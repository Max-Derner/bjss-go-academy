package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestDataStoreDirectAdd(t *testing.T) {
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(item)

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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(item)
		err = store.add(item)

		if err != ErrCannotAdd {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreChannelAdd(t *testing.T) {
	t.Run("Add value to store", func(t *testing.T) {
		dataKey := "Keep sanity"
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)
		errChan := make(chan error)
		dataChan := make(chan []ToDoItem)
		err, data := toDoMapper([]ToDoItem{item})
		if err != nil {
			t.Fatalf("setup failed! -> %v", err)
		}
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.requests <- dbRequest{
			Add,
			item,
			errChan,
			dataChan,
		}

		err = <-errChan
		if err != nil {
			t.Errorf("Unexpected error thrown! %v", err)
		}
		if !equalData(want.data, store.data) {
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
		errChan := make(chan error)
		dataChan := make(chan []ToDoItem)
		dataErr, data := toDoMapper([]ToDoItem{item})
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(item)
		store.requests <- dbRequest{
			Add,
			item,
			errChan,
			dataChan,
		}

		err := <-errChan
		if err != ErrCannotAdd {
			t.Errorf("Unexpected error thrown! got %v want %v", err, ErrCannotAdd)
		}
		if !equalData(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreAPIAdd(t *testing.T) {
	t.Run("Add value to store", func(t *testing.T) {
		dataKey := "Keep sanity"
		item := ConstructToDoItem(
			Title(dataKey),
			"high",
			false,
		)

		dataErr, data := toDoMapper([]ToDoItem{item})
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		err := store.Add(item)

		if err != nil {
			t.Errorf("Unexpected error! -> %v", err)
		}
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.Add(item)
		err = store.Add(item)

		if err != ErrCannotAdd {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreAPIAddConcurrency(t *testing.T) {
	db := NewEmptyDataStore()
	title := Title("Something")
	priority := Priority("something else")
	expectedNumItems := 1_000
	finChan := make(chan struct{})
	for i := 0; i < expectedNumItems; i++ {
		go func(database dataStore, fin chan struct{}) {
			item := ConstructToDoItem(
				title,
				priority,
				false,
			)
			db.Add(item)
			fin <- struct{}{}
		}(db, finChan)
	}
	for i := 0; i < expectedNumItems; i++ {
		<- finChan
	}

	if len(db.data) != expectedNumItems {
		t.Errorf("Incorrect number of items found in DB! got %v want %v", len(db.data), expectedNumItems)
	}
}
func TestDataStoreDirectUpdate(t *testing.T) {
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(initialItem)
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		err := store.update(item)

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreChannelUpdate(t *testing.T) {
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(initialItem)
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- dbRequest{
			Update,
			updateItem,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- dbRequest{
			Update,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreAPIUpdate(t *testing.T) {
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(initialItem)
		err := store.Update(updateItem)

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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		err := store.Update(item)

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreDirectDelete(t *testing.T) {
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(item)
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		err := store.delete(item)

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreChannelDelete(t *testing.T) {
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(item)
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- dbRequest{
			Delete,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan

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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- dbRequest{
			Delete,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreAPIDelete(t *testing.T) {
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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		store.add(item)
		err := store.Delete(item)

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
		want := NewDataStore(data)

		store := NewEmptyDataStore()
		err := store.Delete(item)

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !reflect.DeepEqual(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
}
func TestDataStoreDirectRead(t *testing.T) {
	t.Run("Reading datastore", func(t *testing.T) {
		items := populatedToDoList()
		dataErr, data := toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		store := NewDataStore(data)

		got := store.read()

		if !equalSlicesNoOrder(items, got) {
			t.Errorf("want %v, got %v", items, got)
		}
	})
}
func TestDataStoreChannelRead(t *testing.T) {
	t.Run("Reading datastore", func(t *testing.T) {
		items := populatedToDoList()
		dataErr, data := toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		store := NewDataStore(data)

		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		dbr := dbRequest{
			Read,
			ToDoItem{},
			errReturnChan,
			dataReturnChan,
		}
		store.requests <- dbr
		err := <-errReturnChan
		got := <-dataReturnChan

		if err != nil {
			t.Errorf("Unexpected error thrown! Got: %v", err)
		}
		if got == nil || reflect.DeepEqual(got, []ToDoItem{}) {
			t.Fatal("No data returned!")
		}
		if !equalSlicesNoOrder(items, got) {
			t.Errorf("want %v, got %v", items, got)
		}
	})
}
func TestDataStoreAPIRead(t *testing.T) {
	t.Run("Reading datastore", func(t *testing.T) {
		items := populatedToDoList()
		dataErr, data := toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		store := NewDataStore(data)

		got := store.Read()

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

var ErrOverWritten = errors.New("item overwritten")

func toDoMapper(data []ToDoItem) (error, map[Id]ToDoItem) {
	dataMap := make(map[Id]ToDoItem)
	var err error
	for _, item := range data {
		_, itemExists := dataMap[item.Id]
		if itemExists {
			err = ErrOverWritten
		}
		dataMap[item.Id] = item
	}
	return err, dataMap
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
		a := NewDataStore(data)
		dataErr, data = toDoMapper(items)
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		b := NewDataStore(data)
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
		a := NewDataStore(data)
		dataErr, data = toDoMapper(items[:2])
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		b := NewDataStore(data)
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
		a := NewDataStore(data)
		dataErr, data = toDoMapper(append([]ToDoItem{items[1]}, items[:1]...))
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		b := NewDataStore(data)

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
