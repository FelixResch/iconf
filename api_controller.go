package main

import (
	"github.com/kataras/iris/mvc"
	"github.com/go-redis/redis"
	"net"
	"github.com/kataras/iris"
)

type ApiController struct {
	mvc.C

	Client redis.Client
}

func (c* ApiController) PostRegister() mvc.Result {
	ctx := c.Ctx
	client := c.Client
	log := ctx.Application().Logger()
	ip := net.ParseIP(ctx.RemoteAddr())
	mac, err := client.HGet("machine:" + ip.String(), "mac").Result()
	if err == redis.Nil {
		return mvc.Response{
			Code: iris.StatusNotFound,
			Content: []byte("Your device wasn't found on this server, is your address assigned statically?\n"),
		}
	}

	log.Debug("Incoming request from ", ip, " with HWAddr ", mac)

	//TODO check mac
	visitor := Visitor{}
	err = ctx.ReadForm(&visitor)
	if err != nil {
		log.Warn("Could not parse IP", err)
		return mvc.Response{
			Code: iris.StatusInternalServerError,
			Content: []byte("Internal Error\n"),
		}
	}
	name := visitor.Name
	if name == "master" {
		return mvc.Response{
			Code: iris.StatusForbidden,
			Content: []byte("Name \"master\" can't be chosen\n"),
		}
	}
	_, err = client.HSet("machine:" + ip.String(), "name", name + ".").Result()
	if  err != nil {
		log.Warn("Could not write to redis", err)
	}
	_, err = client.HSet("record:" + ip.String() + ":" + name + ":A", "type", "A").Result()
	if  err != nil {
		log.Warn("Could not write to redis", err)
	}
	_, err = client.HSet("record:" + ip.String() + ":" + name + ":A", "host", ip.String()).Result()
	if  err != nil {
		log.Warn("Could not write to redis", err)
	}
	return mvc.Response{
		Code: iris.StatusOK,
		Content: []byte("Ok\n"),
	}
}

func (c* ApiController) GetActiveDevices() uint32 {
	return activeDeviceCount(c.Client)
}

func (c* ApiController) GetActiveServers() uint32 {
	return activeServers(c.Client)
}

func (c* ApiController) GetActiveRecords() uint32 {
	return storedNameRecords(c.Client)
}

func (c* ApiController) GetActiveGames() uint32 {
	return activeGames(c.Client)
}

func (c* ApiController) GetDevices() []Device {
	return findActiveDevices(c.Client)
}

func (c* ApiController) GetDevicesByRecords(identifier string) uint32 {
	ip := net.ParseIP(identifier).To4()
	return findNumOfRecordsForDevice(c.Client, ip)
}

func (c* ApiController) GetDevicesByServers(identifier string) uint32 {
	ip := net.ParseIP(identifier).To4()
	return findNumOfServersForDevice(c.Client, ip)
}

func (c* ApiController) GetRecords() []Record {
	return findRecords(c.Client)
}

func (c* ApiController) GetRecordsBy(key string) Record {
	return findRecord(c.Client, key)
}

type CrudResult struct {
	State bool
	Identifier string
	Reason string
}

func (c* ApiController) DeleteRecordsBy(key string) CrudResult {
	return CrudResult{
		State:false,
		Identifier:key,
		Reason:"not implemented!",
	}
}

func (c* ApiController) PostRecords() CrudResult {
	ctx := c.Ctx
	raw := new(RawRecord)
	err := ctx.ReadJSON(raw)
	if err != nil {
		ctx.Application().Logger().Debug(err)
		ctx.StatusCode(iris.StatusNotAcceptable)
		return CrudResult{
			State:false,
		}
	}
	ctx.Application().Logger().Debug(raw)
	key, err := createRecord(c.Client, *raw)
	if err == nil {
		return CrudResult{
			State:true,
			Identifier:key,
		}
	} else {
		return CrudResult{
			State:false,
			Reason: err.Error(),
		}
	}
}