package repository

import (
	"os"
	"testing"
)

const test_db_path = "./test_db.bin"

func TestNewRepository(t *testing.T) {
	repo := NewRepository(test_db_path)

	if repo == nil {
		t.Error("NewRepository should not return nil")
	}

	if repo.data_file_path != test_db_path {
		t.Error("data_file_path should be set")
	}
}

func TestListWithoutQuery(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	secrets := []Entry{
		Entry{Name: "first entry", UserName: "first", Password: "frstpwd", Note: "first value"},
		Entry{Name: "second entry", UserName: "second", Password: "scndpwd", Note: "second value"},
	}
	repo.entries = secrets

	repo.dump("testpassword")

	defer os.Remove(test_db_path)

	entries := repo.List("", "testpassword")

	if len(entries) != 2 {
		t.Error("List should return 2 entries")
	}

	if entries[0].Name != "first entry" {
		t.Error("First entry name should be 'first entry'")
	}

	if entries[1].Name != "second entry" {
		t.Error("Second entry name should be 'second entry'")
	}
}

func TestListWithQuery(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	secrets := []Entry{
		Entry{Name: "first entry", UserName: "first", Password: "frstpwd", Note: "first value"},
		Entry{Name: "second entry", UserName: "second", Password: "scndpwd", Note: "second value"},
	}
	repo.entries = secrets

	repo.dump("testpassword")

	defer os.Remove(test_db_path)

	entries := repo.List("first", "testpassword")

	if len(entries) != 1 {
		t.Error("List should return 1 entries, but got", len(entries))
	}

	if entries[0].Name != "first entry" {
		t.Error("First entry name should be 'first entry'")
	}
}

func TestListQueryCaseInsensitive(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	secrets := []Entry{
		Entry{Name: "first entry", UserName: "first", Password: "frstpwd", Note: "first value"},
		Entry{Name: "second entry", UserName: "second", Password: "scndpwd", Note: "second value"},
	}
	repo.entries = secrets

	repo.dump("testpassword")

	defer os.Remove(test_db_path)

	entries := repo.List("First", "testpassword")

	if len(entries) != 1 {
		t.Error("List should return 1 entries, but got", len(entries))
	}

	if entries[0].Name != "first entry" {
		t.Error("First entry name should be 'first entry'")
	}
}

func TestAdd(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	repo.load("testpassword")

	defer os.Remove(test_db_path)

	entry := Entry{Name: "first entry", UserName: "first", Password: "frstpwd", Note: "first value"}

	repo.entries = []Entry{entry}
	repo.dump("testpassword")

	entries := repo.List("", "testpassword")

	if len(entries) != 1 {
		t.Error("List should return 1 entry")
	}

	if entries[0].Name != "first entry" {
		t.Error("First entry name should be 'first entry'")
	}

	file_size, _ := os.Stat(test_db_path)
	if file_size.Size() == 0 {
		t.Error("File should not be empty")
	}
}

func TestGet(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	repo.load("testpassword")

	defer os.Remove(test_db_path)

	entry := Entry{Name: "first entry", UserName: "first realy", Password: "frstpwd", Note: "first value"}
	one_more_first := Entry{Name: "one more first entry", UserName: "first realy", Password: "frstpwd", Note: "first value"}
	second := Entry{Name: "second entry", UserName: "first", Password: "frstpwd", Note: "first value"}

	repo.Add(entry, "testpassword")
	repo.Add(one_more_first, "testpassword")
	repo.Add(second, "testpassword")

	firstID := repo.entries[0].ID

	found, err := repo.Get(firstID, "testpassword")

	if err != nil {
		t.Error("Get should not return error")
	}

	if found.Name != "first entry" {
		t.Error("First entry name should be 'first entry'")
	}
}

func TestCheckPasswordWhenFileIsEmpty(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	err := repo.CheckPassword("testpassword")
	defer os.Remove(test_db_path)

	if err != nil {
		t.Error("Load should not return error")
	}
}

func TestCheckPasswordWhenPasswordIsWrong(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	repo.load("testpassword")
	repo.Add(Entry{Name: "first entry", UserName: "first", Password: "frstpwd", Note: "first value"}, "testpassword")

	defer os.Remove(test_db_path)

	err := repo.CheckPassword("wrongpassword")

	if err == nil {
		t.Error("CheckPassword should return error")
	}
}

func TestCheckPasswordWhenPasswordIsCorrect(t *testing.T) {
	os.Create(test_db_path)
	repo := NewRepository(test_db_path)

	repo.load("testpassword")
	repo.Add(Entry{Name: "first entry", UserName: "first", Password: "frstpwd", Note: "first value"}, "testpassword")

	defer os.Remove(test_db_path)

	err := repo.CheckPassword("testpassword")

	if err != nil {
		t.Error("CheckPassword should not return error")
	}
}
