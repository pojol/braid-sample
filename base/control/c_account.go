package control

import (
	"braid-game/api"
	"braid-game/errcode"
	"braid-game/proto"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pojol/braid/module/rpc/client"
)

// Rename 游客登录
func Rename(ctx context.Context, reqBody []byte) (interface{}, error) {
	res := proto.RenameRes{}
	req := proto.RenameReq{}
	mailRes := &api.SendMailRes{}
	var err error

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, err
	}

	fmt.Println("rename")

	client.Invoke(ctx, "mail", "/api.mail/send", &api.SendMailReq{
		Accountid: "testaccountid",
		Body: &api.MailBody{
			Title: "testTitle",
		},
	}, mailRes)

	if mailRes.Errcode != int32(errcode.Succ) {
		//
	}

	return res, err
}
