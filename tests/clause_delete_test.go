package tests

import (
	"testing"
)

func TestDeleteFrom(t *testing.T) {
	qb := sqlQb.Delete().
		From("my_table")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `delete from "my_table"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestDeleteWhere(t *testing.T) {
	qb := sqlQb.Delete().
		From("my_table").
		Where("id", "=", 42)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `delete from "my_table" where "id" = $1`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != 42 {
		t.Errorf("expected args[0] = 42, got %v", args[0])
	}
}

func TestDeleteClone(t *testing.T) {
	base := sqlQb.Delete().From("tbl").Where("status", "=", "active")

	del1 := base.Clone().Where("id", "=", 1)
	del2 := base.Clone().Where("id", "=", 2)

	sql1, args1, err := del1.ToSql()
	if err != nil {
		t.Fatalf("del1 error: %v", err)
	}
	sql2, args2, err := del2.ToSql()
	if err != nil {
		t.Fatalf("del2 error: %v", err)
	}

	want1 := `delete from "tbl" where "status" = $1 and "id" = $2`
	want2 := `delete from "tbl" where "status" = $1 and "id" = $2`
	if sql1 != want1 {
		t.Errorf("del1: got %q, want %q", sql1, want1)
	}
	if sql2 != want2 {
		t.Errorf("del2: got %q, want %q", sql2, want2)
	}
	if args1[1] != 1 {
		t.Errorf("del1: expected args[1]=1, got %v", args1[1])
	}
	if args2[1] != 2 {
		t.Errorf("del2: expected args[1]=2, got %v", args2[1])
	}

	// base must have only the original where
	sqlBase, _, err := base.ToSql()
	if err != nil {
		t.Fatalf("base error: %v", err)
	}
	wantBase := `delete from "tbl" where "status" = $1`
	if sqlBase != wantBase {
		t.Errorf("base: got %q, want %q", sqlBase, wantBase)
	}
}

func TestCteDelete(t *testing.T) {
	cte := sqlQb.Select().Columns("id").From("source_table").Where("active", "=", true)

	qb := sqlQb.With("to_delete", cte).
		Delete().
		From("my_table").
		Where("id", "=", 99)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `with "to_delete" as (select "id" from "source_table" where "active" = $1) delete from "my_table" where "id" = $2`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != true {
		t.Errorf("expected args[0] = true, got %v", args[0])
	}
	if args[1] != 99 {
		t.Errorf("expected args[1] = 99, got %v", args[1])
	}
}

func TestDeleteWhereNull(t *testing.T) {
	qb := sqlQb.Delete().
		From("my_table").
		WhereNull("deleted_at")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `delete from "my_table" where "deleted_at" is null`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}
