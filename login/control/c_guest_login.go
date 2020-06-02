package control

import (
	"braid-game/proto"
	"encoding/json"
	"strconv"
	"time"
)

// GuestLogin 游客登录
func GuestLogin(reqBody []byte) (interface{}, error) {
	res := proto.GuestLoginRes{}
	req := proto.GuestLoginReq{}
	var err error

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, err
	}

	res.Token = "test_token_" + strconv.Itoa(int(time.Now().Unix()))

	return res, err
}
