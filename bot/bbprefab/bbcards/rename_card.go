package bbcards

import (
	"braid-game/bot/bbprefab"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// RenameCard 改名
type RenameCard struct {
	URL   string
	delay time.Duration
	md    *bbprefab.BotData
}

// RenameReq 更新用户名请求
type RenameReq struct {
	Token    string
	Nickname string
}

// RenameRes 更新用户名返回
type RenameRes struct {
	nickname string
}

// NewRenameCard 修改昵称
func NewRenameCard(md *bbprefab.BotData) *RenameCard {
	return &RenameCard{
		URL:   "/v1/base/rename",
		delay: time.Millisecond,
		md:    md,
	}
}

// GetURL 获取服务器地址
func (card *RenameCard) GetURL() string { return card.URL }

// GetHeader get card header
func (card *RenameCard) GetHeader() map[string]string {
	return map[string]string{
		"token": card.md.AccToken,
	}
}

// SetDelay 设置卡片之间调用的延迟
func (card *RenameCard) SetDelay(delay time.Duration) { card.delay = delay }

// GetDelay 获取卡片之间调用的延迟
func (card *RenameCard) GetDelay() time.Duration { return card.delay }

// Marshal 序列化传入消息体
func (card *RenameCard) Marshal() []byte {

	req := RenameReq{
		Token:    card.md.AccToken,
		Nickname: "newname",
	}

	b, err := json.Marshal(&req)
	if err != nil {
		fmt.Println(card.GetURL(), "proto.Marshal err", err)
	}

	return b
}

// Unmarshal 反序列化返回消息
func (card *RenameCard) Unmarshal(res *http.Response) {

	errcode, _ := strconv.Atoi(res.Header["Errcode"][0])
	if errcode != 0 {
		fmt.Println(card.GetURL(), "request err", errcode)
	}
}
