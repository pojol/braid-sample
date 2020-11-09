package main

import (
	"braid-game/api"
	"braid-game/common"
	"braid-game/mail/handle"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pojol/braid"
	"github.com/pojol/braid/module/tracer"
	"github.com/pojol/braid/modules/electorconsul"
	"github.com/pojol/braid/modules/grpcserver"
	"google.golang.org/grpc"
)

var (
	help bool

	consulAddr string
	jaegerAddr string
	localPort  int

	// NodeName 节点名
	NodeName = "mail"
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&consulAddr, "consul", "http://127.0.0.1:8500", "set consul address")
	flag.StringVar(&jaegerAddr, "jaeger", "http://127.0.0.1:9411/api/v2/spans", "set jaeger address")
	flag.IntVar(&localPort, "localPort", 0, "run locally")
}

func main() {
	initFlag()

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	b, err := braid.New(NodeName)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	var rpcserver braid.Module
	if localPort == 0 {
		rpcserver = braid.GRPCServer(grpcserver.Name)
	} else {
		addr := ":" + strconv.Itoa(localPort)
		rpcserver = braid.GRPCServer(grpcserver.Name, grpcserver.WithListen(addr))

		id := strconv.Itoa(int(time.Now().UnixNano())) + addr
		err := common.Regist(common.ConsulRegistReq{
			Name:    NodeName,
			ID:      id,
			Tags:    []string{"braid", NodeName},
			Address: "127.0.0.1",
			Port:    localPort,
		}, consulAddr)
		if err != nil {
			panic(err.Error())
		}

		defer common.Deregist(id, consulAddr)
	}

	b.RegistModule(
		rpcserver,
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
