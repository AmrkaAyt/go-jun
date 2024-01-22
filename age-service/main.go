package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/enrich/age", enrichAge)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func enrichAge(w http.ResponseWriter, r *http.Request) {
	// Implement age enrichment logic here
	// Fetch data from the agify API
	// ...
	// Respond with enriched data
	json.NewEncoder(w).Encode(map[string]interface{}{"age": 30})
}
