package main

import (
	bm "braid-game/gate/middleware"
	"braid-game/gate/routes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/pojol/braid"
	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/3rd/redis"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/plugin/balancerswrr"
	"github.com/pojol/braid/plugin/discoverconsul"
	"github.com/pojol/braid/plugin/electorconsul"
	"github.com/pojol/braid/plugin/linkerredis"
	"github.com/pojol/braid/plugin/pubsubnsq"
)

var (
	help bool

	consulAddr    string
	redisAddr     string
	jaegerAddr    string
	nsqLookupAddr string
	nsqdAddr      string
)

const (
	// NodeName box 节点名
	NodeName = "gateway"
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&consulAddr, "consul", "http://127.0.0.1:8900", "set consul address")
	flag.StringVar(&redisAddr, "redis", "redis://127.0.0.1:6379/0", "set redis address")
	flag.StringVar(&jaegerAddr, "jaeger", "http://127.0.0.1:9411/api/v2/spans", "set jaeger address")
	flag.StringVar(&nsqLookupAddr, "nsqlookup", "127.0.0.1:4161", "set nsq lookup address")
	flag.StringVar(&nsqdAddr, "nsqd", "127.0.0.1:4150", "set nsqd address")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			stack := string(debug.Stack())
			fmt.Printf("stack: %v\n", stack)
			fmt.Println(fmt.Errorf("error: %v", err).Error())
		}
	}()

	initFlag()
	//var kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//var nodeID = flag.String("node-id", "", "node id used for leader election")

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	log.New(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/gateway",
		Suffex: ".log",
	}, log.WithSys(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/gateway",
		Suffex: ".sys",
	}))

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
		braid.Elector(
			electorconsul.Name,
			electorconsul.WithConsulAddr(consulAddr),
		),
		braid.LinkCache(linkerredis.Name),
		braid.Pubsub(
			pubsubnsq.Name,
			pubsubnsq.WithLookupAddr([]string{nsqLookupAddr}),
			pubsubnsq.WithNsqdAddr([]string{nsqdAddr})),
		braid.JaegerTracing(tracer.WithHTTP(jaegerAddr), tracer.WithProbabilistic(0.01)))

	b.Run()
	defer b.Close()

	e := echo.New()
	e.Use(bm.ReqTrace())
	e.Use(bm.ReqLimit())
	e.POST("/*", routes.PostRouting)

	//go gatemid.Tick()
	/*
		go func() {
			fmt.Println(http.ListenAndServe(":6060", nil))
		}()
	*/
	err = e.Start(":14222")
	if err != nil {
		log.Fatalf("start echo err", err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	<-ch

	rc.Close()
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
