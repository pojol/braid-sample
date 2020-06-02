package main

import (
	"braid-game/api"
	"braid-game/mail/handle"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/pojol/braid/log"
	"github.com/pojol/braid/rpc/server"
	"github.com/pojol/braid/service/election"
	"github.com/pojol/braid/tracer"
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

	l := log.New(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/mail",
		Suffex: ".log",
	}, log.WithSys(log.Config{
		Mode:   log.DebugMode,
		Path:   "/var/log/mail",
		Suffex: ".sys",
	}))
	defer l.Close()

	elec, err := election.New(NodeName, consulAddr)
	if err != nil {
		log.Fatalf(err.Error())
	}
	elec.Run()

	tr := tracer.New(NodeName, jaegerAddr)
	tr.Init()

	s := server.New(NodeName, server.WithListen(":1201"), server.WithTracing())
	api.RegisterMailServer(server.Get(), &handle.MailServer{})

	s.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	elec.Close()
	s.Close()
}
