package clog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanonical_SetInt(t *testing.T) {
	ctx := context.Background()
	ctx = Init(ctx)
	SetInt(ctx, "foo.bar", 1)
	SetInt(ctx, "foo.baz", 2)
	require.Equal(t, `{"foo":{"bar":1,"baz":2}}`, MarshalJSON(ctx))
}

func TestCanonical_SetFloat64(t *testing.T) {
	ctx := context.Background()
	ctx = Init(ctx)
	SetFloat64(ctx, "foo.bar", 1.1)
	SetFloat64(ctx, "foo.baz", 2.1)
	require.Equal(t, `{"foo":{"bar":1.1,"baz":2.1}}`, MarshalJSON(ctx))
}

func TestCanonical_Nested(t *testing.T) {
	ctx := context.Background()
	ctx = Init(ctx)
	SetString(ctx, "request_id", "req-123")

	func() {
		ctx1 := context.WithValue(ctx, "foo", "bar")

		SetString(ctx1, "user_id", "123")
	}()

	require.Equal(t, `{"request_id":"req-123","user_id":"123"}`, MarshalJSON(ctx))
}

func TestCanonical_AddInt(t *testing.T) {
	ctx := context.Background()
	ctx = Init(ctx)
	AddInt(ctx, "foo.bar", 1)
	AddInt(ctx, "foo.bar", 1)

	require.Equal(t, `{"foo":{"bar":2}}`, MarshalJSON(ctx))
}

func TestCanonical_AddFloat64(t *testing.T) {
	ctx := context.Background()
	ctx = Init(ctx)
	AddFloat64(ctx, "foo.bar", 1.1)
	AddFloat64(ctx, "foo.bar", 1.2)

	require.Equal(t, `{"foo":{"bar":2.3}}`, MarshalJSON(ctx))
}

func TestCanonical_MultiNested(t *testing.T) {
	ctx := context.Background()
	ctx = Init(ctx)
	SetString(ctx, "http.request.path", "/foo")
	SetString(ctx, "http.request.code", "200")
	SetInt(ctx, "http.response.duration_ms", 10)

	require.Equal(t, `{"http":{"request":{"path":"/foo","code":"200"},"response":{"duration_ms":10}}}`, MarshalJSON(ctx))
}
