package main

type SearchResult struct {
	Date        string `json:"日期"`
	Code        string `json:"股票代码"`
	Name        string `json:"名称"`
	ClosePrice  string `json:"收盘价"`
	HighPrice   string `json:"最高价"`
	LowPrice    string `json:"最低价"`
	OpenPrice   string `json:"开盘价"`
	LastPrice   string `json:"前收盘"`
	Quota       string `json:"涨跌额"`
	Percent     string `json:"涨跌幅"`
	Rate        string `json:"换手率"`
	Amount      string `json:"成交量"`
	MoneyAmount string `json:"成交金额"`
	TotalValue  string `json:"总市值"`
	MarketValue string `json:"流通市值"`
}

func main() {

}
