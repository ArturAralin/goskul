package internal

type Relation struct {
	rel string
}

func NewRelation(rel string) *Relation {
	return &Relation{rel: rel}
}

// plain value binding
type BindingValue struct {
	Binding interface{}
}

func NewBindingValue(binding interface{}) *BindingValue {
	return &BindingValue{Binding: binding}
}

func (bv *BindingValue) clone() *BindingValue {
	return &BindingValue{Binding: bv.Binding}
}

type QbSettings struct {
	BindingSymbol          byte
	RelationSymbol         byte
	BindingValidatorFn     *func(binding interface{}) error
	InsertReturningFeature bool
}

func (qbSettings *QbSettings) ValidateBinding(binding interface{}) error {
	if qbSettings.BindingValidatorFn != nil {
		return (*qbSettings.BindingValidatorFn)(binding)
	}

	return nil
}

func NewQbSettings() *QbSettings {
	return &QbSettings{}
}

type QbRaw struct {
	settings *QbSettings
	sql      string
	bindings []interface{}
}

func (raw *QbRaw) ExtendSql(ctx *SqlBuildingCtx) error {
	bindingIdx := 0

	for _, ch := range raw.sql {
		if ch == '?' {
			if err := ctx.WriteArg(raw.bindings[bindingIdx], false); err != nil {
				return err
			}

			bindingIdx++
		} else {
			ctx.Sql.WriteRune(ch)
		}
	}

	return nil
}

func WrapValue(v interface{}) interface{} {
	switch v.(type) {
	case *QbRaw, QbRaw:
		return v
	}

	return NewBindingValue(v)
}

func (raw *QbRaw) ToSql() (string, []interface{}, error) {
	return raw.sql, raw.bindings, nil
}

func NewQbRaw(settings *QbSettings, sql string, bindings ...interface{}) *QbRaw {
	wrappedBindings := make([]interface{}, len(bindings))

	for i, binding := range bindings {
		wrappedBindings[i] = WrapValue(binding)
	}

	return &QbRaw{
		settings: settings,
		sql:      sql,
		bindings: wrappedBindings,
	}
}

// todo: raw clone
