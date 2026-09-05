package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"module-name/api"
	"module-name/bins"
)

func TestParseOptions(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	testCases := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{name: "create", args: []string{"--create", "--file=data.json", "--name=my-bin"}},
		{name: "update", args: []string{"--update", "--file=data.json", "--id=bin-id"}},
		{name: "delete", args: []string{"--delete", "--id=bin-id"}},
		{name: "get", args: []string{"--get", "--id=bin-id"}},
		{name: "list", args: []string{"--list"}},
		{name: "missing action", expectedError: "select exactly one action"},
		{name: "several actions", args: []string{"--get", "--list", "--id=bin-id"}, expectedError: "select exactly one action"},
		{name: "create missing name", args: []string{"--create", "--file=data.json"}, expectedError: "--create requires"},
		{name: "update missing ID", args: []string{"--update", "--file=data.json"}, expectedError: "--update requires"},
		{name: "get missing ID", args: []string{"--get"}, expectedError: "require --id"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act - выполняем функцию
			_, err := parseOptions(tc.args)

			// Assert - проверка результата с expected
			if tc.expectedError == "" && err != nil {
				t.Fatalf("parseOptions() error = %v", err)
			}
			if tc.expectedError != "" && (err == nil || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("parseOptions() error = %v, want substring %q", err, tc.expectedError)
			}
		})
	}
}

func TestExecuteCreateBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	appOptions := parseTestOptions(t, "--create", "--file=data.json", "--name=my-bin")
	store := &memoryStorage{}
	binList, err := bins.NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	client := &fakeAPI{created: bins.Bin{ID: "new-id", Private: true, Name: "my-bin"}}
	files := memoryFiles{"data.json": []byte(`{"value":1}`)}
	var output bytes.Buffer

	// Act - выполняем функцию
	err = execute(
		context.Background(),
		appOptions,
		&output,
		client,
		binList,
		files,
	)

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("execute(create) error = %v", err)
	}
	if got, want := strings.Join(client.calls, ","), "create"; got != want {
		t.Fatalf("API calls = %q, want %q", got, want)
	}
	if got := string(client.createDocument); got != `{"value":1}` || client.createName != "my-bin" {
		t.Fatalf("Create arguments = %q, %q", got, client.createName)
	}
	if len(store.saved) != 1 || store.saved[0] != client.created {
		t.Fatalf("saved bins = %+v", store.saved)
	}

	var printed bins.Bin
	if err := json.Unmarshal(output.Bytes(), &printed); err != nil {
		t.Fatalf("decode create output: %v", err)
	}
	if printed != client.created {
		t.Fatalf("create output = %+v, want %+v", printed, client.created)
	}
}

func TestExecuteUpdateBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	appOptions := parseTestOptions(t, "--update", "--file=data.json", "--id=bin-id")
	store := &memoryStorage{}
	binList, err := bins.NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	client := &fakeAPI{updated: []byte(`{"record":{"value":2},"metadata":{"id":"bin-id"}}`)}
	files := memoryFiles{"data.json": []byte(`{"value":2}`)}
	var output bytes.Buffer

	// Act - выполняем функцию
	err = execute(
		context.Background(),
		appOptions,
		&output,
		client,
		binList,
		files,
	)

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("execute(update) error = %v", err)
	}
	if got, want := strings.Join(client.calls, ","), "update"; got != want {
		t.Fatalf("API calls = %q, want %q", got, want)
	}
	if got := string(client.updateDocument); got != `{"value":2}` || client.updatedID != "bin-id" {
		t.Fatalf("Update arguments = %q, %q", client.updatedID, got)
	}
	if got, want := output.String(), "{\n  \"record\": {\n    \"value\": 2\n  },\n  \"metadata\": {\n    \"id\": \"bin-id\"\n  }\n}\n"; got != want {
		t.Fatalf("update output = %q, want %q", got, want)
	}
	if len(store.saved) != 0 || store.deletedID != "" {
		t.Fatalf("update changed local storage: %+v", store)
	}
}

func TestExecuteGetBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	appOptions := parseTestOptions(t, "--get", "--id=bin-id")
	store := &memoryStorage{}
	binList, err := bins.NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	client := &fakeAPI{got: []byte(`{"value":1}`)}
	var output bytes.Buffer

	// Act - выполняем функцию
	err = execute(
		context.Background(),
		appOptions,
		&output,
		client,
		binList,
		memoryFiles{},
	)

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("execute(get) error = %v", err)
	}
	if got, want := strings.Join(client.calls, ","), "get"; got != want {
		t.Fatalf("API calls = %q, want %q", got, want)
	}
	if client.gotID != "bin-id" {
		t.Fatalf("Get ID = %q, want %q", client.gotID, "bin-id")
	}
	if got, want := output.String(), "{\n  \"value\": 1\n}\n"; got != want {
		t.Fatalf("get output = %q, want %q", got, want)
	}
}

func TestExecuteDeleteBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	appOptions := parseTestOptions(t, "--delete", "--id=bin-id")
	store := &memoryStorage{saved: []bins.Bin{{ID: "bin-id", Name: "my-bin"}}}
	binList, err := bins.NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	client := &fakeAPI{}
	var output bytes.Buffer

	// Act - выполняем функцию
	err = execute(context.Background(), appOptions, &output, client, binList, memoryFiles{})

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("execute(delete) error = %v", err)
	}
	if got, want := strings.Join(client.calls, ","), "delete"; got != want {
		t.Fatalf("API calls = %q, want %q", got, want)
	}
	if client.deletedID != "bin-id" || store.deletedID != "bin-id" {
		t.Fatalf("deleted IDs: remote=%q local=%q", client.deletedID, store.deletedID)
	}
	if got := binList.Bins(); len(got) != 0 {
		t.Fatalf("BinList.Bins() = %+v, want empty", got)
	}
	if got, want := output.String(), "deleted bin bin-id\n"; got != want {
		t.Fatalf("delete output = %q, want %q", got, want)
	}
}

func TestReadJSONFileRejectsInvalidDocument(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expected := api.ErrInvalidJSON
	files := memoryFiles{"broken.json": []byte(`{"broken"`)}
	path := "broken.json"

	// Act - выполняем функцию
	_, err := readJSONFile(files, path)

	// Assert - проверка результата с expected
	if !errors.Is(err, expected) {
		t.Fatalf("readJSONFile() error = %v, want %v", err, expected)
	}
}

func parseTestOptions(t *testing.T, args ...string) options {
	t.Helper()

	appOptions, err := parseOptions(args)
	if err != nil {
		t.Fatalf("parseOptions(%q) error = %v", args, err)
	}
	return appOptions
}

type fakeAPI struct {
	calls          []string
	created        bins.Bin
	createDocument []byte
	createName     string
	got            []byte
	gotID          string
	updated        []byte
	updatedID      string
	updateDocument []byte
	deletedID      string
}

func (api *fakeAPI) Create(_ context.Context, document []byte, name string) (bins.Bin, error) {
	api.calls = append(api.calls, "create")
	api.createDocument = append([]byte(nil), document...)
	api.createName = name
	return api.created, nil
}

func (api *fakeAPI) Get(_ context.Context, id string) ([]byte, error) {
	api.calls = append(api.calls, "get")
	api.gotID = id
	return append([]byte(nil), api.got...), nil
}

func (api *fakeAPI) Update(_ context.Context, id string, document []byte) ([]byte, error) {
	api.calls = append(api.calls, "update")
	api.updatedID = id
	api.updateDocument = append([]byte(nil), document...)
	return append([]byte(nil), api.updated...), nil
}

func (api *fakeAPI) Delete(_ context.Context, id string) error {
	api.calls = append(api.calls, "delete")
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
