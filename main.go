package main

import (
	"flag"
	"spider/service"
	"time"
)

func main() {
	startTime := time.Now().AddDate(0, 0, -7)
	startDate := startTime.Format("20060102")
	endDate := time.Now().Format("20060102")

	start := flag.String("start", startDate, "start date")
	end := flag.String("end", endDate, "end date")
	flag.Parse()

	_, errS := time.Parse("20060102", *start)
	if errS != nil {
		panic("开始日期格式错误")
	}

	_, errE := time.Parse("20060102", *end)
	if errE != nil {
		panic("结束日期格式错误")
	}

	new(service.Task).Task(*start, *end)
}
