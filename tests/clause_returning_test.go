package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
)

func TestDeleteReturning(t *testing.T) {
	qb := sqlQb.Delete().
		From("my_table").
		Where("id", "=", 1).
		Returning("id", "name")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `delete from "my_table" where "id" = $1 returning "id", "name"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
}

func TestDeleteReturningWildcard(t *testing.T) {
	qb := sqlQb.Delete().
		From("my_table").
		Returning("*")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `delete from "my_table" returning *`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestUpdateReturning(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("name", "alice").
		Where("id", "=", 1).
		Returning("id", "name")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "name" = $1 where "id" = $2 returning "id", "name"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
}

func TestUpdateReturningRaw(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("val", 42).
		Returning(sqlQb.Raw(`coalesce("id", 0)`))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "val" = $1 returning coalesce("id", 0)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
}

func TestInsertReturningSingleColumn(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(1, 2).
		Returning("id")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values ($1, $2) returning "id"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
}

func TestInsertReturningMultipleColumns(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(1, 2).
		Returning("id", "created_at")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values ($1, $2) returning "id", "created_at"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestInsertReturningWildcard(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		Returning("*")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a") values ($1) returning *`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestInsertReturningWithOnConflict(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(1, 2).
		OnConflict("a").DoNothing().
		Returning("id")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values ($1, $2) on conflict ("a") do nothing returning "id"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestInsertReturningMixedWildcard(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		Returning("id", "*")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a") values ($1) returning "id", *`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestInsertReturningPanicNoColumns(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for Returning with no columns, got none")
		}
	}()
	sqlQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		Returning()
}

func TestInsertReturningPanicFeatureDisabled(t *testing.T) {
	noReturnQb := goskul.SetupQueryBuilder(goskul.NewQbSettings())
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when InsertReturningFeature is disabled, got none")
		}
	}()
	noReturnQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		Returning("id")
}

func TestInsertReturningRaw(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		Returning(sqlQb.Raw(`coalesce("id", 0)`))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a") values ($1) returning coalesce("id", 0)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
}

func TestInsertReturningSubquery(t *testing.T) {
	sub := sqlQb.Select().Columns("id").From("other_table")
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		Returning(sub)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a") values ($1) returning (select "id" from "other_table")`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
}
