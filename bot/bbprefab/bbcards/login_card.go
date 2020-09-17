package bbcards

import (
	"braid-game/bot/bbprefab"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// GuestLoginCard 游客登录
type GuestLoginCard struct {
	URL   string
	delay time.Duration
	md    *bbprefab.BotData
}

// GuestLoginRes 游客登录返回
type GuestLoginRes struct {
	Token string
}

// NewGuestLoginCard 生成账号创建预制
func NewGuestLoginCard(md *bbprefab.BotData) *GuestLoginCard {
	return &GuestLoginCard{
		URL:   "/v1/login/guest",
		delay: time.Millisecond,
		md:    md,
	}
}

// GetURL 获取服务器地址
func (card *GuestLoginCard) GetURL() string { return card.URL }

// GetHeader get card header
func (card *GuestLoginCard) GetHeader() map[string]string { return nil }

// SetDelay 设置卡片之间调用的延迟
func (card *GuestLoginCard) SetDelay(delay time.Duration) { card.delay = delay }

// GetDelay 获取卡片之间调用的延迟
func (card *GuestLoginCard) GetDelay() time.Duration { return card.delay }

// Marshal 序列化传入消息体
func (card *GuestLoginCard) Marshal() []byte {

	b := []byte{}

	return b
}

// Unmarshal 反序列化返回消息
func (card *GuestLoginCard) Unmarshal(res *http.Response) {

	errcode, _ := strconv.Atoi(res.Header["Errcode"][0])
	if errcode != 0 {
		fmt.Println(card.GetURL(), "request err", errcode)
	}

	cres := GuestLoginRes{}
	b, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal(b, &cres)
	if err != nil {
		fmt.Println(card.GetURL(), "json.Unmarshal", errcode)
	}

	card.md.AccToken = cres.Token
}
