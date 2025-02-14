package gormexpect

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

// ExecExpectation is returned by Expecter. It exposes a narrower API than
// Execer to limit footguns.
type ExecExpectation interface {
	WillSucceed(lastReturnedID, rowsAffected int64) ExecExpectation
	WillFail(err error) ExecExpectation
}

// SqlmockExecExpectation implement ExecExpectation with gosqlmock
type SqlmockExecExpectation struct {
	parent *Expecter
}

// WillSucceed sets the exec to be successful with the passed ID and rows.
// This method may also call Query, if there are default columns.
func (e *SqlmockExecExpectation) WillSucceed(lastReturnedID, rowsAffected int64) ExecExpectation {
	exec, _ := e.parent.recorder.GetFirst()
	e.parent.adapter.ExpectExec(exec).WillSucceed(lastReturnedID, rowsAffected)

	// for now, just return empty row
	if len(e.parent.recorder.stmts) >= 1 {
		// follow-up query
		query, _ := e.parent.recorder.GetFirst()

		switch query.kind {
		case "query":
			if len(e.parent.recorder.blankColumns) >= 1 {
				e.parent.adapter.ExpectQuery(query).Returns(sqlmock.NewRows(e.parent.recorder.blankColumns))
			}
		case "exec":
			e.parent.adapter.ExpectExec(query).WillSucceed(1, 1)
		default:
			return e
		}
	}

	return e
}

// WillFail sets the exec to fail with the passed error
func (e *SqlmockExecExpectation) WillFail(err error) ExecExpectation {
	query, _ := e.parent.recorder.GetFirst()
	e.parent.adapter.ExpectExec(query).WillFail(err)

	return e
}
