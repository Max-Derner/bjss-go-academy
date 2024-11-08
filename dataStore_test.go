package main

import (
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

		want := testDataStore([]ToDoItem{item})
		
		store := NewDataStore()
		store.Add(item)

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

		want := testDataStore([]ToDoItem{item})
		
		store := NewDataStore()
		store.Add(item)
		err := store.Add(item)

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
		want := testDataStore([]ToDoItem{item})
		
		store := NewDataStore()
		store.requests <- DbRequest{
			Add,
			item,
			errChan,
			dataChan,
		}

		err := <- errChan
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
		want := testDataStore([]ToDoItem{item})
		
		store := NewDataStore()
		store.Add(item)
		store.requests <- DbRequest{
			Add,
			item,
			errChan,
			dataChan,
		}

		err := <- errChan
		if err != ErrCannotAdd {
			t.Errorf("Unexpected error thrown! got %v want %v", err, ErrCannotAdd)
		}
		if !equalData(want.data, store.data) {
			t.Errorf("want %v, got %v", want, store)
		}
	})
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

		want := testDataStore([]ToDoItem{updateItem})
		
		store := NewDataStore()
		store.Add(initialItem)
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

		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		err := store.Update(item)

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
		want := testDataStore([]ToDoItem{updateItem})
		
		store := NewDataStore()
		store.Add(initialItem)
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- DbRequest{
			Update,
			updateItem,
			errReturnChan,
			dataReturnChan,
		}
		err := <- errReturnChan
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

		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- DbRequest{
			Update,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <- errReturnChan

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
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		store.Add(item)
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
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		err := store.Delete(item)

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
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		store.Add(item)
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- DbRequest{
			Delete,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <- errReturnChan
		
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
		want := testDataStore([]ToDoItem{})
		
		store := NewDataStore()
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		store.requests <- DbRequest{
			Delete,
			item,
			errReturnChan,
			dataReturnChan,
		}
		err := <- errReturnChan

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
		store := testDataStore(items)

		got := store.Read()

		if !equalSlicesNoOrder(items, got) {
			t.Errorf("want %v, got %v", items, got)
		}
	})
}
func TestDataStoreChannelRead(t *testing.T) {
	t.Run("Reading datastore", func(t *testing.T) {
		items := populatedToDoList()
		store := testDataStore(items)
		
		errReturnChan := make(chan error)
		dataReturnChan := make(chan []ToDoItem)
		dbr := DbRequest{	// Whole thing hangs up here for fuck knows why
			Read,			//
			ToDoItem{},		//  X    X
			errReturnChan,	//    __  '
			dataReturnChan,	//  /    \
		}
		store.requests <- dbr
		err := <- errReturnChan
		got := <- dataReturnChan

		if err != nil {
			t.Errorf("Unexpected error thrown! Got: %v", err)
		}
		if got == nil {
			t.Fatal("No data returned!")
		}
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
		make(chan DbRequest),
	}
}

func equalData(a, b map[Id]ToDoItem) bool {
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

func TestEqualData(t *testing.T) {
	t.Run("data is equal", func(t *testing.T) {
		items := populatedToDoList()
		a := testDataStore(items)
		b := testDataStore(items)
		if !equalData(a.data, b.data) {
			t.Errorf("Data was not considered equal! %v != %v", a.data, b.data)
		}
	})
	t.Run("data is not equal", func(t *testing.T) {
		items := populatedToDoList()
		a := testDataStore(items)
		b := testDataStore(items[:2])
		if equalData(a.data, b.data) {
			t.Errorf("Data was considered equal! %v == %v", a.data, b.data)
		}
	})
	t.Run("data is still not equal", func(t *testing.T) {
		items := populatedToDoList()
		a := testDataStore(items)
		b := testDataStore(append([]ToDoItem{items[1]}, items[:1]...))
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