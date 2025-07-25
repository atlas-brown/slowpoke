package boutique

import (
	"context"
	// "github.com/eniac/mucache/pkg/invoke"
	"github.com/eniac/mucache/pkg/slowpoke"
	"fmt"
)

const (
	debug_frontend = false
)

func Home(ctx context.Context, request HomeRequest) HomeResponse {
	if (debug_frontend) { fmt.Println("Home: ", request) }
	req1 := GetSupportedCurrenciesRequest{}
	currenciesRes := slowpoke.Invoke[GetSupportedCurrenciesResponse](ctx, "currency", "ro_get_currencies", req1)
	//http.HandleFunc("/ro_get_currencies", wrappers.ROWrapper[boutique.GetSupportedCurrenciesRequest, boutique.GetSupportedCurrenciesResponse](getCurrencies))

	req2 := GetCartRequest{UserId: request.Userid}
	cartRes := slowpoke.Invoke[GetCartResponse](ctx, "cart", "ro_get_cart", req2)
	//http.HandleFunc("/ro_get_cart", wrappers.ROWrapper[boutique.GetCartRequest, boutique.GetCartResponse](getCart))

	req3 := FetchCatalogRequest{CatalogSize: request.CatalogSize}
	catalogRes := slowpoke.Invoke[FetchCatalogResponse](ctx, "productcatalog", "ro_fetch_catalog", req3)
	//http.HandleFunc("/ro_fetch_catalog", wrappers.ROWrapper[boutique.FetchCatalogRequest, boutique.FetchCatalogResponse](fetchCatalog))

	res := HomeResponse{
		Products:   catalogRes.Catalog,
		UserCart:   cartRes.Cart,
		Currencies: currenciesRes.Currencies,
	}
	return res
}

func FrontendSetCurrency(ctx context.Context, currency Currency) {
	if (debug_frontend) { fmt.Println("FrontendSetCurrency: ", currency) }
	req := SetCurrencySupportRequest{Currency: currency}
	slowpoke.Invoke[SetCurrencySupportResponse](ctx, "currency", "set_currency", req)
}

func BrowseProduct(ctx context.Context, productId string) BrowseProductResponse {
	if (debug_frontend) { fmt.Println("BrowseProduct: ", productId) }
	req := GetProductRequest{ProductId: productId}
	res := slowpoke.Invoke[GetProductResponse](ctx, "productcatalog", "ro_get_product", req)
	return BrowseProductResponse{res.Product}
}

func AddToCart(ctx context.Context, request AddToCartRequest) AddToCartResponse {
	if (debug_frontend) { fmt.Println("AddToCart: ", request) }
	req := AddItemRequest{
		UserId:    request.UserId,
		ProductId: request.ProductId,
		Quantity:  request.Quantity,
	}
	res := slowpoke.Invoke[AddItemResponse](ctx, "cart", "add_item", req)
	//http.HandleFunc("/add_item", wrappers.NonROWrapper[boutique.AddItemRequest, boutique.AddItemResponse](addItemToCart))
	return AddToCartResponse{OK: res.Ok}
}

func ViewCart(ctx context.Context, request ViewCartRequest) ViewCartResponse {
	if (debug_frontend) { fmt.Println("ViewCart: ", request) }
	req := GetCartRequest{
		UserId: request.UserId,
	}
	res := slowpoke.Invoke[GetCartResponse](ctx, "cart", "ro_get_cart", req)
	//http.HandleFunc("/ro_get_cart", wrappers.ROWrapper[boutique.GetCartRequest, boutique.GetCartResponse](getCart))
	return ViewCartResponse{C: res.Cart}
}

func Checkout(ctx context.Context, request CheckoutRequest) CheckoutResponse {
	if (debug_frontend) { fmt.Println("Checkout: ", request) }
	req := PlaceOrderRequest{
		UserId:       request.UserId,
		UserCurrency: request.UserCurrency,
		Address:      request.Address,
		Email:        request.Email,
		CreditCard:   request.CreditCard,
	}
	res := slowpoke.Invoke[PlaceOrderResponse](ctx, "checkout", "place_order", req)
	//http.HandleFunc("/place_order", wrappers.NonROWrapper[boutique.PlaceOrderRequest, boutique.PlaceOrderResponse](placeOrder))
	return CheckoutResponse{
		Res: res.Order,
	}
}
