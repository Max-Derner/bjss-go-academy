package main

import (
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("Add value to store", func(t *testing.T) {
		dataKey := "Keep sanity"
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)

		err, data := toDoMapper([]ToDoItem{item})
		if err != nil {
			t.Fatalf("setup failed! -> %v", err)
		}
		want := inMemoryDataStore{
			data,
		}

		store := inMemoryDataStore{make(map[Id]ToDoItem)}
		store.create(item)

		if !reflect.DeepEqual(want.data, store.data) {
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

		err, data := toDoMapper([]ToDoItem{item})
		if err != nil {
			t.Fatalf("setup failed! -> %v", err)
		}
		want := inMemoryDataStore{
			data,
		}

		store := inMemoryDataStore{make(map[Id]ToDoItem)}
		store.create(item)
		err = store.create(item)

		if err != ErrCannotCreate {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update value in store", func(t *testing.T) {
		dataKey := "Keep sanity"
		initialItem := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)
		updateItem := initialItem
		updateItem.Complete = true

		dataErr, data := toDoMapper([]ToDoItem{updateItem})
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		want := inMemoryDataStore{
			data,
		}

		store := inMemoryDataStore{make(map[Id]ToDoItem)}
		store.create(initialItem)
		err := store.update(updateItem)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !reflect.DeepEqual(want.data, store.data) {
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

		dataErr, data := toDoMapper([]ToDoItem{})
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		want := inMemoryDataStore{
			data,
		}

		store := inMemoryDataStore{make(map[Id]ToDoItem)}
		err := store.update(item)

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Deleting item", func(t *testing.T) {
		item := ConstructToDoItem(
			"Keep sanity",
			"high",
			false,
		)
		dataErr, data := toDoMapper([]ToDoItem{})
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		want := inMemoryDataStore{
			data,
		}

		store := inMemoryDataStore{make(map[Id]ToDoItem)}
		store.create(item)
		err := store.delete(item)

		if err != nil {
			t.Errorf("Unexpected error thrown! Got: %v", err)
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
	t.Run("Deleting non-existent item", func(t *testing.T) {
		item := ConstructToDoItem(
			"Keep sanity",
			"high",
			false,
		)
		dataErr, data := toDoMapper([]ToDoItem{})
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		want := inMemoryDataStore{
			data,
		}

		store := inMemoryDataStore{make(map[Id]ToDoItem)}
		err := store.delete(item)

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}

func TestRead(t *testing.T) {
	t.Run("Reading datastore", func(t *testing.T) {
		items := populatedToDoList()
		dataErr, data := toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		store := inMemoryDataStore{
			data,
		}

		got := store.read()

		if !equalSlicesNoOrder(items, got) {
			t.Errorf("want %v, got %v", items, got)
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
	for _, title := range titles {
		item := ConstructToDoItem(
			Title(title),
			"high",
			false,
		)
		items = append(items, item)
	}
	return items
}

func equalData(a, b map[Id]ToDoItem) bool {
	if len(a) != len(b) {
		return false
	}
	present := make(chan bool)
	for _, item := range a {
		go func(i ToDoItem, c chan bool) {
			for _, bItem := range b {
				if bItem == i {
					c <- true
					return
				}
			}
			c <- false
		}(item, present)
	}
	for _, item := range b {
		go func(i ToDoItem, c chan bool) {
			for _, aItem := range a {
				if aItem == i {
					c <- true
					return
				}
			}
			c <- false
		}(item, present)
	}

	for i := 0; i < len(a)+len(b); i++ {
		if !<-present {
			return false
		}
	}
	return true
}

func TestEqualData(t *testing.T) {
	t.Run("data is equal", func(t *testing.T) {
		items := populatedToDoList()
		dataErr, data := toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		a := inMemoryDataStore{data}
		dataErr, data = toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		b := inMemoryDataStore{data}
		if !equalData(a.data, b.data) {
			t.Errorf("Data was not considered equal! %v != %v", a.data, b.data)
		}
	})
	t.Run("data is not equal", func(t *testing.T) {
		items := populatedToDoList()
		dataErr, data := toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		a := inMemoryDataStore{data}
		dataErr, data = toDoMapper(items[:2])
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		b := inMemoryDataStore{data}
		if equalData(a.data, b.data) {
			t.Errorf("Data was considered equal! %v == %v", a.data, b.data)
		}
	})
	t.Run("data is still not equal", func(t *testing.T) {
		items := populatedToDoList()
		dataErr, data := toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		a := inMemoryDataStore{data}
		dataErr, data = toDoMapper(append([]ToDoItem{items[1]}, items[:1]...))
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		b := inMemoryDataStore{data}

		if equalData(a.data, b.data) {
			t.Errorf("Data was considered equal! %v == %v", a.data, b.data)
		}
	})
}

func equalSlicesNoOrder(a, b []ToDoItem) bool {
	if len(a) != len(b) {
		return false
	}
	present := make(chan bool)
	for _, item := range a {
		go func(i ToDoItem, c chan bool) {
			for _, bItem := range b {
				if bItem == i {
					c <- true
					return
				}
			}
			c <- false
		}(item, present)
	}
	for _, item := range b {
		go func(i ToDoItem, c chan bool) {
			for _, aItem := range a {
				if aItem == i {
					c <- true
					return
				}
			}
			c <- false
		}(item, present)
	}

	for i := 0; i < len(a)+len(b); i++ {
		if !<-present {
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
		b := append(a[splitPoint:], a[:splitPoint+1]...)
		equal := equalSlicesNoOrder(a, b)

		if equal {
			t.Error("Incorrectly stated non-equality")
		}
	})
	t.Run("Still Not Equal", func(t *testing.T) {
		a := populatedToDoList()
		splitPoint := 3
		b := append(a[splitPoint:], a[:splitPoint-1]...)
		equal := equalSlicesNoOrder(a, b)

		if equal {
			t.Error("Incorrectly stated non-equality")
		}
	})
}
