package parser

import (
	"strings"

	"github.com/xwb1989/sqlparser"
)

func collectTableNames(from sqlparser.TableExprs) []string {
	var tables []string
	
	for _, te := range from {
		tables = append(tables, tablesFromExpr(te)...)
	}
	
	return dedupeStrings(tables)
}

func tablesFromExpr(te sqlparser.TableExpr) []string {
	switch t := te.(type) {
		case *sqlparser.AliasedTableExpr:
			if tn, ok := t.Expr.(sqlparser.TableName); ok {
				return []string{formatTableName(tn)}
			}
		
		case *sqlparser.JoinTableExpr:
			return append(tablesFromExpr(t.LeftExpr), tablesFromExpr(t.RightExpr)...)
		
		case *sqlparser.ParenTableExpr:
			return collectTableNames(t.Exprs)
	}

	return nil
}

func ddlTables(d *sqlparser.DDL) []string {
	if !d.Table.Name.IsEmpty() {
		return []string{formatTableName(d.Table)}
	}
	
	if !d.NewName.Name.IsEmpty() {
		return []string{formatTableName(d.NewName)}
	}
	
	return nil
}

func formatTableName(t sqlparser.TableName) string {
	if !t.Qualifier.IsEmpty() {
		return t.Qualifier.String() + "." + t.Name.String()
	}
	
	return t.Name.String()
}

func tableNamesToStrings(names sqlparser.TableNames) []string {
	out := make([]string, 0, len(names))
	
	for _, t := range names {
		out = append(out, formatTableName(t))
	}
	
	return out
}

func extractTablesFromText(sql string, typ StatementType) []string {
	norm := strings.ToUpper(sql)
	
	var table string
	
	switch typ {
		case StmtUpdate:
			if setIdx := strings.Index(norm, " SET "); setIdx > 6 {
				table = strings.TrimSpace(sql[6:setIdx])
			}
	
		case StmtDelete:
			fromIdx := strings.Index(norm, "FROM")
			whereIdx := strings.Index(norm, " WHERE ")
			if fromIdx >= 0 {
				end := len(sql)
				if whereIdx > fromIdx {
					end = whereIdx
				}
				table = strings.TrimSpace(sql[fromIdx+4 : end])
			}
	
		case StmtInsert:
			if intoIdx := strings.Index(norm, "INTO"); intoIdx >= 0 {
				rest := strings.TrimSpace(sql[intoIdx+4:])
				if sp := strings.Fields(rest); len(sp) > 0 {
					table = sp[0]
				}
			}
	}
	
	if table == "" {
		return nil
	}
	
	return []string{table}
}

func dedupeStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	
	return out
}
