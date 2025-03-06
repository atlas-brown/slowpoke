package hotel

import (
	"context"
	"github.com/eniac/mucache/pkg/slowpoke"
	"os"
	"github.com/goccy/go-json"
	"fmt"
)

func InitHotelAvailability() {
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
		// addr, _ := item["address"].(map[string]interface{})
		// location, _ := addr["city"].(string)
		// rate := 100
		capacity := 11
		// info := getRandomString(1000)

		AddHotelAvailability(ctx, hotelId, capacity)
	}
	fmt.Printf("Initialized %d avalibilities\n", len(data))
}

func datesIntersect(inDate1 string, outDate1 string, inDate2 string, outDate2 string) bool {
	// Note: This is a little hacky since there is no check that the dates are in the same format.
	// Assumes: date is in YYYY-MM-DD format
	// Assumes: outDates are after inDates
	if (inDate2 > outDate1) || (inDate1 > outDate2) {
		return false
	} else {
		return true
	}
}

func checkAvailability(availability HotelAvailability, inDate string, outDate string, numberOfRooms int) bool {
	capacity := availability.Capacity
	reservationsTheseDays := 0
	for _, reservation := range availability.Reservations {
		if datesIntersect(inDate, outDate, reservation.InDate, reservation.OutDate) {
			reservationsTheseDays += 1
		}
	}

	return reservationsTheseDays+numberOfRooms <= capacity
}

func CheckAvailability(ctx context.Context, customerName string, hotelIds []string, inDate string, outDate string, numberOfRooms int) []string {
	// Get all reservations for that hotel
	availableHotelIds := []string{}
	for _, hotelId := range hotelIds {
		availability, err := slowpoke.GetState[HotelAvailability](ctx, hotelId)
		if err != nil {
			panic(err)
		}

		isAvailable := checkAvailability(availability, inDate, outDate, numberOfRooms)
		//fmt.Printf("Availability of hotel: %v for dates: %v-%v and %v rooms: %v\n", hotelId, inDate, outDate, numberOfRooms, isAvailable)
		if isAvailable {
			availableHotelIds = append(availableHotelIds, hotelId)
		}
	}
	return availableHotelIds
}

func MakeReservation(ctx context.Context, customerName string, hotelId string, inDate string, outDate string, numberOfRooms int) bool {
	availability, err := slowpoke.GetState[HotelAvailability](ctx, hotelId)
	if err != nil {
		panic(err)
	}

	if !checkAvailability(availability, inDate, outDate, numberOfRooms) {
		return false
	}

	// Note: When we make a reservation, make sure that there are at most 10 reservations
	//       for the hotel so that we get predictable latency when fetching the state.
	if len(availability.Reservations) >= 10 {
		availability.Reservations = availability.Reservations[1:]
	}

	newReservation := Reservation{
		CustomerName: customerName,
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   numberOfRooms,
	}
	availability.Reservations = append(availability.Reservations, newReservation)
	slowpoke.SetState(ctx, hotelId, availability)
	return true
}

func AddHotelAvailability(ctx context.Context, hotelId string, capacity int) string {
	slowpoke.SetState(ctx, hotelId, HotelAvailability{Reservations: []Reservation{}, Capacity: capacity})
	return hotelId
}
