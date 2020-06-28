package main

import (
	"braid-game/base/handle"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/module/election"
	"github.com/pojol/braid/module/rpc/client"
	"github.com/pojol/braid/module/rpc/client/bproto"
	"github.com/pojol/braid/module/rpc/server"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/plugin/election/consulelection"
)

var (
	help bool

	consulAddr string
	jaegerAddr string

	// NodeName 节点名
	NodeName = "base"
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&consulAddr, "consul", "http://127.0.0.1:8900", "set consul address")
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
		Path:   "/var/log/base",
		Suffex: ".log",
	}, log.WithSys(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/base",
		Suffex: ".sys",
	}))
	defer l.Close()

	elec, err := election.GetBuilder(consulelection.ElectionName).Build(consulelection.Cfg{
		Address:           consulAddr,
		Name:              NodeName,
		LockTick:          time.Second * 2,
		RefushSessionTick: time.Second * 2,
	})
	if err != nil {
		log.Fatalf("elector build err", err)
	}
	elec.Run()

	tr := tracer.New(NodeName, jaegerAddr)
	tr.Init()

	rpcClient := client.New(NodeName, consulAddr, client.WithTracing())
	rpcClient.Discover()
	defer rpcClient.Close()

	s := server.New(NodeName, server.WithListen(":1201"), server.WithTracing())
	bproto.RegisterListenServer(server.Get(), &handle.RouteServer{})

	s.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	elec.Close()
	s.Close()
}
