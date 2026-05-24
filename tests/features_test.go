package tests

import (
	"errors"
	"testing"

	goskul "github.com/ArturAralin/goskul"
	"github.com/ArturAralin/goskul/internal"
)

var mysqlQb = goskul.SetupQueryBuilder(goskul.MySQLSettings())

// MySQL settings — binding symbol

func TestMySQLSettingsScalarBinding(t *testing.T) {
	sql, args, err := mysqlQb.Select().From("t").Where("id", "=", 42).ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "select * from `t` where `id` = ?"
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 || args[0] != 42 {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestMySQLSettingsSliceBinding(t *testing.T) {
	sql, args, err := mysqlQb.Select().From("t").Where("id", "in", []int{1, 2}).ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "select * from `t` where `id` in (?, ?)"
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}

func TestMySQLSettingsRelationSymbol(t *testing.T) {
	sql, _, err := mysqlQb.Select().Columns("t.id").From("t").ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "select `t`.`id` from `t`"
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestMySQLSettingsMultipleBindingsOrdered(t *testing.T) {
	sql, args, err := mysqlQb.Select().From("t").
		Where("a", "=", 1).
		Where("b", "=", 2).
		ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "select * from `t` where `a` = ? and `b` = ?"
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != 2 {
		t.Errorf("unexpected args: %v", args)
	}
}

// MySQL settings — disabled feature flags

func TestMySQLSettingsLockForKeySharePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for LockForKeyShareFeature disabled in MySQL")
		}
	}()
	mysqlQb.Select().LockForKeyShare()
}

func TestMySQLSettingsLockForNoKeyUpdatePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for LockForNoKeyUpdateFeature disabled in MySQL")
		}
	}()
	mysqlQb.Select().LockForNoKeyUpdate()
}

// BindingValidatorFn

func TestBindingValidatorFnNilIsNoop(t *testing.T) {
	settings := internal.NewQbSettings()
	settings.BindingSymbol = '$'
	settings.BindingNumeration = true
	settings.RelationSymbol = '"'
	qb := goskul.SetupQueryBuilder(settings)

	_, _, err := qb.Select().From("t").Where("id", "=", 1).ToSql()
	if err != nil {
		t.Errorf("unexpected error with nil validator: %v", err)
	}
}

func TestBindingValidatorFnValidValuePasses(t *testing.T) {
	fn := func(v interface{}) error { return nil }
	settings := goskul.PostgreSQLSettings()
	settings.BindingValidatorFn = &fn
	qb := goskul.SetupQueryBuilder(settings)

	_, args, err := qb.Select().From("t").Where("id", "=", 99).ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 1 || args[0] != 99 {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestBindingValidatorFnErrorPropagates(t *testing.T) {
	validationErr := errors.New("invalid binding")
	fn := func(v interface{}) error { return validationErr }
	settings := goskul.PostgreSQLSettings()
	settings.BindingValidatorFn = &fn
	qb := goskul.SetupQueryBuilder(settings)

	_, _, err := qb.Select().From("t").Where("id", "=", 1).ToSql()
	if err == nil {
		t.Fatal("expected error from validator, got nil")
	}
	if !errors.Is(err, validationErr) {
		t.Errorf("got %v, want %v", err, validationErr)
	}
}
