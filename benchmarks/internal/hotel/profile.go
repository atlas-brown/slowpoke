package hotel

import (
	"context"
	"github.com/eniac/mucache/pkg/slowpoke"
	"fmt"
	"os"
	"github.com/goccy/go-json"
)

func getRandomString(sz int) string {
	// Generate a random string of length sz
	randomString := make([]byte, sz)
	for i := range randomString {
		randomString[i] = byte('a' + i%26)
	}
	return string(randomString)
}

func InitProfiles() {
	ctx := context.Background()
	catalogJSON, err := os.ReadFile("/app/internal/hotel/data/hotels.json")
	if err != nil {
		panic(err)
	}
	var data []map[string]interface{}
	err = json.Unmarshal(catalogJSON, &data)
	if err != nil {
		panic(err)
	}
	for _, item := range data {
		hotelId, _ := item["id"].(string)
		name, _ := item["name"].(string)
		phone, _ := item["phone"].(string)
		// addr, _ := item["address"].(map[string]interface{})
		// location, _ := addr["city"].(string)
		// rate := 100
		// capacity := 11
		info := getRandomString(1000)
		profile := HotelProfile{
			HotelId:   hotelId,
			Name:      name,
			Phone:     phone,
			Info:  info,
		}
		StoreProfile(ctx, profile)
	}
	fmt.Printf("Initialized %d profiles\n", len(data))
}


func StoreProfile(ctx context.Context, profile HotelProfile) string {
	slowpoke.SetState(ctx, profile.HotelId, profile)
	return profile.HotelId
}

func GetProfiles(ctx context.Context, hotelIds []string) []HotelProfile {
	//fmt.Printf("[ReviewStorage] Asked for: %v\n", reviewIds)
	//profiles := make([]HotelProfile, len(hotelIds))
	//for i, hotelId := range hotelIds {
	//	profile, err := slowpoke.GetState[HotelProfile](ctx, hotelId)
	//	if err != nil {
	//		panic(err)
	//	}
	//	profiles[i] = profile
	//}

	// Bulk
	var profiles []HotelProfile
	if len(hotelIds) > 0 {
		profiles = slowpoke.GetBulkStateDefault[HotelProfile](ctx, hotelIds, HotelProfile{})
	} else {
		profiles = make([]HotelProfile, len(hotelIds))
	}
	//fmt.Printf("[ReviewStorage] Returning: %v\n", reviews)
	return profiles
}
