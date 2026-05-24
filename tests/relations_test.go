package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
)

func TestFromSimpleName(t *testing.T) {
	qb := sqlQb.Select().From("orders")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestFromLeadingSpaces(t *testing.T) {
	qb := sqlQb.Select().From("  orders")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestFromWithAlias(t *testing.T) {
	qb := sqlQb.Select().From("orders as o")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" as "o"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestFromWithAliasExtraSpaces(t *testing.T) {
	qb := sqlQb.Select().From(" orders     as o ")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" as "o"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestFromWithLongAlias(t *testing.T) {
	qb := sqlQb.Select().From("orders as alias_name")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" as "alias_name"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestJoinWithAlias(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users as u", func(j *goskul.JoinClause) {
			j.On("u.id", "=", "orders.user_id")
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" as "u" on "u"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestFromEmptyString(t *testing.T) {
	qb := sqlQb.Select().From("")

	_, _, err := qb.ToSql()
	if err == nil {
		t.Fatalf("expected error for empty relation name, got none")
	}
}
