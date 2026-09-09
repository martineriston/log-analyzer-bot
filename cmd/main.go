package main

import (
	"encoding/json"
	"fmt"
	"log-analyzer-bot/internal/model"
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

	prompt := fmt.Sprintf(`Analisa log berikut. Kembalikan HANYA dalam format JSON, isinya harus murni tanpa tambahan apapun dan wajib dikembalikan sesuai struktur ini:
	{
		"issues": [
			{
			"error": "pesan error singkat",
			"occurrences": jumlah_kemunculan_sebagai_angka,
			"root_cause": "kemungkinan penyebab",
			"recommendation": "saran tindakan"
			}
		]
	}
	Log:
	%s`, string(logContent))

	result, err := service.AskGemini(apiKey, prompt)

	if err != nil {
		fmt.Println("Error saat memanggil Gemini API:", err)
		return
	}

	var resp model.AnalysisResult

	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		fmt.Println("Error saat parsing hasil JSON:", err)
		return
	}

	printAnalysis(resp)
}

func printAnalysis(resp model.AnalysisResult) {
	fmt.Println("=== Hasil Analisa ===")
	for i, issue := range resp.Issues {
		fmt.Printf("\n[%d] %s\n", i+1, issue.Error)
		fmt.Printf("	Terjadi: %d kali\n", issue.Occurrences)
		fmt.Printf("	Root Cause: %s\n", issue.RootCause)
		fmt.Printf("	Rekomendasi: %s\n", issue.Recommendation)
	}
}
