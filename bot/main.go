package main

import (
	"braid-game/bot/bbprefab"
	"braid-game/bot/bbprefab/bbsteps"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pojol/gobot"
	"github.com/pojol/gobot/botreport"
)

var (
	help bool

	// target server addr
	target string

	// Strategy
	strategy string

	// robot number
	num int
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&target, "target", "http://localhost:14001", "set target server address")
	flag.StringVar(&strategy, "strategy", "login", "set strategy `login` & `loginout` ")
	flag.IntVar(&num, "num", 1, "robot number")
}

/*
        +---->login1+--+--->mail1
        |              |
        +---->login2+--+--->mail2
gate1+->+
        |
        +---->base1
        |
        +---->base2
*/

var tokenm map[string]int

func main() {

	initFlag()

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	tokenm = make(map[string]int)

	rand.Seed(time.Now().UnixNano())
	ports := []string{"14001", "14002" /*, "14003"*/}

	var bots []*gobot.Bot

	for i := 0; i < num; i++ {

		md := &bbprefab.BotData{}
		bot := gobot.New(gobot.BotConfig{
			Addr:   "http://localhost:" + ports[rand.Intn(len(ports))],
			Report: false,
		}, md)

		bot.Timeline.AddStep(bbsteps.NewGuestLoginStep(md))
		bot.Timeline.AddLoopStep(bbsteps.NewRenameStep(md))

		bot.Run()

		bots = append(bots, bot)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				minfo := make(map[string][]botreport.Info)
				for _, v := range bots {
					r := v.GetReport()
					for rk := range r.Info {
						minfo[rk] = append(minfo[rk], r.Info[rk]...)
					}

					r.Clear()
				}

				for k := range minfo {
					t, err := botreport.GetAverageTime(minfo[k])
					if err != nil {
						continue
					}

					rate := botreport.GetSuccRate(minfo[k])

					fmt.Printf("%-30s Req count %-5d Average time %-5s Succ rate %-10s \n", k, len(minfo[k]), strconv.Itoa(t)+"ms", rate)
				}

			default:
			}
		}

	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	for _, v := range bots {
		v.Close()
	}
	time.Sleep(time.Second)
}
