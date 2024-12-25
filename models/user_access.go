package models

import (
	"fmt"
	"strconv"
)

// UserAccess представляет права доступа пользователя
type UserAccess struct {
	UserID    int `gorm:"primaryKey"`
	FeatureID int
}

func UserFeatureAccess(userID, featureID int) bool {

	// Проверяем Bloom filter
	key := strconv.Itoa(userID) + ":" + strconv.Itoa(featureID)
	if !bloomFilter.Test([]byte(key)) {
		return false
	}

	// Запрос в БД для подтверждения
	var access UserAccess
	result := db.Where("user_id = ? AND feature_id = ?",
		userID, featureID).First(&access)
	if result.Error != nil {
		fmt.Printf("БД: Пользователь %d не имеет доступа к %d.\n",
			userID, featureID)
		return false
	}

	fmt.Printf("БД: Пользователь %d имеет доступ к %d.\n", userID, featureID)
	return true
}
