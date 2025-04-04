package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/hotel"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/wrappers"
	"net/http"
	"runtime"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func searchHotels(ctx context.Context, req *hotel.SearchHotelsRequest) *hotel.SearchHotelsResponse {
	hotels := hotel.SearchHotels(ctx, req.InDate, req.OutDate, req.Location)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := hotel.SearchHotelsResponse{Profiles: hotels}
	//fmt.Printf("[ReviewStorage] Response: %v\n", resp)
	return &resp
}

func storeHotel(ctx context.Context, req *hotel.StoreHotelRequest) *hotel.StoreHotelResponse {
	hotelId := hotel.StoreHotel(ctx, req.HotelId, req.Name, req.Phone, req.Location, req.Rate, req.Capacity, req.Info)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := hotel.StoreHotelResponse{HotelId: hotelId}
	//fmt.Printf("[ReviewStorage] Response: %v\n", resp)
	return &resp
}

func reservation(ctx context.Context, req *hotel.FrontendReservationRequest) *hotel.FrontendReservationResponse {
	success := hotel.FrontendReservation(ctx, req.HotelId, req.InDate, req.OutDate, req.Rooms, req.Username, req.Password)
	resp := hotel.FrontendReservationResponse{Success: success}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	slowpoke.SlowpokeInit()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_search_hotels", wrappers.SlowpokeWrapper[hotel.SearchHotelsRequest, hotel.SearchHotelsResponse](searchHotels, "ro_search_hotels"))
	http.HandleFunc("/store_hotel", wrappers.SlowpokeWrapper[hotel.StoreHotelRequest, hotel.StoreHotelResponse](storeHotel, "store_hotel"))
	http.HandleFunc("/reservation", wrappers.SlowpokeWrapper[hotel.FrontendReservationRequest, hotel.FrontendReservationResponse](reservation, "reservation"))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
