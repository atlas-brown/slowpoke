package hotel

import (
	"context"
	"github.com/eniac/mucache/pkg/slowpoke"
	"os"
	"github.com/goccy/go-json"
	"fmt"
)

func InitLocations() {
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
		// name, _ := item["name"].(string)
		// phone, _ := item["phone"].(string)
		addr, _ := item["address"].(map[string]interface{})
		location, _ := addr["city"].(string)
		// rate := 100
		// capacity := 11
		// info := getRandomString(1000)

		StoreHotelLocation(ctx, hotelId, location)
	}
	fmt.Printf("Initialized %d locations \n", len(data))
}

func Nearby(ctx context.Context, inDate string, outDate string, location string) []Rate {
	// Find the hotel ids in that location
	hotelIds := getHotelIdsForLocation(ctx, location)

	// Get the rates for these hotels
	req := GetRatesRequest{HotelIds: hotelIds}
	ratesRes := slowpoke.Invoke[GetRatesResponse](ctx, "rate", "ro_get_rates", req)
	return ratesRes.Rates
}

func StoreHotelLocation(ctx context.Context, hotelId string, location string) string {
	hotelIds := getHotelIdsForLocation(ctx, location)
	// Keep saved reviews bounded to 10 for consistent performance measurements
	if len(hotelIds) >= 10 {
		hotelIds = hotelIds[1:]
	}
	hotelIds = append(hotelIds, hotelId)
	slowpoke.SetState(ctx, location, hotelIds)
	return hotelId
}

func getHotelIdsForLocation(ctx context.Context, location string) []string {
	hotelIds, err := slowpoke.GetState[[]string](ctx, location)
	// If err != nil then the key does not exist
	if err != nil {
		return []string{}
	} else {
		return hotelIds
	}
}
