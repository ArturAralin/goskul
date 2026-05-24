package internal

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"

	"github.com/modern-go/reflect2"
)

var aliasRegex regexp.Regexp = *regexp.MustCompile(`.\s+as\s+(?P<alias>\w+)\s*$`)

type SqlBuildingCtx struct {
	settings *QbSettings
	Sql      bytes.Buffer
	bindings []interface{}
}

func (ctx *SqlBuildingCtx) WriteBinding(v interface{}) (int, error) {
	if err := ctx.settings.ValidateBinding(v); err != nil {
		return -1, err
	}

	ctx.bindings = append(ctx.bindings, v)

	return len(ctx.bindings), nil
}

func NewSqlBuildingCtx(settings *QbSettings) SqlBuildingCtx {
	return SqlBuildingCtx{
		settings: settings,
		Sql:      *bytes.NewBuffer(make([]byte, 0, settings.BuildingSqlDefaultBufferSize)),
	}
}

func (ctx *SqlBuildingCtx) WriteRelation(rel string, aliasing bool) error {
	if len(rel) == 0 {
		return fmt.Errorf("relation name cannot be empty")
	}

	aliasIndices := aliasRegex.FindStringSubmatchIndex(rel)

	if !aliasing && aliasIndices != nil {
		return fmt.Errorf("alias is not allowed here")
	}

	relStartIdx := 0
	relEndIdx := len(rel)

	if aliasIndices != nil {
		// +1 for regex first char compensation
		relEndIdx = aliasIndices[0] + 1
	}

	// skip leading spaces
	for idx, ch := range rel {
		if ch != ' ' {
			relStartIdx = idx
			break
		}
	}

	if relStartIdx >= relEndIdx {
		return fmt.Errorf("relation name cannot be blank")
	}

	partStart := relStartIdx
	for i := relStartIdx; i <= relEndIdx; i++ {
		if i == relEndIdx || rel[i] == '.' {
			if err := ctx.Sql.WriteByte(ctx.settings.RelationSymbol); err != nil {
				return err
			}
			for j := partStart; j < i; j++ {
				if err := ctx.Sql.WriteByte(rel[j]); err != nil {
					return err
				}
			}
			if err := ctx.Sql.WriteByte(ctx.settings.RelationSymbol); err != nil {
				return err
			}
			if i < relEndIdx {
				if err := ctx.Sql.WriteByte('.'); err != nil {
					return err
				}
				partStart = i + 1
			}
		}
	}

	if aliasIndices != nil {
		alias := rel[aliasIndices[2]:aliasIndices[3]]
		if _, err := ctx.Sql.WriteString(" as "); err != nil {
			return err
		}
		if err := ctx.Sql.WriteByte(ctx.settings.RelationSymbol); err != nil {
			return err
		}
		if _, err := ctx.Sql.WriteString(alias); err != nil {
			return err
		}
		if err := ctx.Sql.WriteByte(ctx.settings.RelationSymbol); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *SqlBuildingCtx) WriteArg(rel interface{}, relationAliasing bool) error {
	switch rel := rel.(type) {
	case *Relation:
		if err := ctx.WriteRelation(rel.rel, relationAliasing); err != nil {
			return err
		}
	case *BindingValue:
		{
			v := reflect2.TypeOf(rel.Binding)

			switch v.Kind() {
			case reflect.Slice, reflect.Array:
				v := reflect.ValueOf(rel.Binding)

				ctx.Sql.WriteByte('(')

				for i := 0; i < v.Len(); i++ {
					if i > 0 {
						ctx.Sql.WriteByte(',')
						ctx.Sql.WriteByte(' ')
					}

					n, err := ctx.WriteBinding(v.Index(i).Interface())

					if err != nil {
						return err
					}

					ctx.Sql.WriteByte(ctx.settings.BindingSymbol)

					if ctx.settings.BindingNumeration {
						ctx.Sql.WriteString(fmt.Sprintf("%d", n))
					}
				}

				ctx.Sql.WriteByte(')')
			default:
				{
					n, err := ctx.WriteBinding(rel.Binding)

					if err != nil {
						return err
					}

					ctx.Sql.WriteByte(ctx.settings.BindingSymbol)
					if ctx.settings.BindingNumeration {
						ctx.Sql.WriteString(fmt.Sprintf("%d", n))
					}
				}
			}
		}
	case *QbRaw:
		if err := rel.ExtendSql(ctx); err != nil {
			return err
		}
	case *SelectQb:
		if err := ctx.Sql.WriteByte('('); err != nil {
			return err
		}
		if err := rel.ExtendSql(ctx); err != nil {
			return err
		}
		return ctx.Sql.WriteByte(')')
	default:
		panic(fmt.Sprintf("Unsupported type %v", reflect2.TypeOf(rel)))
	}

	return nil
}
