# Goskul 💀

SQL query builder for PostgreSQL. Produces parameterized queries with `$N` placeholders and double-quoted identifiers.

## Setup

```go
qb := goskul.SetupQueryBuilder(goskul.PostgreSQLSettings())
```

Create one instance and reuse it across the application.

## SELECT

```go
sql, args, err := qb.Select().
    From("users").
    Where("active", "=", true).
    ToSql()
// select * from "users" where "active" = $1
// args: [true]
```

### Columns

```go
qb.Select().Column("id").Column("users.name").From("users")
// select "id", "users"."name" from "users"

qb.Select().Columns("id", "name", "email").From("users")
// select "id", "name", "email" from "users"
```

Without `.Column`/`.Columns` the builder emits `select *`.

### Subquery as FROM

```go
inner := qb.Select().From("raw_events").Where("type", "=", "click").Alias("ev")

qb.Select().Columns("ev.id", "ev.ts").From(inner)
// select "ev"."id", "ev"."ts" from (select * from "raw_events" where "type" = $1) as "ev"
```

`.Alias` is required when using a subquery as a FROM source.

## UPDATE

```go
sql, args, err := qb.Update().
    Table("users").
    Set("name", "Alice").
    Set("age", 30).
    Where("id", "=", 1).
    ToSql()
// update "users" set "name" = $1, "age" = $2 where "id" = $3
// args: ["Alice", 30, 1]
```

Both `.Table` and at least one `.Set` are required; `ToSql` returns an error otherwise.

## DELETE

```go
sql, args, err := qb.Delete().
    From("sessions").
    Where("expires_at", "<", now).
    ToSql()
// delete from "sessions" where "expires_at" < $1
```

## INSERT

```go
sql, args, err := qb.Insert().
    Into("users").
    Columns("name", "email").
    Values("Alice", "alice@example.com").
    ToSql()
// insert into "users" ("name", "email") values ($1, $2)
// args: ["Alice", "alice@example.com"]
```

Both `.Into` and at least one `.Values` are required; `ToSql` returns an error otherwise. When `.Columns` is provided, every `.Values` call must supply the same number of values.

### Multiple rows

```go
qb.Insert().Into("tags").Columns("name").
    Values("go").
    Values("sql")
// insert into "tags" ("name") values ($1), ($2)
```

### Without column list

```go
qb.Insert().Into("tags").Values("go", 1)
// insert into "tags" values ($1, $2)
```

### ON CONFLICT

`.OnConflict(cols...)` returns an intermediate builder. Chain `.DoNothing()` or `.DoUpdate(fn)` on it.

```go
// Do nothing on conflict
qb.Insert().Into("users").Columns("id", "email").Values(1, "a@b.com").
    OnConflict("id").DoNothing()
// insert into "users" ("id", "email") values ($1, $2) on conflict ("id") do nothing

// Without a conflict target
qb.Insert().Into("users").Columns("id").Values(1).
    OnConflict().DoNothing()
// insert into "users" ("id") values ($1) on conflict do nothing
```

`.DoUpdate` receives a callback with a `*DoUpdateQb`. Inside the callback, `.Set(col, val)` builds the `SET` list:

| Value type | Result |
|---|---|
| `string` | relation identifier — `"excluded.col"` → `"excluded"."col"` |
| `*QbRaw` | raw SQL fragment inlined as-is |
| any other | bound parameter `$N` |

```go
// Relation reference (most common — mirror the excluded row)
qb.Insert().Into("users").Columns("id", "score").Values(1, 10).
    OnConflict("id").DoUpdate(func(u *goskul.DoUpdateQb) {
        u.Set("score", "excluded.score")
    })
// ... on conflict ("id") do update set "score" = "excluded"."score"

// Raw expression
qb.Insert().Into("users").Columns("id", "score").Values(1, 10).
    OnConflict("id").DoUpdate(func(u *goskul.DoUpdateQb) {
        u.Set("score", qb.Raw(`coalesce("excluded"."score", "users"."score")`))
    })
// ... on conflict ("id") do update set "score" = coalesce("excluded"."score", "users"."score")

// Bound value
qb.Insert().Into("users").Columns("id", "score").Values(1, 10).
    OnConflict("id").DoUpdate(func(u *goskul.DoUpdateQb) {
        u.Set("score", 0)
    })
// ... on conflict ("id") do update set "score" = $3
```

### Subquery as value

Any `*SelectQb` can be passed as a value in `Where`, `WhereIn`, etc. — it is rendered as `(select ...)` inline.

```go
sub := qb.Select().Column("user_id").From("banned_users")

qb.Delete().From("sessions").WhereIn("user_id", sub)
// delete from "sessions" where "user_id" in (select "user_id" from "banned_users")

qb.Select().From("orders").Where("user_id", "=", sub)
// select * from "orders" where "user_id" = (select "user_id" from "banned_users")
```

## WHERE conditions

All three query types share the same WHERE API.

| Method | SQL |
|---|---|
| `.Where(col, op, val)` | `col op $N` (AND) |
| `.OrWhere(col, op, val)` | `col op $N` (OR) |
| `.WhereNull(col)` | `col is null` (AND) |
| `.OrWhereNull(col)` | `col is null` (OR) |
| `.WhereNotNull(col)` | `col is not null` (AND) |
| `.OrWhereNotNull(col)` | `col is not null` (OR) |
| `.WhereIn(col, values)` | `col in ($1, $2, ...)` or `col in (subquery)` (AND) |
| `.WhereNotIn(col, values)` | `col not in ($1, $2, ...)` or `col not in (subquery)` (AND) |
| `.OrWhereIn(col, values)` | `col in ($1, $2, ...)` or `col in (subquery)` (OR) |
| `.OrWhereNotIn(col, values)` | `col not in ($1, $2, ...)` or `col not in (subquery)` (OR) |
| `.WhereRaw(raw)` | arbitrary raw fragment (AND) |
| `.OrWhereRaw(raw)` | arbitrary raw fragment (OR) |

