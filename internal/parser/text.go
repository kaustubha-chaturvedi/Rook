package parser

import (
	"strings"
	"github.com/xwb1989/sqlparser"
)

func analyzeFromText(raw, sanitized string, forceBypass bool) Analysis {
	a := newBaseAnalysis(raw, sanitized, forceBypass)
	a.Type = stmtTypeFromPreview(sanitized)
	norm := strings.ToUpper(sanitized)

	switch a.Type {
		case StmtSelect:
			a.HasLimitOrTop = strings.Contains(norm, " LIMIT ") || strings.Contains(norm, " TOP ")
			a.UsesSelectStar = strings.Contains(norm, "SELECT *") || strings.Contains(norm, "SELECT  *")
		
		case StmtUpdate, StmtDelete:
			a.HasWhere = strings.Contains(norm, " WHERE ")
			a.EstimatedCountSQL = estimateCountFromText(sanitized)
	}

	a.Tables = extractTablesFromText(sanitized, a.Type)
	applyTextSafetyChecks(&a, sanitized)
	
	return a
}

func applyTextSafetyChecks(a *Analysis, sanitized string) {
	norm := strings.ToUpper(sanitized)
	
	a.HasNonSelectiveWhere = a.HasNonSelectiveWhere || hasRiskyWhere(norm)

	a.HasCartesianJoin = a.HasCartesianJoin ||
		strings.Contains(norm, "CROSS JOIN") ||
		(strings.Contains(norm, " JOIN ") && !strings.Contains(norm, " ON "))

	a.HasDynamicExec = a.HasDynamicExec ||
		strings.Contains(norm, "EXEC(@") ||
		strings.Contains(norm, "SP_EXECUTESQL") ||
		strings.Contains(norm, "||")

	for k, v := range detectDangerousKeywords(norm) {
		a.DangerousKeywords[k] = v
	}
	
	if a.Type == StmtSelect && !a.HasLimitOrTop {
		a.HasLimitOrTop = strings.Contains(norm, " TOP ")
	}
}

func stmtTypeFromPreview(sql string) StatementType {
	switch sqlparser.Preview(sql) {
		case sqlparser.StmtSelect, sqlparser.StmtStream:
			return StmtSelect
		
		case sqlparser.StmtInsert, sqlparser.StmtReplace:
			return StmtInsert
		
		case sqlparser.StmtUpdate:
			return StmtUpdate
		
		case sqlparser.StmtDelete:
			return StmtDelete
		
		case sqlparser.StmtDDL:
			return stmtTypeFromPrefix(sql)
		
		default:
			return stmtTypeFromPrefix(sql)
	}
}

func stmtTypeFromPrefix(sql string) StatementType {
	norm := strings.ToUpper(strings.TrimSpace(sql))
	switch {
		case strings.HasPrefix(norm, "SELECT"):
			return StmtSelect

		case strings.HasPrefix(norm, "INSERT"):
			return StmtInsert

		case strings.HasPrefix(norm, "UPDATE"):
			return StmtUpdate

		case strings.HasPrefix(norm, "DELETE"):
			return StmtDelete

		case strings.HasPrefix(norm, "ALTER"):
			return StmtAlter

		case strings.HasPrefix(norm, "DROP"):
			return StmtDrop

		case strings.HasPrefix(norm, "TRUNCATE"):
			return StmtTruncate

		case strings.HasPrefix(norm, "CREATE"):
			return StmtCreate

		default:
			return StmtUnknown
	}
}

func detectDangerousKeywords(norm string) map[string]string {
	rules := map[string]string{
		"DROP": "Danger", "TRUNCATE": "Danger", "SHUTDOWN": "Block", "KILL": "Danger",
		"EXEC": "Warning", "XP_CMDSHELL": "Block", "BULK INSERT": "Warning",
		"OPENROWSET": "Warning", "LOAD DATA": "Warning",
	}

	m := make(map[string]string, len(rules))
	
	for k, v := range rules {
		if strings.Contains(norm, k) {
			m[k] = v
		}
	}
	
	return m
}

func hasRiskyWhere(norm string) bool {
	return strings.Contains(norm, "WHERE 1=1") ||
		strings.Contains(norm, "WHERE TRUE") ||
		strings.Contains(norm, "IS NOT NULL") ||
		strings.Contains(norm, "!= 0")
}
