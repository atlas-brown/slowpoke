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

func addProduct(ctx context.Context, req *boutique.AddProductRequest) *boutique.AddProductResponse {
	// slowpoke.SlowpokeCheck("addProduct")
	productId := boutique.AddProduct(ctx, req.Product)
	resp := boutique.AddProductResponse{ProductId: productId}
	return &resp
}

func getProduct(ctx context.Context, req *boutique.GetProductRequest) *boutique.GetProductResponse {
	// slowpoke.SlowpokeCheck("getProduct")
	product := boutique.GetProduct(ctx, req.ProductId)
	//fmt.Printf("Product read: %+v\n", product)
	resp := boutique.GetProductResponse{Product: product}
	return &resp
}

func searchProducts(ctx context.Context, req *boutique.SearchProductsRequest) *boutique.SearchProductsResponse {
	// slowpoke.SlowpokeCheck("searchProducts")
	products := boutique.SearchProducts(ctx, req.Query)
	//fmt.Printf("Products read: %+v\n", products)
	resp := boutique.SearchProductsResponse{Products: products}
	return &resp
}

func fetchCatalog(ctx context.Context, req *boutique.FetchCatalogRequest) *boutique.FetchCatalogResponse {
	// slowpoke.SlowpokeCheck("fetchCatalog")
	products := boutique.FetchCatalog(ctx, req.CatalogSize)
	resp := boutique.FetchCatalogResponse{Catalog: products}
	return &resp
}

func addProducts(ctx context.Context, req *boutique.AddProductsRequest) *boutique.AddProductsResponse {
	// slowpoke.SlowpokeCheck("addProducts")
	boutique.AddProducts(ctx, req.Products)
	resp := boutique.AddProductsResponse{OK: "OK"}
	return &resp
}

func loadProducts(ctx context.Context) []boutique.Product {

	// List directory
	var products []boutique.Product
	catalogJSON, err := os.ReadFile("/app/cmd/boutique/products.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(catalogJSON, &products)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded products: ", len(products))
	// for _, product := range products {
	// 	fmt.Println(product)
	// }
	return products
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	// go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	// http.HandleFunc("/add_product", wrappers.NonROWrapper[boutique.AddProductRequest, boutique.AddProductResponse](addProduct))
	http.HandleFunc("/add_product", wrappers.SlowpokeWrapper[boutique.AddProductRequest, boutique.AddProductResponse](addProduct, "addProduct"))
	// http.HandleFunc("/add_products", wrappers.NonROWrapper[boutique.AddProductsRequest, boutique.AddProductsResponse](addProducts))
	http.HandleFunc("/add_products", wrappers.SlowpokeWrapper[boutique.AddProductsRequest, boutique.AddProductsResponse](addProducts, "addProducts"))
	// http.HandleFunc("/ro_get_product", wrappers.ROWrapper[boutique.GetProductRequest, boutique.GetProductResponse](getProduct))
	http.HandleFunc("/ro_get_product", wrappers.SlowpokeWrapper[boutique.GetProductRequest, boutique.GetProductResponse](getProduct, "getProduct"))
	// http.HandleFunc("/ro_search_products", wrappers.ROWrapper[boutique.SearchProductsRequest, boutique.SearchProductsResponse](searchProducts))
	http.HandleFunc("/ro_search_products", wrappers.SlowpokeWrapper[boutique.SearchProductsRequest, boutique.SearchProductsResponse](searchProducts, "searchProducts"))
	// http.HandleFunc("/ro_fetch_catalog", wrappers.ROWrapper[boutique.FetchCatalogRequest, boutique.FetchCatalogResponse](fetchCatalog))
	http.HandleFunc("/ro_fetch_catalog", wrappers.SlowpokeWrapper[boutique.FetchCatalogRequest, boutique.FetchCatalogResponse](fetchCatalog, "fetchCatalog"))
	boutique.InitAllProducts(context.Background(), loadProducts(context.Background()))
	slowpoke.SlowpokeInit()
	// boutique.GetAllProducts(context.Background())
	fmt.Println("Server started on port 3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	slowpokeListener := &slowpoke.SlowpokeListener{listener}
	panic(http.Serve(slowpokeListener, nil))
}
