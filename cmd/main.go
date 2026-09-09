package main

import (
	"encoding/json"
	"fmt"
	"log-analyzer-bot/internal/model"
	"log-analyzer-bot/internal/service"
	"os"
	"strings"

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

	filteredContent := FilterLog(string(logContent))

	chunks := service.ChunkLog(filteredContent, 2)

	var allResp []model.AnalysisResult

	for i, chunk := range chunks {
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
		%s`, chunk)

		result, err := service.AskGemini(apiKey, prompt)
		if err != nil {
			fmt.Printf("Error saat memanggil Gemini untuk chunk %d: %v\n", i+1, err)
			continue
		}

		var analysis model.AnalysisResult
		if err := json.Unmarshal([]byte(result), &analysis); err != nil {
			fmt.Printf("Error saat parsing hasil chunk %d: %v\n", i+1, err)
			continue
		}
		allResp = append(allResp, analysis)
	}

	printAnalysis(allResp)
}

func printAnalysis(resp []model.AnalysisResult) {
	if len(resp) == 0 {
		fmt.Println("Tidak ada hasil analisa yang ditemukan.")
		return
	}
	fmt.Println("=== Hasil Analisa ===")
	for _, issues := range resp {
		for i, issue := range issues.Issues {
			fmt.Printf("\n[%d] %s\n", i+1, issue.Error)
			fmt.Printf("	Terjadi: %d kali\n", issue.Occurrences)
			fmt.Printf("	Root Cause: %s\n", issue.RootCause)
			fmt.Printf("	Rekomendasi: %s\n", issue.Recommendation)
		}

	}
}
func FilterLog(logContent string) string {
	var filteredLines []string

	lines := strings.Split(logContent, "\n")

	for _, line := range lines {
		upperLine := strings.ToUpper(line)
		if !strings.Contains(upperLine, "ERROR") && !strings.Contains(upperLine, "WARN") {
			continue
		}
		filteredLines = append(filteredLines, line)
	}

	return strings.Join(filteredLines, "\n")
}
