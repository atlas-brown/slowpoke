package wrappers

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/goccy/go-json"
	"hash/fnv"
	"net/http"
	"strconv"
	"io"
)

var DEBUG_CA = false

func HashCallArgs(app string, method string, input []byte) string {
	if DEBUG_CA {
		return fmt.Sprintf("%s.%s.%s", app, method, input)
	}
	h := fnv.New32a()
	h.Write([]byte(app))
	h.Write([]byte(method))
	h.Write(input)
	inputHash := h.Sum32()
	return fmt.Sprintf("%x", inputHash)
}

func HashCallArgsWithShard(app string, method string, input []byte) (string, int) {
	utility.Assert(common.ShardEnabled)
	shard, err := strconv.Atoi(common.ShardCount)
	if err != nil {
		panic(err)
	}
	h := fnv.New32a()
	h.Write([]byte(app))
	h.Write([]byte(method))
	h.Write(input)
	inputHash := h.Sum32()
	// shard is in [1, shardCount]
	if DEBUG_CA {
		return fmt.Sprintf("%s.%s.%s", app, method, input), int(inputHash%uint32(shard) + 1)
	}
	return fmt.Sprintf("%x", inputHash), int(inputHash%uint32(shard) + 1)
}

func PreReqStart(ctx context.Context) {
	if ReadOnlyContext(ctx) && CtxCaller(ctx) != "client" {
		// Initialize the dependencies for this call
		callId := CtxCallId(ctx)
		deps.InitDep(callId)
		// Inform the cache manager that this call has started
		// TODO: Do we need to send a callId to the cache-manager too?
		callArgs := CtxCallArgs(ctx)
		if common.ZMQ {
			cm.SendRequestZmq(&cm.StartRequest{CallArgs: callArgs}, cm.TypeStartRequest)
		} else {
			cm.SendStartRequestHttp(&cm.StartRequest{CallArgs: callArgs}, common.CMUrl)
		}
	}
}

// Could also happen after read (or at the same time)
func PreRead(ctx context.Context, key cm.Key) {
	if ReadOnlyContext(ctx) && CtxCaller(ctx) != "client" {
		// Initialize the dependencies for this call
		callId := CtxCallId(ctx)
		deps.AddKeyDep(callId, key)
	}
}

func PreWrite(ctx context.Context, key cm.Key) {
	if common.ZMQ {
		cm.SendRequestZmq(&cm.InvalidateKeyRequest{Key: key}, cm.TypeInvRequest)
	} else {
		cm.SendInvRequestHttp(&cm.InvalidateKeyRequest{Key: key}, common.CMUrl)
	}
}

func PostWrite(ctx context.Context, key cm.Key) {
	if common.ZMQ {
		cm.SendRequestZmq(&cm.InvalidateKeyRequest{Key: key}, cm.TypeInvRequest)
	} else {
		cm.SendInvRequestHttp(&cm.InvalidateKeyRequest{Key: key}, common.CMUrl)
	}
}

func PreReqEnd(ctx context.Context, retVal cm.ReturnVal) {
	if ReadOnlyContext(ctx) && CtxCaller(ctx) != "client" {
		// Get the dependencies
		callId := CtxCallId(ctx)
		keyDeps, callDeps := deps.PopDeps(callId)

		callArgs := CtxCallArgs(ctx)
		currServiceName := CtxCaller(ctx)
		endReq := cm.EndRequest{CallArgs: callArgs, KeyDeps: keyDeps, CallDeps: callDeps, Caller: currServiceName, ReturnVal: retVal}
		if common.ZMQ {
			cm.SendRequestZmq(&endReq, cm.TypeEndRequest)
		} else {
			cm.SendEndRequestHttp(&endReq, common.CMUrl)
		}
	}
}

// This returns true if cache was hit. The caller of this method needs to avoid the call if that
// is the case.
func PreCall(ctx context.Context, ca cm.CallArgs) (cm.ReturnVal, bool) {
	// Add call dependencies
	if ReadOnlyContext(ctx) && CtxCaller(ctx) != "client" {
		// Initialize the dependencies for this call
		callId := CtxCallId(ctx)
		deps.AddCallDep(callId, ca)
	}
	//start := time.Now()
	mc := cm.GetOrCreateCacheClient()
	ret, exists := cm.CacheGet(mc, ca)
	//if time.Since(start) > 1*time.Millisecond {
	//	glog.Info("CacheGet took ", time.Since(start), ", ", runtime.NumGoroutine())
	//}
	return ret, exists
}

func ROWrapper[ReqType interface{}, RespType interface{}](handler func(context.Context, *ReqType) *RespType) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		input, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		var req ReqType
		err = json.Unmarshal(input, &req)
		if err != nil {
			panic(err)
		}
		resp := handler(ctx, &req)
		utility.DumpJson(resp, w)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		slowpoke.SlowpokeDelay()
	}
}
func NonROWrapper[ReqType interface{}, RespType interface{}](handler func(context.Context, *ReqType) *RespType) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, input := SetupCtxFromHTTPReq(r, false)
		var req ReqType
		err := json.Unmarshal(input, &req)
		if err != nil {
			panic(err)
		}
		resp := handler(ctx, &req)
		//respByte, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		utility.DumpJson(resp, w)
		if f, ok := w.(http.Flusher); ok {
			f.Flush() // Force buffer to flush
		}
		slowpoke.SlowpokeDelay()
	}
}

func SlowpokeWrapper[ReqType interface{}, RespType interface{}](handler func(context.Context, *ReqType) *RespType, endpointName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slowpoke.SlowpokeCheck(endpointName)
		ctx := r.Context()
		input, err := io.ReadAll(r.Body)
		r.Body.Close()
		var req ReqType
		err = json.Unmarshal(input, &req)
		if err != nil {
			panic(err)
		}
		resp := handler(ctx, &req)
		utility.DumpJson(resp, w)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		slowpoke.SlowpokeDelay()
	}
}