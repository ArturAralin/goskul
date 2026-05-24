package internal

type ReturningQb struct {
	columns []interface{}
}

func (qb *ReturningQb) add(col interface{}) {
	qb.columns = append(qb.columns, col)
}

func (qb *ReturningQb) clone() *ReturningQb {
	cloned := &ReturningQb{}
	if qb.columns != nil {
		cloned.columns = make([]interface{}, len(qb.columns))
		copy(cloned.columns, qb.columns)
	}
	return cloned
}

func (qb *ReturningQb) ExtendSql(ctx *SqlBuildingCtx) error {
	if _, err := ctx.Sql.WriteString(" returning "); err != nil {
		return err
	}
	for i, col := range qb.columns {
		if i > 0 {
			if _, err := ctx.Sql.WriteString(", "); err != nil {
				return err
			}
		}
		if _, ok := col.(string); ok { // only "*" is stored as a plain string
			if err := ctx.Sql.WriteByte('*'); err != nil {
				return err
			}
		} else {
			if err := ctx.WriteArg(col, false); err != nil {
				return err
			}
		}
	}
	return nil
}

type ReturningBuilder[T any] struct {
	self            *T
	settings        *QbSettings
	returningClause *ReturningQb
}

func NewReturningBuilder[T any](self *T, settings *QbSettings) ReturningBuilder[T] {
	return ReturningBuilder[T]{self: self, settings: settings}
}

func (b *ReturningBuilder[T]) Returning(cols ...interface{}) *T {
	if !b.settings.InsertReturningFeature {
		panic("InsertReturningFeature is not supported by database settings")
	}
	if len(cols) == 0 {
		panic("Returning requires at least one column or *")
	}
	if b.returningClause == nil {
		b.returningClause = &ReturningQb{}
	}
	for _, col := range cols {
		if s, ok := col.(string); ok {
			if s == "*" {
				b.returningClause.add("*")
			} else {
				b.returningClause.add(NewRelation(s))
			}
		} else {
			b.returningClause.add(col)
		}
	}
	return b.self
}

func (b *ReturningBuilder[T]) cloneReturningClause() *ReturningQb {
	if b.returningClause == nil {
		return nil
	}
	return b.returningClause.clone()
}
