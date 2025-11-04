package migrate

import (
	"gorm.io/gorm"
	"sync"
)

var (
	models     []interface{}
	modelsLock sync.Mutex
)

// RegisterModel 注册模型
func RegisterModel(model interface{}) {
	modelsLock.Lock()
	defer modelsLock.Unlock()
	models = append(models, model)
}

// GetModels 获取所有注册的模型
func GetModels() []interface{} {
	modelsLock.Lock()
	defer modelsLock.Unlock()
	return models
}

// AutoMigrate 执行自动迁移
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(GetModels()...)
}
