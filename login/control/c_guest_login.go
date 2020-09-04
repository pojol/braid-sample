package control

import (
	"braid-game/api"
	"braid-game/errcode"
	"braid-game/proto"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pojol/braid"
	"github.com/pojol/braid/module/pubsub"
	"github.com/pojol/braid/plugin/linkerredis"
)

// GuestLogin 游客登录
func GuestLogin(ctx context.Context, token string, reqBody []byte) (interface{}, error) {
	res := proto.GuestLoginRes{}
	mailRes := &api.SendMailRes{}

	var err error
	token = "t" + strconv.Itoa(int(time.Now().Unix()))

	braid.Client().Invoke(ctx, "mail", "/api.mail/send", token, &api.SendMailReq{
		Accountid: "testaccountid",
		Body: &api.MailBody{
			Title: "hello,braid.",
		},
	}, mailRes)

	if mailRes.Errcode != int32(errcode.Succ) {
		fmt.Println("send mail err", mailRes.Errcode)
	}

	fmt.Println("guest login", token)

	res.Token = token
	return res, err
}

// Loginout 登出
func Loginout(ctx context.Context, token string, reqBody []byte) (interface{}, error) {

	fmt.Println(token, "login out")

	//
	braid.Pubsub().Pub(linkerredis.LinkerTopicUnlink, &pubsub.Message{
		Body: []byte(token),
	})

	return []byte{}, nil
}
