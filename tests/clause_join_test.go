package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
)

func TestInnerJoin(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				AndOn("users.active", "=", true)
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" and "users"."active" = $1`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 1 || args[0] != true {
		t.Errorf("expected args=[true], got %v", args)
	}
}

func TestInnerJoinOn(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoinOn("users", "users.id", "orders.user_id")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestLeftJoin(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		LeftJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id")
		})

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" left join "users" on "users"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestLeftJoinOn(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		LeftJoinOn("users", "users.id", "orders.user_id")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" left join "users" on "users"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestRightJoin(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		RightJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id")
		})

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" right join "users" on "users"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestRightJoinOn(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		RightJoinOn("users", "users.id", "orders.user_id")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" right join "users" on "users"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestFullJoin(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		FullJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id")
		})

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" full join "users" on "users"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestFullJoinOn(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		FullJoinOn("users", "users.id", "orders.user_id")

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" full join "users" on "users"."id" = "orders"."user_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestCrossJoin(t *testing.T) {
	qb := sqlQb.Select().
		From("products").
		CrossJoin("colors")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "products" cross join "colors"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestJoinOrOn(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrOn("users.guest_id", "=", "orders.guest_id")
		})

	sql, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" or "users"."guest_id" = "orders"."guest_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
}

func TestMultipleJoins(t *testing.T) {
	qb := sqlQb.Select().
		Column("orders.id").
		Column("users.name").
		Column("statuses.label").
		From("orders").
		InnerJoinOn("users", "users.id", "orders.user_id").
		LeftJoinOn("statuses", "statuses.id", "orders.status_id")

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select "orders"."id", "users"."name", "statuses"."label" from "orders" inner join "users" on "users"."id" = "orders"."user_id" left join "statuses" on "statuses"."id" = "orders"."status_id"`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestJoinOnNull(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		LeftJoin("users", func(j *goskul.JoinClause) {
			j.OnNull("orders.deleted_at")
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" left join "users" on "orders"."deleted_at" is null`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestJoinAndOnNull(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				AndOnNull("orders.deleted_at")
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" and "orders"."deleted_at" is null`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestJoinOrOnNull(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrOnNull("orders.deleted_at")
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" or "orders"."deleted_at" is null`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestJoinOnRaw(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.OnRaw(sqlQb.Raw("users.id = orders.user_id and users.active = ?", true))
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on users.id = orders.user_id and users.active = $1`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 1 || args[0] != true {
		t.Errorf("expected args=[true], got %v", args)
	}
}

func TestJoinOrOnRaw(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrOnRaw(sqlQb.Raw("users.guest_id = ?", 42))
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" or users.guest_id = $1`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 1 || args[0] != 42 {
		t.Errorf("expected args=[42], got %v", args)
	}
}

func TestJoinOnSubCond(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OnSubCond(func(s *goskul.SubCond) {
					s.Where("users.active", "=", true)
				})
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// sub-cond has 1 inner condition → no parens
	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" and "users"."active" = $1`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 1 || args[0] != true {
		t.Errorf("expected args=[true], got %v", args)
	}
}

func TestJoinAndOnSubCond(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				AndOnSubCond(func(s *goskul.SubCond) {
					s.Where("users.active", "=", true).OrWhere("users.role", "=", "admin")
				})
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" and ("users"."active" = $1 or "users"."role" = $2)`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 2 || args[0] != true || args[1] != "admin" {
		t.Errorf("expected args=[true, admin], got %v", args)
	}
}

func TestJoinOrOnSubCond(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrOnSubCond(func(s *goskul.SubCond) {
					s.Where("users.type", "=", "guest").
						Where("users.active", "=", true)
				})
		})

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" or ("users"."type" = $1 and "users"."active" = $2)`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 2 || args[0] != "guest" || args[1] != true {
		t.Errorf("expected args=[guest, true], got %v", args)
	}
}

func TestJoinWithWhere(t *testing.T) {
	qb := sqlQb.Select().
		From("orders").
		InnerJoin("users", func(j *goskul.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				AndOn("users.role_id", "=", 5)
		}).
		Where("orders.total", ">", 100)

	sql, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "orders" inner join "users" on "users"."id" = "orders"."user_id" and "users"."role_id" = $1 where "orders"."total" > $2`
	if sql != want {
		t.Errorf("got  %q\nwant %q", sql, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != 5 {
		t.Errorf("expected args[0]=5, got %v", args[0])
	}
	if args[1] != 100 {
		t.Errorf("expected args[1]=100, got %v", args[1])
	}
}
