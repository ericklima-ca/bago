package cachingservice

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func connect() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	rdb = redis.NewClient(opt)
}

func SetToken(subKey string, userID uint) {
	connect()
	token := Token{
		ID: generateRandomToken(),
		Expiration: time.Hour,
		UserID: userID,
	}
	key := fmt.Sprintf("token:%v:%v", subKey, token.UserID)
	rdb.Set(context.TODO(), key, token.ID, token.Expiration )
}
func GetToken(subKey string, userID uint) string {
	connect()
	key := fmt.Sprintf("token:%v:%v", subKey, userID)
	str, err := rdb.Get(context.TODO(), key).Result()
	
	if err != nil {
		log.Fatal(err)
	}
	return str
}

type Token struct {
	ID             string
	Expiration time.Duration
	UserID         uint
}

func generateRandomToken() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 64)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}