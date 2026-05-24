package internal

type cteEntry struct {
	alias string
	qb    *SelectQb
}

type CteQb struct {
	settings *QbSettings
	ctes     []cteEntry
}

func NewCteQb(settings *QbSettings, alias string, qb *SelectQb) *CteQb {
	return &CteQb{
		settings: settings,
		ctes:     []cteEntry{{alias: alias, qb: qb}},
	}
}

func (qb *CteQb) With(alias string, selectQb *SelectQb) *CteQb {
	qb.ctes = append(qb.ctes, cteEntry{alias: alias, qb: selectQb})
	return qb
}

func (qb *CteQb) Select() *SelectQb {
	sel := NewSelectQueryBuilder(qb.settings)
	sel.cteQb = qb
	return sel
}

func (qb *CteQb) Delete() *DeleteQb {
	del := NewDeleteQueryBuilder(qb.settings)
	del.cteQb = qb
	return del
}

func (qb *CteQb) Insert() *InsertQb {
	ins := NewInsertQueryBuilder(qb.settings)
	ins.cteQb = qb
	return ins
}

func (qb *CteQb) Update() *UpdateQb {
	upd := NewUpdateQueryBuilder(qb.settings)
	upd.cteQb = qb
	return upd
}

func (qb *CteQb) ExtendSql(ctx *SqlBuildingCtx) error {
	ctx.Sql.WriteString("with ")
	for i, cte := range qb.ctes {
		if i > 0 {
			ctx.Sql.WriteString(", ")
		}
		if err := ctx.WriteRelation(cte.alias, false); err != nil {
			return err
		}
		ctx.Sql.WriteString(" as (")
		if err := cte.qb.ExtendSql(ctx); err != nil {
			return err
		}
		ctx.Sql.WriteByte(')')
	}
	return nil
}
