package service

import "strings"

func ChunkLog(logContent string, linesPerChunk int) []string {
	var chunks []string

	lines := strings.Split(logContent, "\n")

	for i := 0; i < len(lines); i += linesPerChunk {
		end := i + linesPerChunk
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, strings.Join(lines[i:end], "\n"))
	}
	return chunks
}
