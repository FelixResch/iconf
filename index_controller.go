package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"net"
	"github.com/go-redis/redis"
)

// Section: Views

var loginView = mvc.View{
	Name: "login.pug",
}

type IndexController struct {
	mvc.C

	Client redis.Client
}

//
// GET: /

func (c* IndexController) Get() {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr()).To4()
	ipS := ip.String()
	request := c.Client.HGet("machine:" + ipS, "name")
	_, e := request.Result()
	if e == redis.Nil {
		ctx.Redirect("/login")
	} else {
		ctx.Redirect("/dashboard", iris.StatusPermanentRedirect)
	}
}

//
// GET: /login

func (c* IndexController) GetLogin() mvc.Result {
	c.Ctx.Gzip(true)
	return loginView
}

//
// GET: /leave
func (c* IndexController) GetLeave() mvc.Result {
	ctx := c.Ctx
	client := c.Client
	ip := net.ParseIP(ctx.RemoteAddr())
	log := ctx.Application().Logger()
	name, e := client.HGet("machine:" + ip.String(), "name").Result()
	if e == nil {
		log.Debug("Deleting machine information")
		_, e := client.HDel("machine:" + ip.String(), "name").Result()
		if e == nil {
			log.Debug("Deleting dns entries")
			go deleteAllKeysMatching("record:" + name[0:len(name)-1] + ":*", client)
			ctx.Redirect("/login", iris.StatusTemporaryRedirect)
			ctx.Header("Pragma", "no-cache")
			return mvc.Response{
				Code: 200,
				Content: []byte("OK"),
			}
		}
	}
	return mvc.Response{
		Code: iris.StatusNotFound,
		Content: []byte("No element found for this machine!"),
	}
}