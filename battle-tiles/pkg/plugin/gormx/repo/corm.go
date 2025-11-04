package repo

import (
	utils2 "battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/lang/reflectx"
	"battle-tiles/pkg/utils/lang/stringx"
	"context"
	"github.com/pkg/errors"
	"go/ast"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
	"reflect"
	"strings"
)

// 只要实现了DataImpl接口，那么就可以使用DBWithContext()
type DataImpl interface {
	GetDBWithContext(ctx context.Context) *gorm.DB
}

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

type CORMImplRepo[T any] struct {
	data DataImpl
	// 泛型参数代表的struct名称，例如：BaseRepo[Pop]
	StructName string
	PrimaryKey string
}

func (b *CORMImplRepo[T]) DB(ctx context.Context) *gorm.DB {
	return b.data.GetDBWithContext(ctx)
}

func (b *CORMImplRepo[T]) recursiveDeleteAutoTime(v reflect.Value, updateData map[string]interface{}) {
	t := reflectx.IndirectType(v.Type())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			b.recursiveDeleteAutoTime(v.Field(i), updateData)
		} else {
			tag := field.Tag.Get("gorm")
			autoCreateTime := strings.Contains(tag, "autoCreateTime") &&
				!strings.Contains(tag, "autoCreateTime:false")

			autoUpdateTime := strings.Contains(tag, "autoUpdateTime") &&
				!strings.Contains(tag, "autoUpdateTime:false")
			if autoCreateTime || autoUpdateTime {
				delete(updateData, field.Name)
			}
		}
	}
}

func (b *CORMImplRepo[T]) deleteAutoTime(updateData map[string]interface{}) {
	var m T
	b.recursiveDeleteAutoTime(reflect.ValueOf(m), updateData)
}

func (b *CORMImplRepo[T]) _select(ctx context.Context, condition interface{}) ([]*T, error) {
	var (
		m   T
		res []*T
	)
	if err := b.DB(ctx).Model(&m).Where(condition).Find(&res).Error; err != nil {
		return nil, errors.Wrapf(err, "dbx: select %s error, condition: %+v", b.StructName, condition)
	}
	return res, nil
}

func (b *CORMImplRepo[T]) Insert(ctx context.Context, m *T) (err error) {
	if err = b.DB(ctx).Create(m).Error; err != nil {
		err = errors.Wrapf(err, "dbx: insert %s error, param: %+v", b.StructName, m)
	}
	return
}

func (b *CORMImplRepo[T]) BatchInsert(ctx context.Context, m []*T, batchSize int) (int64, error) {
	tx := b.DB(ctx).CreateInBatches(m, batchSize)
	if tx.Error != nil {
		//tx.Error

		return 0, errors.Wrapf(tx.Error, "dbx: batch insert %s error, param: %+v", b.StructName, m)
	}
	return tx.RowsAffected, nil
}

func (b *CORMImplRepo[T]) DeleteByPK(ctx context.Context, pks interface{}) (int64, error) {
	var m T
	tx := b.DB(ctx).Where(map[string]interface{}{
		b.PrimaryKey: pks,
	}).Delete(&m)
	if err := tx.Error; err != nil {
		return 0, errors.Wrapf(err, "dbx: delete %s by pks error, pks: %v", b.StructName, pks)
	}
	return tx.RowsAffected, nil
}

func (b *CORMImplRepo[T]) DeleteByMap(ctx context.Context, condition map[string]interface{}) (int64, error) {
	c := camel2SnakeForMapKey(condition)
	var m T
	tx := b.DB(ctx).Where(c).Delete(&m)
	if err := tx.Error; err != nil {
		return 0, errors.Wrapf(err, "dbx: delete %s by map error, condition: %v", b.StructName, condition)
	}
	return tx.RowsAffected, nil
}

func (b *CORMImplRepo[T]) UpdateByPK(ctx context.Context, t *T) (int64, error) {
	tx := b.DB(ctx).Model(t).Updates(t)
	if err := tx.Error; err != nil {
		return 0, errors.Wrapf(err, "dbx: update %s by pk error, param: %+v", b.StructName, t)
	}
	return tx.RowsAffected, nil
}

func (b *CORMImplRepo[T]) UpdateByPKWithMap(ctx context.Context, pk interface{}, updateData map[string]interface{}) (int64, error) {
	return b.UpdateByMap(ctx, map[string]interface{}{b.PrimaryKey: pk}, updateData)
}

func (b *CORMImplRepo[T]) UpdateByMap(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	c := camel2SnakeForMapKey(condition)
	b.deleteAutoTime(updateData)

	var m T
	tx := b.DB(ctx).Model(&m).Where(c).Updates(updateData)
	if err := tx.Error; err != nil {
		return 0, errors.Wrapf(err, "dbx: update %s by map error, condition: %v, updateData: %v", b.StructName, c, updateData)
	}
	return tx.RowsAffected, nil
}

func (b *CORMImplRepo[T]) SelectOne(ctx context.Context, condition *T) (*T, error) {
	return b.selectOne(ctx, condition)
}

func (b *CORMImplRepo[T]) SelectOneByPK(ctx context.Context, pk interface{}) (*T, error) {
	return b.SelectOneByMap(ctx, map[string]interface{}{b.PrimaryKey: pk})
}

func (b *CORMImplRepo[T]) SelectOneByMap(ctx context.Context, condition map[string]interface{}) (*T, error) {
	c := camel2SnakeForMapKey(condition)
	return b.selectOne(ctx, c)
}

