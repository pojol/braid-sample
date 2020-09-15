package control

import (
	"braid-game/proto"
	"context"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pojol/braid"
	"github.com/pojol/braid/3rd/redis"
	"github.com/pojol/braid/module/pubsub"
	"github.com/pojol/braid/plugin/linkerredis"
)

// unixtime / 1000 + redis incr
// 60 * 60 * 24 = 86400 + 100 = 85600
// 使用原子操作保证，expire必定被调用
var uniqueIncrScript = redigo.NewScript(1, `
local current
current = redis.call("incr", KEYS[1])
if tonumber(current) == 1 then
    redis.call("expire", KEYS[1], 86500)
end`)

// GetUniqueID 基于redis获取唯一id
func GetUniqueID() string {
	uniqueid := ""
	val := ""
	ival := 0

	// 2006（golang创立年份） 年以来的天数
	totalDay := ((time.Now().Year() - 2006) * 365) + time.Now().YearDay()

	conn := redis.Get().Conn()
	defer conn.Close()

	field := "uniqueid_field" + strconv.Itoa(totalDay)
	err := uniqueIncrScript.Send(conn, field)
	if err != nil {
		//plog.BSystemError("redis", "GetUniqueID", "send", err.Error())
		goto TAG
	}

	val, err = redis.ConnGet(conn, field)
	if err != nil || val == "" {
		//plog.BSystemError("redis", "GetUniqueID", "ConnGet", err.Error())
		goto TAG
	}

	ival, _ = strconv.Atoi(val)
	ival += 100000
	uniqueid = strconv.FormatInt((int64(totalDay)*1000000)+int64(ival), 10)

TAG:
	return uniqueid
}

// GuestLogin 游客登录
func GuestLogin(ctx context.Context, token string, reqBody []byte) (interface{}, error) {
	res := proto.GuestLoginRes{}

	var err error
	token = "token" + GetUniqueID()

	time.AfterFunc(time.Minute, func() {
		braid.Pubsub().Pub(linkerredis.LinkerTopicUnlink, &pubsub.Message{
			Body: []byte(token),
		})
	})

	res.Token = token
	return res, err
}

// Loginout 登出
func Loginout(ctx context.Context, token string, reqBody []byte) (interface{}, error) {

	return []byte{}, nil
}
