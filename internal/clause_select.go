package internal

import (
	"fmt"

	"github.com/modern-go/reflect2"
)

type SelectQb struct {
	WhereBuilder[SelectQb]
	selectColumns []interface{}
	fromClause    interface{}
	alias         *string
	cteQb         *CteQb
	joinsClause   []*JoinClause
	orderClauses  []orderClause
	limitClause   *uint64
	offsetClause  *uint64
	lockClause    *selectLockClause
}

func NewSelectQueryBuilder(settings *QbSettings) *SelectQb {
	qb := &SelectQb{}
	qb.WhereBuilder = NewWhereBuilder(qb, settings)
	return qb
}

func (qb *SelectQb) Column(col interface{}) *SelectQb {
	switch col := col.(type) {
	case string:
		// todo: validate format "col" or "table.col"
		qb.selectColumns = append(qb.selectColumns, NewRelation(col))
	case *string:
		// todo: validate format "col" or "table.col"
		qb.selectColumns = append(qb.selectColumns, NewRelation(*col))
	case *QbRaw:
		qb.selectColumns = append(qb.selectColumns, col)
	case QbRaw:
		qb.selectColumns = append(qb.selectColumns, &col)
	default:
		panic(fmt.Errorf("unknown type in Column %s", reflect2.TypeOf(col)))
	}
	return qb
}

func (qb *SelectQb) Columns(cols ...interface{}) *SelectQb {
	for _, col := range cols {
		qb.Column(col)
	}
	return qb
}

// string, *string or Raw
func (qb *SelectQb) From(from interface{}) *SelectQb {
	qb.fromClause = from
	return qb
}

func (qb *SelectQb) Alias(alias string) *SelectQb {
	qb.alias = &alias
	return qb
}

func (qb *SelectQb) addJoin(joinType, tbl string, f func(*JoinClause)) *SelectQb {
	join := NewJoinClause(qb.settings, joinType, &tbl)
	f(join)
	qb.joinsClause = append(qb.joinsClause, join)
	return qb
}

func (qb *SelectQb) addJoinOn(joinType, tbl, left, right string) *SelectQb {
	join := NewJoinClause(qb.settings, joinType, &tbl)
	join.On(left, "=", right)
	qb.joinsClause = append(qb.joinsClause, join)
	return qb
}

func (qb *SelectQb) InnerJoin(tbl string, f func(*JoinClause)) *SelectQb {
	return qb.addJoin("inner join", tbl, f)
}

func (qb *SelectQb) InnerJoinOn(tbl, left, right string) *SelectQb {
	return qb.addJoinOn("inner join", tbl, left, right)
}

func (qb *SelectQb) LeftJoin(tbl string, f func(*JoinClause)) *SelectQb {
	return qb.addJoin("left join", tbl, f)
}

func (qb *SelectQb) LeftJoinOn(tbl, left, right string) *SelectQb {
	return qb.addJoinOn("left join", tbl, left, right)
}

func (qb *SelectQb) RightJoin(tbl string, f func(*JoinClause)) *SelectQb {
	return qb.addJoin("right join", tbl, f)
}

func (qb *SelectQb) RightJoinOn(tbl, left, right string) *SelectQb {
	return qb.addJoinOn("right join", tbl, left, right)
}

func (qb *SelectQb) FullJoin(tbl string, f func(*JoinClause)) *SelectQb {
	return qb.addJoin("full join", tbl, f)
}

func (qb *SelectQb) FullJoinOn(tbl, left, right string) *SelectQb {
	return qb.addJoinOn("full join", tbl, left, right)
}

func (qb *SelectQb) OrderBy(t interface{}, direction string, nulls ...string) *SelectQb {
	var nullsVal string
	if len(nulls) > 0 {
		nullsVal = nulls[0]
	}
	qb.orderClauses = append(qb.orderClauses, orderClause{term: t, direction: direction, nulls: nullsVal})
	return qb
}

func (qb *SelectQb) Limit(n uint64) *SelectQb {
	qb.limitClause = &n
	return qb
}

func (qb *SelectQb) Offset(n uint64) *SelectQb {
	qb.offsetClause = &n
	return qb
}

func (qb *SelectQb) LockForUpdate() *SelectQb {
	if !qb.settings.LockForUpdateFeature {
		panic("LockForUpdateFeature is not supported by database settings")
	}
	if qb.lockClause == nil {
		qb.lockClause = &selectLockClause{}
	}
	qb.lockClause.lockType = selectLockForUpdate
	return qb
}

func (qb *SelectQb) LockForShare() *SelectQb {
	if !qb.settings.LockForShareFeature {
		panic("LockForShareFeature is not supported by database settings")
	}
	if qb.lockClause == nil {
		qb.lockClause = &selectLockClause{}
	}
	qb.lockClause.lockType = selectLockForShare
	return qb
}

func (qb *SelectQb) LockForNoKeyUpdate() *SelectQb {
	if !qb.settings.LockForNoKeyUpdateFeature {
		panic("LockForNoKeyUpdateFeature is not supported by database settings")
	}
	if qb.lockClause == nil {
		qb.lockClause = &selectLockClause{}
	}
	qb.lockClause.lockType = selectLockForNoKeyUpdate
	return qb
}

