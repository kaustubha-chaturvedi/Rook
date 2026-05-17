package parser

import "github.com/xwb1989/sqlparser"

func analyzeAST(stmt sqlparser.Statement, raw, sanitized string, forceBypass bool) Analysis {
	a := newBaseAnalysis(raw, sanitized, forceBypass)

	switch s := stmt.(type) {		
		case *sqlparser.Select:
			fillSelectAnalysis(&a, s)
		
		case *sqlparser.Union:
			a.Type = StmtSelect
			a.HasLimitOrTop = s.Limit != nil
			a = mergeSelectBranch(a, s.Left)
			a = mergeSelectBranch(a, s.Right)
		
		case *sqlparser.Insert:
			a.Type = StmtInsert
			a.Tables = []string{formatTableName(s.Table)}
		
		case *sqlparser.Update:
			a.Type = StmtUpdate
			a.HasWhere = s.Where != nil
			a.Tables = collectTableNames(s.TableExprs)
			a.EstimatedCountSQL = estimateCountFromUpdate(s)
		
		case *sqlparser.Delete:
			a.Type = StmtDelete
			a.HasWhere = s.Where != nil
			a.Tables = collectTableNames(s.TableExprs)
			if len(s.Targets) > 0 {
				a.Tables = append(a.Tables, tableNamesToStrings(s.Targets)...)
			}
			a.EstimatedCountSQL = estimateCountFromDelete(s)
		
		case *sqlparser.DDL:
			a.Tables = ddlTables(s)
			a.Type = ddlStatementType(s.Action, sanitized)
		
		case *sqlparser.DBDDL:
			a.Type, a.Tables = dbddlTypeAndTables(s)
		
		default:
		a.Type = stmtTypeFromPreview(sanitized)
	}
	
	return a
}

func fillSelectAnalysis(a *Analysis, s *sqlparser.Select) {
	a.Type = StmtSelect
	a.HasWhere = s.Where != nil
	a.HasLimitOrTop = s.Limit != nil
	a.UsesSelectStar = hasSelectStar(s.SelectExprs)
	a.Tables = collectTableNames(s.From)
	a.HasCartesianJoin = hasCartesianJoin(s.From)
}

func mergeSelectBranch(a Analysis, sel sqlparser.SelectStatement) Analysis {
	switch s := sel.(type) {
		case *sqlparser.Select:
			if s.Where != nil {
				a.HasWhere = true
			}
			if s.Limit != nil {
				a.HasLimitOrTop = true
			}
			if hasSelectStar(s.SelectExprs) {
				a.UsesSelectStar = true
			}
			a.Tables = append(a.Tables, collectTableNames(s.From)...)
			if hasCartesianJoin(s.From) {
				a.HasCartesianJoin = true
			}
	
		case *sqlparser.Union:
			a = mergeSelectBranch(a, s.Left)
			a = mergeSelectBranch(a, s.Right)
			if s.Limit != nil {
				a.HasLimitOrTop = true
			}
	
		case *sqlparser.ParenSelect:
			a = mergeSelectBranch(a, s.Select)
	}

	return a
}

func ddlStatementType(action, sanitized string) StatementType {
	switch action {
	case sqlparser.AlterStr:
		return StmtAlter
	case sqlparser.DropStr:
		return StmtDrop
	case sqlparser.TruncateStr:
		return StmtTruncate
	case sqlparser.CreateStr:
		return StmtCreate
	default:
		return stmtTypeFromPreview(sanitized)
	}
}

func dbddlTypeAndTables(s *sqlparser.DBDDL) (StatementType, []string) {
	switch s.Action {
	case sqlparser.CreateStr:
		return StmtCreate, []string{s.DBName}
	case sqlparser.DropStr:
		return StmtDrop, []string{s.DBName}
	default:
		return StmtUnknown, nil
	}
}

func hasSelectStar(exprs sqlparser.SelectExprs) bool {
	for _, e := range exprs {
		if _, ok := e.(*sqlparser.StarExpr); ok {
			return true
		}
	}
	return false
}

func hasCartesianJoin(from sqlparser.TableExprs) bool {
	for _, te := range from {
		if hasCartesianInExpr(te) {
			return true
		}
	}
	return false
}

func hasCartesianInExpr(te sqlparser.TableExpr) bool {
	join, ok := te.(*sqlparser.JoinTableExpr)
	if !ok {
		return false
	}
	if join.Join == sqlparser.JoinStr && join.Condition.On == nil && len(join.Condition.Using) == 0 {
		return true
	}
	return hasCartesianInExpr(join.LeftExpr) || hasCartesianInExpr(join.RightExpr)
}
