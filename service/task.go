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

type Task struct {
	StartDate     string
	EndDate       string
	ExecStartTime int64
}

var RequestLimit = 300

func (t *Task) Run() {
	t.ExecStartTime = time.Now().Unix()
	db, err := model.GormOpenDB()
	if err != nil {
		return
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	defer db.Close()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	t.Task(db, ticker)

	return
}

func (t *Task) Task(db *gorm.DB, ticker *time.Ticker) {
	sql := fmt.Sprintf("TRUNCATE TABLE %s", stock.StockTrend{}.TableName())
	if err := db.Exec(sql).Error; err != nil {
		return
	}

	stocks, _ := new(stock.Stock).FindAll()
	stocksLen := len(stocks)
	taskNum := math.Ceil(float64(stocksLen) / float64(RequestLimit))
	taskIndex := 1

	fmt.Printf("Total requests:%d，Total tasks:%d，Each task handing requests:%d...\n\n", stocksLen, int(taskNum), RequestLimit)

	for i := 0; i < stocksLen; i += RequestLimit {
		fmt.Printf("Task %d start...\n", taskIndex)

		var sliceStocks []stock.Stock
		if taskIndex == int(taskNum) {
			// last task
			sliceStocks = stocks[i:stocksLen]
		} else {
			sliceStocks = stocks[i : taskIndex*RequestLimit]
		}

		sliceStocksLen := len(sliceStocks)
		spiderChan := make(chan []Result, sliceStocksLen)
		doneChan := make(chan int, sliceStocksLen)
		go t.Store(db, spiderChan, doneChan)

		for _, value := range sliceStocks {
			go t.Do(db, value, spiderChan)
		}

	Loop:
		for {
			select {
			case <-ticker.C:
				chanLen := len(doneChan)
				if chanLen == sliceStocksLen {
					fmt.Printf("Task %d completed,Total spend:%ds\n\n", taskIndex, time.Now().Unix()-t.ExecStartTime)
					break Loop
				}

				progressString := fmt.Sprintf("%.2f", float64(chanLen)/float64(sliceStocksLen))
				progressFloat, _ := strconv.ParseFloat(progressString, 64)
				progress := int(progressFloat * 100)

				fmt.Printf("Handing requests %d,Current progress:%d%s\n", chanLen, progress, "%")
			}
		}

		taskIndex++

		close(spiderChan)
		close(doneChan)
	}

	fmt.Printf("All task done,Total spend:%ds\n\n", time.Now().Unix()-t.ExecStartTime)

	return
}

func (t Task) Store(db *gorm.DB, spiderChan chan []Result, doneChan chan int) {
	for {
		select {
		case s, ok := <-spiderChan:
			if !ok {
				return
			}
			go t.BulkStorage(db, s, doneChan)
		}
	}
}

func (t Task) Do(db *gorm.DB, s stock.Stock, spiderChan chan []Result) {
	results, _ := Result{}.Request(s.Code, t.StartDate, t.EndDate)
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

func (t Task) BulkStorage(db *gorm.DB, results []Result, doneChan chan int) error {
	if len(results) == 0 {
		doneChan <- 1
		return fmt.Errorf("results empty")
	}

	trends := t.ConvertStockTrend(results)
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT INTO %s (`code`, `name`, `open_price`, `close_price`, `open_percent`, `close_percent`, `high_percent`, `low_percent`, `shock`, `amount`, `amount_format`, `close_color`, `date`, `created_at`, `updated_at`) VALUES", stock.StockTrend{}.TableName())
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
