package Repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
	"github.com/SiddharthSharechat/CRUDGo/Models"
	"log"
	"time"
)

var ctx = context.Background()

func getKeyType(key string) string {
	keyType, e := Initializers.RDb.Type(ctx, key).Result()
	if e != nil {
		log.Fatalf("Type check failed: %v\n", e)
	}
	return keyType
}

func SetValue(key string, value interface{}) {
	binaryVal, _ := json.Marshal(&value)
	res := Initializers.RDb.Set(ctx, key, binaryVal, time.Minute)
	log.Printf("set value res: %v\n", res)
}

func GetValue(key string, ptr interface{}) bool {
	keyType := getKeyType(key)
	val, er := Initializers.RDb.Get(ctx, key).Result()
	if keyType != "none" {
		if er != nil {
			log.Fatalf("could not get for key %s got error %s for keyType %s", key, er, keyType)
		} else {
			err := json.Unmarshal([]byte(val), ptr)
			if err != nil {
				log.Printf("could not unmarshal value for key %s got error %s for keyType %s", key, er, err)
			}
			return true
		}
	}
	return false
}

func Expire(key string) {
	res := Initializers.RDb.Expire(ctx, key, 0)
	log.Printf("expire res: %v\n", res)
}

func RPush(key string, value interface{}) {
	val, e := json.Marshal(&value)
	if e != nil {
		log.Printf("could not marshal value for key %s got error %s", key, e)
	}
	response := Initializers.RDb.RPush(ctx, key, val)
	log.Printf("rpush res: %v\n", response)
	res := Initializers.RDb.Expire(ctx, key, time.Minute)
	log.Printf("expire set for Rpush res: %v\n", res)
}

func LGet(key string, storedUsers *[]Models.UserResponse) bool {
	fmt.Printf("%s is the key for pagination api\n", key)
	storedUsersJSON, e := Initializers.RDb.LRange(ctx, key, 0, -1).Result()
	if e != nil {
		log.Printf("LRange failed: %v\n", e)
	}

	if len(storedUsersJSON) == 0 {
		return false
	}

	for _, userJSON := range storedUsersJSON {
		var user Models.UserResponse
		err := json.Unmarshal([]byte(userJSON), &user)
		if err != nil {
			log.Printf("could not unmarshal user response for key %s got error %s", key, err)
			return false
		}
		*storedUsers = append(*storedUsers, user)
	}
	return true
}

func ClearPaginationCache(key string) {
	storedCachekeys, err := Initializers.RDb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Printf("could not get cache keys: %v\n", err)
	}
	for _, storedCachekey := range storedCachekeys {
		k := storedCachekey[1 : len(storedCachekey)-1]
		er := Initializers.RDb.Expire(ctx, k, time.Second)
		if er != nil {
			log.Printf("could not expire cache with key %v: %v\n", k, er)
		}
	}
	e := Initializers.RDb.Expire(ctx, key, 0)
	if e != nil {
		log.Printf("could not expire cache with key %v: %v\n", key, e)
	}
}
