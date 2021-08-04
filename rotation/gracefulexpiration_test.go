package rotation

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"testing"
	"time"
)

func TestNoKeysCreates(t *testing.T) {
	harnessRunPlan(t, newEmptyMock, assertCreateKeyOnly)
}

func TestCurrentKeyDoesNothing(t *testing.T) {
	harnessRunPlan(t, newMockWithKey, assertPlanNoOp)
}

func TestGracePeriodKeyCreates(t *testing.T) {
	harnessRunPlan(t, newMockInGrace, func(t *testing.T, plan *KeyRotationPlan) {
		plan.assertCreating(t)
		plan.assertNotDestroying(t)
		plan.assertNoGoodKeys(t)
	})
}

func TestExpiredKeyCreatesNewDestroysOld(t *testing.T) {
	harnessRunPlan(t, newMockWithExpiredKey, func(t *testing.T, plan *KeyRotationPlan) {
		plan.assertCreating(t)
		plan.assertDestroying(t, 1)
		plan.assertNoGoodKeys(t)
	})
}

func TestDestroyMultipleExpiredKeys(t *testing.T) {
	harnessRunPlan(t, newMockWithMultipleExpiredKey, func(t *testing.T, plan *KeyRotationPlan) {
		plan.assertCreating(t)
		plan.assertDestroying(t, 2)
		plan.assertNoGoodKeys(t)
	})
}

func TestMaximumInGrace(t *testing.T) {
	harnessRunPlan(t, func() KeyStore {
		store := newMock()
		store.maximumCount = 2
		store.mockInGrace()
		store.mockInGrace()
		return store
	}, func(t *testing.T, plan *KeyRotationPlan) {
		plan.assertCreating(t)
		plan.assertDestroying(t, 1)
		plan.assertNoGoodKeys(t)
	})
}

func testContext(t *testing.T) (context.Context, func()) {
	//todo: include background
	return context.WithCancel(context.Background())
}

func assert(t *testing.T, depth int, condition bool, format string, args ...interface{}) {
	if !condition {
		message := fmt.Sprintf(format, args...)
		_, file, line, _ := runtime.Caller(depth)
		_, err := fmt.Printf("\t%s:%d -- FAILED -- %s\n", path.Base(file), line, message)
		if err != nil {
			panic(err)
		}
		t.Fail()
	}
}

func (plan *KeyRotationPlan) assertCreating(t *testing.T) {
	assert(t, 2, plan.CreateKey, "Expected plan to create key, did not")
}

func (plan *KeyRotationPlan) assertNotCreating(t *testing.T) {
	assert(t, 2, !plan.CreateKey, "Expected plan not create key, creating")
}

func (plan *KeyRotationPlan) assertNoGoodKeys(t *testing.T) {
	count := len(plan.goodKeys)
	assert(t, 2, count == 0, "Expected 0 good keys, got %d", count)
}

func (plan *KeyRotationPlan) assertGoodKeysCount(t *testing.T, expected int) {
	count := len(plan.goodKeys)
	assert(t, 2, count != expected, "Expected %d good keys, got %d", expected, count)
}

func (plan *KeyRotationPlan) assertNotDestroying(t *testing.T) {
	count := len(plan.DestroyKeys)
	if count > 0 {
		t.Errorf("Expected 0 keys to be destroyed, got %d", count)
	}
}

func (plan *KeyRotationPlan) assertDestroying(t *testing.T, expected int) {
	count := len(plan.DestroyKeys)
	assert(t, 2, count == expected, "Expecting %d keys to be destroyed, got %d", expected, count)
}

func assertCreateKeyOnly(t *testing.T, plan *KeyRotationPlan) {
	plan.assertCreating(t)
	plan.assertNoGoodKeys(t)
}

func assertPlanNoOp(t *testing.T, plan *KeyRotationPlan) {
	if plan.CreateKey {
		t.Error("Expected to do nothing, creating key instead")
	}

	destroyCount := len(plan.DestroyKeys)
	if destroyCount > 0 {
		t.Errorf("Expected to delete nothing, instead attempting to deleted %d", destroyCount)
	}

	goodKeyCount := len(plan.goodKeys)
	if goodKeyCount == 0 {
		t.Error("Expected good keys, got none")
	}
	for i, k := range plan.goodKeys {
		if k == nil {
			t.Errorf("goodKeys[%d] is nil", i)
		}
	}
}

func harnessRunPlan(t *testing.T, given func() KeyStore, expect func(*testing.T, *KeyRotationPlan)) {
	ctx, done := testContext(t)
	defer done()

	rotation := &GracefulExpiration{
		maximumAge: 1 * time.Minute,
		graceAge:   30 * time.Second,
	}

	store := given()
	plan, err := rotation.Plan(ctx, store)
	if err != nil {
		t.Fatalf("Failed planning because %s", err.Error())
		return
	}
	expect(t, plan)
}
