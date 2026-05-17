package safety

import (
	"fmt"

	"github.com/kaustubha-chaturvedi/Rook/internal/parser"
)

func Validate(statements []parser.Analysis, estimatedRows int64, cfg RulesConfig) Result {
	res := Result{Status: Safe}
	destructiveCount := 0
	types := map[parser.StatementType]bool{}

	for _, st := range statements {
		types[st.Type] = true

		if cfg.ReadOnlyConnection && isWriteType(st.Type) {
			return blockResult("Write operations blocked on read-only/replica connection")
		}

		if st.ForceBypassSafety {
			if cfg.IsProduction && !cfg.AllowForceBypassInProd {
				return blockResult("FORCE bypass is disabled in PROD")
			}
			res = elevate(res, Warning, "Safety checks bypassed via --FORCE")
			res.Bypassed = true
			continue
		}

		res = applyPatternChecks(res, st)

		var early *Result
		res, destructiveCount, early = applyStatementRules(res, st, destructiveCount, cfg)
		if early != nil {
			return *early
		}
	}

	if destructiveCount > 1 {
		return blockResult("Multiple destructive statements blocked")
	}
	if hasMixedBatch(types) {
		res = elevate(res, Warning, "Mixed safe + destructive statements in same cell")
	}

	res = applyRowThresholds(res, estimatedRows, cfg)
	if cfg.IsProduction && anyDestructive(types) {
		res = elevate(res, Danger, "PROD operation: confirmation and reason required")
		res.RequiresTypedConfirm = true
	}
	return res
}

func applyPatternChecks(res Result, st parser.Analysis) Result {
	for kw, act := range st.DangerousKeywords {
		res = elevate(res, mapKeywordAction(act), fmt.Sprintf("Dangerous keyword detected: %s", kw))
	}

	if st.HasCartesianJoin {
		res = elevate(res, Danger, "Potential cartesian join detected")
	}
	
	if st.HasNonSelectiveWhere {
		res = elevate(res, Warning, "Non-selective WHERE detected")
	}
	
	if st.HasDynamicExec {
		res = elevate(res, Warning, "Potential unsafe dynamic SQL detected")
	}
	
	return res
}

func applyStatementRules(res Result, st parser.Analysis, destructiveCount int, cfg RulesConfig) (Result, int, *Result) {
	switch st.Type {
		case parser.StmtUpdate:
			destructiveCount++
			if !st.HasWhere {
				r := blockResult("Unsafe UPDATE detected", "Missing WHERE clause", "Execution blocked")
				return res, destructiveCount, &r
			}
			res = elevate(res, Danger, "UPDATE detected, review affected-row estimate before commit")
		
		case parser.StmtDelete:
			destructiveCount++
			if !st.HasWhere {
				r := blockResult("Unsafe DELETE detected", "Missing WHERE clause", "Execution blocked")
				return res, destructiveCount, &r
			}
			res = elevate(res, Danger, "DELETE detected, review affected-row estimate before commit")
		
		case parser.StmtDrop:
			destructiveCount++
			res = elevate(res, Danger, "DROP detected - typed confirmation required")
			res.RequiresTypedConfirm = true
		
		case parser.StmtTruncate:
			destructiveCount++
			res = elevate(res, Danger, "TRUNCATE detected - confirmation required")
			if cfg.IsProduction && cfg.StrictProd {
				res = elevate(res, Blocked, "TRUNCATE blocked in strict PROD mode")
			}
		
		case parser.StmtAlter:
			res = elevate(res, Warning, "ALTER TABLE may cause locks/downtime")
			res.RequiresConfirmation = true
		
		case parser.StmtSelect:
			if !st.HasLimitOrTop {
				res = elevate(res, Warning, "SELECT without LIMIT/TOP may be expensive")
			}
			if st.UsesSelectStar {
				res = elevate(res, Info, "SELECT * detected")
			}
	}

	return res, destructiveCount, nil
}

func blockResult(msgs ...string) Result {
	return Result{Status: Blocked, Messages: msgs}
}

func applyRowThresholds(res Result, estimatedRows int64, cfg RulesConfig) Result {
	switch {
		case estimatedRows >= cfg.LargeRowsBlockThreshold:
			res = elevate(res, Danger, "Very large row impact estimated")
			res.RequiresTypedConfirm = true
		
		case estimatedRows >= cfg.LargeRowsDangerThreshold:
			res = elevate(res, Danger, "Large row impact estimated")
			res.RequiresConfirmation = true
		
		case estimatedRows >= cfg.LargeRowsWarnThreshold:
			res = elevate(res, Warning, "Moderate row impact estimated")
	}
	
	return res
}
