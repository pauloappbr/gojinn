package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Order struct {
	ID    string  `json:"id"`
	Value float64 `json:"value"`
}

func main() {
	http.HandleFunc("/api/bench", func(w http.ResponseWriter, r *http.Request) {
		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		tax := order.Value * 0.15
		total := order.Value + tax

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"order_id": "%s", "total_final": %.2f, "engine": "docker-native"}`, order.ID, total)
	})

	fmt.Println("Docker server running on port 8081...")
	http.ListenAndServe(":8081", nil)
}
