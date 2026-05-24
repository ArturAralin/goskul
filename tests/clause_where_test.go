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
