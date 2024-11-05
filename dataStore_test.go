package main
import (
	"testing"
	"reflect"
	// "sort"
)

func TestDataStore(t *testing.T) {
	t.Run("Add value to store", func(t *testing.T) {
		dataKey := "Keep sanity"
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)

		want := testDataStore([]ToDoItem{item})
		
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

		want := testDataStore([]ToDoItem{item})
		
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

		want := testDataStore([]ToDoItem{updateItem})
		
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

		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		err := store.Update(item)

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Deleting item", func(t *testing.T) {
		dataKey := Title("Keep sanity")
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		store.Add(item)
		store.Delete(dataKey)

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Deleting non-existent item", func(t *testing.T) {
		dataKey := Title("Keep sanity")
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		err := store.Delete(dataKey)

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Reading datastore", func(t *testing.T) {
		titles := []string{
			"keep sanity",
			"lose sanity",
			"cry",
			"regain sanity",
			"kill",
		}
		var items []ToDoItem
		for _, title := range(titles) {
			item := ConstructToDoItem(
				Title(title),
				"high",
				false,
			)
			items = append(items, item)
		}
		store := testDataStore(items)

		got := store.Read()

		if !reflect.DeepEqual(items, got) {
			t.Errorf("want %v, got %v", items, got)
		}
	})
}

func testDataStore(items []ToDoItem) DataStore {
	dataMap := make(map[Title]ToDoItem)
	for _, item := range(items) {
		dataKey := item.Title
		dataMap[Title(dataKey)] = item
	}
	return DataStore{
		dataMap,
	}
}