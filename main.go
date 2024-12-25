package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up HTTP server
	http.HandleFunc("/feature_access", UserFeatureCheck)
	http.HandleFunc("/estimate_fp", EstimateFPHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func UserFeatureCheck(w http.ResponseWriter, r *http.Request) {

}

func EstimateFPHandler(w http.ResponseWriter, r *http.Request) {

}
