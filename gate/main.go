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
	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/module/election"
	"github.com/pojol/braid/module/rpc/client"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/plugin/election/k8selector"
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
	var kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	var nodeID = flag.String("node-id", "", "node id used for leader election")

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

	tr := tracer.New(NodeName, jaegerAddr)
	tr.Init()

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

	rpcClient := client.New(NodeName, consulAddr, client.WithTracing())
	rpcClient.Discover()
	defer rpcClient.Close()

	e := echo.New()
	e.Use(middleware.ReqTrace())
	e.Use(middleware.ReqLimit())
	e.POST("/*", routes.PostRouting)

	//go gatemid.Tick()

	err = e.Start(":1202")
	if err != nil {
		log.Fatalf("start echo err", err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	<-ch

	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
