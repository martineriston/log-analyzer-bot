# Log Analyzer Bot

CLI tool written in Go that uses Google's Gemini AI to analyze application logs — automatically summarizing errors, identifying patterns, and suggesting root causes.

## Why This Project

Built as part of learning how to integrate LLM APIs into practical internal automation tools. This solves a real problem: manually reading through log files to spot recurring errors is slow and error-prone — this tool gives an instant AI-generated summary instead.

## Features

- Reads log files and sends them to Gemini for analysis
- Identifies error patterns and frequency
- Suggests possible root causes and recommended actions

## Tech Stack

- Go (standard library only for HTTP — no external HTTP frameworks)
- Google Gemini API (`gemini-2.5-flash`)

## Project Structure

\```
log-analyzer-bot/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── model/
│   │   └── gemini.go         # Request/response structs
│   └── service/
│       └── gemini_service.go # Gemini API integration logic
└── Logs/
    └── sample.log             # Example log file for testing
\```

## Setup

1. Clone the repo
2. Copy `.env.example` to `.env` and add your Gemini API key ([get one free here](https://aistudio.google.com))
3. Install dependencies:
   \```bash
   go mod tidy
   \```
4. Run:
   \```bash
   go run cmd/main.go
   \```

## Example Output

\```
=== Hasil Analisa ===
Berikut adalah hasil analisa dari log yang diberikan:

---
### 1. Error yang Muncul dan Frekuensinya
* **Error:** `Failed to connect to database connection timeout`
  * **Frekuensi:** **3 kali** (pukul 10:16:01, 10:16:05, dan 10:18:12).
* *(Catatan Tambahan)* Terdapat **1 peringatan (WARN)**: `High memory usage detected 87%` pada pukul 10:17:30 sebelum koneksi akhirnya pulih (*restored*) pada pukul 10:20:00.

### 2. Kemungkinan *Root Cause* (Penyebab Utama)
Ada dua skenario utama yang kemungkinan saling berhubungan:

1. **Gangguan Sementara pada Server Database / Jaringan (Paling Kuat):**
   * Database sempat mengalami *downtime*, *restart*, atau gangguan jaringan selama kurang lebih 4 menit (10:16 – 10:20).
   * Gejala lonjakan memori (87%) kemungkinan merupakan **akibat** dari antrean *request* (*connection pool exhaustion* / *thread pile-up*) yang tertahan dan menumpuk di memori aplikasi karena menunggu respon database yang *timeout*.
2. **Keterbatasan Resource Aplikasi (Resource Starvation):**
   * Server aplikasi mengalami lonjakan penggunaan RAM (87%) yang menyebabkan latensi tinggi atau kegagalan I/O jaringan, sehingga aplikasi gagal menyambung ke database tepat waktu (*timeout*).

---
### 3. Rekomendasi Tindakan
1. **Periksa Log Database:**
   * Cek log pada server database di rentang waktu **10:15 – 10:20** untuk memastikan apakah ada *restart*, kehabisan kapasitas (*CPU/RAM spikes*), atau *deadlock*....
\```

## Roadmap

- [x] v1: Basic log analysis via CLI