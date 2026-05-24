package goskul

import "github.com/ArturAralin/goskul/internal"

type JoinClause = internal.JoinClause
type SubCond = internal.SubCond
type DoUpdateQb = internal.DoUpdateQb
type SelectQb = internal.SelectQb
type QbRaw = internal.QbRaw

var NewQbSettings = internal.NewQbSettings

type Goskul struct {
	settings *internal.QbSettings
}

func PostgreSQLSettings() *internal.QbSettings {
	settings := internal.NewQbSettings()
	settings.BindingSymbol = '$'
	settings.BindingNumeration = true
	settings.RelationSymbol = '"'
	settings.ReturningFeature = true
	settings.LockForUpdateFeature = true
	settings.LockForShareFeature = true
	settings.LockForKeyShareFeature = true
	settings.LockForNoKeyUpdateFeature = true
	settings.LockSkipLockedFeature = true
	settings.LockNoWaitFeature = true

	return settings
}

func MySQLSettings() *internal.QbSettings {
	settings := internal.NewQbSettings()
	settings.BindingSymbol = '?'
	settings.BindingNumeration = false
	settings.RelationSymbol = '`'
	settings.ReturningFeature = true
	settings.LockForUpdateFeature = true
	settings.LockForShareFeature = true
	settings.LockForKeyShareFeature = false
	settings.LockForNoKeyUpdateFeature = false
	settings.LockSkipLockedFeature = true
	settings.LockNoWaitFeature = true

	return settings
}

func SetupQueryBuilder(settings *internal.QbSettings) *Goskul {
	return &Goskul{
		settings: settings,
	}
}

func (s *Goskul) Select() *internal.SelectQb {
	return internal.NewSelectQueryBuilder(s.settings)
}

func (s *Goskul) Delete() *internal.DeleteQb {
	return internal.NewDeleteQueryBuilder(s.settings)
}

func (s *Goskul) Update() *internal.UpdateQb {
	return internal.NewUpdateQueryBuilder(s.settings)
}

func (s *Goskul) Insert() *internal.InsertQb {
	return internal.NewInsertQueryBuilder(s.settings)
}

func (s *Goskul) With(alias string, qb *internal.SelectQb) *internal.CteQb {
	return internal.NewCteQb(s.settings, alias, qb)
}

func (s *Goskul) Raw(raw string, args ...interface{}) *internal.QbRaw {
	return internal.NewQbRaw(s.settings, raw, args...)
}

func Rel(rel string) *internal.Relation {
	return internal.NewRelation(rel)
}
