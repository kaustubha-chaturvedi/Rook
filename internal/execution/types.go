package execution

import (
	"context"
	"database/sql"
	"sync"

	"github.com/kaustubha-chaturvedi/Rook/internal/safety"
)

type ConnectRequest struct {
	Label    string `json:"label"`
	Category string `json:"category"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}

type ConnectResponse struct {
	Connected bool   `json:"connected"`
	Label     string `json:"label"`
	Category  string `json:"category"`
	Message   string `json:"message"`
}

type ExecRequest struct {
	CellID      string             `json:"cellId"`
	SQL         string             `json:"sql"`
	AutoAction  string             `json:"autoAction"`
	AutoTimeout int64              `json:"autoTimeoutSeconds"`
	Rules       safety.RulesConfig `json:"rules"`
	Environment string             `json:"environment"`
}

type ExecResponse struct {
	Safety           safety.Result       `json:"safety"`
	Duration         int64               `json:"durationMs"`
	TxState          string              `json:"txState"`
	EstimatedCount   string              `json:"estimatedCountSql,omitempty"`
	StatementCount   int                 `json:"statementCount"`
	RollbackOnCancel bool                `json:"rollbackOnCancel"`
	Columns          []string            `json:"columns"`
	Rows             []map[string]string `json:"rows"`
	RowsAffected     int64               `json:"rowsAffected"`
	Message          string              `json:"message"`
}

type Service struct {
	mu     sync.Mutex
	db     *sql.DB
	conn   ConnectRequest
	cancel map[string]context.CancelFunc
}

func NewService() *Service {
	return &Service{cancel: map[string]context.CancelFunc{}}
}
