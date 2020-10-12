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

	// increase
	increase bool

	lifetime int
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&target, "target", "http://localhost:14001", "set target server address")
	flag.StringVar(&strategy, "strategy", "login", "set strategy `login` & `loginout` ")
	flag.IntVar(&num, "num", 1, "robot number")
	flag.BoolVar(&increase, "increase", false, "incremental robot in every second")
	flag.IntVar(&lifetime, "lifetime", 60, "life time by second")
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

func report(info map[string][]botreport.Info) {
	for k := range info {
		t, err := botreport.GetAverageTime(info[k])
		if err != nil {
			continue
		}

		rate := botreport.GetSuccRate(info[k])

		fmt.Printf("%-30s Req count %-5d Average time %-5s Succ rate %-10s \n", k, len(info[k]), strconv.Itoa(t)+"ms", rate)
	}
}

func reportByMap(m map[string]*gobot.Bot) {
	minfo := make(map[string][]botreport.Info)
	for _, v := range m {
		info := v.GetReportInfo()
		for rk := range info {
			minfo[rk] = append(minfo[rk], info[rk]...)
		}

		v.ClearReportInfo()
	}
	report(minfo)
}

func reportByArr(arr []*gobot.Bot) {
	minfo := make(map[string][]botreport.Info)
	for _, v := range arr {
		info := v.GetReportInfo()
		for rk := range info {
			minfo[rk] = append(minfo[rk], info[rk]...)
		}

		v.ClearReportInfo()
	}
	report(minfo)
}

func staticBot(ports []string, ch chan os.Signal) {

	var bots []*gobot.Bot

	for i := 0; i < num; i++ {

		bot := createBot(ports[rand.Intn(len(ports))])
		bot.Run()

		bots = append(bots, bot)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				reportByArr(bots)
			default:
			}
		}

	}()

	<-ch
	fmt.Println("closed")
	for _, bot := range bots {
		bot.Close()
	}
}

func createBot(port string) *gobot.Bot {
	md := &bbprefab.BotData{}
	bot := gobot.New(gobot.BotConfig{
		Addr:   "http://localhost:" + port,
		Report: false,
	}, md)

	bot.Timeline.AddStep(bbsteps.NewGuestLoginStep(md))
	bot.Timeline.AddStep(bbsteps.NewLoginOutStep(md))
	bot.Timeline.AddLoopStep(bbsteps.NewRenameStep(md))

	return bot
}

func increaseBot(ports []string, ch chan os.Signal) {

	botm := make(map[string]*gobot.Bot)
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:

				for i := 0; i < num; i++ {
					bot := createBot(ports[rand.Intn(len(ports))])
					bot.Run()
					time.AfterFunc(time.Duration(lifetime)*time.Second, func() {
						bot.Close()

						reportByMap(botm)

						delete(botm, bot.ID())
					})

					botm[bot.ID()] = bot
				}
			default:
			}
		}
	}()

	<-ch
	fmt.Println("closed")
}

func main() {

	initFlag()

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	fmt.Println("===== run bot =====")
	fmt.Println("target", target)
	fmt.Println("num", num)
	fmt.Println("increase", increase)
	fmt.Println("lifetime", lifetime)

	rand.Seed(time.Now().UnixNano())
	ports := []string{"14001", "14002" /*, "14003"*/}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	if increase {
		increaseBot(ports, ch)
	} else {
		staticBot(ports, ch)
	}

	time.Sleep(time.Second)
}
