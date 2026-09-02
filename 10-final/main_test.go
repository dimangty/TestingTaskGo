package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"module-name/bins"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "create", args: []string{"--create", "--file=data.json", "--name=my-bin"}},
		{name: "update", args: []string{"--update", "--file=data.json", "--id=bin-id"}},
		{name: "delete", args: []string{"--delete", "--id=bin-id"}},
		{name: "get", args: []string{"--get", "--id=bin-id"}},
		{name: "list", args: []string{"--list"}},
		{name: "missing action", wantErr: "select exactly one action"},
		{name: "several actions", args: []string{"--get", "--list", "--id=bin-id"}, wantErr: "select exactly one action"},
		{name: "create missing name", args: []string{"--create", "--file=data.json"}, wantErr: "--create requires"},
		{name: "update missing ID", args: []string{"--update", "--file=data.json"}, wantErr: "--update requires"},
		{name: "get missing ID", args: []string{"--get"}, wantErr: "require --id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseOptions(test.args)
			if test.wantErr == "" && err != nil {
				t.Fatalf("parseOptions() error = %v", err)
			}
			if test.wantErr != "" && (err == nil || !strings.Contains(err.Error(), test.wantErr)) {
				t.Fatalf("parseOptions() error = %v, want substring %q", err, test.wantErr)
			}
		})
	}
}

func TestExecuteCreateStoresMetadataAndListPrintsIDAndName(t *testing.T) {
	store := &memoryStorage{}
	binList, err := bins.NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	client := &fakeAPI{created: bins.Bin{ID: "new-id", Private: true, Name: "my-bin"}}
	files := memoryFiles{"data.json": []byte(`{"value":1}`)}
	var output bytes.Buffer

	err = execute(
		context.Background(),
		options{create: true, file: "data.json", name: "my-bin"},
		&output,
		client,
		binList,
		files,
	)
	if err != nil {
		t.Fatalf("execute(create) error = %v", err)
	}
	if got := string(client.createDocument); got != `{"value":1}` || client.createName != "my-bin" {
		t.Fatalf("Create arguments = %q, %q", got, client.createName)
	}
	if len(store.saved) != 1 || store.saved[0].ID != "new-id" {
		t.Fatalf("saved bins = %+v", store.saved)
	}

	output.Reset()
	if err := execute(context.Background(), options{list: true}, &output, nil, binList, files); err != nil {
		t.Fatalf("execute(list) error = %v", err)
	}
	if got, want := output.String(), "new-id\tmy-bin\n"; got != want {
		t.Fatalf("list output = %q, want %q", got, want)
	}
}

func TestExecuteDeleteRemovesRemoteAndLocalBin(t *testing.T) {
	store := &memoryStorage{saved: []bins.Bin{{ID: "bin-id", Name: "my-bin"}}}
	binList, err := bins.NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	client := &fakeAPI{}

	err = execute(context.Background(), options{delete: true, id: "bin-id"}, &bytes.Buffer{}, client, binList, memoryFiles{})
	if err != nil {
		t.Fatalf("execute(delete) error = %v", err)
	}
	if client.deletedID != "bin-id" || store.deletedID != "bin-id" {
		t.Fatalf("deleted IDs: remote=%q local=%q", client.deletedID, store.deletedID)
	}
	if got := binList.Bins(); len(got) != 0 {
		t.Fatalf("BinList.Bins() = %+v, want empty", got)
	}
}

func TestReadJSONFileRejectsInvalidDocument(t *testing.T) {
	_, err := readJSONFile(memoryFiles{"broken.json": []byte(`{"broken"`)}, "broken.json")
	if err == nil {
		t.Fatal("readJSONFile() error = nil, want invalid JSON error")
	}
}

type fakeAPI struct {
	created        bins.Bin
	createDocument []byte
	createName     string
	deletedID      string
}

func (api *fakeAPI) Create(_ context.Context, document []byte, name string) (bins.Bin, error) {
	api.createDocument = append([]byte(nil), document...)
	api.createName = name
	return api.created, nil
}

func (*fakeAPI) Get(context.Context, string) ([]byte, error) {
	return nil, errors.New("unexpected Get call")
}

func (*fakeAPI) Update(context.Context, string, []byte) ([]byte, error) {
	return nil, errors.New("unexpected Update call")
}

func (api *fakeAPI) Delete(_ context.Context, id string) error {
	api.deletedID = id
	return nil
}

type memoryFiles map[string][]byte

func (files memoryFiles) Read(path string) ([]byte, error) {
	document, ok := files[path]
	if !ok {
		return nil, errors.New("file not found")
	}
	return append([]byte(nil), document...), nil
}

func (memoryFiles) IsJSON(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".json")
}

type memoryStorage struct {
	saved     []bins.Bin
	deletedID string
}

func (storage *memoryStorage) ReadBins() ([]bins.Bin, error) {
	return append([]bins.Bin(nil), storage.saved...), nil
}

func (storage *memoryStorage) SaveBin(bin bins.Bin) error {
	storage.saved = append(storage.saved, bin)
	return nil
}

func (storage *memoryStorage) DeleteBin(id string) error {
	storage.deletedID = id
	for index := range storage.saved {
		if storage.saved[index].ID == id {
			storage.saved = append(storage.saved[:index], storage.saved[index+1:]...)
			break
		}
	}
	return nil
}
