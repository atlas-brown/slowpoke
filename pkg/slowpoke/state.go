package slowpoke

import (
    "sync"
    "context"
    "errors"
    "github.com/goccy/go-json"
)
 
var internalStates = sync.Map{}

func GetState[T interface{}](ctx context.Context, key string) (T, error) {
    value, ok := internalStates.Load(key)
    var result T
    if !ok {
        return result, errors.New("key not found")
    }
    err := json.Unmarshal(value.([]byte), &result)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func GetBulkState[T interface{}](ctx context.Context, keys []string) ([]T, error) {
    var results []T
    for _, key := range keys {
        value, err := GetState[T](ctx, key)
        if err != nil {
            return nil, err
        }
        results = append(results, value)
    }
    return results, nil
}

func SetState(ctx context.Context, key string, value interface{}) {
    valueBytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
    internalStates.Store(key, valueBytes)
}

func SetBulkState(ctx context.Context, values map[string]interface{}) {
    for key, value := range values {
        SetState(ctx, key, value)
    }
}