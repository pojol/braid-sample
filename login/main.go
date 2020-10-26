package main

import (
	"braid-game/login/handle"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pojol/braid"
	"github.com/pojol/braid/3rd/redis"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/plugin/balancerswrr"
	"github.com/pojol/braid/plugin/discoverconsul"
	"github.com/pojol/braid/plugin/electorconsul"
	"github.com/pojol/braid/plugin/grpcclient"
	"github.com/pojol/braid/plugin/grpcclient/bproto"
	"github.com/pojol/braid/plugin/grpcserver"
	"github.com/pojol/braid/plugin/linkerredis"
	"github.com/pojol/braid/plugin/mailboxnsq"
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
		log.Fatalf("redis init %s\n", err)
	}

	b, _ := braid.New(
		NodeName,
		mailboxnsq.WithLookupAddr([]string{nsqLookupAddr}),
		mailboxnsq.WithNsqdAddr([]string{nsqdAddr}))

	b.RegistPlugin(
		braid.GRPCServer(
			grpcserver.Name,
			grpcserver.WithListen(":14222"),
		),
		braid.Discover(
			discoverconsul.Name,
			discoverconsul.WithConsulAddr(consulAddr),
			discoverconsul.WithBlacklist([]string{"gateway"})),
		braid.Balancer(balancerswrr.Name),
		braid.GRPCClient(grpcclient.Name),
		braid.Elector(
			electorconsul.Name,
			electorconsul.WithConsulAddr(consulAddr),
		),
		braid.LinkCache(linkerredis.Name, linkerredis.WithRedisAddr(redisAddr)),
		braid.JaegerTracing(tracer.WithHTTP(jaegerAddr), tracer.WithProbabilistic(0.01)))

	bproto.RegisterListenServer(braid.Server().(*grpc.Server), &handle.RouteServer{})

	b.Run()
	defer b.Close()
	defer rc.Close()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
