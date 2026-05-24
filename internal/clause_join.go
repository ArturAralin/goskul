package internal

type JoinClause struct {
	joinType string
	table    *string
	cond     *ConditionQb
	rawSql   *QbRaw
}

func prepareRightArg(value interface{}) interface{} {
	switch value := value.(type) {
	case *QbRaw, *Relation, *SelectQb:
		return value
	case string:
		return NewRelation(value)
	case *string:
		return NewRelation(*value)
	default:
		return NewBindingValue(value)
	}
}

func prepareLeftArg(value interface{}) interface{} {
	switch value := value.(type) {
	case *QbRaw:
		return value
	case string:
		return NewRelation(value)
	case *string:
		return NewRelation(*value)
	default:
		return NewBindingValue(value)
	}
}

func NewJoinClause(settings *QbSettings, joinType string, table *string) *JoinClause {
	return &JoinClause{
		joinType: joinType,
		table:    table,
		cond:     NewConditionQb(settings),
	}
}

func NewCrossJoinClause(table string) *JoinClause {
	return &JoinClause{
		joinType: "cross join",
		table:    &table,
	}
}

func NewRawJoinClause(raw *QbRaw) *JoinClause {
	return &JoinClause{rawSql: raw}
}

func (qb *JoinClause) On(left interface{}, op string, right interface{}) *JoinClause {
	qb.cond.PushBinaryCond("and", prepareLeftArg(left), op, prepareRightArg(right))

	return qb
}

func (qb *JoinClause) AndOn(left interface{}, op string, right interface{}) *JoinClause {
	qb.cond.PushBinaryCond("and", prepareLeftArg(left), op, prepareRightArg(right))

	return qb
}

func (qb *JoinClause) OrOn(left interface{}, op string, right interface{}) *JoinClause {
	qb.cond.PushBinaryCond("or", prepareLeftArg(left), op, prepareRightArg(right))

	return qb
}

func (qb *JoinClause) OnNull(col interface{}) *JoinClause {
	qb.cond.PushUnaryCondLeft("and", "is null", prepareLeftArg(col))
	return qb
}

func (qb *JoinClause) AndOnNull(col interface{}) *JoinClause {
	qb.cond.PushUnaryCondLeft("and", "is null", prepareLeftArg(col))
	return qb
}

func (qb *JoinClause) OrOnNull(col interface{}) *JoinClause {
	qb.cond.PushUnaryCondLeft("or", "is null", prepareLeftArg(col))
	return qb
}

func (qb *JoinClause) OnRaw(raw *QbRaw) *JoinClause {
	qb.cond.PushRawCond("and", raw)
	return qb
}

func (qb *JoinClause) OrOnRaw(raw *QbRaw) *JoinClause {
	qb.cond.PushRawCond("or", raw)
	return qb
}

func (qb *JoinClause) onSubCondOp(unionOp string, fn func(*SubCond)) *JoinClause {
	sub := newSubCond(qb.cond.settings)
	fn(sub)
	if sub.whereClause != nil {
		qb.cond.PushSubCond(unionOp, sub.whereClause)
	}
	return qb
}

func (qb *JoinClause) OnSubCond(fn func(*SubCond)) *JoinClause {
	return qb.onSubCondOp("and", fn)
}

func (qb *JoinClause) AndOnSubCond(fn func(*SubCond)) *JoinClause {
	return qb.onSubCondOp("and", fn)
}

func (qb *JoinClause) OrOnSubCond(fn func(*SubCond)) *JoinClause {
	return qb.onSubCondOp("or", fn)
}

// todo: add OnBetween

func (qb *JoinClause) Clone() *JoinClause {
	if qb.rawSql != nil {
		return &JoinClause{rawSql: qb.rawSql}
	}
	clone := &JoinClause{
		joinType: qb.joinType,
		cond:     qb.cond.clone(),
	}
	if qb.table != nil {
		t := *qb.table
		clone.table = &t
	}
	return clone
}

func (qb *JoinClause) ExtendSql(ctx *SqlBuildingCtx) error {
	if qb.rawSql != nil {
		return qb.rawSql.ExtendSql(ctx)
	}

	if _, err := ctx.Sql.WriteString(qb.joinType); err != nil {
		return err
	}

	if err := ctx.Sql.WriteByte(' '); err != nil {
		return err
	}

	if qb.table != nil {
		if err := ctx.WriteRelation(*qb.table, true); err != nil {
			return err
		}
	}

	if qb.cond != nil {
		if _, err := ctx.Sql.WriteString(" on "); err != nil {
			return err
		}

		if err := qb.cond.ExtendSql(ctx); err != nil {
			return err
		}
	}

	return nil
}
