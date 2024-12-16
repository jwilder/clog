package clog

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/wk8/go-ordered-map/v2"
)

// Package clog provides a canonical logging context for Go applications for easily
// constructing wide events.  Wide events, also called canonical logs or events, are
// a way to log structured data that can be easily queried and analyzed.  This package
// makes it easy to construct wide events by providing a simple API for setting and
// retrieving structured logging data.  Using canonical complements (and can sometimes replace)
// typical places where you might use metrics, logs or traces.  This approach is central
// to the Observability 2.0 ideas.
//
// A canonical log is a JSON object that contains a set of key-value pairs.  The keys
// are strings that are dot-separated.  The values are strings, integers, or floats.  The
// keys are case-insensitive.  The values are case-sensitive. The logging context is propagated through
// a context.Context so that different parts of the application can add to the same logging context.
// At the end of the unit of work, the logging context can be marshaled to JSON and logged by the application.
//
//	// Initialize the canonical logging context
//	ctx := clog.Init(context.Background())
//
//	// Set some values in the logging context.  As long as the context is passed around, these values
//	// will be aggregated into the same logging context.
//	clog.SetString(ctx, "http.request.method", "GET")
//	clog.SetString(ctx, "http.request.path", "/example")
//	clog.SetInt(ctx, "http.response.status_code", 200)
//	clog.SetFloat64(ctx, "http.response.duration_ms", 123.45)
//
//	// Add values to existing keys
//	clog.AddInt(ctx, "http.response.status_code", 1)
//	clog.AddFloat64(ctx, "http.response.duration_ms", 10.5)
//
//	// Marshal the logging context to JSON
//	log := clog.MarshalJSON(ctx)
//	fmt.Println(log)

const (
	contextKey = "__clog__"
)

type canonical struct {
	values *orderedmap.OrderedMap[string, any] //nolint:typecheck
}

func newCanonical() *canonical {
	return &canonical{
		values: orderedmap.New[string, any](), //nolint:typecheck
	}
}

// Init initializes the canonical logging context.  This must be called before any other canonical logging functions
// are called.  This is typically called at the beginning of a request handler or the beginning of a background task.
func Init(ctx context.Context) context.Context {
	v := ctx.Value(contextKey)
	if v == nil {
		ctx = context.WithValue(ctx, contextKey, newCanonical())
	}
	return ctx
}

// MarshalJSON returns the canonical logging context as a JSON string.
func MarshalJSON(ctx context.Context) string {
	if c, ok := ctx.Value(contextKey).(*canonical); ok {
		return c.string()
	}
	return ""
}

// SetString sets a string value in the canonical logging context.  If the string exists, it will be overwritten.
func SetString(ctx context.Context, key, value string) {
	if c, ok := ctx.Value(contextKey).(*canonical); ok {
		c.setString(key, value)
	}
}

// SetInt sets an int value in the canonical logging context.  If the int exists, it will be overwritten.
func SetInt(ctx context.Context, key string, value int) {
	if c, ok := ctx.Value(contextKey).(*canonical); ok {
		c.setInt(key, value)
	}
}

// SetFloat64 sets a float64 value in the canonical logging context.  If the float64 exists, it will be overwritten.
func SetFloat64(ctx context.Context, key string, value float64) {
	if c, ok := ctx.Value(contextKey).(*canonical); ok {
		c.setFloat64(key, value)
	}
}

// AddInt adds an int value to the canonical logging context.  If the int does not exist, it will be created.
func AddInt(ctx context.Context, key string, value int) {
	if c, ok := ctx.Value(contextKey).(*canonical); ok {
		c.addInt(key, value)
	}
}

// AddFloat64 adds a float64 value to the canonical logging context.  If the float64 does not exist, it will be created.
func AddFloat64(ctx context.Context, key string, value float64) {
	if c, ok := ctx.Value(contextKey).(*canonical); ok {
		c.addFloat64(key, value)
	}
}

func (c *canonical) normalizeKey(key string) []string {
	return strings.Split(strings.ToLower(key), ".")
}

func (c *canonical) setString(key string, value string) {
	c.set(c.normalizeKey(key), c.values, value)
}

func (c *canonical) setInt(key string, value int) {
	c.set(c.normalizeKey(key), c.values, value)
}

func (c *canonical) setFloat64(key string, value float64) {
	c.set(c.normalizeKey(key), c.values, value)
}

func (c *canonical) set(parts []string, state *orderedmap.OrderedMap[string, any], value any) { //nolint:typecheck
	if len(parts) == 1 {
		state.Set(parts[0], value)
		return
	}

	val, ok := state.Get(parts[0])
	if !ok {
		val = orderedmap.New[string, any]() //nolint:typecheck
		state.Set(parts[0], val)
	}
	c.set(parts[1:], val.(*orderedmap.OrderedMap[string, any]), value)
}

func (c *canonical) addInt(key string, value int) {
	c.add(c.normalizeKey(key), c.values, value)
}

func (c *canonical) add(parts []string, state *orderedmap.OrderedMap[string, any], value int) { //nolint:typecheck
	if len(parts) == 1 {
		val, ok := state.Get(parts[0])
		if !ok {
			state.Set(parts[0], value)
		}
		if vv, ok := val.(int); ok {
			state.Set(parts[0], vv+value)
		}
		return
	}

	val, ok := state.Get(parts[0])
	if !ok {
		val = orderedmap.New[string, any]()
		state.Set(parts[0], val)
	}
	c.add(parts[1:], val.(*orderedmap.OrderedMap[string, any]), value)
}

func (c *canonical) addFloat64(key string, value float64) {
	c.addFloat(c.normalizeKey(key), c.values, value)
}

func (c *canonical) addFloat(parts []string, state *orderedmap.OrderedMap[string, any], value float64) {
	if len(parts) == 1 {
		val, ok := state.Get(parts[0])
		if !ok {
			state.Set(parts[0], value)
		}
		if vv, ok := val.(float64); ok {
			state.Set(parts[0], vv+value)
		}
		return
	}

	val, ok := state.Get(parts[0])
	if !ok {
		val = orderedmap.New[string, any]()
		state.Set(parts[0], val)
	}
	c.addFloat(parts[1:], val.(*orderedmap.OrderedMap[string, any]), value)
}

func (c *canonical) string() string {
	b, _ := json.Marshal(c.values)
	return string(b)
}
