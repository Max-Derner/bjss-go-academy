package main

import (
	"reflect"
	"testing"
)


func TestChannelCreate(t *testing.T) {
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
			Create,
			item,
			errChan,
			dataChan,
		}

		err = <-errChan
		if err != nil {
			t.Errorf("Unexpected error thrown! %v", err)
		}
		if !equalSlicesNoOrder(want.db.read(), store.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), store.db.read())
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

		dal := NewDataStore(data)
		dal.requests <- dbRequest{
			Create,
			item,
			errChan,
			dataChan,
		}

		err := <-errChan
		if err != ErrCannotCreate {
			t.Errorf("Unexpected error thrown! got %v want %v", err, ErrCannotCreate)
		}
		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
		}
	})
}
func TestAPICreate(t *testing.T) {
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

		dal := NewEmptyDataStore()
		err := dal.Create(item)

		if err != nil {
			t.Errorf("Unexpected error! -> %v", err)
		}
		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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
		store.Create(item)
		err = store.Create(item)

		if err != ErrCannotCreate {
			t.Fatal("Error not thrown")
		}

		if !equalSlicesNoOrder(want.db.read(), store.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), store.db.read())
		}
	})
}
func BenchmarkAPICreateConcurrency(t *testing.B) {
	dal := NewEmptyDataStore()
	title := Title("Something")
	priority := Priority("something else")
	expectedNumItems := t.N
	finChan := make(chan struct{})
	for i := 0; i < expectedNumItems; i++ {
		go func(database DataAccessLayer, fin chan struct{}) {
			item := ConstructToDoItem(
				title,
				priority,
				false,
			)
			dal.Create(item)
			fin <- struct{}{}
		}(dal, finChan)
	}
	for i := 0; i < expectedNumItems; i++ {
		<-finChan
	}

	if len(dal.db.read()) != expectedNumItems {
		t.Errorf("Incorrect number of items found in DB! got %v want %v", len(dal.db.read()), expectedNumItems)
	}
}

func TestChannelUpdate(t *testing.T) {
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

		dal := NewDataStore(map[Id]ToDoItem{initialItem.Id: initialItem})
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		dal.requests <- dbRequest{
			Update,
			updateItem,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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

		dal := NewEmptyDataStore()
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		dal.requests <- dbRequest{
			Update,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
		}
	})
}
func TestAPIUpdate(t *testing.T) {
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

		dal := NewDataStore(map[Id]ToDoItem{initialItem.Id: initialItem})
		err := dal.Update(updateItem)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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

		if !equalSlicesNoOrder(want.db.read(), store.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), store.db.read())
		}
	})
}
func TestChannelDelete(t *testing.T) {
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

		store := NewDataStore(map[Id]ToDoItem{item.Id: item})
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

		if !equalSlicesNoOrder(want.db.read(), store.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), store.db.read())
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

		if !equalSlicesNoOrder(want.db.read(), store.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), store.db.read())
		}
	})
}
func TestAPIDelete(t *testing.T) {
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

		store := NewDataStore(map[Id]ToDoItem{item.Id: item})
		err := store.Delete(item)

		if err != nil {
			t.Errorf("Unexpected error thrown! Got: %v", err)
		}

		if !equalSlicesNoOrder(want.db.read(), store.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), store.db.read())
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

		if !equalSlicesNoOrder(want.db.read(), store.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), store.db.read())
		}
	})
}
func TestChannelRead(t *testing.T) {
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
func TestAPIRead(t *testing.T) {
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
