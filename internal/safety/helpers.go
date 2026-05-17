package safety

import (
	"strings"

	"github.com/kaustubha-chaturvedi/Rook/internal/parser"
)

func hasMixedBatch(types map[parser.StatementType]bool) bool {
	return types[parser.StmtSelect] &&
		(types[parser.StmtDelete] || types[parser.StmtDrop] ||
			(types[parser.StmtAlter] && types[parser.StmtUpdate]))
}

func anyDestructive(types map[parser.StatementType]bool) bool {
	return types[parser.StmtUpdate] || types[parser.StmtDelete] ||
		types[parser.StmtAlter] || types[parser.StmtDrop] || types[parser.StmtTruncate]
}

func isWriteType(t parser.StatementType) bool {
	switch t {
		case parser.StmtUpdate, parser.StmtDelete, parser.StmtInsert,
			parser.StmtAlter, parser.StmtDrop, parser.StmtTruncate, parser.StmtCreate:
			return true
		
		default:
			return false
	}
}

func elevate(current Result, status Status, message string) Result {
	if severity(status) > severity(current.Status) {
		current.Status = status
	}

	current.Messages = append(current.Messages, strings.TrimSpace(message))
	if status == Danger {
		current.RequiresConfirmation = true
	}
	
	if status == Blocked {
		current.RequiresTypedConfirm = true
	}
	
	return current
}

func severity(s Status) int {
	switch s {
		case Safe:
			return 0
	
		case Info:
			return 1
	
		case Warning:
			return 2
	
		case Danger:
			return 3
	
		case Blocked:
			return 4
	
		default:
			return 0
	}
}

func mapKeywordAction(action string) Status {
	switch action {
		case "Block":
			return Blocked
	
		case "Danger":
			return Danger
	
		case "Warning":
			return Warning
	
		default:
			return Info
	}
}
