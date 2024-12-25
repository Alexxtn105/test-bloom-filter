package models

import (
	"fmt"
	"strconv"
)

func SeedDatabase() {
	var accessList []UserAccess

	// генерируем 500 записей со 100 пользователями и 10 фичами
	for i := 1; i <= 500; i++ {
		userID := (i-1)%100 + 1
		featureID := (i-1)%10 + 101
		accessList = append(accessList, UserAccess{
			UserID:    userID,
			FeatureID: featureID,
		})
	}

	for _, access := range accessList {
		db.FirstOrCreate(&access)
		// Добавляем комбинацию UserID и FeatureID в Bloom filter
		key := strconv.Itoa(access.UserID) + ":" + strconv.Itoa(access.FeatureID)
		bloomFilter.Add([]byte(key))
	}

	fmt.Println("База данных заполнена, bloom-фильтр инициализирован.")
}
