package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type jsonDataStore struct {
	data     map[Id]ToDoItem
	fileName string
}

func newJSONDataStore() jsonDataStore {
	ds := jsonDataStore{
		make(map[Id]ToDoItem),
		"data.json",
	}
	ds.lift()
	return ds
}

func (d *jsonDataStore) lift() {
	f, err := os.Open(d.fileName)
	if err == nil {
		defer f.Close()
		jsonData, err := io.ReadAll(f)
		if err != nil {
			fmt.Printf("ERROR! %q\n", err)
		}
		todos := make(map[Id]ToDoItem)
		err = json.Unmarshal(jsonData, &todos)
		if err != nil {
			fmt.Printf("x: %q\n", err)
			fmt.Printf("ERROR! %q\n", err)
		}
		if err != nil {
			fmt.Printf("ERROR! %q\n", err)
		}
		d.data = todos
	} else {
		os.Create(d.fileName)
	}
}

func (d *jsonDataStore) place() {
	f, err := os.Create(d.fileName)
	if err == nil {
		defer f.Close()
		jsonNibbles, err := json.MarshalIndent(d.data, "", "  ")
		if err == nil {
			_, err = f.Write(jsonNibbles)
			if err != nil {
				fmt.Printf("ERROR: %q\n", err)
			}
		} else {
			fmt.Printf("ERROR: %q\n", err)
		}
	} else {
		fmt.Printf("ERROR! %q\n", err)
	}
}

func (d jsonDataStore) read() []ToDoItem {
	d.lift()
	var dataSlice []ToDoItem
	for _, item := range d.data {
		dataSlice = append(dataSlice, item)
	}
	return dataSlice
}

func (d *jsonDataStore) delete(item ToDoItem) error {
	d.lift()
	defer d.place()
	_, keyExists := d.data[item.Id]
	if !keyExists {
		return ErrCannotDelete
	}
	delete(d.data, item.Id)
	return nil
}

func (d *jsonDataStore) update(item ToDoItem) error {
	d.lift()
	defer d.place()
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if !keyExists {
		return ErrCannotUpdate
	}
	d.data[dataKey] = item
	return nil
}

func (d *jsonDataStore) create(item ToDoItem) error {
	d.lift()
	defer d.place()
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if keyExists {
		return ErrCannotCreate
	}
	d.data[dataKey] = item
	return nil
}
