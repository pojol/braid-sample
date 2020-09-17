package bbcards

import (
	"braid-game/bot/bbprefab"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// LoginOutCard 游客登录
type LoginOutCard struct {
	URL   string
	delay time.Duration
	md    *bbprefab.BotData
}

// NewLoginOutCard 生成账号创建预制
func NewLoginOutCard(md *bbprefab.BotData) *LoginOutCard {
	return &LoginOutCard{
		URL:   "/v1/login/out",
		delay: time.Millisecond,
		md:    md,
	}
}

// GetURL 获取服务器地址
func (card *LoginOutCard) GetURL() string { return card.URL }

// GetHeader get card header
func (card *LoginOutCard) GetHeader() map[string]string {
	return map[string]string{
		"token": card.md.AccToken,
	}
}

// SetDelay 设置卡片之间调用的延迟
func (card *LoginOutCard) SetDelay(delay time.Duration) { card.delay = delay }

// GetDelay 获取卡片之间调用的延迟
func (card *LoginOutCard) GetDelay() time.Duration { return card.delay }

// Marshal 序列化传入消息体
func (card *LoginOutCard) Marshal() []byte {

	return []byte{}
}

// Unmarshal 反序列化返回消息
func (card *LoginOutCard) Unmarshal(res *http.Response) {

	errcode, _ := strconv.Atoi(res.Header["Errcode"][0])
	if errcode != 0 {
		fmt.Println(card.GetURL(), "request err", errcode)
	}

}
