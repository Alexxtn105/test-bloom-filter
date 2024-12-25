package main

import (
	"encoding/json"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"log"
	"net/http"
	"strconv"
	"test-bloom-filter/models"
)

func init() {
	models.DBInit()

	models.BloomFilterSetup()

	// заполнение БД
	models.SeedDatabase()
}

func main() {
	// настраиваем сервер
	// эта ручка определяет, обладает ли пользователь указанной способностью
	http.HandleFunc("/feature_access", UserFeatureCheck)
	http.HandleFunc("/estimate_fp", EstimateFPHandler)

	fmt.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func UserFeatureCheck(w http.ResponseWriter, r *http.Request) {
	// Извлекаем userID из параметров запроса
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w,
			"Отсутствует параметр 'user_id'",
			http.StatusBadRequest)
		return
	}

	// Преобразуем userID в целое число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w,
			"Неверный 'userID'. Должен быть целым числом.",
			http.StatusBadRequest)
		return
	}

	// Извлекаем featureID из параметров запроса
	featureIDStr := r.URL.Query().Get("feature_id")
	if featureIDStr == "" {
		http.Error(w,
			"Отсутствует параметр 'feature_id'",
			http.StatusBadRequest)
		return
	}

	// Преобразуем featureID в целое число
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		http.Error(w,
			"Неверный 'featureID'. Должен быть целым числом.",
			http.StatusBadRequest)
		return
	}
	access := models.UserFeatureAccess(userID, featureID)

	// Создаем API-ответ
	response := map[string]bool{
		"access": access,
	}

	// Пишем ответ в виде JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EstimateAPIResponse структура для возврата ответа в виде json
type EstimateAPIResponse struct {
	EstimatedFPRate  float64 `json:"estimated_false_positive_rate,omitempty"`
	EstimatedM       uint    `json:"estimated_m,omitempty"`
	EstimatedK       uint    `json:"estimated_k,omitempty"`
	DesiredFPRate    float64 `json:"desired_fp_rate,omitempty"`
	ExpectedAccuracy string  `json:"expected_accuracy,omitempty"`
}

func EstimateFPHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим параметры запроса
	query := r.URL.Query()
	nStr := query.Get("n")
	fpStr := query.Get("desired_fp_rate")

	if nStr == "" || fpStr == "" {
		http.Error(w, "Отсутствует параметр запроса 'n' (число элементов) или 'desired_fp_rate'", http.StatusBadRequest)
		return
	}

	// Преобразуем в целые числа
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		http.Error(w, "'n' должно быть положительным целым", http.StatusBadRequest)
		return
	}

	fp, err := strconv.ParseFloat(fpStr, 64)
	if err != nil || fp <= 0 || fp >= 1 {
		http.Error(w, "'desired_fp_rate' должен быть числом с плавающей точкой в диапазоне от 0 до 1", http.StatusBadRequest)
		return
	}

	// Подсчет параметров Bloom filter
	m, k := bloom.EstimateParameters(uint(n), fp)

	// Валидация посчитанных параметров
	estimatedFP := bloom.EstimateFalsePositiveRate(m, k, uint(n))

	// Готовим ответ
	response := EstimateAPIResponse{
		EstimatedFPRate:  estimatedFP,
		EstimatedM:       m,
		EstimatedK:       k,
		DesiredFPRate:    fp,
		ExpectedAccuracy: fmt.Sprintf("Посчитаненый false positive rate в пределах %.2f%% желаемого.", (estimatedFP-fp)/fp*100),
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
