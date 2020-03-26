<p>本项目是用go语言写的抓取股票数据爬虫。</p>

## 特点

1. 使用go协程、连接池特性，30s秒内抓取所有股票历史数据
2. 程序执行时进度条显示
3. 可以自定义日期查询股票数据

## 使用方式：

```go
//默认抓取近七天数据
go run main.go

//抓取自义定时间段数据
go run main.go -start=20200323 -end=20200324
```
## 运行效果

<img src="http://blog.herozw.com/wp-content/uploads/2020/03/20200326161150_64608.png" height="750" />
<img src="http://blog.herozw.com/wp-content/uploads/2020/03/20200326160419_27952.png"  />

