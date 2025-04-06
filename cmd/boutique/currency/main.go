package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/boutique"
	// "github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/wrappers"
	"github.com/eniac/mucache/pkg/slowpoke"
	"net"
	"net/http"
	"runtime"
	"os"
	"github.com/goccy/go-json"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func setCurrency(ctx context.Context, req *boutique.SetCurrencySupportRequest) *boutique.SetCurrencySupportResponse {
	// slowpoke.SlowpokeCheck("setCurrency")
	ok := boutique.SetCurrencySupport(ctx, req.Currency)
	resp := boutique.SetCurrencySupportResponse{Ok: ok}
	return &resp
}

func getCurrencies(ctx context.Context, req *boutique.GetSupportedCurrenciesRequest) *boutique.GetSupportedCurrenciesResponse {
	// slowpoke.SlowpokeCheck("getCurrencies")
	currencies := boutique.GetSupportedCurrencies(ctx)
	resp := boutique.GetSupportedCurrenciesResponse{Currencies: currencies}
	return &resp
}

func convertCurrency(ctx context.Context, req *boutique.ConvertCurrencyRequest) *boutique.ConvertCurrencyResponse {
	// slowpoke.SlowpokeCheck("convertCurrency")
	amount := boutique.ConvertCurrency(ctx, req.Amount, req.ToCurrency)
	resp := boutique.ConvertCurrencyResponse{Amount: amount}
	return &resp
}

func initCurrencies(ctx context.Context, req *boutique.InitCurrencyRequest) *boutique.InitCurrencyResponse {
	// slowpoke.SlowpokeCheck("initCurrencies")
	boutique.InitCurrencies(ctx, req.Currencies)
	resp := boutique.InitCurrencyResponse{Ok: "OK"}
	return &resp
}

func loadCurrencies(ctx context.Context) []boutique.Currency {
	// List directory
	var currencies []boutique.Currency
	catalogJSON, err := os.ReadFile("/app/cmd/boutique/currency_conversion.json")
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(catalogJSON, &data); err != nil {
		panic(err)
	}
	for key, value := range data {
		currency := boutique.Currency{
			CurrencyCode: key,
			Rate:         fmt.Sprintf("%v", value),
		}
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	// go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	// http.HandleFunc("/set_currency", wrappers.NonROWrapper[boutique.SetCurrencySupportRequest, boutique.SetCurrencySupportResponse](setCurrency))
	http.HandleFunc("/set_currency", wrappers.SlowpokeWrapper[boutique.SetCurrencySupportRequest, boutique.SetCurrencySupportResponse](setCurrency, "setCurrency"))
	// http.HandleFunc("/init_currencies", wrappers.NonROWrapper[boutique.InitCurrencyRequest, boutique.InitCurrencyResponse](initCurrencies))
	http.HandleFunc("/init_currencies", wrappers.SlowpokeWrapper[boutique.InitCurrencyRequest, boutique.InitCurrencyResponse](initCurrencies, "initCurrencies"))
	// http.HandleFunc("/ro_get_currencies", wrappers.ROWrapper[boutique.GetSupportedCurrenciesRequest, boutique.GetSupportedCurrenciesResponse](getCurrencies))
	http.HandleFunc("/ro_get_currencies", wrappers.SlowpokeWrapper[boutique.GetSupportedCurrenciesRequest, boutique.GetSupportedCurrenciesResponse](getCurrencies, "getCurrencies"))
	// http.HandleFunc("/ro_convert_currency", wrappers.ROWrapper[boutique.ConvertCurrencyRequest, boutique.ConvertCurrencyResponse](convertCurrency))
	http.HandleFunc("/ro_convert_currency", wrappers.SlowpokeWrapper[boutique.ConvertCurrencyRequest, boutique.ConvertCurrencyResponse](convertCurrency, "convertCurrency"))
	boutique.InitAllCurrencies(context.Background(), loadCurrencies(context.Background()))
	slowpoke.SlowpokeInit()
	fmt.Println("Server started on :3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	slowpokeListener := &slowpoke.SlowpokeListener{listener}
	panic(http.Serve(slowpokeListener, nil))
}
