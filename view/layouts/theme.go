package layouts

import "context"

type themeContextKey struct{}

func WithTheme(ctx context.Context, isDark bool) context.Context {
	return context.WithValue(ctx, themeContextKey{}, isDark)
}

func isDarkFromContext(ctx context.Context) bool {
	isDark, _ := ctx.Value(themeContextKey{}).(bool)
	return isDark
}
