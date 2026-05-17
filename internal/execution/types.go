package execution

import (
	"context"
	"sync"

	"github.com/kaustubha-chaturvedi/Rook/internal/safety"
)

type ExecRequest struct {
	CellID      string             `json:"cellId"`
	SQL         string             `json:"sql"`
	AutoAction  string             `json:"autoAction"` // none|commit|rollback
	AutoTimeout int64              `json:"autoTimeoutSeconds"`
	Rules       safety.RulesConfig `json:"rules"`
}

type ExecResponse struct {
	Safety           safety.Result `json:"safety"`
	Duration         int64         `json:"durationMs"`
	TxState          string        `json:"txState"`
	EstimatedCount   string        `json:"estimatedCountSql,omitempty"`
	StatementCount   int           `json:"statementCount"`
	RollbackOnCancel bool          `json:"rollbackOnCancel"`
}

type Service struct {
	mu     sync.Mutex
	cancel map[string]context.CancelFunc
}

func NewService() *Service {
	return &Service{cancel: map[string]context.CancelFunc{}}
}
