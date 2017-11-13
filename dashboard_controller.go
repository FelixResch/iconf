package main

import (
	"github.com/kataras/iris/mvc"
	"github.com/go-redis/redis"
	"net"
)

type DashboardInfo struct {
	Visitor Visitor
}

type DashboardController struct {
	mvc.C

	Client redis.Client
}

func (c* DashboardController) Get() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard.pug",
		Data: dashboardInfo,
	}
}

func (c* DashboardController) GetDevices() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard/devices.pug",
		Data: dashboardInfo,
	}
}

func (c* DashboardController) GetServers() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard/servers.pug",
		Data: dashboardInfo,
	}
}

func (c* DashboardController) GetRecords() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard/records.pug",
		Data: dashboardInfo,
	}
}

type RecordDetailInfo struct {
	DashboardInfo

	RecordKey string
}

func (c* DashboardController) GetRecordsAdd() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard/records/add.pug",
		Data: dashboardInfo,
	}
}

func (c* DashboardController) GetRecordsBy(key string) mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := RecordDetailInfo{DashboardInfo: DashboardInfo{Visitor:Visitor{Name:name}}, RecordKey: key}
	return mvc.View{
		Name: "dashboard/records/view.pug",
		Data: dashboardInfo,
	}
}

func (c* DashboardController) GetGames() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard/games.pug",
		Data: dashboardInfo,
	}
}

func (c* DashboardController) GetAdmin() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard/admin.pug",
		Data: dashboardInfo,
	}
}

func (c* DashboardController) GetSettings() mvc.Result {
	ctx := c.Ctx
	ctx.Gzip(true)
	ip := net.ParseIP(ctx.RemoteAddr())
	name, e := c.Client.HGet("machine:" + ip.String(), "name").Result()
	if e == redis.Nil {
		panic("User deleted while operation was performed!")
	}
	dashboardInfo := DashboardInfo{Visitor{name}}
	return mvc.View{
		Name: "dashboard/settings.pug",
		Data: dashboardInfo,
	}
}