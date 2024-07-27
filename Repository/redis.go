package Repository

import (
	"context"
	"encoding/json"
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
	Initializers.RDb.Set(ctx, key, binaryVal, time.Minute)
}

func GetValue(key string, ptr interface{}) bool {
	keyType := getKeyType(key)
	val, er := Initializers.RDb.Get(ctx, key).Result()
	if keyType != "none" {
		if er != nil {
			log.Fatalf("could not get for key %s got error %s for keyType %s", key, er, keyType)
		} else {
			json.Unmarshal([]byte(val), ptr)
			return true
		}
	}
	return false
}

func Expire(key string) {
	Initializers.RDb.Expire(ctx, key, 0)
}

func RPush(key string, value interface{}) {
	val, _ := json.Marshal(&value)
	Initializers.RDb.RPush(ctx, key, val)
	Initializers.RDb.Expire(ctx, key, time.Minute)
}

func LGet(key string, storedUsers *[]Models.UserResponse) bool {
	storedUsersJSON, _ := Initializers.RDb.LRange(ctx, key, 0, -1).Result()

	if len(storedUsersJSON) == 0 {
		return false
	}

	for _, userJSON := range storedUsersJSON {
		var user Models.UserResponse
		err := json.Unmarshal([]byte(userJSON), &user)
		if err != nil {
			return false
		}
		*storedUsers = append(*storedUsers, user)
	}
	return true
}
