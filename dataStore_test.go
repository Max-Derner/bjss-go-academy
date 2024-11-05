package main
import (
	"testing"
	"reflect"
)

func TestDataStore(t *testing.T) {
	t.Run("Add value to store", func(t *testing.T) {
		dataKey := "Keep sanity"
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)

		want := mockDataStore([]ToDoItem{item})
		
		store := NewDataStore()
		store.Add(item)

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Add existing value to store", func(t *testing.T) {
		dataKey := "Keep sanity"
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)

		want := mockDataStore([]ToDoItem{item})
		
		store := NewDataStore()
		store.Add(item)
		err := store.Add(item)

		if err != ErrCannotAdd {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Update value in store", func(t *testing.T) {
		dataKey := "Keep sanity"
		initialItem := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)
		updateItem := ConstructToDoItem(
			Title(dataKey),
			"high",
			true,
		)

		want := mockDataStore([]ToDoItem{updateItem})
		
		store := NewDataStore()
		store.Add(initialItem)
		store.Update(updateItem)

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Update non-existent item in store", func(t *testing.T) {
		dataKey := "Keep sanity"
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			true,
		)

		want := mockDataStore([]ToDoItem{})
		
		store := NewDataStore()
		err := store.Update(item)

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}

func mockDataStore(items []ToDoItem) DataStore {
	dataMap := make(map[Title]ToDoItem)
	for _, item := range(items) {
		dataKey := item.Title
		dataMap[Title(dataKey)] = item
	}
	return DataStore{
		dataMap,
	}
}