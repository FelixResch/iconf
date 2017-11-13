package main

import (
	"github.com/go-redis/redis"
	"log"
	"strings"
	"net"
)

type Device struct {
	Name string
	Ip string
	Mac string
}

type Record struct {
	Device string
	Type string
	Name string
	Description string
	Key string
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

func deleteAllKeysMatching(pattern string, client redis.Client) {
	keys, e := client.Keys(pattern).Result()
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

func findActiveDevices(client redis.Client) []Device {
	keys, _ := client.Keys("machine:*").Result()
	devices := make([]Device, len(keys))
	for i, key := range keys {
		deviceData, e := client.HGetAll(key).Result()
		if e == redis.Nil {
			panic("No data found for " + key)
		} else {
			ip := key[strings.Index(key, ":") + 1:]
			devices[i] = Device{
				Name:deviceData["name"][:len(deviceData["name"])-1],
				Mac:deviceData["mac"],
				Ip:ip,
			}
		}
	}
	return devices
}

func findNumOfRecordsForDevice(client redis.Client, ip net.IP) uint32 {
	keys, _ := client.Keys("record:" + ip.String() + ":*").Result()
	return uint32(len(keys))
}

func findNumOfServersForDevice(client redis.Client, ip net.IP) uint32 {
	keys, _ := client.Keys("server:" + ip.String() + ":*").Result()
	return uint32(len(keys))
}

func findRecords(client redis.Client) []Record {
	keys, _ := client.Keys("record:*").Result()
	records := make([]Record, len(keys))
	for i, key := range keys {
		records[i] = findRecord(client, key)
	}
	return records
}

func findRecord(client redis.Client, key string) Record {
	record, e := client.HGetAll(key).Result()
	if e == redis.Nil {
		panic("Inconsistent redis state!")
	}
	parts := strings.Split(key, ":")
	ip := parts[1]
	name := parts[2]
	typeN := record["type"]
	var descr string
	switch typeN {
	case "A":
		descr = record["host"]
	case "CNAME":
		descr = record["host"]
	case "SRV":
		descr = record["host"] + ":" + record["port"]
	}
	return Record{
		Device: ip,
		Name: name,
		Type: typeN,
		Description: descr,
		Key: key,
	}
}

type RawRecord struct {
	Machine string
	Name string
	Type string
	Host string
	Port string
}

func createRecord(client redis.Client, record RawRecord) (string, error) {
	//Attempting to guess the machine, if not possible this record will be stored as a rogue record
	ip := net.ParseIP(record.Host)
	if record.Machine == "" {
		if ip != nil {
			exists, _ := client.Exists("machine:" + record.Host).Result()
			if exists == 1 {
				record.Machine = record.Host
				goto verify
			}
		}
		keys, _ := client.Keys("record:*:" + record.Host + ":*").Result()
		if len(keys) > 0 {
			for _, key := range keys {
				mName := strings.Split(key, ":")[1]
				if mName != "_" {
					record.Machine = mName
					goto verify
				}
			}
		}

		record.Machine = "_"
	}

	verify:
		key := "record:" + record.Machine + ":" + record.Name + ":" + record.Type
		switch record.Type {
		case "A":
			if ip != nil && record.Name != "" {
				data := make(map[string]interface{})
				data["type"] = "A"
				data["host"] = record.Host
				_, err := client.HMSet(key, data).Result()
				if err != nil {
					//TODO generate errors
				}
			} else {
				//TODO generate errors
			}
		case "CNAME":
			if ip == nil && record.Host != "" && record.Name != "" {
				data := make(map[string]interface{})
				data["type"] = "CNAME"
				data["host"] = record.Host
				_, err := client.HMSet(key, data).Result()
				if err != nil {
					//TODO generate errors
				}
			} else {
				//TODO generate errors
			}
		case "SRV":
			if ip == nil && record.Host != "" && record.Name != "" && record.Port != "" {
				data := make(map[string]interface{})
				data["type"] = "SRV"
				data["host"] = record.Host
				data["port"] = record.Port
				_, err := client.HMSet(key, data).Result()
				if err != nil {
					//TODO generate errors
				}
			} else {
				//TODO generate errors
			}
		}
		return key, nil
}