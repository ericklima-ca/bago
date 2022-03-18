package cachingservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
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

func SetToken(subkey string, moreParams ...interface{}) {
	connect()
	token := Token{
		ID:         generateRandomToken(),
		Expiration: time.Hour,
		UserID:     moreParams[0].(uint),
	}
	var key string
	switch subkey {
	case "signup":
		key = fmt.Sprintf("token:signup:%v", token.UserID)
		rdb.Set(context.TODO(), key, token.ID, token.Expiration)
	case "recovery":
		key = fmt.Sprintf("token:recovery:%v", token.UserID)
		pass := encryptPass(moreParams[1].(string))
		rdb.Set(context.TODO(), key, token.ID+":"+pass, token.Expiration)
	}

}
func GetToken(subkey string, userID uint) (str string) {
	connect()
	key := fmt.Sprintf("token:%v:%v", subkey, userID)
	str, err := rdb.Get(context.TODO(), key).Result()
	if err != nil {
		return
	}
	return str
}

func DeleteToken(subkey string, userID uint) {
	connect()
	key := fmt.Sprintf("token:%v:%v", subkey, userID)
	rdb.Del(context.TODO(), key)
}

type Token struct {
	ID         string
	Expiration time.Duration
	UserID     uint
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

func encryptPass(pass string) string {
	userPassBytes := []byte(pass)
	passBytes, err := bcrypt.GenerateFromPassword(userPassBytes, bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("unable to hash password")
	}
	return string(passBytes)
}
