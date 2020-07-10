package main

import (
	"braid-game/base/handle"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/pojol/braid"
	"github.com/pojol/braid/3rd/log"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/plugin/rpc/grpcclient/bproto"
	"github.com/pojol/braid/plugin/rpc/grpcserver"
	"google.golang.org/grpc"
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

	tr, err := tracer.New(NodeName, jaegerAddr)
	if err != nil {
		log.Fatalf("tracer init", err)
	}

	b := braid.New(NodeName)
	b.RegistPlugin(braid.DiscoverByConsul(consulAddr),
		braid.BalancerBySwrr(),
		braid.GRPCClient(),
		braid.GRPCServer(grpcserver.WithListen(":14222"), grpcserver.WithTracing()),
		braid.ElectorByConsul())

	bproto.RegisterListenServer(braid.Server().Server().(*grpc.Server), &handle.RouteServer{})

	b.Run()
	defer b.Close()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	tr.Close()
}
