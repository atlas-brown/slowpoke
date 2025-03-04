package hotel

import (
	"bytes"
	"context"
	"crypto/sha256"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/lithammer/shortuuid"
	"fmt"
	"strconv"
)

func InitUsers() {
	ctx := context.Background()
	// Initialize the slowpoke state for users
	for i := 0; i < 100; i++ {
		username := "username" + strconv.Itoa(i)
		password := "password" + strconv.Itoa(i)
		RegisterUser(ctx, username, password)
	}
	fmt.Println("Initialized 100 users")
}

func RegisterUser(ctx context.Context, username string, password string) bool {
	fmt.Println("Registering user: ", username)
	userId := shortuuid.New()
	salt := shortuuid.New()
	hashPass := hash(password + salt)
	user := User{
		UserId:   userId,
		Username: username,
		Password: hashPass,
		Salt:     salt,
	}

	slowpoke.SetState(ctx, username, user)
	return true
}

func hash(str string) []byte {
	h := sha256.New()
	h.Write([]byte(str))
	val := h.Sum(nil)
	return val
}

func Login(ctx context.Context, username string, password string) string {
	fmt.Println("Logging in user: ", username)
	user, err := slowpoke.GetState[User](ctx, username)
	if err != nil {
		panic(err)
	}
	salt := user.Salt
	givenPass := hash(password + salt)
	if bytes.Equal(givenPass, user.Password) {
		return "OK"
	}
	return "NOT-OK"
}

func GetUserId(ctx context.Context, username string) string {
	user, err := slowpoke.GetState[User](ctx, username)
	if err != nil {
		panic(err)
	}
	return user.UserId
}
