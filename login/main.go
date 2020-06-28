package main

import (
	"braid-game/login/handle"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/module/rpc/client/bproto"
	"github.com/pojol/braid/module/rpc/server"
	"github.com/pojol/braid/module/tracer"
)

var (
	help bool

	consulAddr string
	jaegerAddr string

	// NodeName 节点名
	NodeName = "login"
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
		Path:   "/var/log/login",
		Suffex: ".log",
	}, log.WithSys(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/login",
		Suffex: ".sys",
	}))
	defer l.Close()

	tr := tracer.New(NodeName, jaegerAddr)
	tr.Init()

	s := server.New(NodeName, server.WithListen(":1201"), server.WithTracing())
	bproto.RegisterListenServer(server.Get(), &handle.RouteServer{})

	s.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	s.Close()
}
