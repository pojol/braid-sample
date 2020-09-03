package control

import (
	"braid-game/proto"
	"context"
	"encoding/json"
	"fmt"
)

// Rename 游客登录
func Rename(ctx context.Context, token string, reqBody []byte) (interface{}, error) {
	res := proto.RenameRes{}
	req := proto.RenameReq{}
	var err error

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, err
	}

	fmt.Println("rename", "token", token, "name:", req.Nickname)
	res.Nickname = req.Nickname

	return res, err
}
