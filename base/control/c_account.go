package control

import (
	"braid-game/api"
	"braid-game/errcode"
	"braid-game/proto"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pojol/braid/rpc/client"
)

// Rename 游客登录
func Rename(ctx context.Context, reqBody []byte) (interface{}, error) {
	res := proto.RenameRes{}
	req := proto.RenameReq{}
	var err error

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, err
	}

	fmt.Println("rename")

	conn, err := client.GetConn("mail")
	if err != nil {
		return nil, err
	}
	defer conn.Put()

	cc := api.NewMailClient(conn.ClientConn)
	rres, err := cc.Send(ctx, &api.SendMailReq{
		Accountid: "testaccountid",
		Body: &api.MailBody{
			Title: "testTitle",
		},
	})
	if err != nil {
		conn.Unhealthy()
	}

	if rres.Errcode != int32(errcode.Succ) {
		//
	}

	return res, err
}
