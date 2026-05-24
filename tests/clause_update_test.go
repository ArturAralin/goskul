package tests

import "testing"

func TestUpdateSet(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("name", "john")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "name" = $1`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != "john" {
		t.Errorf("expected args[0] = \"john\", got %v", args[0])
	}
}

func TestUpdateMultipleSet(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("name", "john").
		Set("age", 30)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "name" = $1, "age" = $2`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != "john" {
		t.Errorf("expected args[0] = \"john\", got %v", args[0])
	}
	if args[1] != 30 {
		t.Errorf("expected args[1] = 30, got %v", args[1])
	}
}

func TestUpdateSetWhere(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("name", "john").
		Where("id", "=", 25)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "name" = $1 where "id" = $2`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != "john" {
		t.Errorf("expected args[0] = \"john\", got %v", args[0])
	}
	if args[1] != 25 {
		t.Errorf("expected args[1] = 25, got %v", args[1])
	}
}

func TestUpdateMultipleSetWhere(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("name", "john").
		Set("age", 30).
		Where("id", "=", 25)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "name" = $1, "age" = $2 where "id" = $3`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != "john" {
		t.Errorf("expected args[0] = \"john\", got %v", args[0])
	}
	if args[1] != 30 {
		t.Errorf("expected args[1] = 30, got %v", args[1])
	}
	if args[2] != 25 {
		t.Errorf("expected args[2] = 25, got %v", args[2])
	}
}

func TestUpdateWhereNull(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("active", false).
		WhereNull("deleted_at")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "active" = $1 where "deleted_at" is null`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d: %v", len(args), args)
	}
	if args[0] != false {
		t.Errorf("expected args[0] = false, got %v", args[0])
	}
}

func TestUpdateClone(t *testing.T) {
	base := sqlQb.Update().Table("tbl").Set("active", true)

	upd1 := base.Clone().Where("id", "=", 1)
	upd2 := base.Clone().Where("id", "=", 2)

	sql1, args1, err := upd1.ToSql()
	if err != nil {
		t.Fatalf("upd1 error: %v", err)
	}
	sql2, args2, err := upd2.ToSql()
	if err != nil {
		t.Fatalf("upd2 error: %v", err)
	}

	want1 := `update "tbl" set "active" = $1 where "id" = $2`
	want2 := `update "tbl" set "active" = $1 where "id" = $2`
	if sql1 != want1 {
		t.Errorf("upd1: got %q, want %q", sql1, want1)
	}
	if sql2 != want2 {
		t.Errorf("upd2: got %q, want %q", sql2, want2)
	}
	if args1[1] != 1 {
		t.Errorf("upd1: expected args[1]=1, got %v", args1[1])
	}
	if args2[1] != 2 {
		t.Errorf("upd2: expected args[1]=2, got %v", args2[1])
	}

	// base must not have a where clause
	sqlBase, argsBase, err := base.ToSql()
	if err != nil {
		t.Fatalf("base error: %v", err)
	}
	wantBase := `update "tbl" set "active" = $1`
	if sqlBase != wantBase {
		t.Errorf("base: got %q, want %q", sqlBase, wantBase)
	}
	if len(argsBase) != 1 {
		t.Errorf("base: expected 1 arg, got %v", argsBase)
	}
}

func TestCteUpdate(t *testing.T) {
	cte := sqlQb.Select().Columns("id").From("source_table").Where("active", "=", true)

	qb := sqlQb.With("to_update", cte).
		Update().
		Table("my_table").
		Set("status", "archived").
		Where("id", "=", 99)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `with "to_update" as (select "id" from "source_table" where "active" = $1) update "my_table" set "status" = $2 where "id" = $3`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != true {
		t.Errorf("expected args[0] = true, got %v", args[0])
	}
	if args[1] != "archived" {
		t.Errorf("expected args[1] = \"archived\", got %v", args[1])
	}
	if args[2] != 99 {
		t.Errorf("expected args[2] = 99, got %v", args[2])
	}
}

func TestUpdateSetRaw(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table").
		Set("views", sqlQb.Raw("views + ?", 1)).
		Where("id", "=", 5)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `update "my_table" set "views" = views + $1 where "id" = $2`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != 1 {
		t.Errorf("expected args[0] = 1, got %v", args[0])
	}
	if args[1] != 5 {
		t.Errorf("expected args[1] = 5, got %v", args[1])
	}
}

func TestUpdateNoTable(t *testing.T) {
	qb := sqlQb.Update().
		Set("name", "john")

	_, _, err := qb.ToSql()
	if err == nil {
		t.Fatal("expected error for missing table, got nil")
	}
}

func TestUpdateNoSet(t *testing.T) {
	qb := sqlQb.Update().
		Table("my_table")

	_, _, err := qb.ToSql()
	if err == nil {
		t.Fatal("expected error for missing set clauses, got nil")
	}
}
