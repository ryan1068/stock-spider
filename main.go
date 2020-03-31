package main

import (
	"flag"
	"spider/service"
	"time"
)

func main() {
	startDate := time.Now().AddDate(0, 0, -7).Format("20060102")
	endDate := time.Now().Format("20060102")

	start := flag.String("start", startDate, "start date")
	end := flag.String("end", endDate, "end date")
	flag.Parse()

	_, errS := time.Parse("20060102", *start)
	if errS != nil {
		panic("start date format error")
	}

	_, errE := time.Parse("20060102", *end)
	if errE != nil {
		panic("end date format error")
	}

	task := &service.Task{StartDate: *start, EndDate: *end}
	task.Run()
}
