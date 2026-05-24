package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
)

func TestWhereSimpleCondition(t *testing.T) {
	qb := sqlQb.Select().Where("col", "=", "my str")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * where "col" = $1`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 1 || args[0] != "my str" {
		t.Errorf("expected args [\"my str\"], got %v", args)
	}
}

func TestWhereRelationCondition(t *testing.T) {
	qb := sqlQb.Select().Where("col", "=", goskul.Rel("another.col"))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * where "col" = "another"."col"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestWhereInSubquery(t *testing.T) {
	sub := sqlQb.Select().Column("id").From("t")
	qb := sqlQb.Select().WhereIn("x", sub)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * where "x" in (select "id" from "t")`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestWhereEqual(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where("a", "=", 10)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "a" = $1`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}

	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != 10 {
		t.Errorf("expected args[0] = 10, got %v", args[0])
	}
}

func TestWhereNullAndWhereNull(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		WhereNull("my_col_1").
		WhereNull("my_col_2")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "my_col_1" is null and "my_col_2" is null`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestWhereNullOrWhereNull(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		WhereNull("my_col_1").
		OrWhereNull("my_col_2")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "my_col_1" is null or "my_col_2" is null`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestWhereRawNotEqual(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where(sqlQb.Raw("custom + ?", 10), "<>", sqlQb.Raw("custom2 - ? + ?", 20, 30))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where custom + $1 <> custom2 - $2 + $3`
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
	if args[2] != 30 {
		t.Errorf("expected args[2] = 30, got %v", args[2])
	}
}

func TestWhereEqualOrWhereNotEqual(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where("a", "=", 10).
		OrWhere("b", "<>", 30)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "a" = $1 or "b" <> $2`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}

	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != 10 {
		t.Errorf("expected args[0] = 10, got %v", args[0])
	}
	if args[1] != 30 {
		t.Errorf("expected args[1] = 30, got %v", args[1])
	}
}

func TestWhereIn(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		WhereIn("a", []int{1, 2, 3})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "a" in ($1, $2, $3)`
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

func TestWhereNotIn(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		WhereNotIn("a", []int{1, 2, 3})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "a" not in ($1, $2, $3)`
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

func TestOrWhereIn(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where("b", "=", 10).
		OrWhereIn("a", []int{1, 2, 3})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "b" = $1 or "a" in ($2, $3, $4)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}

	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d: %v", len(args), args)
	}
	if args[0] != 10 {
		t.Errorf("expected args[0] = 10, got %v", args[0])
	}
	if args[1] != 1 {
		t.Errorf("expected args[1] = 1, got %v", args[1])
	}
	if args[2] != 2 {
		t.Errorf("expected args[2] = 2, got %v", args[2])
	}
	if args[3] != 3 {
		t.Errorf("expected args[3] = 3, got %v", args[3])
	}
}

func TestOrWhereNotIn(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where("b", "=", 10).
		OrWhereNotIn("a", []int{1, 2, 3})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "b" = $1 or "a" not in ($2, $3, $4)`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}

	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d: %v", len(args), args)
	}
	if args[0] != 10 {
		t.Errorf("expected args[0] = 10, got %v", args[0])
	}
	if args[1] != 1 {
		t.Errorf("expected args[1] = 1, got %v", args[1])
	}
	if args[2] != 2 {
		t.Errorf("expected args[2] = 2, got %v", args[2])
	}
	if args[3] != 3 {
		t.Errorf("expected args[3] = 3, got %v", args[3])
	}
}

func TestWhereRaw(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		WhereRaw(sqlQb.Raw("lower(name) = ?", "alice"))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where lower(name) = $1`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != "alice" {
		t.Errorf("expected args[0] = \"alice\", got %v", args[0])
	}
}

func TestWhereRawCombinedWithWhere(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where("status", "=", 1).
		WhereRaw(sqlQb.Raw("score(?) > 0.5", "term"))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "status" = $1 and score($2) > 0.5`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != 1 {
		t.Errorf("expected args[0] = 1, got %v", args[0])
	}
	if args[1] != "term" {
		t.Errorf("expected args[1] = \"term\", got %v", args[1])
	}
}

func TestOrWhereRaw(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where("a", "=", 1).
		OrWhereRaw(sqlQb.Raw("b > ?", 100))

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "a" = $1 or b > $2`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != 1 {
		t.Errorf("expected args[0] = 1, got %v", args[0])
	}
	if args[1] != 100 {
		t.Errorf("expected args[1] = 100, got %v", args[1])
	}
}
