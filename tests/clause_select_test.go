package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
)

var sqlQb = goskul.SetupQueryBuilder(goskul.PostgreSQLSettings())

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
