package main

import (
	"fmt"
	"log-analyzer-bot/internal/service"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file tidak ditemukan")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		fmt.Println("Error: GEMINI_API_KEY tidak di set di env variable")
		return
	}

	logContent, err := os.ReadFile("Logs/sample.log")

	if err != nil {
		fmt.Println("Error saat membaca file:", err)
		return
	}

	prompt := fmt.Sprintf(`Analisa log berikut, kasih ringkasan:
		- Ada error apa aja dan berapa kali muncul
		- Kemungkinan root cause
		- Rekomendasi tindakan

	Log:
	%s`, string(logContent))

	result, err := service.AskGemini(apiKey, prompt)

	if err != nil {
		fmt.Println("Error saat memanggil Gemini API:", err)
		return
	}
	fmt.Println("=== Hasil Analisa ===")
	fmt.Println(result)
}
