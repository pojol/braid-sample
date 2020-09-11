module braid-game

go 1.13

require (
	github.com/garyburd/redigo v1.6.2
	github.com/golang/protobuf v1.3.3
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/labstack/echo/v4 v4.1.6
	github.com/pojol/braid v1.1.44
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1
	google.golang.org/grpc v1.27.1
)

// 切换到本地包进行测试
replace github.com/pojol/braid => /Users/pojol/work/gohome/src/braid
