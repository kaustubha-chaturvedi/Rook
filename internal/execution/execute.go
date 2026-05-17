package execution

import (
	"context"
	"errors"
	"time"

	"github.com/kaustubha-chaturvedi/Rook/internal/parser"
	"github.com/kaustubha-chaturvedi/Rook/internal/safety"
)

func (s *Service) Execute(req ExecRequest) (ExecResponse, error) {
	start := time.Now()
	cfg := req.Rules

	if cfg.MaxQueryTimeSeconds == 0 {
		cfg = safety.DefaultRulesConfig()
	}

	analyses, est := collectAnalyses(req.SQL)
	safe := safety.Validate(analyses, 0, cfg)
	
	base := ExecResponse{
		Safety:           safe,
		EstimatedCount:   est,
		StatementCount:   len(analyses),
		RollbackOnCancel: true,
	}

	if safe.Status == safety.Blocked {
		base.TxState = "blocked"
		base.Duration = time.Since(start).Milliseconds()
		return base, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.MaxQueryTimeSeconds)*time.Second)
	
	s.mu.Lock()
	s.cancel[req.CellID] = cancel
	s.mu.Unlock()
	
	defer func() {
		s.mu.Lock()
		delete(s.cancel, req.CellID)
		s.mu.Unlock()
	}()

	select {
		case <-ctx.Done():
			base.TxState = "cancelled_rollback"
			base.Duration = time.Since(start).Milliseconds()
			return base, ctx.Err()
	
		case <-time.After(50 * time.Millisecond):
	}

	txState, err := runAutoAction(ctx, req)
	base.TxState = txState
	base.Duration = time.Since(start).Milliseconds()
	
	return base, err
}

func collectAnalyses(sql string) ([]parser.Analysis, string) {
	cells := parser.SplitCells(sql)
	analyses := make([]parser.Analysis, 0, len(cells))
	est := ""
	
	appendCell := func(cell string) {
		for _, a := range parser.AnalyzeStatements(cell) {
			if est == "" {
				est = a.EstimatedCountSQL
			}
			analyses = append(analyses, a)
		}
	}
	
	for _, cell := range cells {
		appendCell(cell)
	}
	
	if len(analyses) == 0 {
		appendCell(sql)
	}
	
	return analyses, est
}

func runAutoAction(ctx context.Context, req ExecRequest) (string, error) {
	action := req.AutoAction
	
	if action == "" {
		action = "rollback"
	}
	
	if action != "commit" && action != "rollback" {
		return "pending_commit", nil
	}
	
	timeout := req.AutoTimeout
	
	if timeout <= 0 {
		timeout = 5
	}
	
	select {
		case <-ctx.Done():
			return "cancelled_rollback", errors.New("cancelled before auto action")
		case <-time.After(time.Duration(timeout) * time.Second):
			return "auto_" + action, nil
	}
}

func (s *Service) Cancel(cellID string) {
	s.mu.Lock()
	c := s.cancel[cellID]
	s.mu.Unlock()
	
	if c != nil {
		c()
	}
}

func (s *Service) Shutdown(context.Context) {}
