package model

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Response struct {
	Candidates    []Candidate   `json:"candidates"`
	UsageMetadata UsageMetadata `json:"usageMetadata"`
	ModelVersion  string        `json:"modelVersion"`
	ResponseID    string        `json:"responseId"`
}

type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason"`
	Index        int     `json:"index"`
}

type Part struct {
	Text             string `json:"text"`
	ThoughtSignature string `json:"thoughtSignature,omitempty"`
}

type UsageMetadata struct {
	PromptTokenCount     int                  `json:"promptTokenCount"`
	CandidatesTokenCount int                  `json:"candidatesTokenCount"`
	TotalTokenCount      int                  `json:"totalTokenCount"`
	PromptTokensDetails  []PromptTokensDetail `json:"promptTokensDetails"`
	ThoughtsTokenCount   int                  `json:"thoughtsTokenCount"`
	ServiceTier          string               `json:"serviceTier"`
}

type PromptTokensDetail struct {
	Modality   string `json:"modality"`
	TokenCount int    `json:"tokenCount"`
}

type AnalysisResult struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Error          string `json:"error"`
	Occurrences    int    `json:"occurrences"`
	RootCause      string `json:"root_cause"`
	Recommendation string `json:"recommendation"`
}
