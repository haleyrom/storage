module github.com/haleyrom/storage

go 1.14

replace github.com/go-xorm/core v0.6.3 => xorm.io/core v0.6.3

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/builder v0.3.4
	github.com/go-xorm/core v0.6.3
	github.com/go-xorm/xorm v0.7.9
	go.uber.org/zap v1.15.0
	xorm.io/builder v0.3.7 // indirect
	xorm.io/core v0.7.3 // indirect
)
