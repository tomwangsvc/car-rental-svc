package context

import "context"

type contextKey int

const (
	ContextKeyIamTokenClientId           contextKey = iota
	ContextKeyIamTokenClientName         contextKey = iota
	ContextKeyIamTokenEmail              contextKey = iota
	ContextKeyIamTokenIdentityExternalId contextKey = iota
	ContextKeyIamTokenIdentityId         contextKey = iota
	ContextKeyIamTokenIdentityProvider   contextKey = iota
	ContextKeyIamTokenName               contextKey = iota
	ContextKeyIamTokenRoles              contextKey = iota
)

func iamNewContext(newContext, ctx context.Context) context.Context {
	newContext = WithIamTokenClientId(newContext, IamTokenClientId(ctx))
	newContext = WithIamTokenClientName(newContext, IamTokenClientName(ctx))
	newContext = WithIamTokenEmail(newContext, IamTokenEmail(ctx))
	newContext = WithIamTokenRoles(newContext, IamTokenRoles(ctx))
	newContext = WithIamTokenIdentityExternalId(newContext, IamTokenIdentityExternalId(ctx))
	newContext = WithIamTokenIdentityId(newContext, IamTokenIdentityId(ctx))
	newContext = WithIamTokenIdentityProvider(newContext, IamTokenIdentityProvider(ctx))
	newContext = WithIamTokenName(newContext, IamTokenName(ctx))
	return newContext
}

func IamTokenClientId(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyIamTokenClientId).(string)
	if !ok {
		return ""
	}
	return v
}

func WithIamTokenClientId(ctx context.Context, iamTokenClientId string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenClientId, iamTokenClientId)
}

func IamTokenClientName(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyIamTokenClientName).(string)
	if !ok {
		return ""
	}
	return v
}

func WithIamTokenClientName(ctx context.Context, iamTokenClientName string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenClientName, iamTokenClientName)
}

func IamTokenIdentityExternalId(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyIamTokenIdentityExternalId).(string)
	if !ok {
		return ""
	}
	return v
}

func IamTokenEmail(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyIamTokenEmail).(string)
	if !ok {
		return ""
	}
	return v
}

func WithIamTokenEmail(ctx context.Context, iamTokenEmail string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenEmail, iamTokenEmail)
}

func WithIamTokenIdentityExternalId(ctx context.Context, iamTokenIdentityExternalId string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenIdentityExternalId, iamTokenIdentityExternalId)
}

func IamTokenIdentityId(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyIamTokenIdentityId).(string)
	if !ok {
		return ""
	}
	return v
}

func WithIamTokenIdentityId(ctx context.Context, iamTokenIdentityId string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenIdentityId, iamTokenIdentityId)
}

func IamTokenIdentityProvider(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyIamTokenIdentityProvider).(string)
	if !ok {
		return ""
	}
	return v
}

func WithIamTokenIdentityProvider(ctx context.Context, iamTokenIdentityProvider string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenIdentityProvider, iamTokenIdentityProvider)
}

func IamTokenName(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyIamTokenName).(string)
	if !ok {
		return ""
	}
	return v
}

func WithIamTokenName(ctx context.Context, iamTokenName string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenName, iamTokenName)
}

func IamTokenRoles(ctx context.Context) []string {
	v, ok := ctx.Value(ContextKeyIamTokenRoles).([]string)
	if !ok {
		return nil
	}
	return v
}

func WithIamTokenRoles(ctx context.Context, iamTokenRoles []string) context.Context {
	return context.WithValue(ctx, ContextKeyIamTokenRoles, iamTokenRoles)
}
