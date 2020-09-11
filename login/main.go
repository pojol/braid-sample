package main

import (
	"braid-game/login/handle"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pojol/braid"
	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/3rd/redis"
	"github.com/pojol/braid/plugin/discoverconsul"
	"github.com/pojol/braid/plugin/grpcclient/bproto"
	"github.com/pojol/braid/plugin/grpcserver"
	"google.golang.org/grpc"
)

var (
	help bool

	consulAddr    string
	jaegerAddr    string
	redisAddr     string
	nsqLookupAddr string
	nsqdAddr      string

	// NodeName 节点名
	NodeName = "login"
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&consulAddr, "consul", "http://127.0.0.1:8900", "set consul address")
	flag.StringVar(&nsqLookupAddr, "nsqlookup", "127.0.0.1:4161", "set nsq lookup address")
	flag.StringVar(&nsqdAddr, "nsqd", "127.0.0.1:4150", "set nsqd address")
	flag.StringVar(&redisAddr, "redis", "redis://127.0.0.1:6379/0", "set redis address")
	flag.StringVar(&jaegerAddr, "jaeger", "http://127.0.0.1:9411/api/v2/spans", "set jaeger address")
}

func main() {
	initFlag()

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	l := log.New(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/login",
		Suffex: ".log",
	}, log.WithSys(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/login",
		Suffex: ".sys",
	}))
	defer l.Close()

	rc := redis.New()
	err := rc.Init(redis.Config{
		Address:        "redis://192.168.50.201:6379/0",
		ReadTimeOut:    5 * time.Second,
		WriteTimeOut:   5 * time.Second,
		ConnectTimeOut: 2 * time.Second,
		MaxIdle:        16,
		MaxActive:      128,
		IdleTimeout:    0,
	})
	if err != nil {
		log.Fatalf("redis init", err)
	}

	b := braid.New(NodeName)
	b.RegistPlugin(braid.GRPCServer(grpcserver.WithListen(":14222")),
		braid.DiscoverByConsul(consulAddr, discoverconsul.WithBlacklist([]string{"gateway"})),
		braid.BalancerBySwrr(),
		braid.GRPCClient(),
		braid.ElectorByConsul(consulAddr),
		braid.PubsubByNsq([]string{nsqLookupAddr}, []string{nsqdAddr}),
		braid.LinkerByRedis(),
		braid.JaegerTracing(jaegerAddr))

	bproto.RegisterListenServer(braid.Server().Server().(*grpc.Server), &handle.RouteServer{})

	b.Run()
	defer b.Close()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
