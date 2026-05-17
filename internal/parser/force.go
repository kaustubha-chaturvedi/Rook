package parser

import (
	"regexp"
	"strings"
)

var forceTokenRegex = regexp.MustCompile(`(?i)--\s*FORCE\b`)

func hasForceBypass(sql string) bool { return forceTokenRegex.MatchString(sql) }

func sanitizeSQL(sql string) (sanitized string, forceBypass bool) {
	forceBypass = hasForceBypass(sql)
	sanitized = strings.TrimSpace(forceTokenRegex.ReplaceAllString(sql, ""))
	
	return sanitized, forceBypass
}

func newBaseAnalysis(raw, sanitized string, forceBypass bool) Analysis {
	return Analysis{
		Raw:               raw,
		SanitizedSQL:      sanitized,
		Type:              StmtUnknown,
		ForceBypassSafety: forceBypass,
		DangerousKeywords: map[string]string{},
	}
}
