package control

import (
	"braid-game/api"
	"braid-game/errcode"
	"braid-game/proto"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pojol/braid"
)

// Rename 游客登录
func Rename(ctx context.Context, token string, reqBody []byte) (interface{}, error) {
	res := proto.RenameRes{}
	req := proto.RenameReq{}
	mailRes := &api.SendMailRes{}
	var err error

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, err
	}

	res.Nickname = req.Nickname

	braid.Invoke(ctx, "mail", "/api.mail/send", token, &api.SendMailReq{
		Accountid: "testaccountid",
		Body: &api.MailBody{
			Title: "hello,braid.",
		},
	}, mailRes)
	if mailRes.Errcode != int32(errcode.Succ) {
		fmt.Println("send mail err", mailRes.Errcode, token)
	}

	return res, err
}
