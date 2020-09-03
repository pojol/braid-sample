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
)

// GuestLogin 游客登录
func GuestLogin(ctx context.Context, token string, reqBody []byte) (interface{}, error) {
	res := proto.GuestLoginRes{}
	mailRes := &api.SendMailRes{}

	var err error

	braid.Client().Invoke(ctx, "mail", "/api.mail/send", "", &api.SendMailReq{
		Accountid: "testaccountid",
		Body: &api.MailBody{
			Title: "hello,braid.",
		},
	}, mailRes)

	if mailRes.Errcode != int32(errcode.Succ) {
		fmt.Println("send mail err", mailRes.Errcode)
	}

	res.Token = "test_token_" + strconv.Itoa(int(time.Now().Unix()))
	fmt.Println("guest login", res.Token)

	return res, err
}

// Loginout 登出
func Loginout(ctx context.Context, token string, reqBody []byte) (interface{}, error) {

	fmt.Println(token, "login out")

	return []byte{}, nil
}
