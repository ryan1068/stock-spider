package model

type StockTrend struct {
	Id           uint32 `gorm:"column:id;type:int(11);PRIMARY_KEY;AUTO_INCREMENT;"`
	Code         uint32 `gorm:"column:code;type:varchar(10);DEFAULT:0;NOT NULL;"`
	Name         string `gorm:"column:name;type:varchar(20);NOT NULL;"`
	OpenPrice    uint8  `gorm:"column:open_price;type:float;DEFAULT:0;NOT NULL;"`
	ClosePrice   uint32 `gorm:"column:close_price;type:float;DEFAULT:0;NOT NULL;"`
	OpenPercent  uint32 `gorm:"column:open_percent;type:float;DEFAULT:0;NOT NULL;"`
	ClosePercent uint32 `gorm:"column:close_percent;type:float;DEFAULT:0;NOT NULL;"`
	HighPercent  string `gorm:"column:high_percent;type:float;NOT NULL;"`
	LowPercent   uint8  `gorm:"column:low_percent;type:float;DEFAULT:0;NOT NULL;"`
	Shock        uint32 `gorm:"column:shock;type:float;DEFAULT:0;NOT NULL;"`
	Amount       uint32 `gorm:"column:amount;type:varchar(10);DEFAULT:0;NOT NULL;"`
	AmountFormat string `gorm:"column:amount_format;type:varchar(20);NOT NULL;"`
	CloseColor   uint8  `gorm:"column:close_color;type:tinyint(1);NOT NULL;"`
	CreateAt     uint32 `gorm:"column:create_at;type:int(11);DEFAULT:0;NOT NULL;"`
	UpdatedAt    uint32 `gorm:"column:updated_at;type:int(11);DEFAULT:0;NOT NULL;"`
}

func (s StockTrend) DBName() string {
	return "stock"
}

func (s StockTrend) TableName() string {
	return s.DBName() + ".stock_trend"
}
