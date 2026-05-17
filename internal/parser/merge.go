package parser

func mergeAnalyses(stmts []Analysis) Analysis {
	merged := stmts[0]
	
	for _, s := range stmts[1:] {
		if isDestructiveType(s.Type) {
			merged.Type = s.Type
		}
		if (s.Type == StmtUpdate || s.Type == StmtDelete) && !s.HasWhere {
			merged.HasWhere = false
		}
		merged.Tables = append(merged.Tables, s.Tables...)
		merged.HasLimitOrTop = merged.HasLimitOrTop || s.HasLimitOrTop
		merged.UsesSelectStar = merged.UsesSelectStar || s.UsesSelectStar
		merged.HasNonSelectiveWhere = merged.HasNonSelectiveWhere || s.HasNonSelectiveWhere
		merged.HasCartesianJoin = merged.HasCartesianJoin || s.HasCartesianJoin
		merged.HasDynamicExec = merged.HasDynamicExec || s.HasDynamicExec
		merged.ForceBypassSafety = merged.ForceBypassSafety || s.ForceBypassSafety
		if merged.EstimatedCountSQL == "" {
			merged.EstimatedCountSQL = s.EstimatedCountSQL
		}
		for k, v := range s.DangerousKeywords {
			merged.DangerousKeywords[k] = v
		}
	}
	
	return merged
}

func isDestructiveType(t StatementType) bool {
	switch t {
		case StmtUpdate, StmtDelete, StmtDrop, StmtTruncate, StmtAlter:
			return true
	
		default:
			return false
	}
}
