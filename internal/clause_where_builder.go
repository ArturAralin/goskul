package internal

type WhereBuilder[T any] struct {
	self        *T
	settings    *QbSettings
	whereClause *ConditionQb
}

func NewWhereBuilder[T any](self *T, settings *QbSettings) WhereBuilder[T] {
	return WhereBuilder[T]{self: self, settings: settings}
}

func (b *WhereBuilder[T]) initWhere() {
	if b.whereClause == nil {
		b.whereClause = NewConditionQb(b.settings)
	}
}

func rawOrSubQueryOrBindingV(value interface{}) interface{} {
	switch value := value.(type) {
	case *QbRaw, *SelectQb, *Relation:
		return value
	default:
		return NewBindingValue(value)
	}
}

func rawOrPass(value interface{}) interface{} {
	switch value := value.(type) {
	case *QbRaw:
		return value
	case string:
		return NewRelation(value)
	case *string:
		return NewRelation(*value)
	default:
		return value
	}
}

func prepareLeftWhereArg(value interface{}) interface{} {
	switch value := value.(type) {
	case *QbRaw, *SelectQb, *Relation:
		return value
	case string:
		return NewRelation(value)
	case *string:
		return NewRelation(*value)
	default:
		return NewBindingValue(value)
	}
}

func (b *WhereBuilder[T]) Where(col interface{}, op string, value interface{}) *T {
	b.initWhere()
	b.whereClause.PushBinaryCond("and", prepareLeftWhereArg(col), op, rawOrSubQueryOrBindingV(value))
	return b.self
}

func (b *WhereBuilder[T]) OrWhere(col interface{}, op string, value interface{}) *T {
	b.initWhere()
	b.whereClause.PushBinaryCond("or", prepareLeftWhereArg(col), op, rawOrSubQueryOrBindingV(value))
	return b.self
}

func (b *WhereBuilder[T]) WhereNull(v interface{}) *T {
	b.initWhere()
	b.whereClause.PushUnaryCondLeft("and", "is null", rawOrPass(v))
	return b.self
}

func (b *WhereBuilder[T]) OrWhereNull(v interface{}) *T {
	b.initWhere()
	b.whereClause.PushUnaryCondLeft("or", "is null", rawOrPass(v))
	return b.self
}

func (b *WhereBuilder[T]) WhereNotNull(v interface{}) *T {
	b.initWhere()
	b.whereClause.PushUnaryCondLeft("and", "is not null", rawOrPass(v))
	return b.self
}

func (b *WhereBuilder[T]) OrWhereNotNull(v interface{}) *T {
	b.initWhere()
	b.whereClause.PushUnaryCondLeft("or", "is not null", rawOrPass(v))
	return b.self
}

func (b *WhereBuilder[T]) WhereRaw(raw *QbRaw) *T {
	b.initWhere()
	b.whereClause.PushRawCond("and", raw)
	return b.self
}

func (b *WhereBuilder[T]) OrWhereRaw(raw *QbRaw) *T {
	b.initWhere()
	b.whereClause.PushRawCond("or", raw)
	return b.self
}

func (b *WhereBuilder[T]) WhereIn(col interface{}, values interface{}) *T {
	b.initWhere()
	b.whereClause.PushBinaryCond("and", prepareLeftWhereArg(col), "in", rawOrSubQueryOrBindingV(values))
	return b.self
}

func (b *WhereBuilder[T]) WhereNotIn(col interface{}, values interface{}) *T {
	b.initWhere()
	b.whereClause.PushBinaryCond("and", prepareLeftWhereArg(col), "not in", rawOrSubQueryOrBindingV(values))
	return b.self
}

func (b *WhereBuilder[T]) OrWhereIn(col interface{}, values interface{}) *T {
	b.initWhere()
	b.whereClause.PushBinaryCond("or", prepareLeftWhereArg(col), "in", rawOrSubQueryOrBindingV(values))
	return b.self
}

func (b *WhereBuilder[T]) OrWhereNotIn(col interface{}, values interface{}) *T {
	b.initWhere()
	b.whereClause.PushBinaryCond("or", prepareLeftWhereArg(col), "not in", rawOrSubQueryOrBindingV(values))
	return b.self
}

type SubCond struct {
	WhereBuilder[SubCond]
}

func newSubCond(settings *QbSettings) *SubCond {
	s := &SubCond{}
	s.WhereBuilder = NewWhereBuilder[SubCond](s, settings)
	return s
}

func (b *WhereBuilder[T]) whereSubCondOp(unionOp string, fn func(*SubCond)) *T {
	b.initWhere()
	sub := newSubCond(b.settings)
	fn(sub)
	if sub.whereClause != nil {
		b.whereClause.PushSubCond(unionOp, sub.whereClause)
	}
	return b.self
}

func (b *WhereBuilder[T]) WhereSubCond(fn func(*SubCond)) *T {
	return b.whereSubCondOp("and", fn)
}

func (b *WhereBuilder[T]) AndWhereSubCond(fn func(*SubCond)) *T {
	return b.whereSubCondOp("and", fn)
}

func (b *WhereBuilder[T]) OrWhereSubCond(fn func(*SubCond)) *T {
	return b.whereSubCondOp("or", fn)
}

func (b *WhereBuilder[T]) cloneWhereClause() *ConditionQb {
	if b.whereClause == nil {
		return nil
	}
	return b.whereClause.clone()
}
