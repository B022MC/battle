package dbx

import "context"

func NewCtxWithDB(ctx context.Context) context.Context {
	if ctx.Value(CtxDBKey) == nil {
		return ctx
	}
	dbKey, ok := ctx.Value(CtxDBKey).(string)
	if !ok {
		return ctx
	}

	return context.WithValue(context.Background(), CtxDBKey, dbKey)
}

func GetDBKeyFromCtx(ctx context.Context) string {
	dbKey, ok := ctx.Value(CtxDBKey).(string)
	if !ok {
		return ""
	}
	return dbKey
}
