package model

type Stock struct {
	Id        uint32 `gorm:"column:id;type:int(11);PRIMARY_KEY;AUTO_INCREMENT;"`
	Code      uint32 `gorm:"column:code;type:varchar(10);DEFAULT:0;NOT NULL;"`
	Name      string `gorm:"column:name;type:varchar(20);NOT NULL;"`
	Type      uint8  `gorm:"column:type;type:tinyint(1);DEFAULT:0;NOT NULL;"`
	CreateAt  uint32 `gorm:"column:create_at;type:int(11);DEFAULT:0;NOT NULL;"`
	UpdatedAt uint32 `gorm:"column:updated_at;type:int(11);DEFAULT:0;NOT NULL;"`
}

func (s Stock) DBName() string {
	return "stock"
}

func (s Stock) TableName() string {
	return s.DBName() + ".stock"
}
