#  系统云平台



## 如何使用上下文切换机构

1. 登录云平台，基于账号名称作为前置判断，自动选择切换机构 指向存在该用户的机构（即数据库）。

2. 前端 请求过程中，Header携带Token/platform 根据Token/platform自动切换机构（即数据库）。
### 请求头

| 参数名      | 类型     | 是否必填 | 描述                    | 示例值                         |
|----------|--------|------|-----------------------|-----------------------------|
| Token    | string | 必填   | 认证信息，通常是 Bearer Token | Bearer abcdef123            |
| platform | string | 非    | 机构信息                  | H65402100014 |



## BIZ 示例


### 数据库使用示例
```go
// c *gin.Context

ctx := c.Request.Context()

// gorm
db := uc.repo.GetDB(ctx)

// redis
rdb := uc.repo.GetRDB()

// redis/v8 原生方法
rdb.Set(ctx, "key", "value", 0)

// pkg/db 封装方法 若需要补充方法 添加到 pkg/db/rdb.go
err := rdb.SetWithContext(ctx, "key", "value", 0)

```
jias
## 功能覆盖对比（waiter/plaza/const.go 清单）

- 已覆盖（battle-tiles）:
  - 统计接口：`/stats/today`, `/stats/yesterday`, `/stats/week`, `/stats/lastweek`
  - 钱包/流水查询：`/members/wallet/get`, `/members/wallet/list`, `/members/ledger/list`
  - 店铺管理：管理员分配/撤销/列表；房间/成员管理（踢人、解散、列表等）

- 新增中：
  - 店铺设置：运费设置/查询、分运开关设置/取消、推送额度查询/设置（/shops/fees/*, /shops/sharefee/set, /shops/pushcredit/*）

- 待确认或未覆盖：
  - 约战战绩明细回放（需对接 foxuc HTTP 战绩接口并落库/映射）
