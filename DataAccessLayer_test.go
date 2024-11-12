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
		want := NewDataAccessLayer(&inMemoryDataStore{data})

		db := newEmptyInMemoryDataStore()
		dal := NewDataAccessLayer(&db)
		dal.requests <- dbRequest{
			Create,
			item,
			errChan,
			dataChan,
		}

		err = <-errChan
		if err != nil {
			t.Errorf("Unexpected error thrown! %v", err)
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
		errChan := make(chan error)
		dataChan := make(chan []ToDoItem)
		dataErr, data := toDoMapper([]ToDoItem{item})
		if dataErr != nil {
			t.Fatalf("setup failed! -> %v", dataErr)
		}
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{data}
		dal := NewDataAccessLayer(&db2)
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{make(map[Id]ToDoItem)}
		dal := NewDataAccessLayer(&db2)
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{make(map[Id]ToDoItem)}
		dal := NewDataAccessLayer(&db2)
		dal.Create(item)
		err = dal.Create(item)

		if err != ErrCannotCreate {
			t.Fatal("Error not thrown")
		}

		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
		}
	})
}
func BenchmarkAPICreateConcurrency(t *testing.B) {
	db := newEmptyInMemoryDataStore()
	dal := NewDataAccessLayer(&db)
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{map[Id]ToDoItem{initialItem.Id: initialItem}}
		dal := NewDataAccessLayer(&db2)
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{make(map[Id]ToDoItem)}
		dal := NewDataAccessLayer(&db2)
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{map[Id]ToDoItem{initialItem.Id: initialItem}}
		dal := NewDataAccessLayer(&db2)
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{make(map[Id]ToDoItem)}
		dal := NewDataAccessLayer(&db2)
		err := dal.Update(item)

		if err != ErrCannotUpdate {
			t.Fatal("Error not thrown")
		}

		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{map[Id]ToDoItem{item.Id: item}}
		dal := NewDataAccessLayer(&db2)
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		dal.requests <- dbRequest{
			Delete,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan

		if err != nil {
			t.Errorf("Unexpected error thrown! Got: %v", err)
		}

		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{make(map[Id]ToDoItem)}
		dal := NewDataAccessLayer(&db2)
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		dal.requests <- dbRequest{
			Delete,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <-errReturnChan

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{map[Id]ToDoItem{item.Id: item}}
		dal := NewDataAccessLayer(&db2)
		err := dal.Delete(item)

		if err != nil {
			t.Errorf("Unexpected error thrown! Got: %v", err)
		}

		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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
		db := inMemoryDataStore{data}
		want := NewDataAccessLayer(&db)

		db2 := inMemoryDataStore{make(map[Id]ToDoItem)}
		dal := NewDataAccessLayer(&db2)
		err := dal.Delete(item)

		if err != ErrCannotDelete {
			t.Fatal("Error not thrown")
		}

		if !equalSlicesNoOrder(want.db.read(), dal.db.read()) {
			t.Errorf("want %v, got %v", want.db.read(), dal.db.read())
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
		db := inMemoryDataStore{data}
		dal := NewDataAccessLayer(&db)

		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		dbr := dbRequest{
			Read,
			ToDoItem{},
			errReturnChan,
			dataReturnChan,
		}
		dal.requests <- dbr
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
		db := inMemoryDataStore{data}
		dal := NewDataAccessLayer(&db)

		got := dal.Read()

		if !equalSlicesNoOrder(items, got) {
			t.Errorf("want %v, got %v", items, got)
		}
	})
}
