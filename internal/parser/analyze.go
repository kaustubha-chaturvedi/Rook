package parser

import (
	"fmt"
	"io"

	"github.com/xwb1989/sqlparser"
)

func Analyze(sql string) Analysis {
	stmts := AnalyzeStatements(sql)
	switch len(stmts) {
		case 0:
			sanitized, force := sanitizeSQL(sql)
			return analyzeFromText(sql, sanitized, force)
		
		case 1:
			return stmts[0]
		
		default:
			return mergeAnalyses(stmts)
	}
}

func AnalyzeStatements(sql string) []Analysis {
	sanitized, forceBypass := sanitizeSQL(sql)

	if sanitized == "" {
		return nil
	}

	tokenizer := sqlparser.NewStringTokenizer(sanitized)
	var out []Analysis
	for {
		stmt, err := sqlparser.ParseNext(tokenizer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return []Analysis{analyzeFromText(sql, sanitized, forceBypass)}
		}

		a := analyzeAST(stmt, sql, sanitized, forceBypass)
		applyTextSafetyChecks(&a, sanitized)
		out = append(out, a)
	}

	if len(out) == 0 {
		return []Analysis{analyzeFromText(sql, sanitized, forceBypass)}
	}

	return out
}

func ValidateSyntax(sql string) error {

	tokenizer := sqlparser.NewStringTokenizer(sql)

	for {
		_, err := sqlparser.ParseNext(tokenizer)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}
	}
}
