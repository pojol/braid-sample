package main

import (
	"braid-game/gate/middleware"
	"braid-game/gate/routes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pojol/braid"
	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/3rd/redis"
	"github.com/pojol/braid/module/tracer"
)

var (
	help bool

	consulAddr string
	redisAddr  string
	jaegerAddr string
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

	tr, err := tracer.New(NodeName, jaegerAddr)
	if err != nil {
		log.Fatalf("tracer init", err)
	}

	rc := redis.New()
	err = rc.Init(redis.Config{
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

	/*
		elec, err := election.GetBuilder(k8selector.ElectionName).Build(k8selector.Cfg{
			KubeCfg:     *kubeconfig,
			NodID:       *nodeID,
			Namespace:   "default",
			RetryPeriod: time.Second * 2,
		})
		if err != nil {
			log.Fatalf("elector build err", err)
		}

		elec.Run()
	*/

	b := braid.New(NodeName)
	b.RegistPlugin(braid.DiscoverByConsul(consulAddr),
		braid.BalancerBySwrr(),
		braid.GRPCClient())

	b.Run()
	defer b.Close()

	e := echo.New()
	e.Use(middleware.ReqTrace())
	e.Use(middleware.ReqLimit())
	e.POST("/*", routes.PostRouting)

	//go gatemid.Tick()

	err = e.Start(":14222")
	if err != nil {
		log.Fatalf("start echo err", err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	<-ch

	tr.Close()
	rc.Close()
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
