package dbx

func FetchScan[T any](sql string, values ...interface{}) ([]T, error) {
	ConnPool.mu.Lock()
	result := make([]T, 0)
	for _, db := range ConnPool.dbMap {
		list := make([]T, 0)
		if err := db.Raw(sql, values...).Scan(&list).Error; err != nil {
			continue
		}
		result = append(result, list...)
	}
	ConnPool.mu.Unlock()
	return result, nil
}

func FetchScanByDB[T any](platform, sql string, values ...interface{}) ([]T, error) {
	ConnPool.mu.Lock()
	db := ConnPool.dbMap[platform]
	defer ConnPool.mu.Unlock()

	list := make([]T, 0)
	if err := db.Raw(sql, values...).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