func (qb *SelectQb) LockForKeyShare() *SelectQb {
	if !qb.settings.LockForKeyShareFeature {
		panic("LockForKeyShareFeature is not supported by database settings")
	}
	if qb.lockClause == nil {
		qb.lockClause = &selectLockClause{}
	}
	qb.lockClause.lockType = selectLockForKeyShare
	return qb
}

func (qb *SelectQb) LockSkipLocked() *SelectQb {
	if !qb.settings.LockSkipLockedFeature {
		panic("LockSkipLockedFeature is not supported by database settings")
	}
	if qb.lockClause == nil {
		qb.lockClause = &selectLockClause{}
	}
	qb.lockClause.lockModifier = selectLockSkipLocked
	return qb
}

func (qb *SelectQb) LockNoWait() *SelectQb {
	if !qb.settings.LockNoWaitFeature {
		panic("LockNoWaitFeature is not supported by database settings")
	}
	if qb.lockClause == nil {
		qb.lockClause = &selectLockClause{}
	}
	qb.lockClause.lockModifier = selectLockNoWait
	return qb
}

func (qb *SelectQb) CrossJoin(tbl string) *SelectQb {
	qb.joinsClause = append(qb.joinsClause, NewCrossJoinClause(tbl))
	return qb
}

func (qb *SelectQb) JoinRaw(raw *QbRaw) *SelectQb {
	qb.joinsClause = append(qb.joinsClause, NewRawJoinClause(raw))
	return qb
}

func (qb *SelectQb) Clone() *SelectQb {
	clone := &SelectQb{
		fromClause: qb.fromClause,
		alias:      qb.alias,
		cteQb:      qb.cteQb,
	}
	if qb.selectColumns != nil {
		clone.selectColumns = make([]interface{}, len(qb.selectColumns))
		copy(clone.selectColumns, qb.selectColumns)
	}
	if qb.joinsClause != nil {
		clone.joinsClause = make([]*JoinClause, len(qb.joinsClause))
		for i, j := range qb.joinsClause {
			clone.joinsClause[i] = j.Clone()
		}
	}
	if qb.orderClauses != nil {
		clone.orderClauses = make([]orderClause, len(qb.orderClauses))
		copy(clone.orderClauses, qb.orderClauses)
	}
	clone.limitClause = qb.limitClause
	clone.offsetClause = qb.offsetClause
	if qb.lockClause != nil {
		lc := *qb.lockClause
		clone.lockClause = &lc
	}
	clone.WhereBuilder = NewWhereBuilder(clone, qb.settings)
	clone.whereClause = qb.cloneWhereClause()
	return clone
}

func (qb *SelectQb) ExtendSql(ctx *SqlBuildingCtx) error {
	if _, err := ctx.Sql.WriteString("select"); err != nil {
		return err
	}

	if qb.selectColumns != nil {
		for i, col := range qb.selectColumns {
			if i == 0 {
				if _, err := ctx.Sql.WriteString(" "); err != nil {
					return err
				}
			} else {
				if _, err := ctx.Sql.WriteString(", "); err != nil {
					return err
				}
			}

			if err := ctx.WriteArg(col, true); err != nil {
				return err
			}
		}
	} else {
		if _, err := ctx.Sql.WriteString(" *"); err != nil {
			return err
		}
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
		case *QbRaw:
			if _, err := ctx.Sql.WriteString(" "); err != nil {
				return err
			}

			if err := from.ExtendSql(ctx); err != nil {
				return err
			}
		case *SelectQb:
			if err := ctx.Sql.WriteByte(' '); err != nil {
				return err
			}

			if from.alias != nil {
				if err := ctx.Sql.WriteByte('('); err != nil {
					return err
				}
			}

			if err := from.ExtendSql(ctx); err != nil {
				return err
			}

			if from.alias != nil {
				if _, err := ctx.Sql.WriteString(") as "); err != nil {
					return err
				}

				if err := ctx.WriteRelation(*from.alias, false); err != nil {
					return err
				}
			}
		default:
			panic(fmt.Sprintf("unsupported from clause type %v", reflect2.TypeOf(from)))
		}
	}

	if len(qb.joinsClause) > 0 {
		for _, join := range qb.joinsClause {
			if err := ctx.Sql.WriteByte(' '); err != nil {
				return err
			}

			if err := join.ExtendSql(ctx); err != nil {
				return err
			}
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

	if len(qb.orderClauses) > 0 {
		if _, err := ctx.Sql.WriteString(" order by "); err != nil {
			return err
		}
		for i, order := range qb.orderClauses {
			if i > 0 {
				if _, err := ctx.Sql.WriteString(", "); err != nil {
					return err
				}
			}
			if err := order.ExtendSql(ctx); err != nil {
				return err
			}
		}
	}

	if qb.limitClause != nil {
		if _, err := fmt.Fprintf(&ctx.Sql, " limit %d", *qb.limitClause); err != nil {
			return err
		}
	}

	if qb.offsetClause != nil {
		if _, err := fmt.Fprintf(&ctx.Sql, " offset %d", *qb.offsetClause); err != nil {
			return err
		}
	}

	if qb.lockClause != nil && qb.lockClause.lockType != "" {
		if err := qb.lockClause.ExtendSql(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (qb *SelectQb) ToSql() (string, []interface{}, error) {
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
