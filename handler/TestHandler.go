package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rudithu/CRYPTO-WalletApp/models"
)

type TestHandler struct {
	DB *sql.DB
}

func (h *TestHandler) HandleHello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userIdStr := r.URL.Query().Get("id")

	if userIdStr != "" {

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}

		user, err := models.GetUserById(h.DB, userId)

		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "Hello World %s!", user.Name)
		}
	} else {
		fmt.Fprint(w, "Hello World 2")
	}
}

func (h *TestHandler) HandleEcho(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

}
