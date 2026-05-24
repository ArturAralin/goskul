package internal

type selectLockType string
type selectLockModifier string

const (
	selectLockForUpdate      selectLockType = "for update"
	selectLockForShare       selectLockType = "for share"
	selectLockForKeyShare    selectLockType = "for key share"
	selectLockForNoKeyUpdate selectLockType = "for no key update"

	selectLockSkipLocked selectLockModifier = "skip locked"
	selectLockNoWait     selectLockModifier = "nowait"
)

type selectLockClause struct {
	lockType     selectLockType
	lockModifier selectLockModifier
}

func (l *selectLockClause) ExtendSql(ctx *SqlBuildingCtx) error {
	if _, err := ctx.Sql.WriteString(" "); err != nil {
		return err
	}
	if _, err := ctx.Sql.WriteString(string(l.lockType)); err != nil {
		return err
	}
	if l.lockModifier != "" {
		if _, err := ctx.Sql.WriteString(" "); err != nil {
			return err
		}
		if _, err := ctx.Sql.WriteString(string(l.lockModifier)); err != nil {
			return err
		}
	}
	return nil
}