`col` and `val` accept plain strings (treated as identifiers), `*QbRaw` (inlined as-is), or any value (bound as a parameter).

```go
// WhereIn with a slice
qb.Select().From("orders").WhereIn("status_id", []int{1, 2, 3})
// select * from "orders" where "status_id" in ($1, $2, $3)

// WhereIn with a subquery
sub := qb.Select().Column("id").From("active_users")
qb.Select().From("orders").WhereIn("user_id", sub)
// select * from "orders" where "user_id" in (select "id" from "active_users")

// Combined AND / OR
qb.Select().From("t").Where("a", "=", 1).OrWhere("b", ">", 2)
// select * from "t" where "a" = $1 or "b" > $2

// Raw fragment in WHERE
qb.Select().From("users").
    Where("status", "=", 1).
    WhereRaw(qb.Raw("score(?) > 0.5", "term"))
// select * from "users" where "status" = $1 and score($2) > 0.5
```

## Joins

Every join type has two variants: a callback form for complex `ON` conditions and a shorthand `On(left, right)` form.

```go
// Shorthand: single equality condition
qb.Select().From("orders").InnerJoinOn("users", "users.id", "orders.user_id")
// ... inner join "users" on "users"."id" = "orders"."user_id"

// Callback: multiple / mixed conditions
qb.Select().From("orders").
    InnerJoin("users", func(j *goskul.JoinClause) {
        j.On("users.id", "=", "orders.user_id").
            AndOn("users.active", "=", true)
    })
// ... inner join "users" on "users"."id" = "orders"."user_id" and "users"."active" = $1
```

Available join types:

| Method | SQL |
|---|---|
| `InnerJoin` / `InnerJoinOn` | `inner join` |
| `LeftJoin` / `LeftJoinOn` | `left join` |
| `RightJoin` / `RightJoinOn` | `right join` |
| `FullJoin` / `FullJoinOn` | `full join` |
| `CrossJoin(tbl)` | `cross join` (no ON clause) |

`JoinClause` methods: `.On`, `.AndOn`, `.OrOn` — all accept `(left, op, right)`.

## ORDER BY

`.OrderBy(term, direction, nulls...)` appends an ordering term. Call it multiple times for multiple terms.

```go
// Single column
qb.Select().From("posts").OrderBy("created_at", "desc")
// select * from "posts" order by "created_at" desc

// Multiple columns
qb.Select().From("tasks").
    OrderBy("priority", "asc").
    OrderBy("created_at", "desc")
// select * from "tasks" order by "priority" asc, "created_at" desc

// Raw expression
qb.Select().From("docs").
    OrderBy(qb.Raw("score(?, title)", "golang"), "desc")
// select * from "docs" order by score($1, title) desc

// NULLS FIRST / NULLS LAST (optional third argument)
qb.Select().From("scores").OrderBy("score", "desc", "first")
// select * from "scores" order by "score" desc nulls first

qb.Select().From("scores").OrderBy("score", "asc", "last")
// select * from "scores" order by "score" asc nulls last
```

## LIMIT

```go
qb.Select().From("events").
    Where("active", "=", true).
    OrderBy("created_at", "desc").
    Limit(10)
// select * from "events" where "active" = $1 order by "created_at" desc limit 10
```

`Limit` takes a `uint64`. `ORDER BY` and `LIMIT` are also preserved through `.Clone()`.

## CTE (WITH)

```go
active := qb.Select().Columns("id").From("source_table").Where("active", "=", true)
other  := qb.Select().Column("x").From("other_table")

sql, args, err := qb.With("cte1", active).
    With("cte2", other).
    Select().
    Columns("cte1.id", "cte2.x").
    From("cte1").
    ToSql()
// with "cte1" as (...), "cte2" as (...) select "cte1"."id", "cte2"."x" from "cte1"
```

`.With` is also available on `CteQb` to chain multiple CTEs. The CTE block can terminate with `.Select()`, `.Update()`, `.Delete()`, or `.Insert()`.

## Raw expressions

Use `.Raw(sql, args...)` when you need to embed an arbitrary SQL fragment. Placeholders inside the raw string use `?` and are re-indexed automatically.

```go
// In a column list
qb.Select().
    Column(qb.Raw("count(*) as total")).
    From("events")
// select count(*) as total from "events"

// With bindings
qb.Select().
    Column(qb.Raw("coalesce(price, ?)", 0)).
    From("products")
// select coalesce(price, $1) from "products"

// In a WHERE clause
qb.Select().From("t").
    Where(qb.Raw("lower(email)"), "=", qb.Raw("lower(?)", email))
// ... where lower(email) = lower($1)
```

## Clone

All three builder types support `.Clone()`, which produces a deep copy that shares no state with the original. Useful for building query variants from a common base.

```go
base := qb.Select().Column("a").From("tbl")

q1 := base.Clone().Where("x", "=", 10)
q2 := base.Clone().Where("x", "=", 20)
// base is unchanged; q1 and q2 are independent
```

## ToSql

Every builder terminates with:

```go
sql  string
args []interface{}
err  error
```

Pass `sql` and `args` directly to `pgx` / `database/sql`:

```go
rows, err := db.Query(ctx, sql, args...)
```
