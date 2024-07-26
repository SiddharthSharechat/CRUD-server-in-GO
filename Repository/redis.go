package Repository

import (
	"context"
	"encoding/json"
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
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
