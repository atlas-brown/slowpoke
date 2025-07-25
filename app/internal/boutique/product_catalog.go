package boutique

import (
	"context"
	// "github.com/eniac/mucache/pkg/state"
	"strings"
	"fmt"
)

const (
	debug_product_catalog = false
)

var allProducts []Product

var CatalogSize = 1000

func InitAllProducts(ctx context.Context, products []Product) {
	allProducts = products
	fmt.Println("InitAllProducts: ", len(allProducts))
}

func GetAllProducts(ctx context.Context) {
	for _, product := range allProducts {
		fmt.Println(product)
	}
}

func AddProduct(ctx context.Context, product Product) string {
	// keys, err := state.GetState[[]string](ctx, "KEYS")
	// if err != nil {
	// 	// fmt.fmt.Println("empty db")
	// }
	// keys = append(keys, product.Id)
	// state.SetState(ctx, "KEYS", keys)
	// state.SetState(ctx, product.Id, product)

	allProducts = append(allProducts, product)
	return product.Id
}

func AddProducts(ctx context.Context, products []Product) {
	// keys, err := state.GetState[[]string](ctx, "KEYS")
	// if err != nil {
	// 	// fmt.fmt.Println("empty db")
	// }
	// // If keys are 100 then we don't want to add more to the catalog
	// if len(keys) < CatalogSize {
	// 	rest := CatalogSize - len(keys)
	// 	if len(products) < rest {
	// 		rest = len(products)
	// 	}
	// 	for i := 0; i < rest; i++ {
	// 		keys = append(keys, products[i].Id)
	// 	}
	// 	state.SetState(ctx, "KEYS", keys)
	// }

	// productMap := make(map[string]interface{})
	// for _, product := range products {
	// 	productMap[product.Id] = product
	// }
	// state.SetBulkState(ctx, productMap)

	if len(allProducts) < CatalogSize {
		rest := CatalogSize - len(allProducts)
		if len(products) < rest {
			rest = len(products)
		}
		for i := 0; i < rest; i++ {
			allProducts = append(allProducts, products[i])
		}
	}
	return
}

func GetProduct(ctx context.Context, Id string) Product {
	if debug_product_catalog { fmt.Println("GetProduct: ", Id) }
	// product, err := state.GetState[Product](ctx, Id)
	// if err != nil {
	// 	panic(err)
	// }
	// return product

	var product Product
	for _, p := range allProducts {
		if p.Id == Id {
			product = p
			break
		}
	}
	return product
}

func SearchProducts(ctx context.Context, name string) []Product {
	if debug_product_catalog { fmt.Println("SearchProducts: ", name) }
	// products := make([]Product, 0)
	// keys, err := state.GetState[[]string](ctx, "KEYS")
	// if err != nil {
	// 	panic(err)
	// }
	// for _, id := range keys {
	// 	product, err := state.GetState[Product](ctx, id)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if strings.Contains(strings.ToLower(product.Name), strings.ToLower(name)) ||
	// 		strings.Contains(strings.ToLower(product.Name), strings.ToLower(name)) {
	// 		products = append(products, product)
	// 	}
	// }
	// return products

	var products []Product
	for _, p := range allProducts {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(name)) {
			products = append(products, p)
		}
	}
	return products
}

func FetchCatalog(ctx context.Context, catalogSize int) []Product {
	if debug_product_catalog { fmt.Println("FetchCatalog: ", catalogSize) }
	var products []Product
	if catalogSize < len(allProducts) {
		products = allProducts[:catalogSize]
	} else {
		products = allProducts
	}
	return products
// 	keys, err := state.GetState[[]string](ctx, "KEYS")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Limit fetches to the catalog size
// 	if catalogSize < len(keys) {
// 		keys = keys[:catalogSize]
// 	}
// 	// Bulk
// 	var products []Product
// 	if len(keys) > 0 {
// 		products = state.GetBulkStateDefault[Product](ctx, keys, Product{})
// 	} else {
// 		products = make([]Product, len(keys))
// 	}
// 	// Prior non-bulk implementation
// 	//for _, id := range keys {
// 	//	product, err := state.GetState[Product](ctx, id)
// 	//	if err != nil {
// 	//		panic(err)
// 	//	}
// 	//	products = append(products, product)
// 	//}
// 	return products
}
