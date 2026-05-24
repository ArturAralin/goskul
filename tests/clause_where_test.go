package tests

import "testing"

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
