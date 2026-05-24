package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
)

func TestInsertSingleRow(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b", "c").
		Values(1, 2, 3)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b", "c") values ($1, $2, $3)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != 1 {
		t.Errorf("expected args[0] = 1, got %v", args[0])
	}
	if args[1] != 2 {
		t.Errorf("expected args[1] = 2, got %v", args[1])
	}
	if args[2] != 3 {
		t.Errorf("expected args[2] = 3, got %v", args[2])
	}
}

func TestInsertMultipleRows(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b", "c")

	qb.Values(1, 2, 3)
	qb.Values(4, 5, 6)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b", "c") values ($1, $2, $3), ($4, $5, $6)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 6 {
		t.Fatalf("expected 6 args, got %d: %v", len(args), args)
	}
	if args[0] != 1 || args[3] != 4 {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestInsertNoColumns(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Values("x", "y")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" values ($1, $2)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != "x" {
		t.Errorf("expected args[0] = \"x\", got %v", args[0])
	}
	if args[1] != "y" {
		t.Errorf("expected args[1] = \"y\", got %v", args[1])
	}
}

func TestInsertNoTable(t *testing.T) {
	qb := sqlQb.Insert().
		Columns("a", "b").
		Values(1, 2)

	_, _, err := qb.ToSql()
	if err == nil {
		t.Fatal("expected error for missing table, got nil")
	}
}

func TestInsertNoValues(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b")

	_, _, err := qb.ToSql()
	if err == nil {
		t.Fatal("expected error for missing values, got nil")
	}
}

func TestInsertColumnCountMismatch(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b", "c").
		Values(1, 2)

	_, _, err := qb.ToSql()
	if err == nil {
		t.Fatal("expected error for column/value count mismatch, got nil")
	}
}

func TestInsertClone(t *testing.T) {
	base := sqlQb.Insert().
		Columns("name").
		Values("alice")

	ins1 := base.Clone().Into("table1")
	ins2 := base.Clone().Into("table2")

	sql1, _, err := ins1.ToSql()
	if err != nil {
		t.Fatalf("ins1 error: %v", err)
	}
	sql2, _, err := ins2.ToSql()
	if err != nil {
		t.Fatalf("ins2 error: %v", err)
	}

	want1 := `insert into "table1" ("name") values ($1)`
	want2 := `insert into "table2" ("name") values ($1)`
	if sql1 != want1 {
		t.Errorf("ins1: got %q, want %q", sql1, want1)
	}
	if sql2 != want2 {
		t.Errorf("ins2: got %q, want %q", sql2, want2)
	}

	// base must not have a table set
	_, _, err = base.ToSql()
	if err == nil {
		t.Fatal("base: expected error for missing table, got nil")
	}
}

func TestCteInsert(t *testing.T) {
	cte := sqlQb.Select().Columns("id", "name").From("source_table").Where("active", "=", true)

	qb := sqlQb.With("src", cte).
		Insert().
		Into("my_table").
		Columns("id", "name").
		Values(1, "bob")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `with "src" as (select "id", "name" from "source_table" where "active" = $1) insert into "my_table" ("id", "name") values ($2, $3)`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != true {
		t.Errorf("expected args[0] = true, got %v", args[0])
	}
	if args[1] != 1 {
		t.Errorf("expected args[1] = 1, got %v", args[1])
	}
	if args[2] != "bob" {
		t.Errorf("expected args[2] = \"bob\", got %v", args[2])
	}
}

func TestInsertOnConflictDoNothing(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(1, 2).
		OnConflict("a").DoNothing()

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values ($1, $2) on conflict ("a") do nothing`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
}

func TestInsertOnConflictDoNothingNoTarget(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		OnConflict().DoNothing()

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a") values ($1) on conflict do nothing`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestInsertOnConflictDoUpdateRelation(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(1, 2).
		OnConflict("a").DoUpdate(func(u *goskul.DoUpdateQb) {
		u.Set("a", "excluded.a").Set("b", "excluded.b")
	})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values ($1, $2) on conflict ("a") do update set "a" = "excluded"."a", "b" = "excluded"."b"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	// string values are relations, not bindings — only the VALUES bindings remain
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
}

func TestInsertOnConflictDoUpdateRaw(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a").
		Values(1).
		OnConflict("a").
		DoUpdate(func(u *goskul.DoUpdateQb) {
			u.Set("a", sqlQb.Raw(`coalesce("excluded"."a", "my_table"."a")`))
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a") values ($1) on conflict ("a") do update set "a" = coalesce("excluded"."a", "my_table"."a")`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
}

func TestInsertOnConflictDoUpdateBinding(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(1, 2).
		OnConflict("a").DoUpdate(func(u *goskul.DoUpdateQb) {
		u.Set("b", 42)
	})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values ($1, $2) on conflict ("a") do update set "b" = $3`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[2] != 42 {
		t.Errorf("expected args[2] = 42, got %v", args[2])
	}
}

func TestInsertValuesRawNoBindings(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(1, sqlQb.Raw("NOW()"))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values ($1, NOW())`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != 1 {
		t.Errorf("expected args[0] = 1, got %v", args[0])
	}
}

func TestInsertValuesRawWithBindings(t *testing.T) {
	qb := sqlQb.Insert().
		Into("my_table").
		Columns("a", "b").
		Values(sqlQb.Raw("coalesce(?, ?)", 10, 20), 3)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `insert into "my_table" ("a", "b") values (coalesce($1, $2), $3)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != 10 {
		t.Errorf("expected args[0] = 10, got %v", args[0])
	}
	if args[1] != 20 {
		t.Errorf("expected args[1] = 20, got %v", args[1])
	}
	if args[2] != 3 {
		t.Errorf("expected args[2] = 3, got %v", args[2])
	}
}
