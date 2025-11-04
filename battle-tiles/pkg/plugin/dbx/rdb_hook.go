package dbx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"strings"
)

const SkipPrefixKey string = "skipPrefix"

// 定义一个上下文辅助函数
func WithSkipPrefix(ctx context.Context) context.Context {
	return context.WithValue(ctx, SkipPrefixKey, true)
}

func shouldSkipPrefix(ctx context.Context) bool {
	value := ctx.Value(SkipPrefixKey)
	skip, ok := value.(bool)
	return ok && skip
}

// 允许适配前缀的命令
var commandsWithPrefix = []string{
	"GET", "SET", "EXISTS", "DEL", "TYPE",
	"RPUSH", "LPOP", "RPOP", "LLEN", "LRANGE",
	"SADD", "SREM", "SISMEMBER", "SMEMBERS", "SCARD",
	"HSET", "HMSET", "HGET", "HGETALL",
	"ZADD", "ZRANGE", "ZRANGEBYSCORE", "ZREVRANGEBYSCORE", "ZREM",
	"INCR", "INCRBY", "INCRBYFLOAT",
	"WATCH", "MULTI", "EXEC", "EXPIRE",
}

// 公共前缀处理函数
func addPrefixToArgs(ctx context.Context, cmd redis.Cmder) {
	// 直接改变args变量即可，因为内存地址是一样
	args := cmd.Args()
	if len(args) <= 1 {
		return
	}
	prefix := GetDBKeyFromCtx(ctx)
	name := strings.ToUpper(cmd.Name())
	switch name {
	case "MGET", "DEL":
		for i := 1; i < len(args); i++ {
			args[i] = prefix + "_" + cast.ToString(args[i])
		}
	case "MSET":
		for i := 1; i < len(args); i += 2 {
			args[i] = prefix + "_" + cast.ToString(args[i])
		}
	case "SCAN":
		if len(args) > 2 {
			for i := 2; i < len(args); i += 2 {
				if args[i] == "match" && i+1 < len(args) {
					args[i+1] = prefix + "_" + cast.ToString(args[i+1])
					break
				}
			}
		}
	default:
		if lo.IndexOf[string](commandsWithPrefix, name) != -1 {
			args[1] = prefix + "_" + cast.ToString(args[1])
		}
	}
}

type WithPlatformKeyHook struct {
}

func (h *WithPlatformKeyHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !shouldSkipPrefix(ctx) {
		addPrefixToArgs(ctx, cmd)
	}
	return ctx, nil
}

func (h *WithPlatformKeyHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {

	return nil
}

func (h *WithPlatformKeyHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !shouldSkipPrefix(ctx) {
		for _, cmd := range cmds {
			addPrefixToArgs(ctx, cmd)
		}
	}
	return ctx, nil

}

func (h *WithPlatformKeyHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func NewWithPlatformKeyHook() redis.Hook {
	return &WithPlatformKeyHook{}
}
