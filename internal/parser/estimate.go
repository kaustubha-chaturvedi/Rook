package parser

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func estimateCountFromUpdate(u *sqlparser.Update) string {
	if u.Where == nil {
		return ""
	}

	tables := collectTableNames(u.TableExprs)
	
	if len(tables) == 0 {
		return ""
	}
	
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tables[0], sqlparser.String(u.Where.Expr))
}

func estimateCountFromDelete(d *sqlparser.Delete) string {
	if d.Where == nil {
		return ""
	}

	
	tables := collectTableNames(d.TableExprs)
	
	if len(tables) == 0 {
		return "SELECT COUNT(*) FROM " + sqlparser.String(d.TableExprs) + " WHERE " + sqlparser.String(d.Where.Expr)
	}
	
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tables[0], sqlparser.String(d.Where.Expr))
}

func estimateCountFromText(sql string) string {
	n := strings.ToUpper(sql)
	whereIdx := strings.Index(n, " WHERE ")
	
	if whereIdx == -1 {
		return ""
	}
	
	if strings.HasPrefix(n, "UPDATE") {
		if setIdx := strings.Index(n, " SET "); setIdx > 6 {
			table := strings.TrimSpace(sql[6:setIdx])
			return "SELECT COUNT(*) FROM " + table + sql[whereIdx:]
		}
	}
	
	if strings.HasPrefix(n, "DELETE") {
		if fromIdx := strings.Index(n, "FROM"); fromIdx >= 0 {
			return "SELECT COUNT(*) " + sql[fromIdx:]
		}
	}
	
	return ""
}
