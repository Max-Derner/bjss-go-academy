package main

type inMemoryDataStore struct {
	data map[Id]ToDoItem
}

func newEmptyInMemoryDataStore() inMemoryDataStore {
	return inMemoryDataStore{
		make(map[Id]ToDoItem),
	}
}

func (d inMemoryDataStore) read() []ToDoItem {
	var dataSlice []ToDoItem
	for _, item := range d.data {
		dataSlice = append(dataSlice, item)
	}
	return dataSlice
}

func (d *inMemoryDataStore) delete(item ToDoItem) error {
	_, keyExists := d.data[item.Id]
	if !keyExists {
		return ErrCannotDelete
	}
	delete(d.data, item.Id)
	return nil
}

func (d *inMemoryDataStore) update(item ToDoItem) error {
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if !keyExists {
		return ErrCannotUpdate
	}
	d.data[dataKey] = item
	return nil
}

func (d *inMemoryDataStore) create(item ToDoItem) error {
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if keyExists {
		return ErrCannotCreate
	}
	d.data[dataKey] = item
	return nil
}
