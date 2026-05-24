package internal

import "fmt"

type setClause struct {
	col string
	val interface{}
}

type UpdateQb struct {
	WhereBuilder[UpdateQb]
	tableClause string
	setClauses  []setClause
	cteQb       *CteQb
}

func NewUpdateQueryBuilder(settings *QbSettings) *UpdateQb {
	qb := &UpdateQb{}
	qb.WhereBuilder = NewWhereBuilder(qb, settings)
	return qb
}

func (qb *UpdateQb) Table(table string) *UpdateQb {
	qb.tableClause = table
	return qb
}

func (qb *UpdateQb) Set(col string, val interface{}) *UpdateQb {
	qb.setClauses = append(qb.setClauses, setClause{col: col, val: val})
	return qb
}

func (qb *UpdateQb) Clone() *UpdateQb {
	clone := &UpdateQb{
		tableClause: qb.tableClause,
		cteQb:       qb.cteQb,
	}
	if qb.setClauses != nil {
		clone.setClauses = make([]setClause, len(qb.setClauses))
		copy(clone.setClauses, qb.setClauses)
	}
	clone.WhereBuilder = NewWhereBuilder(clone, qb.settings)
	clone.whereClause = qb.cloneWhereClause()
	return clone
}

func (qb *UpdateQb) ExtendSql(ctx *SqlBuildingCtx) error {
	if _, err := ctx.Sql.WriteString("update"); err != nil {
		return err
	}

	if qb.tableClause == "" {
		return fmt.Errorf("update: table is required")
	}

	if err := ctx.Sql.WriteByte(' '); err != nil {
		return err
	}
	if err := ctx.WriteArg(&qb.tableClause, true); err != nil {
		return err
	}

	if len(qb.setClauses) == 0 {
		return fmt.Errorf("update: at least one set clause is required")
	}

	if _, err := ctx.Sql.WriteString(" set "); err != nil {
		return err
	}

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

		if err := ctx.WriteArg(WrapValue(s.val), false); err != nil {
			return err
		}
	}

	if qb.whereClause != nil {
		if _, err := ctx.Sql.WriteString(" where "); err != nil {
			return err
		}
		if err := qb.whereClause.ExtendSql(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (qb *UpdateQb) ToSql() (string, []interface{}, error) {
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
