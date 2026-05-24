package internal

import "fmt"

type InsertQb struct {
	settings        *QbSettings
	tableClause     string
	columns         []string
	rows            [][]interface{}
	cteQb           *CteQb
	conflictColumns []string
	conflictAction  string // "" | "nothing" | "update"
	conflictUpdate  *DoUpdateQb
}

func NewInsertQueryBuilder(settings *QbSettings) *InsertQb {
	return &InsertQb{settings: settings}
}

func (qb *InsertQb) Into(table string) *InsertQb {
	qb.tableClause = table
	return qb
}

func (qb *InsertQb) Columns(cols ...string) *InsertQb {
	qb.columns = append(qb.columns, cols...)
	return qb
}

func (qb *InsertQb) Values(vals ...interface{}) *InsertQb {
	qb.rows = append(qb.rows, vals)
	return qb
}

func (qb *InsertQb) OnConflict(cols ...string) *insertOnConflictQb {
	return &insertOnConflictQb{insert: qb, columns: cols}
}

func (qb *InsertQb) Clone() *InsertQb {
	clone := &InsertQb{
		settings:       qb.settings,
		tableClause:    qb.tableClause,
		cteQb:          qb.cteQb,
		conflictAction: qb.conflictAction,
	}
	if qb.columns != nil {
		clone.columns = make([]string, len(qb.columns))
		copy(clone.columns, qb.columns)
	}
	if qb.rows != nil {
		clone.rows = make([][]interface{}, len(qb.rows))
		for i, row := range qb.rows {
			r := make([]interface{}, len(row))
			copy(r, row)
			clone.rows[i] = r
		}
	}
	if qb.conflictColumns != nil {
		clone.conflictColumns = make([]string, len(qb.conflictColumns))
		copy(clone.conflictColumns, qb.conflictColumns)
	}
	if qb.conflictUpdate != nil {
		du := &DoUpdateQb{}
		du.setClauses = make([]setClause, len(qb.conflictUpdate.setClauses))
		copy(du.setClauses, qb.conflictUpdate.setClauses)
		clone.conflictUpdate = du
	}
	return clone
}

func (qb *InsertQb) ExtendSql(ctx *SqlBuildingCtx) error {
	if qb.tableClause == "" {
		return fmt.Errorf("insert: table is required")
	}

	if len(qb.rows) == 0 {
		return fmt.Errorf("insert: at least one values row is required")
	}

	if _, err := ctx.Sql.WriteString("insert into "); err != nil {
		return err
	}
	if err := ctx.WriteArg(&qb.tableClause, true); err != nil {
		return err
	}

	if len(qb.columns) > 0 {
		if _, err := ctx.Sql.WriteString(" ("); err != nil {
			return err
		}
		for i, col := range qb.columns {
			if i > 0 {
				if _, err := ctx.Sql.WriteString(", "); err != nil {
					return err
				}
			}
			c := col
			if err := ctx.WriteArg(&c, true); err != nil {
				return err
			}
		}
		if err := ctx.Sql.WriteByte(')'); err != nil {
			return err
		}
	}

	if _, err := ctx.Sql.WriteString(" values "); err != nil {
		return err
	}

	for i, row := range qb.rows {
		if len(qb.columns) > 0 && len(row) != len(qb.columns) {
			return fmt.Errorf("insert: row %d has %d values, expected %d", i, len(row), len(qb.columns))
		}

		if i > 0 {
			if _, err := ctx.Sql.WriteString(", "); err != nil {
				return err
			}
		}

		if err := ctx.Sql.WriteByte('('); err != nil {
			return err
		}
		for j, val := range row {
			if j > 0 {
				if _, err := ctx.Sql.WriteString(", "); err != nil {
					return err
				}
			}
			var arg interface{}
			if raw, ok := val.(*QbRaw); ok {
				arg = raw
			} else {
				arg = NewBindingValue(val)
			}
			if err := ctx.WriteArg(arg, false); err != nil {
				return err
			}
		}
		if err := ctx.Sql.WriteByte(')'); err != nil {
			return err
		}
	}

	if qb.conflictAction != "" {
		if _, err := ctx.Sql.WriteString(" on conflict"); err != nil {
			return err
		}
		if len(qb.conflictColumns) > 0 {
			if _, err := ctx.Sql.WriteString(" ("); err != nil {
				return err
			}
			for i, col := range qb.conflictColumns {
				if i > 0 {
					if _, err := ctx.Sql.WriteString(", "); err != nil {
						return err
					}
				}
				c := col
				if err := ctx.WriteArg(&c, false); err != nil {
					return err
				}
			}
			if err := ctx.Sql.WriteByte(')'); err != nil {
				return err
			}
		}
		switch qb.conflictAction {
		case "nothing":
			if _, err := ctx.Sql.WriteString(" do nothing"); err != nil {
				return err
			}
		case "update":
			if _, err := ctx.Sql.WriteString(" do update set "); err != nil {
				return err
			}
			if err := qb.conflictUpdate.extendSql(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (qb *InsertQb) ToSql() (string, []interface{}, error) {
	ctx := NewSqlBuildingCtx(qb.settings)

	if qb.cteQb != nil {
		if err := qb.cteQb.ExtendSql(&ctx); err != nil {
			return "", nil, err
		}
		ctx.Sql.WriteByte(' ')
	}

	if err := qb.ExtendSql(&ctx); err != nil {
		return "", nil, err
	}

	return ctx.Sql.String(), ctx.bindings, nil
}

// insertOnConflictQb is an intermediate builder returned by OnConflict.
// The user immediately chains DoNothing() or DoUpdate() on it.
type insertOnConflictQb struct {
	insert  *InsertQb
	columns []string
}

func (qb *insertOnConflictQb) DoNothing() *InsertQb {
	qb.insert.conflictColumns = qb.columns
	qb.insert.conflictAction = "nothing"
	return qb.insert
}

func (qb *insertOnConflictQb) DoUpdate(fn func(*DoUpdateQb)) *InsertQb {
	du := &DoUpdateQb{}
	fn(du)
	qb.insert.conflictColumns = qb.columns
	qb.insert.conflictAction = "update"
	qb.insert.conflictUpdate = du
	return qb.insert
}

// DoUpdateQb builds the SET list for ON CONFLICT DO UPDATE.
type DoUpdateQb struct {
	setClauses []setClause
}

func (qb *DoUpdateQb) Set(col string, val interface{}) *DoUpdateQb {
	qb.setClauses = append(qb.setClauses, setClause{col: col, val: val})
	return qb
}

func (qb *DoUpdateQb) extendSql(ctx *SqlBuildingCtx) error {
	for i, s := range qb.setClauses {
		if i > 0 {
			if _, err := ctx.Sql.WriteString(", "); err != nil {
				return err
			}
		}
		if err := ctx.WriteArg(&s.col, false); err != nil {
			return err
		}
		if _, err := ctx.Sql.WriteString(" = "); err != nil {
			return err
		}
		// string  → relation identifier ("excluded.a" → "excluded"."a")
		// *QbRaw → raw SQL fragment
		// other  → bound parameter
		switch v := s.val.(type) {
		case string:
			if err := ctx.WriteArg(&v, false); err != nil {
				return err
			}
		case *string:
			if err := ctx.WriteArg(v, false); err != nil {
				return err
			}
		case *QbRaw:
			if err := ctx.WriteArg(v, false); err != nil {
				return err
			}
		default:
			if err := ctx.WriteArg(NewBindingValue(s.val), false); err != nil {
				return err
			}
		}
	}
	return nil
}
