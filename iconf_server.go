package main

import (
	"log"

	"github.com/go-redis/redis"
	"net"

	"github.com/kataras/iris"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/context"
	"github.com/spf13/viper"
	"fmt"
	"strconv"
	"strings"
	"github.com/kataras/golog"
)

type Visitor struct {
	Name string
}

func main() {

	viper.SetDefault("redisAddr", "127.0.0.1:6379")
	viper.SetDefault("redisPw", "")
	viper.SetDefault("redisDb", 0)

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/iserv/")
	viper.AddConfigPath("$HOME/.iserv")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file %s \n", err))
	}

	client := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redisAddr"),
		Password: viper.GetString("redisPw"),
		DB: viper.GetInt("redisDb"),
	})

	setting, err := client.Get("config:ip").Result()

	if err == redis.Nil {
		setting = "172.16.0.1"
		log.Print("Could not read config from redis")
	} else {
		log.Print("Config read from redis")
	}

	portS, err := client.Get("config:http:port").Result()

	var port uint64

	if err == redis.Nil {
		fmt.Print("Port was not readable ", port)
		port = 80
	} else {
		port, err = strconv.ParseUint(portS, 10, 32)
		if err != nil {
			fmt.Print("Could not parse port string")
			port = 80
		}
	}

	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Logger().Level = golog.DebugLevel

	pugEngine := iris.Pug("./templates", ".pug").Binary(Asset, AssetNames)
	pugEngine.Reload(true)
	app.RegisterView(pugEngine)

	app.StaticEmbedded("/static", "./assets", Asset, AssetNames)

	app.Use(func(ctx context.Context) {
		ip := net.ParseIP(ctx.RemoteAddr())
		key := "machine:" + ip.String()
		num, e := client.Exists(key).Result()
		if num == 0 {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.View("static.pug")
			return
		}
		_, e = client.HGet(key, "mac").Result()
		if e == redis.Nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.View("static.pug")
			return
		}
		path := ctx.Request().RequestURI
		if !(strings.Index(path, "/api") == 0 || strings.Index(path, "/login") == 0) {
			_, err := client.HGet(key, "name").Result()
			if err == redis.Nil {
				ctx.Redirect("/login", iris.StatusTemporaryRedirect)
				return
			}
		}
		ctx.Next()
	})

	app.Controller("/", new(IndexController), *client)
	app.Controller("/dashboard", new(DashboardController), *client)
	app.Controller("/api", new(ApiController), *client)

	app.Run(iris.Addr(setting + ":" + strconv.FormatUint(port, 10)), iris.WithoutVersionChecker, iris.WithCharset("UTF-8"))
}

