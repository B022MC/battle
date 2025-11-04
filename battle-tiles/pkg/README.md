# go-utils 工具库

## 目录结构
```shell
├─cmd // 命令行
│  └─template // 模版
├─logger
├─plugin
│  ├─dbx  // 基于数据库扩展
│  ├─ginx
│  ├─gormx
│  │  ├─base // 基础
│  │  ├─migrate // 迁移
│  │  └─repo
│  └─middleware
└─utils
    ├─ecode
    ├─lang  // 提供的一些常用方法
    │  ├─browserx
    │  ├─mapx
    │  ├─reflectx
    │  ├─stringx
    │  └─timex
    ├─render
    ├─request
    ├─response
    └─timeutil

```

## 如何使用基础存储库
CORMImpl 默认实现：基础的增加、删除、更新、查询、分页、批量插入、批量删除、批量更新、批量查询、批量查询分页、批量查询总数、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、批量查询总数分页、
```shell

type CORMImpl[T any] interface {
	DB(ctx context.Context) *gorm.DB
	_select(ctx context.Context, condition interface{}) ([]*T, error)
	Insert(ctx context.Context, m *T) (err error)
	BatchInsert(ctx context.Context, m []*T, batchSize int) (int64, error) // BatchInsert 批量插入
	DeleteByPK(ctx context.Context, pks interface{}) (int64, error)        // DeleteByPK 根据主键删除，支持单个主键或者一个主键数组
	DeleteByMap(ctx context.Context, condition map[string]interface{}) (int64, error)
	UpdateByPK(ctx context.Context, t *T) (int64, error)
	UpdateByPKWithMap(ctx context.Context, pk interface{}, updateData map[string]interface{}) (int64, error)             // UpdateByPKWithMap 根据id更新，支持零值
	UpdateByMap(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) (int64, error) // UpdateByMap 根据条件更新，支持零值
	SelectOne(ctx context.Context, condition *T) (*T, error)
	SelectOneByPK(ctx context.Context, pk interface{}) (*T, error)                    // SelectOne 条件不能是零值，如果要查零值，请用 SelectOneByMap
	SelectOneByMap(ctx context.Context, condition map[string]interface{}) (*T, error) // SelectOneByPK 根据主键查找
	selectOne(ctx context.Context, condition interface{}) (*T, error)                 // SelectOneByMap 根据条件查找，支持零值
	Select(ctx context.Context, condition *T) ([]*T, error)                           // Select 根据非空字段查询
	SelectAll(ctx context.Context) ([]*T, error)                                      // SelectAll 查询所有
	SelectByPK(ctx context.Context, pks interface{}) ([]*T, error)                    // SelectByPK 根据主键查找，支持单个主键或者一个主键数组
	SelectByMap(ctx context.Context, condition map[string]interface{}) ([]*T, error)  // SelectByMap 根据条件查找，支持零值
	ListPage(ctx context.Context, query *gorm.DB, page *utils2.PageParam) ([]*T, int64, error)
	ListPageByOption(ctx context.Context, query *gorm.DB, page *utils2.PageParam, target interface{}) (int64, error)
	PageSelect(ctx context.Context, page *utils2.PageParam, query interface{}, args ...interface{}) ([]*T, int64, error)
}
```
Repo层使用
```shell

type BasicUserRepo interface {
	repox.CORMImpl[basicModel.BasicUser]
}

type basicUserRepo struct {
	repox.CORMImpl[basicModel.BasicUser]
	data *infra.Data
	log  *log.Helper
}

// NewBasicUserRepo .
func NewBasicUseRepo(data *infra.Data, logger log.Logger) BasicUserRepo {
	return &basicUserRepo{
		CORMImpl: repox.NewCORMImplRepo[basicModel.BasicUser](data),
		data:     data,
		log:      log.NewHelper(log.With(logger, "module", "repo/basicUser")),
	}
}

```


## 如何基于Model 生成 biz/repo/req/vo/service 模版文件
指定模板包名、模板包下要生成的model名称 即可 详细参考 cmd/tmpl.go 文件
```shell
generate tmpl --pkg "basic" --models "BasicUser"
```
