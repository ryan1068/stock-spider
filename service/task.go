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

// 一次并发请求的数量
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

	//分批处理任务，控制并发量，防止目标服务器压力过大
	stocksLen := len(stocks)
	batchTimes := math.Ceil(float64(stocksLen) / float64(RequestLimit))
	batchIndex := 1

	fmt.Printf("总共抓取%d条股票数据，分%d批抓取，每批抓取%d条数据...\n\n", stocksLen, int(batchTimes), RequestLimit)

	for i := 0; i < stocksLen; i += RequestLimit {
		fmt.Printf("第%d批任务开始...\n", batchIndex)

		var sliceStocks []stock.Stock
		if batchIndex == int(batchTimes) {
			// 最后一批任务
			sliceStocks = stocks[i:stocksLen]
		} else {
			sliceStocks = stocks[i : batchIndex*RequestLimit]
		}

		sliceStocksLen := len(sliceStocks)
		spiderChan := make(chan []Result, sliceStocksLen)
		doneChan := make(chan int, sliceStocksLen)
		for _, value := range sliceStocks {
			go t.Do(db, value, spiderChan, startDate, endDate)
		}

	Loop:
		for {
			select {
			case <-ticker.C:
				chanLen := len(doneChan)
				if chanLen == sliceStocksLen {
					fmt.Printf("任务完成100%s，累计耗时：%ds\n\n", "%", time.Now().Unix()-startTime)
					break Loop
				}

				progressString := fmt.Sprintf("%.2f", float64(chanLen)/float64(sliceStocksLen))
				progressFloat, _ := strconv.ParseFloat(progressString, 64)
				progress := int(progressFloat * 100)

				fmt.Printf("已经抓取%d条股票数据, 当前进度为%d%s\n", chanLen, progress, "%")

			case s, ok := <-spiderChan:
				if !ok {
					return
				}
				go t.Store(db, s, doneChan)
			}
		}

		batchIndex++
	}

	fmt.Printf("全部抓取完成，累计耗时：%ds\n\n", time.Now().Unix()-startTime)

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

func (t Task) Store(db *gorm.DB, results []Result, doneChan chan int) error {
	if len(results) == 0 {
		doneChan <- 1
		return fmt.Errorf("results数据为空")
	}

	trends := t.ConvertStockTrend(results)
	var buffer bytes.Buffer
	sql := "insert into " + stock.StockTrend{}.TableName() + " (`code`, `name`, `open_price`, `close_price`, `open_percent`, `close_percent`, `high_percent`, `low_percent`, `shock`, `amount`, `amount_format`, `close_color`, `date`, `created_at`, `updated_at`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	for i, v := range trends {
		if i == len(trends)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, '%s', '%s', '%s', %d, %d);", v.Code, v.Name, v.OpenPrice, v.ClosePrice, v.OpenPercent, v.ClosePercent, v.HighPercent, v.LowPercent, v.Shock, v.Amount, v.AmountFormat, v.CloseColor, v.Date, v.CreatedAt, v.UpdatedAt))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, '%s', '%s', '%s', %d, %d),", v.Code, v.Name, v.OpenPrice, v.ClosePrice, v.OpenPercent, v.ClosePercent, v.HighPercent, v.LowPercent, v.Shock, v.Amount, v.AmountFormat, v.CloseColor, v.Date, v.CreatedAt, v.UpdatedAt))
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
		var closeColor string
		if v.ClosePrice-v.OpenPrice >= 0 {
			closeColor = "收红"
		} else {
			closeColor = "收绿"
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
