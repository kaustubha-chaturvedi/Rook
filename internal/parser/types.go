package parser

type StatementType string

const (
	StmtUnknown  StatementType = "UNKNOWN"
	StmtSelect   StatementType = "SELECT"
	StmtInsert   StatementType = "INSERT"
	StmtUpdate   StatementType = "UPDATE"
	StmtDelete   StatementType = "DELETE"
	StmtAlter    StatementType = "ALTER"
	StmtDrop     StatementType = "DROP"
	StmtTruncate StatementType = "TRUNCATE"
	StmtCreate   StatementType = "CREATE"
)

type Analysis struct {
	Raw                  string            `json:"raw"`
	SanitizedSQL         string            `json:"sanitizedSql"`
	Type                 StatementType     `json:"type"`
	Tables               []string          `json:"tables,omitempty"`
	HasWhere             bool              `json:"hasWhere"`
	HasLimitOrTop        bool              `json:"hasLimitOrTop"`
	UsesSelectStar       bool              `json:"usesSelectStar"`
	EstimatedCountSQL    string            `json:"estimatedCountSQL,omitempty"`
	ForceBypassSafety    bool              `json:"forceBypassSafety"`
	HasNonSelectiveWhere bool              `json:"hasNonSelectiveWhere"`
	HasCartesianJoin     bool              `json:"hasCartesianJoin"`
	HasDynamicExec       bool              `json:"hasDynamicExec"`
	DangerousKeywords    map[string]string `json:"dangerousKeywords"`
}

func (a Analysis) StatementType() StatementType { return a.Type }

func (a Analysis) HasWhereClause() bool { return a.HasWhere }

func (a Analysis) IsDestructive() bool {
	switch a.Type {
		case StmtUpdate, StmtDelete, StmtInsert, StmtAlter, StmtDrop, StmtTruncate, StmtCreate:
			return true
	
		default:
			return false
	}
}

func (a Analysis) ExtractTables() []string {
	if len(a.Tables) == 0 {
		return nil
	}

	out := make([]string, len(a.Tables))
	copy(out, a.Tables)
	
	return out
}