func (b *CORMImplRepo[T]) selectOne(ctx context.Context, condition interface{}) (*T, error) {
	res, err := b._select(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	if len(res) > 1 {
		return nil, errors.Errorf("dbx: select one %s error, result must be one, now it is %d, condition %+v", b.StructName, len(res), condition)
	}
	return res[0], err
}

func (b *CORMImplRepo[T]) Select(ctx context.Context, condition *T) ([]*T, error) {
	return b._select(ctx, condition)
}

func (b *CORMImplRepo[T]) SelectAll(ctx context.Context) ([]*T, error) {
	return b.SelectByMap(ctx, map[string]interface{}{})
}

func (b *CORMImplRepo[T]) SelectByPK(ctx context.Context, pks interface{}) ([]*T, error) {
	return b.SelectByMap(ctx, map[string]interface{}{b.PrimaryKey: pks})
}

func (b *CORMImplRepo[T]) SelectByMap(ctx context.Context, condition map[string]interface{}) ([]*T, error) {
	c := camel2SnakeForMapKey(condition)
	return b._select(ctx, c)
}

func (b *CORMImplRepo[T]) ListPage(ctx context.Context, query *gorm.DB, page *utils2.PageParam) ([]*T, int64, error) {
	var (
		total int64
		res   []*T
	)
	if page != nil {
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, errors.Wrapf(err, "dbx: select count %s error", b.StructName)
		}
	}
	var orders []string
	if page != nil {
		if !page.NotPage {
			query = query.Offset(int(page.PageNo-1) * int(page.PageSize)).Limit(int(page.PageSize))
		}
		if page.OrderBy != "" {
			query = query.Order(page.OrderBy)
		}
	}
	if err := query.Find(&res).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "dbx: select %s error, orders: %s", b.StructName, strings.Join(orders, ","))
	}
	if total == 0 {
		total = int64(len(res))
	}
	return res, total, nil
}

func (b *CORMImplRepo[T]) ListPageByOption(ctx context.Context, query *gorm.DB, page *utils2.PageParam, target interface{}) (int64, error) {
	var (
		total int64
	)
	if page != nil {
		if err := query.Count(&total).Error; err != nil {
			return 0, errors.Wrapf(err, "dbx: select count %s error", b.StructName)
		}
	}
	var orders []string
	if page != nil {
		if !page.NotPage {
			query = query.Offset(int(page.PageNo-1) * int(page.PageSize)).Limit(int(page.PageSize))
		}
		if page.OrderBy != "" {
			query = query.Order(page.OrderBy)
		}
	}
	if err := query.Scan(target).Error; err != nil {
		return 0, errors.Wrapf(err, "dbx: select %s error, orders: %s", b.StructName, strings.Join(orders, ","))
	}
	return total, nil
}

func (b *CORMImplRepo[T]) PageSelect(ctx context.Context, page *utils2.PageParam, query interface{}, args ...interface{}) ([]*T, int64, error) {
	var (
		m     T
		total int64
		res   []*T
	)
	if page != nil {
		if err := b.DB(ctx).Model(&m).Where(query, args...).Count(&total).Error; err != nil {
			return nil, 0, errors.Wrapf(err, "dbx: select count %s error, query: %+v, args: %+v", b.StructName, query, args)
		}
	}
	q := b.DB(ctx).Model(&m).Where(query, args...)
	if page != nil {
		if !page.NotPage {
			q = q.Offset(int(page.PageNo-1) * int(page.PageSize)).Limit(int(page.PageSize))
		}

		if page.OrderBy != "" {
			q = q.Order(page.OrderBy)
		}
	}
	if err := q.Find(&res).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "dbx: select %s error, query: %+v, args: %+v", b.StructName, query, args)
	}
	if total == 0 {
		total = int64(len(res))
	}
	return res, total, nil
}

// NewCORMImplRepo 这个函数的意义在于不暴露db进行初始化，外部只能通过函数DB()获取
func NewCORMImplRepo[T any](data DataImpl) CORMImpl[T] {
	b := CORMImplRepo[T]{
		data: data,
	}
	var m T
	b.StructName = reflect.ValueOf(m).Type().Name()
	b.PrimaryKey = b.parsePrimaryKey()
	return &b
}

func (b *CORMImplRepo[T]) parsePrimaryKey() string {
	var m T
	return recursiveParsePrimaryKey(reflect.ValueOf(m))
}

func recursiveParsePrimaryKey(reflectValue reflect.Value) string {
	reflectType := reflectx.IndirectType(reflectValue.Type())
	var hasId bool
	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
			if fieldStruct.Anonymous {
				res := recursiveParsePrimaryKey(reflectValue.Field(i))
				if res != "" {
					return res
				}
			} else {
				tagSetting := schema.ParseTagSetting(fieldStruct.Tag.Get("gorm"), ";")

				// 数据库字段名
				columnName := tagSetting["COLUMN"]
				if columnName == "" {
					columnName = schema.NamingStrategy{}.ColumnName("", fieldStruct.Name)
				}

				if utils.CheckTruth(tagSetting["PRIMARYKEY"], tagSetting["PRIMARY_KEY"]) {
					return columnName
				}

				if columnName == "id" {
					hasId = true
				}
			}
		}
	}
	if hasId {
		return "id"
	}
	return ""
}

// camel2SnakeForMapKey 将map的key转换为下划线格式
func camel2SnakeForMapKey(condition map[string]interface{}) map[string]interface{} {
	c := make(map[string]interface{})
	for k, v := range condition {
		c[stringx.Camel2Snake(k)] = v
	}
	return c
}
