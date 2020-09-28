package main

import (
	"braid-game/base/handle"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pojol/braid"
	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/3rd/redis"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/plugin/balancerswrr"
	"github.com/pojol/braid/plugin/discoverconsul"
	"github.com/pojol/braid/plugin/electorconsul"
	"github.com/pojol/braid/plugin/grpcclient/bproto"
	"github.com/pojol/braid/plugin/grpcserver"
	"github.com/pojol/braid/plugin/linkerredis"
	"github.com/pojol/braid/plugin/pubsubnsq"
	"google.golang.org/grpc"
)

var (
	help bool

	consulAddr    string
	jaegerAddr    string
	nsqLookupAddr string
	nsqdAddr      string
	redisAddr     string

	// NodeName 节点名
	NodeName = "base"
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&consulAddr, "consul", "http://127.0.0.1:8900", "set consul address")
	flag.StringVar(&jaegerAddr, "jaeger", "http://127.0.0.1:9411/api/v2/spans", "set jaeger address")
	flag.StringVar(&nsqLookupAddr, "nsqlookup", "127.0.0.1:4161", "set nsq lookup address")
	flag.StringVar(&nsqdAddr, "nsqd", "127.0.0.1:4150", "set nsqd address")
	flag.StringVar(&redisAddr, "redis", "redis://127.0.0.1:6379/0", "set redis address")

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
		Path:   "/var/log/base",
		Suffex: ".log",
	}, log.WithSys(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/base",
		Suffex: ".sys",
	}))
	defer l.Close()

	rc := redis.New()
	err := rc.Init(redis.Config{
		Address:        redisAddr,
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
	b.RegistPlugin(
		braid.Discover(
			discoverconsul.Name,
			discoverconsul.WithConsulAddr(consulAddr)),
		braid.Balancer(balancerswrr.Name),
		braid.GRPCClient(),
		braid.GRPCServer(grpcserver.WithListen(":14222")),
		braid.Elector(
			electorconsul.Name,
			electorconsul.WithConsulAddr(consulAddr),
		),
		braid.Pubsub(
			pubsubnsq.Name,
			pubsubnsq.WithLookupAddr([]string{nsqLookupAddr}),
			pubsubnsq.WithNsqdAddr([]string{nsqdAddr})),
		braid.LinkCache(linkerredis.Name),
		braid.JaegerTracing(tracer.WithHTTP(jaegerAddr), tracer.WithProbabilistic(0.01)))

	bproto.RegisterListenServer(braid.Server().Server().(*grpc.Server), &handle.RouteServer{})

	b.Run()
	defer b.Close()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
