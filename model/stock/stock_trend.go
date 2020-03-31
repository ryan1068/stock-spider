package stock

type StockTrend struct {
	Id           uint    `gorm:"column:id;type:int(11);PRIMARY_KEY;AUTO_INCREMENT;"`
	Code         string  `gorm:"column:code;type:varchar(10);DEFAULT:0;NOT NULL;"`
	Name         string  `gorm:"column:name;type:varchar(20);NOT NULL;"`
	OpenPrice    float64 `gorm:"column:open_price;type:float;DEFAULT:0;NOT NULL;"`
	ClosePrice   float64 `gorm:"column:close_price;type:float;DEFAULT:0;NOT NULL;"`
	OpenPercent  float64 `gorm:"column:open_percent;type:float;DEFAULT:0;NOT NULL;"`
	ClosePercent float64 `gorm:"column:close_percent;type:float;DEFAULT:0;NOT NULL;"`
	HighPercent  float64 `gorm:"column:high_percent;type:float;NOT NULL;"`
	LowPercent   float64 `gorm:"column:low_percent;type:float;DEFAULT:0;NOT NULL;"`
	Shock        float64 `gorm:"column:shock;type:float;DEFAULT:0;NOT NULL;"`
	Amount       float64 `gorm:"column:amount;type:float;DEFAULT:0;NOT NULL;"`
	AmountFormat string  `gorm:"column:amount_format;type:varchar(20);NOT NULL;"`
	CloseColor   int8    `gorm:"column:close_color;type:tinyint(1);NOT NULL;"`
	Date         string  `gorm:"column:date;type:varchar(20);NOT NULL;"`
	CreatedAt    int64   `gorm:"column:created_at;type:int(11);DEFAULT:0;NOT NULL;"`
	UpdatedAt    int64   `gorm:"column:updated_at;type:int(11);DEFAULT:0;NOT NULL;"`
}

func (s StockTrend) TableName() string {
	return "stock_trend"
}
