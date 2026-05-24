package tests

import (
	"testing"

	goskul "github.com/ArturAralin/goskul"
	"github.com/ArturAralin/goskul/internal"
)

func TestLockForUpdate(t *testing.T) {
	sql, _, err := sqlQb.Select().From("t").LockForUpdate().ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t" for update`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockForShare(t *testing.T) {
	sql, _, err := sqlQb.Select().From("t").LockForShare().ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t" for share`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockForNoKeyUpdate(t *testing.T) {
	sql, _, err := sqlQb.Select().From("t").LockForNoKeyUpdate().ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t" for no key update`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockForUpdateSkipLocked(t *testing.T) {
	sql, _, err := sqlQb.Select().From("t").LockForUpdate().LockSkipLocked().ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t" for update skip locked`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockForUpdateNoWait(t *testing.T) {
	sql, _, err := sqlQb.Select().From("t").LockForUpdate().LockNoWait().ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t" for update nowait`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockTypeLastWins(t *testing.T) {
	sql, _, err := sqlQb.Select().From("t").LockForUpdate().LockForNoKeyUpdate().ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t" for no key update`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockIsAfterOffset(t *testing.T) {
	sql, _, err := sqlQb.Select().
		From("t").
		Where("id", ">", 0).
		OrderBy("id", "asc").
		Limit(10).
		Offset(5).
		LockForUpdate().
		ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t" where "id" > $1 order by "id" asc limit 10 offset 5 for update`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockClonePreservesLock(t *testing.T) {
	base := sqlQb.Select().From("t").LockForUpdate().LockSkipLocked()
	clone := base.Clone()

	baseSql, _, err := base.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cloneSql, _, err := clone.ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `select * from "t" for update skip locked`
	if baseSql != want {
		t.Errorf("base: got %q, want %q", baseSql, want)
	}
	if cloneSql != want {
		t.Errorf("clone: got %q, want %q", cloneSql, want)
	}
}

func TestLockModifierWithoutLockTypeNotEmitted(t *testing.T) {
	sql, _, err := sqlQb.Select().From("t").LockSkipLocked().ToSql()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `select * from "t"`
	if sql != want {
		t.Errorf("got %q, want %q", sql, want)
	}
}

func TestLockForUpdateFeatureDisabledPanics(t *testing.T) {
	qb := goskul.SetupQueryBuilder(internal.NewQbSettings())
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for LockForUpdateFeature disabled")
		}
	}()
	qb.Select().LockForUpdate()
}

func TestLockForShareFeatureDisabledPanics(t *testing.T) {
	qb := goskul.SetupQueryBuilder(internal.NewQbSettings())
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for LockForShareFeature disabled")
		}
	}()
	qb.Select().LockForShare()
}

func TestLockForNoKeyUpdateFeatureDisabledPanics(t *testing.T) {
	qb := goskul.SetupQueryBuilder(internal.NewQbSettings())
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for LockForNoKeyUpdateFeature disabled")
		}
	}()
	qb.Select().LockForNoKeyUpdate()
}

func TestLockSkipLockedFeatureDisabledPanics(t *testing.T) {
	qb := goskul.SetupQueryBuilder(internal.NewQbSettings())
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for LockSkipLockedFeature disabled")
		}
	}()
	qb.Select().LockSkipLocked()
}

func TestLockNoWaitFeatureDisabledPanics(t *testing.T) {
	qb := goskul.SetupQueryBuilder(internal.NewQbSettings())
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for LockNoWaitFeature disabled")
		}
	}()
	qb.Select().LockNoWait()
}
