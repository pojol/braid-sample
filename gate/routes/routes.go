package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pojol/braid"
	"github.com/pojol/braid/plugin/grpcclient/bproto"

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

	braid.Client().Invoke(ctx.Request().Context(), nodName, "/bproto.listen/routing", token, &bproto.RouteReq{
		Nod:     nodName,
		Service: serviceName,
		Token:   token,
		ReqBody: in,
	}, res)

EXT:
	if err != nil {
		fmt.Println("routes", "routing", err.Error())
		ctx.Response().Header().Set("Errcode", "-1")
		ctx.Response().Header().Set("Errmsg", err.Error())
	} else {
		ctx.Response().Header().Set("Errcode", "0")
	}

	ctx.Blob(http.StatusOK, "text/plain; charset=utf-8", res.ResBody)
	return err
}

func parsingURL(url string) (string, string, string, error) {

	var ver string
	var box string
	var service string
	var err error
	var urls []string

	if url == "" {
		err = echo.ErrNotFound
		goto TAG
	}

	if url[0] == '/' {
		url = url[1:]
	}

	urls = strings.Split(url, "/")
	if len(urls) != 3 {
		err = echo.ErrNotFound
		goto TAG
	}

	err = nil
	ver = urls[0]
	box = urls[1]
	service = urls[2]

TAG:
	if err != nil {
		fmt.Println("routes", "parsingUrl", err.Error())
	}

	return ver, box, service, err
}

// PostRouting 路由
func PostRouting(ctx echo.Context) error {

	token := ctx.Request().Header.Get("token")
	_, boxName, serviceName, err := parsingURL(ctx.Request().RequestURI)
	if err != nil {
		return err
	}

	err = routing(ctx, boxName, serviceName, token)

	return err
}
