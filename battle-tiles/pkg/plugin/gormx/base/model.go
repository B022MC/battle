package base

import (
	"battle-tiles/pkg/utils/lang/reflectx"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"reflect"
	"time"
)

// Model 实体类的基类
// 泛型参数 T 是 Id 字段的类型
type Model[T string | int32] struct {
	Id        T         `gorm:"primaryKey"  json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;type:timestamp with time zone;not null" json:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"false" format:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;type:timestamp with time zone;not null" json:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"false" format:"2006-01-02 15:04:05"` // 更新时间
	DeletedAt time.Time
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt" json:"is_del"` //是否软删除
}

// BeforeCreate 利用gorm的钩子，在insert前给string类型的Id填充uuid
func (r *Model[T]) BeforeCreate(tx *gorm.DB) (err error) {
	idValue := reflectx.Indirect(reflect.ValueOf(&r.Id))
	idInterface := idValue.Interface()
	if idStr, ok := idInterface.(string); ok {
		if idStr == "" {
			uuidStr := uuid.NewString()
			idValue.SetString(uuidStr)
		}
	}
	return
}
