package scontext

import (
	"context"
)

type languageCtxTag struct {
}

// SetLanguage set language to context
func SetLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, languageCtxTag{}, lang)
}

// GetLanguage get language by context
func GetLanguage(ctx context.Context) string {
	if val, ok := ctx.Value(languageCtxTag{}).(string); ok {
		return val
	}
	return "zh"
}

type userUuidCtxTag struct{}

// SetUserUuid set user uuid to context
func SetUserUuid(ctx context.Context, userUuid string) context.Context {
	return context.WithValue(ctx, userUuidCtxTag{}, userUuid)
}

// GetUserUuid get user uuid by context
func GetUserUuid(ctx context.Context) string {
	if val, ok := ctx.Value(userUuidCtxTag{}).(string); ok {
		return val
	}
	return ""
}
