package rotation

import (
	"testing"
)

func (m *mockKeyStore) assertNoKeysCreated(t *testing.T) {
	if m.createdKey {
		t.Error("Expected key not to be created, but did")
	}
}

func (m *mockKeyStore) assertCreatedKey(t *testing.T) {
	if !m.createdKey {
		t.Error("Expected key not to be created, but did")
	}
}

func (m *mockKeyStore) assertKeyCountDeleted(t *testing.T, count int) {
	if len(m.deletedKeys) > count {
		t.Errorf("Expected %d keys to be deleted, got %d", count, len(m.deletedKeys))
	}
}

func assertEmptyKeyList(t *testing.T, k KeyList) {
	assertKeyListSize(t, k, 0)
}
func assertKeyListSize(t *testing.T, k KeyList, size int) {
	if len(k) != size {
		t.Errorf("Expected keylist to have %d, got %d", size, len(k))
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
}

func TestApplyNoOpPlan(t *testing.T) {
	ctx, done := testContext(t)
	defer done()

	store := newMock()
	plan := KeyRotationPlan{
		CreateKey:   false,
		DestroyKeys: nil,
	}
	keys, err := plan.Apply(ctx, store)
	assertEmptyKeyList(t, keys)
	assertNoError(t, err)

	store.assertNoKeysCreated(t)
	store.assertKeyCountDeleted(t, 0)
}

func TestApplyCreateOnly(t *testing.T) {
	ctx, done := testContext(t)
	defer done()

	store := newMock()
	plan := KeyRotationPlan{
		CreateKey:   true,
		DestroyKeys: nil,
	}
	keys, err := plan.Apply(ctx, store)
	assertNoError(t, err)
	assertKeyListSize(t, keys, 1)

	store.assertCreatedKey(t)
	store.assertKeyCountDeleted(t, 0)
}

func TestDestroysTargetKeys(t *testing.T) {
	ctx, done := testContext(t)
	defer done()

	store := newMock()
	keyToDelete := &mockKey{created: InvalidTime()}
	plan := KeyRotationPlan{
		CreateKey:   false,
		DestroyKeys: KeyList{keyToDelete},
	}
	keys, err := plan.Apply(ctx, store)
	assertNoError(t, err)
	assertEmptyKeyList(t, keys)

	store.assertNoKeysCreated(t)
	store.assertKeyCountDeleted(t, 1)
}

func TestCreateAppendsKey(t *testing.T) {
	ctx, done := testContext(t)
	defer done()

	store := newMock()
	existingKey := &mockKey{created: InvalidTime()}
	plan := KeyRotationPlan{
		CreateKey: true,
		goodKeys:  KeyList{existingKey},
	}
	keys, err := plan.Apply(ctx, store)
	assertNoError(t, err)
	assertKeyListSize(t, keys, 2)

	store.assertCreatedKey(t)
	store.assertKeyCountDeleted(t, 0)
}
