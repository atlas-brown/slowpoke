package boutique

import (
	"context"
	"errors"
	"fmt"
	"github.com/eniac/mucache/pkg/state"
	"github.com/redis/go-redis/v9"
	"github.com/goccy/go-json"
	"sync"
)

const (
	debug_cart = false
)

// var local_carts map[string]Cart
var local_carts sync.Map

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func CartInit() {
	// _redisPing(context.Background())
	// local_carts = make(map[string]Cart)
	local_carts = sync.Map{}
}

func _redisPing(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis-cart:6379",
		DB:   0,
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)
}

func _redisGetCart(ctx context.Context, key string) (Cart, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis-cart:6379",
		DB:   0,
	})
	val, err := client.Get(ctx, key).Result()
	if err != redis.Nil {
		fmt.Println("redis error: ", err)
		panic(err)
	}
	var cart Cart
	if len(val) == 0 {
		// fmt.Printf("cart not found: %v\n", key)
		return cart, errors.New("cart not found")
	}
	err = json.Unmarshal([]byte(val), &cart)
	if err != nil {
		fmt.Println("json error: ", err)
		panic(err)
	}
	return cart, nil
}

func _localGetCart(ctx context.Context, key string) (Cart, error) {
	// if cart, ok := local_carts[key]; ok {
	// 	return cart, nil
	// }
	// return Cart{}, errors.New("cart not found")

	cart, ok := local_carts.Load(key)
	if !ok {
		return Cart{}, errors.New("cart not found")
	}
	return cart.(Cart), nil
}

func getCartDefault(ctx context.Context, userId string) Cart {
	cart, err := state.GetState[Cart](ctx, userId)
	if fmt.Sprint(err) == "key not found" {
		cart = Cart{
			UserId: userId,
			Items:  []CartItem{},
		}
	} else if err != nil {
		panic(err)
	}
	return cart
}

func AddItem(ctx context.Context, userId string, productId string, quantity int32) bool {
	if debug_cart { fmt.Println("AddItem: ", userId, productId, quantity) }

	// item := CartItem{
	// 	ProductId: productId,
	// 	Quantity:  quantity,
	// }
	// // cart := getCartDefault(ctx, userId)
	// cart, err := _localGetCart(ctx, userId)
	// if err != nil {
	// 	cart = Cart{
	// 		UserId: userId,
	// 		Items:  []CartItem{},
	// 	}
	// }

	// // Append the new item to the cart
	// cart.Items = append(cart.Items, item)
	// local_carts[userId] = cart
	// // state.SetState(ctx, userId, cart)

	item := CartItem{
		ProductId: productId,
		Quantity:  quantity,
	}
	cart, err := _localGetCart(ctx, userId)
	if err != nil {
		cart = Cart{
			UserId: userId,
			Items:  []CartItem{},
		}
	}
	// Append the new item to the cart
	cart.Items = append(cart.Items, item)
	local_carts.Store(userId, cart)
	return true
}

func GetCart(ctx context.Context, userId string) Cart {
	if debug_cart { fmt.Println("GetCart: ", userId) }
	// cart, err := state.GetState[Cart](ctx, userId)
	// cart, err := _redisGetCart(ctx, userId)
	cart, err := _localGetCart(ctx, userId)
	if err != nil {
		cart = Cart{
			UserId: userId,
			Items:  []CartItem{},
		}
	}
	return cart
}

func EmptyCart(ctx context.Context, userId string) bool {
	// cart := getCartDefault(ctx, userId)
	// cart, err := _localGetCart(ctx, userId)
	// if err != nil {
	// 	cart = Cart{
	// 		UserId: userId,
	// 		Items:  []CartItem{},
	// 	}
	// }
	// cart.Items = []CartItem{}
	// // state.SetState(ctx, userId, cart)
	// local_carts[userId] = cart

	cart, err := _localGetCart(ctx, userId)
	if err != nil {
		cart = Cart{
			UserId: userId,
			Items:  []CartItem{},
		}
	}
	cart.Items = []CartItem{}
	local_carts.Store(userId, cart)
	return true
}
