package service

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leekchan/accounting"
	"math"
	"spider/model"
	"spider/model/stock"
	"strconv"
	"time"
)

type Task struct{}

var RequestLimit = 300

func (t Task) Task(startDate string, endDate string) {
	startTime := time.Now().Unix()
	stocks, _ := stock.Stock{}.FindAll()

	db, err := model.GormOpenDB()
	if err != nil {
		return
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	defer db.Close()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	stocksLen := len(stocks)
	taskNum := math.Ceil(float64(stocksLen) / float64(RequestLimit))
	batchIndex := 1

	fmt.Printf("Total requests:%d，Total tasks:%d，Each task handing requests:%d...\n\n", stocksLen, int(taskNum), RequestLimit)

	for i := 0; i < stocksLen; i += RequestLimit {
		fmt.Printf("Task %d start...\n", batchIndex)

		var sliceStocks []stock.Stock
		if batchIndex == int(taskNum) {
			// last task
			sliceStocks = stocks[i:stocksLen]
		} else {
			sliceStocks = stocks[i : batchIndex*RequestLimit]
		}

		sliceStocksLen := len(sliceStocks)
		spiderChan := make(chan []Result, sliceStocksLen)
		doneChan := make(chan int, sliceStocksLen)
		go t.Store(db, spiderChan, doneChan)
		for _, value := range sliceStocks {
			go t.Do(db, value, spiderChan, startDate, endDate)
		}

	Loop:
		for {
			select {
			case <-ticker.C:
				chanLen := len(doneChan)
				if chanLen == sliceStocksLen {
					fmt.Printf("Task %d completed,Total spend:%ds\n\n", batchIndex, time.Now().Unix()-startTime)
					break Loop
				}

				progressString := fmt.Sprintf("%.2f", float64(chanLen)/float64(sliceStocksLen))
				progressFloat, _ := strconv.ParseFloat(progressString, 64)
				progress := int(progressFloat * 100)

				fmt.Printf("Handing requests %d,Current progress:%d%s\n", chanLen, progress, "%")
			}
		}

		batchIndex++
	}

	fmt.Printf("All task done,Total spend:%ds\n\n", time.Now().Unix()-startTime)

	return
}

func (t Task) Do(db *gorm.DB, s stock.Stock, spiderChan chan []Result, startDate string, endDate string) {
	results, _ := Result{}.Request(s.Code, startDate, endDate)
	spiderChan <- results
	return
}

func (t Task) CalculateOpenPercent(openPrice float64, lastPrice float64) float64 {
	diff := openPrice - lastPrice
	openPercent := fmt.Sprintf("%.3f", diff/lastPrice)
	return ConvertFloat64(openPercent) * 100
}

func (t Task) CalculateHighPercent(highPrice float64, lastPrice float64) float64 {
	diff := highPrice - lastPrice
	highPercent := fmt.Sprintf("%.3f", diff/lastPrice)
	return ConvertFloat64(highPercent) * 100
}

func (t Task) CalculateLowPercent(lowPrice float64, lastPrice float64) float64 {
	diff := lowPrice - lastPrice
	lowPercent := fmt.Sprintf("%.3f", diff/lastPrice)
	return ConvertFloat64(lowPercent) * 100
}

func (t Task) Store(db *gorm.DB, spiderChan chan []Result, doneChan chan int) {
	for {
		select {
		case s, ok := <-spiderChan:
			if !ok {
				return
			}
			go t.BatchSave(db, s, doneChan)
		}
	}
}

func (t Task) BatchSave(db *gorm.DB, results []Result, doneChan chan int) error {
	if len(results) == 0 {
		doneChan <- 1
		return fmt.Errorf("results empty")
	}

	trends := t.ConvertStockTrend(results)
	var buffer bytes.Buffer
	sql := "insert into " + stock.StockTrend{}.TableName() + " (`code`, `name`, `open_price`, `close_price`, `open_percent`, `close_percent`, `high_percent`, `low_percent`, `shock`, `amount`, `amount_format`, `close_color`, `date`, `created_at`, `updated_at`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	for i, v := range trends {
		if i == len(trends)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, '%s', '%d', '%s', %d, %d);", v.Code, v.Name, v.OpenPrice, v.ClosePrice, v.OpenPercent, v.ClosePercent, v.HighPercent, v.LowPercent, v.Shock, v.Amount, v.AmountFormat, v.CloseColor, v.Date, v.CreatedAt, v.UpdatedAt))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, '%s', '%d', '%s', %d, %d),", v.Code, v.Name, v.OpenPrice, v.ClosePrice, v.OpenPercent, v.ClosePercent, v.HighPercent, v.LowPercent, v.Shock, v.Amount, v.AmountFormat, v.CloseColor, v.Date, v.CreatedAt, v.UpdatedAt))
		}
	}

	if err := db.Exec(buffer.String()).Error; err != nil {
		return err
	}

	doneChan <- 1

	return nil
}

func (t Task) ConvertStockTrend(results []Result) []stock.StockTrend {
	var stockTrends []stock.StockTrend
	for _, v := range results {
		highPercent := t.CalculateHighPercent(v.HighPrice, v.LastPrice)
		lowPercent := t.CalculateLowPercent(v.LowPrice, v.LastPrice)
		shock := highPercent - lowPercent
		var closeColor int8
		if v.ClosePrice-v.OpenPrice >= 0 {
			closeColor = 1
		} else {
			closeColor = 2
		}
		stockTrend := stock.StockTrend{
			Code:         v.Code[1:],
			Name:         v.Name,
			OpenPrice:    v.OpenPrice,
			ClosePrice:   v.ClosePrice,
			OpenPercent:  t.CalculateOpenPercent(v.OpenPrice, v.LastPrice),
			ClosePercent: ConvertFloat64(fmt.Sprintf("%.1f", v.Percent)),
			HighPercent:  highPercent,
			LowPercent:   lowPercent,
			Shock:        shock,
			Amount:       v.MoneyAmount,
			AmountFormat: new(accounting.Accounting).FormatMoney(v.MoneyAmount),
			CloseColor:   closeColor,
			Date:         v.Date,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		stockTrends = append(stockTrends, stockTrend)
	}

	return stockTrends
}
