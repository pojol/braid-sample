package main

import (
	"braid-game/api"
	"braid-game/mail/handle"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/pojol/braid"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/plugin/electorconsul"
	"github.com/pojol/braid/plugin/grpcserver"
	"google.golang.org/grpc"
)

var (
	help bool

	consulAddr string
	jaegerAddr string

	// NodeName 节点名
	NodeName = "mail"
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

	b, _ := braid.New(NodeName)

	b.RegistPlugin(
		braid.GRPCServer(
			grpcserver.Name,
			grpcserver.WithListen(":14222"),
		),
		braid.Elector(
			electorconsul.Name,
			electorconsul.WithConsulAddr(consulAddr),
		),
		braid.JaegerTracing(tracer.WithHTTP(jaegerAddr), tracer.WithProbabilistic(0.01)))

	api.RegisterMailServer(braid.Server().(*grpc.Server), &handle.MailServer{})

	b.Init()
	b.Run()
	defer b.Close()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

}
