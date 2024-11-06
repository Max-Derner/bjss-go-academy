package main

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
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
		updateItem := initialItem
		updateItem.Complete = true

		want := testDataStore([]ToDoItem{updateItem})
		
		store := NewDataStore()
		store.Add(initialItem)
		err := store.Update(updateItem)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
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
		item := ConstructToDoItem(
			"Keep sanity",
			"high",
			false,
		)
		dataKey := item.Id
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		store.Add(item)
		store.Delete(dataKey)

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Deleting non-existent item", func(t *testing.T) {
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		err := store.Delete(Id(uuid.New()))

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want, store) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Reading datastore", func(t *testing.T) {
		items := populatedToDoList()
		store := testDataStore(items)

		got := store.Read()

		if !equalSlicesNoOrder(items, got) {
			t.Errorf("want %v, got %v", items, got)
		}
	})
	t.Run("Query datastore", func(t *testing.T) {
		items := populatedToDoList()
		store := testDataStore(items)
		want := items[4]

		err, got := store.Query(want.Title)


		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
		if err != nil {
			t.Errorf("UNexpected exception thrown: %v", err)
		}
	})
	t.Run("Query datastore for non existent item", func(t *testing.T) {
		items := populatedToDoList()
		store := testDataStore(items)

		err, _ := store.Query("being happy")

		if err != ErrCannotQuery {
			t.Errorf("Incorrect error thrown: got %v, want %v", err, ErrCannotQuery)
		}
	})
}

func populatedToDoList() []ToDoItem {
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
	return items
}

func testDataStore(items []ToDoItem) DataStore {
	dataMap := make(map[Id]ToDoItem)
	for _, item := range(items) {
		dataKey := item.Id
		dataMap[dataKey] = item
	}
	return DataStore{
		dataMap,
	}
}

func equalSlicesNoOrder(a, b []ToDoItem) bool {
	if len(a) != len(b) {
		return false
	}
	present := make(chan bool)
	for _, item := range(a) {
		go func(i ToDoItem, c chan bool) {
			for _, bItem := range(b) {
				if bItem == i {
					c <- true
					return
				}
			}
			c <- false
		}(item, present)
	}
	for _, item := range(b) {
		go func(i ToDoItem, c chan bool) {
			for _, aItem := range(a) {
				if aItem == i {
					c <- true
					return
				}
			}
			c <- false
		}(item, present)
	}

	for i := 0; i < len(a) + len(b); i++ {
		if ! <- present {
			return false
		}
	}
	return true
}

func TestEqualSlicesNoOrder(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		a := populatedToDoList()
		splitPoint := 3
		b := append(a[splitPoint:], a[:splitPoint]...)
		equal := equalSlicesNoOrder(a, b)

		if !equal {
			t.Error("Incorrectly stated equality")
		}
	})
	t.Run("Not Equal", func(t *testing.T) {
		a := populatedToDoList()
		splitPoint := 3
		b := append(a[splitPoint:], a[:splitPoint + 1]...)
		equal := equalSlicesNoOrder(a, b)

		if equal {
			t.Error("Incorrectly stated non-equality")
		}
	})
	t.Run("Still Not Equal", func(t *testing.T) {
		a := populatedToDoList()
		splitPoint := 3
		b := append(a[splitPoint:], a[:splitPoint - 1]...)
		equal := equalSlicesNoOrder(a, b)

		if equal {
			t.Error("Incorrectly stated non-equality")
		}
	})
}