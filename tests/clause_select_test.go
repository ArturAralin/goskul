package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
)

var sqlQb = goskul.SetupQueryBuilder(goskul.PostgreSQLSettings())

func TestSubSelectAsFrom(t *testing.T) {
	inner := sqlQb.Select().
		From("inner_table").
		Where("x", "=", 5).
		Alias("sub")

	qb := sqlQb.Select().
		From(inner).
		Where("y", "=", 10)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from (select * from "inner_table" where "x" = $1) as "sub" where "y" = $2`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}

	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != 5 {
		t.Errorf("expected args[0] = 5, got %v", args[0])
	}
	if args[1] != 10 {
		t.Errorf("expected args[1] = 10, got %v", args[1])
	}
}

func TestColumnSingle(t *testing.T) {
	qb := sqlQb.Select().
		Column("my_col").
		From("my_table")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select "my_col" from "my_table"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestColumnTableDotCol(t *testing.T) {
	qb := sqlQb.Select().
		Column("my_table.my_col").
		From("my_table")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select "my_table"."my_col" from "my_table"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestColumnMultiple(t *testing.T) {
	qb := sqlQb.Select().
		Column("a").
		Column("my_table.b").
		From("my_table")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select "a", "my_table"."b" from "my_table"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestSelectClone(t *testing.T) {
	base := sqlQb.Select().Column("a").From("tbl")

	sel1 := base.Clone().Where("x", "=", 10)
	sel2 := base.Clone().Where("x", "=", 20)

	sql1, args1, err := sel1.ToSql()
	if err != nil {
		t.Fatalf("sel1 error: %v", err)
	}
	sql2, args2, err := sel2.ToSql()
	if err != nil {
		t.Fatalf("sel2 error: %v", err)
	}

	want1 := `select "a" from "tbl" where "x" = $1`
	want2 := `select "a" from "tbl" where "x" = $1`
	if sql1 != want1 {
		t.Errorf("sel1: got %q, want %q", sql1, want1)
	}
	if sql2 != want2 {
		t.Errorf("sel2: got %q, want %q", sql2, want2)
	}
	if args1[0] != 10 {
		t.Errorf("sel1: expected args[0]=10, got %v", args1[0])
	}
	if args2[0] != 20 {
		t.Errorf("sel2: expected args[0]=20, got %v", args2[0])
	}

	// base must be unaffected
	sqlBase, argsBase, err := base.ToSql()
	if err != nil {
		t.Fatalf("base error: %v", err)
	}
	wantBase := `select "a" from "tbl"`
	if sqlBase != wantBase {
		t.Errorf("base: got %q, want %q", sqlBase, wantBase)
	}
	if len(argsBase) != 0 {
		t.Errorf("base: expected no args, got %v", argsBase)
	}
}

func TestCte(t *testing.T) {
	cte1 := sqlQb.Select().Columns("a", "b").From("source_table").Where("active", "=", true)
	cte2 := sqlQb.Select().Column("x").From("other_table")

	qb := sqlQb.With("cte1", cte1).
		With("cte2", cte2).
		Select().
		Columns("cte1.a", "cte2.x").
		From("cte1").
		Where("cte1.b", "=", 42)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `with "cte1" as (select "a", "b" from "source_table" where "active" = $1), "cte2" as (select "x" from "other_table") select "cte1"."a", "cte2"."x" from "cte1" where "cte1"."b" = $2`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}

	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != true {
		t.Errorf("expected args[0] = true, got %v", args[0])
	}
	if args[1] != 42 {
		t.Errorf("expected args[1] = 42, got %v", args[1])
	}
}

func TestColumnRaw(t *testing.T) {
	qb := sqlQb.Select().
		Columns(sqlQb.Raw("count(*) as total"), "name").
		From("my_table")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select count(*) as total, "name" from "my_table"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestColumnRawWithBindings(t *testing.T) {
	qb := sqlQb.Select().
		Column(sqlQb.Raw("coalesce(price, ?)", 0)).
		From("my_table")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select coalesce(price, $1) from "my_table"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != 0 {
		t.Errorf("expected args[0] = 0, got %v", args[0])
	}
}

func TestOrderByString(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		OrderBy("created_at", "desc")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" order by "created_at" desc`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestOrderByRaw(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		OrderBy(sqlQb.Raw("score(?, title)", "golang"), "desc")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" order by score($1, title) desc`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != "golang" {
		t.Errorf("expected args[0] = \"golang\", got %v", args[0])
	}
}

func TestOrderByMultiple(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		OrderBy("priority", "asc").
		OrderBy("created_at", "desc")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" order by "priority" asc, "created_at" desc`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestLimit(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Limit(20)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" limit 20`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestOrderByAndLimit(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Where("active", "=", true).
		OrderBy("created_at", "desc").
		Limit(10)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" where "active" = $1 order by "created_at" desc limit 10`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != true {
		t.Errorf("expected args[0] = true, got %v", args[0])
	}
}

func TestOrderByNullsFirst(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		OrderBy("score", "desc", "first")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" order by "score" desc nulls first`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestOrderByNullsLast(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		OrderBy("score", "asc", "last")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" order by "score" asc nulls last`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestOffset(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Offset(40)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" offset 40`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestLimitAndOffset(t *testing.T) {
	qb := sqlQb.Select().
		From("my_table").
		Limit(10).
		Offset(20)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "my_table" limit 10 offset 20`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestSelectClonePreservesOffset(t *testing.T) {
	base := sqlQb.Select().
		From("tbl").
		Limit(5).
		Offset(15)

	clone := base.Clone().Where("x", "=", 1)

	sql, args, err := clone.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "tbl" where "x" = $1 limit 5 offset 15`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if args[0] != 1 {
		t.Errorf("expected args[0]=1, got %v", args[0])
	}
}

func TestSelectClonePreservesOrderAndLimit(t *testing.T) {
	base := sqlQb.Select().
		From("tbl").
		OrderBy("name", "asc").
		Limit(5)

	clone1 := base.Clone().Where("x", "=", 1)
	clone2 := base.Clone().Where("x", "=", 2)

	sql1, args1, err := clone1.ToSql()
	if err != nil {
		t.Fatalf("clone1 error: %v", err)
	}
	sql2, args2, err := clone2.ToSql()
	if err != nil {
		t.Fatalf("clone2 error: %v", err)
	}

	want := `select * from "tbl" where "x" = $1 order by "name" asc limit 5`
	if sql1 != want {
		t.Errorf("clone1: got %q, want %q", sql1, want)
	}
	if sql2 != want {
		t.Errorf("clone2: got %q, want %q", sql2, want)
	}
	if args1[0] != 1 {
		t.Errorf("clone1: expected args[0]=1, got %v", args1[0])
	}
	if args2[0] != 2 {
		t.Errorf("clone2: expected args[0]=2, got %v", args2[0])
	}
}
