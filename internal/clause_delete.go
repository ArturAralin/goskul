package internal

import "fmt"

type DeleteQb struct {
	WhereBuilder[DeleteQb]
	ReturningBuilder[DeleteQb]
	fromClause interface{}
	cteQb      *CteQb
}

func NewDeleteQueryBuilder(settings *QbSettings) *DeleteQb {
	qb := &DeleteQb{}
	qb.WhereBuilder = NewWhereBuilder(qb, settings)
	qb.ReturningBuilder = NewReturningBuilder(qb, settings)
	return qb
}

func (qb *DeleteQb) From(from interface{}) *DeleteQb {
	qb.fromClause = from
	return qb
}

func (qb *DeleteQb) Clone() *DeleteQb {
	clone := &DeleteQb{
		fromClause: qb.fromClause,
		cteQb:      qb.cteQb,
	}
	clone.WhereBuilder = NewWhereBuilder(clone, qb.WhereBuilder.settings)
	clone.whereClause = qb.cloneWhereClause()
	clone.ReturningBuilder = NewReturningBuilder(clone, qb.WhereBuilder.settings)
	clone.ReturningBuilder.returningClause = qb.cloneReturningClause()
	return clone
}

func (qb *DeleteQb) ExtendSql(ctx *SqlBuildingCtx) error {
	if _, err := ctx.Sql.WriteString("delete"); err != nil {
		return err
	}

	if qb.fromClause != nil {
		if _, err := ctx.Sql.WriteString(" from"); err != nil {
			return err
		}

		switch from := qb.fromClause.(type) {
		case string:
			if _, err := ctx.Sql.WriteString(" "); err != nil {
				return err
			}

			if err := ctx.WriteRelation(from, true); err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("unsupported from clause type %v", from))
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

	if qb.returningClause != nil {
		if err := qb.returningClause.ExtendSql(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (qb *DeleteQb) ToSql() (string, []interface{}, error) {
	ctx := NewSqlBuildingCtx(qb.WhereBuilder.settings)

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
