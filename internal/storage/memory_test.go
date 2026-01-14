package storage

import (
	"errors"
	"testing"
)

func TestMemoryStoreCRUD(t *testing.T) {
	store := NewMemoryStore()

	// Create
	task, err := store.CreateTask(42, "first")
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if task.ID != 1 {
		t.Fatalf("expected id 1 got %d", task.ID)
	}
	if task.Text != "first" {
		t.Fatalf("expected text 'first' got %q", task.Text)
	}

	// Get
	got, ok, err := store.GetTask(42, task.ID)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if !ok {
		t.Fatalf("expected task to exist")
	}
	if got.Text != "first" {
		t.Fatalf("expected text 'first' got %q", got.Text)
	}

	// Update
	got.Text = "updated"
	got.Done = true
	if err := store.UpdateTask(42, got); err != nil {
		t.Fatalf("update task: %v", err)
	}
	updated, ok, err := store.GetTask(42, task.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if !ok || updated.Text != "updated" || !updated.Done {
		t.Fatalf("update not persisted")
	}

	// List
	list, err := store.ListTasks(42)
	if err != nil {
		t.Fatalf("list tasks: %v", err)
	}
	if len(list) != 1 || list[0].ID != task.ID {
		t.Fatalf("expected one task")
	}

	// Delete
	if err := store.DeleteTask(42, task.ID); err != nil {
		t.Fatalf("delete task: %v", err)
	}
	_, ok, err = store.GetTask(42, task.ID)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if ok {
		t.Fatalf("expected task to be deleted")
	}
}

func TestMemoryStoreListUserIDs(t *testing.T) {
	store := NewMemoryStore()
	if _, err := store.CreateTask(2, "a"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := store.CreateTask(1, "b"); err != nil {
		t.Fatalf("create: %v", err)
	}

	ids := store.ListUserIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 user ids")
	}
	if ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("expected sorted ids, got %v", ids)
	}
}

func TestMemoryStoreInvalidID(t *testing.T) {
	store := NewMemoryStore()
	if _, err := store.CreateTask(1, "test"); err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := store.DeleteTask(1, 999); err == nil {
		t.Fatalf("expected error for missing task")
	} else if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}

	task, ok, err := store.GetTask(1, 999)
	if err != nil {
		t.Fatalf("get missing: %v", err)
	}
	if ok || task.ID != 0 {
		t.Fatalf("expected missing task")
	}
}
