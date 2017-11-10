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
)

type Visitor struct {
	Name string
}

type DashboardInfo struct {
	Visitor Visitor
	ActiveDevices uint32
	StoredNameRecords uint32
	ActiveServers uint32
	ActiveGames uint32
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
		Addr: viper.GetString("redisArr"),
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

	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New())

	pugEngine := iris.Pug("./templates", ".pug").Binary(Asset, AssetNames)
	pugEngine.Reload(true)
	app.RegisterView(pugEngine)

	app.StaticEmbedded("/static", "./assets", Asset, AssetNames)

	app.Use(func(ctx context.Context) {
		ip := net.ParseIP(ctx.RemoteAddr())
		_, e := client.HGet("machine:" + ip.String(), "mac").Result()
		if e == redis.Nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.View("static.pug")
			return
		} else {
			ctx.Next()
		}
	})

	app.Get("/", func(ctx context.Context) {
		ctx.Gzip(true)
		ip := net.ParseIP(ctx.RemoteAddr())
		_, e := client.HGet("machine:" + ip.String(), "name").Result()
		if e == redis.Nil {
			ctx.Redirect("/login")
		} else {
			ctx.Redirect("/dashboard", iris.StatusPermanentRedirect)
		}
	})

	app.Get("/login", func(ctx context.Context) {
		ctx.Gzip(true)
		ctx.View("login.pug")
	})

	app.Get("/dashboard", func(ctx context.Context) {
		ctx.Gzip(true)
		ip := net.ParseIP(ctx.RemoteAddr())
		name, e := client.HGet("machine:" + ip.String(), "name").Result()
		if e == redis.Nil {
			ctx.Redirect("/login")
			return
		}
		ctx.ViewData("", DashboardInfo{Visitor{name},
			activeDeviceCount(*client),
			storedNameRecords(*client),
			activeServers(*client),
			activeGames(*client)})
		ctx.View("dashboard.pug")
	})

	app.Get("/leave", func(ctx context.Context) {
		ip := net.ParseIP(ctx.RemoteAddr())
		name, e := client.HGet("machine:" + ip.String(), "name").Result()
		if e == nil {
			log.Print("Deleting machine information")
			_, e := client.HDel("machine:" + ip.String(), "name").Result()
			if e == nil {
				log.Print("Deleting dns entries")
				go deleteAllKeysMatching("record:" + name[0:len(name)-1] + ":*", *client)
				ctx.Redirect("/login", iris.StatusTemporaryRedirect)
				ctx.Header("Pragma", "no-cache")
				ctx.WriteString("OK")
			} else {
				ctx.StatusCode(iris.StatusInternalServerError)
				log.Print(e)
			}
		}
	})

	app.Post("/api/register", func(ctx context.Context) {

		ip := net.ParseIP(ctx.RemoteAddr())
		mac, e := client.HGet("machine:" + ip.String(), "mac").Result()
		if e == redis.Nil {
			ctx.StatusCode(404)
			ctx.WriteString("Your device wasn't found on this server, is your address assigned statically?\n")
			return
		}
		log.Print("Incoming request from ", ip, " with HWAddr ", mac)

		//TODO check mac
		visitor := Visitor{}
		err = ctx.ReadForm(&visitor)
		if err != nil {
			log.Fatalln("Could not parse IP", err)
			ctx.StatusCode(500)
			ctx.WriteString("Internal Error\n")
			log.Fatal(err)
			return
		}
		name := visitor.Name
		if name == "master" {
			ctx.StatusCode(403)
			ctx.WriteString("Name \"master\" can't be chosen\n")
			return
		}
		_, err = client.HSet("machine:" + ip.String(), "name", name + ".").Result()
		if  e != nil {
			log.Fatalln("Could not write to redis", e)
		}
		_, e = client.HSet("record:" + name, "type", "A").Result()
		if  e != nil {
			log.Fatalln("Could not write to redis", e)
		}
		_, e = client.HSet("record:" + name, "host", ip.String()).Result()
		if  e != nil {
			log.Fatalln("Could not write to redis", e)
		}
		ctx.WriteString("Okay\n")
	})

	app.Run(iris.Addr(setting + ":80"))
}

func deleteAllKeysMatching(pattern string, client redis.Client) {
	log.Print("Loading keys to delete for pattern: " + pattern)
	keys, e := client.Keys(pattern).Result()
	log.Println(e, keys)
	if e == nil {
		for i := 0; i < len(keys);i++  {
			key := keys[i]
			log.Print("Deleting ", key)
			client.Del(key)
		}
	} else {
		log.Fatalln("Could not clear DNS entries", e)
	}
}

func activeDeviceCount(client redis.Client) uint32 {
	keys, e := client.Keys("machine:*").Result()
	if e == redis.Nil {
		print(e)
		return uint32(0)
	} else {
		return uint32(len(keys))
	}
}

func storedNameRecords(client redis.Client) uint32 {
	keys, e := client.Keys("record:*").Result()
	if e == redis.Nil {
		print(e)
		return uint32(0)
	} else {
		return uint32(len(keys))
	}
}

func activeServers(client redis.Client) uint32 {
	keys, e := client.Keys("server:*").Result()
	if e == redis.Nil {
		print(e)
		return uint32(0)
	} else {
		return uint32(len(keys))
	}
}

func activeGames(client redis.Client) uint32 {
	keys, e := client.Keys("game:*").Result()
	if e == redis.Nil {
		print(e)
		return uint32(0)
	} else {
		return uint32(len(keys))
	}
}