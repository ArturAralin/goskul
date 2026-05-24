package internal

type orderClause struct {
	term      interface{} // string or *QbRaw
	direction string
	nulls     string // "" | "first" | "last"
}

func (o orderClause) ExtendSql(ctx *SqlBuildingCtx) error {
	switch t := o.term.(type) {
	case string:
		if err := ctx.WriteArg(&t, false); err != nil {
			return err
		}
	case *QbRaw:
		if err := t.ExtendSql(ctx); err != nil {
			return err
		}
	}

	if _, err := ctx.Sql.WriteString(" " + o.direction); err != nil {
		return err
	}

	if o.nulls != "" {
		if _, err := ctx.Sql.WriteString(" nulls " + o.nulls); err != nil {
			return err
		}
	}

	return nil
}
