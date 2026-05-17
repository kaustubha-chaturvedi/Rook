package safety

type Status string

const (
	Safe    Status = "Safe"
	Info    Status = "Info"
	Warning Status = "Warning"
	Danger  Status = "Danger"
	Blocked Status = "Blocked"
)

type RulesConfig struct {
	LargeRowsWarnThreshold   int64 `json:"largeRowsWarnThreshold"`
	LargeRowsDangerThreshold int64 `json:"largeRowsDangerThreshold"`
	LargeRowsBlockThreshold  int64 `json:"largeRowsBlockThreshold"`
	LargeResultWarnThreshold int64 `json:"largeResultWarnThreshold"`
	MaxRowsFetched           int64 `json:"maxRowsFetched"`
	MaxQueryTimeSeconds      int64 `json:"maxQueryTimeSeconds"`
	MaxConcurrentQueries     int64 `json:"maxConcurrentQueries"`
	StrictProd               bool  `json:"strictProd"`
	AllowForceBypassInProd   bool  `json:"allowForceBypassInProd"`
	ReadOnlyConnection       bool  `json:"readOnlyConnection"`
	IsProduction             bool  `json:"isProduction"`
}

type Result struct {
	Status               Status   `json:"status"`
	Messages             []string `json:"messages"`
	Bypassed             bool     `json:"bypassed"`
	RequiresConfirmation bool     `json:"requiresConfirmation"`
	RequiresTypedConfirm bool     `json:"requiresTypedConfirm"`
}

func DefaultRulesConfig() RulesConfig {
	return RulesConfig{
		LargeRowsWarnThreshold:   100,
		LargeRowsDangerThreshold: 1000,
		LargeRowsBlockThreshold:  10000,
		LargeResultWarnThreshold: 100000,
		MaxRowsFetched:           10000,
		MaxQueryTimeSeconds:      60,
		MaxConcurrentQueries:     5,
		StrictProd:               true,
		AllowForceBypassInProd:   false,
		ReadOnlyConnection:       false,
		IsProduction:             false,
	}
}
