package report

type ReportAnalysisModerationData struct {
	Toxicity       float64 `json:"toxicity,omitempty"`
	SevereToxicity float64 `json:"severe_toxicity,omitempty"`
	IdentityAttack float64 `json:"identity_attack,omitempty"`
	Insult         float64 `json:"insult,omitempty"`
	Profanity      float64 `json:"profanity,omitempty"`
	Threat         float64 `json:"threat,omitempty"`
}
