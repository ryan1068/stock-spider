package model

import (
	"github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"
	"log"
	"spider/config"
	"sync"
	"time"
)

var mgorusInstance *logrus.Logger
var singleMgorus = sync.Once{}

func Mgoruser() *logrus.Logger {
	singleMgorus.Do(func() {
		mgorusInstance = newMgorus()
	})

	return mgorusInstance
}

func newMgorus() *logrus.Logger {
	logger := logrus.New()
	logger.WithTime(time.Now().In(time.Local))

	mgoHook, err := mgorus.NewHooker(config.Config.MongoDB.Dsn(), "stock", "stock_logs")
	if err != nil {
		log.Fatalf("logrus 新增mongo hook失败：%s", err)
	} else {
		logger.AddHook(mgoHook)
	}

	return logger
}
