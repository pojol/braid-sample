package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/pojol/braid"
	"github.com/pojol/braid/modules/grpcclient/bproto"

	"github.com/labstack/echo/v4"
)

func routing(ctx echo.Context, nodName string, serviceName string, token string) error {
	var err error
	res := &bproto.RouteRes{}
	var in []byte

	in, err = ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		goto EXT
	}

	braid.Invoke(ctx.Request().Context(), nodName, "/bproto.listen/routing", token, &bproto.RouteReq{
		Nod:     nodName,
		Service: serviceName,
		Token:   token,
		ReqBody: in,
	}, res)

EXT:
	if err != nil {
		ctx.Response().Header().Set("Errcode", "-1")
		ctx.Response().Header().Set("Errmsg", err.Error())
	} else {
		ctx.Response().Header().Set("Errcode", "0")
	}

	ctx.Blob(http.StatusOK, "text/plain; charset=utf-8", res.ResBody)
	return err
}

// Regist regist
func Regist(e *echo.Echo) {
	e.POST("/v1/login/guest", guestHandler)
	e.POST("/v1/login/out", outHandler)
	e.POST("/v1/base/rename", renameHandler)
}

func guestHandler(ctx echo.Context) error {
	token := ctx.Request().Header.Get("token")
	return routing(ctx, "login", "guest", token)
}

func outHandler(ctx echo.Context) error {
	token := ctx.Request().Header.Get("token")
	return routing(ctx, "login", "out", token)
}

func renameHandler(ctx echo.Context) error {
	token := ctx.Request().Header.Get("token")
	return routing(ctx, "base", "rename", token)
}
