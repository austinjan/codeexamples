package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pquerna/otp/totp"
	"github.com/rs/cors"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func generateSecret(w http.ResponseWriter, r *http.Request) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MyApp",
		AccountName: "test@example.com",
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := Response{
		Status:  "success",
		Message: "Secret generated",
		Data:    secret.Secret(), // 修改此行，返回 base32 编码的密钥

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateCode(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	secret := data["secret"]
	code := data["code"]

	valid := totp.Validate(code, secret)

	response := Response{
		Status:  "success",
		Message: fmt.Sprintf("Code is valid: %t", valid),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/generate-secret", generateSecret).Methods("GET")
	r.HandleFunc("/validate-code", validateCode).Methods("POST")
	// 添加 CORS 中间件
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsWrapper.Handler(r))
}
