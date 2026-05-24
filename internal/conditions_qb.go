package internal

type binaryCond struct {
	unionOp string
	left    interface{}
	op      string
	right   interface{}
}

type unaryCondRight struct {
	unionOp string
	op      string
	val     interface{}
}

type unaryCondLeft struct {
	unionOp string
	op      string
	val     interface{}
}

type rawCond struct {
	unionOp string
	raw     *QbRaw
}

type subCond struct {
	unionOp    string
	conditions *ConditionQb
}

type ConditionQb struct {
	settings   *QbSettings
	conditions []interface{}
}

func NewConditionQb(settings *QbSettings) *ConditionQb {
	return &ConditionQb{
		settings: settings,
	}
}

func (qb *ConditionQb) PushBinaryCond(unionOp string, left interface{}, op string, right interface{}) {
	// todo: validate op

	qb.conditions = append(qb.conditions, binaryCond{
		unionOp: unionOp,
		left:    left,
		op:      op,
		right:   right,
	})
}

func (qb *ConditionQb) PushUnaryCondRight(unionOp string, op string, val interface{}) {
	// todo: validate op
	qb.conditions = append(qb.conditions, unaryCondRight{
		unionOp: unionOp,
		op:      op,
		val:     val,
	})
}

func (qb *ConditionQb) PushRawCond(unionOp string, raw *QbRaw) {
	qb.conditions = append(qb.conditions, rawCond{unionOp: unionOp, raw: raw})
}

func (qb *ConditionQb) PushSubCond(unionOp string, conditions *ConditionQb) {
	qb.conditions = append(qb.conditions, subCond{
		unionOp:    unionOp,
		conditions: conditions,
	})
}

func (qb *ConditionQb) PushUnaryCondLeft(unionOp string, op string, val interface{}) {
	// todo: validate op
	qb.conditions = append(qb.conditions, unaryCondLeft{
		unionOp: unionOp,
		op:      op,
		val:     val,
	})
}

func (qb *ConditionQb) clone() *ConditionQb {
	cloned := &ConditionQb{
		settings: qb.settings,
	}
	if qb.conditions == nil {
		return cloned
	}
	cloned.conditions = make([]interface{}, len(qb.conditions))
	for i, cond := range qb.conditions {
		switch c := cond.(type) {
		case binaryCond:
			copied := c
			if bv, ok := c.right.(*BindingValue); ok {
				copied.right = bv.clone()
			}
			cloned.conditions[i] = copied
		case unaryCondRight:
			copied := c
			if bv, ok := c.val.(*BindingValue); ok {
				copied.val = bv.clone()
			}
			cloned.conditions[i] = copied
		case unaryCondLeft:
			cloned.conditions[i] = c
		case rawCond:
			cloned.conditions[i] = c
		case subCond:
			cloned.conditions[i] = subCond{
				unionOp:    c.unionOp,
				conditions: c.conditions.clone(),
			}
		default:
			cloned.conditions[i] = cond
		}
	}
	return cloned
}

// todo: handle error
func (qb *ConditionQb) ExtendSql(ctx *SqlBuildingCtx) error {
	for i, cond := range qb.conditions {
		if i > 0 {
			ctx.Sql.WriteByte(' ')
		}

		switch cond := cond.(type) {
		case binaryCond:
			{
				if i > 0 {
					ctx.Sql.WriteString(cond.unionOp)
					ctx.Sql.WriteByte(' ')
				}

				if err := ctx.WriteArg(cond.left, false); err != nil {
					return err
				}
				ctx.Sql.WriteByte(' ')
				ctx.Sql.WriteString(cond.op)
				ctx.Sql.WriteByte(' ')
				if err := ctx.WriteArg(cond.right, false); err != nil {
					return err
				}
			}
		case unaryCondRight:
			{
				if i > 0 {
					ctx.Sql.WriteString(cond.unionOp)
					ctx.Sql.WriteByte(' ')
				}

				ctx.Sql.WriteString(cond.op)
				ctx.Sql.WriteByte(' ')
				if err := ctx.WriteArg(cond.val, false); err != nil {
					return err
				}
			}
		case unaryCondLeft:
			{
				if i > 0 {
					ctx.Sql.WriteString(cond.unionOp)
					ctx.Sql.WriteByte(' ')
				}

				if err := ctx.WriteArg(cond.val, false); err != nil {
					return err
				}
				ctx.Sql.WriteByte(' ')
				ctx.Sql.WriteString(cond.op)
			}
		case rawCond:
			{
				if i > 0 {
					ctx.Sql.WriteString(cond.unionOp)
					ctx.Sql.WriteByte(' ')
				}

				if err := cond.raw.ExtendSql(ctx); err != nil {
					return err
				}
			}
		case subCond:
			{
				if i > 0 {
					ctx.Sql.WriteString(cond.unionOp)
					ctx.Sql.WriteByte(' ')
				}

				needParens := len(cond.conditions.conditions) > 1 && len(qb.conditions) > 1
				if needParens {
					ctx.Sql.WriteByte('(')
				}
				if err := cond.conditions.ExtendSql(ctx); err != nil {
					return err
				}
				if needParens {
					ctx.Sql.WriteByte(')')
				}
			}
		}

	}
	return nil
}
