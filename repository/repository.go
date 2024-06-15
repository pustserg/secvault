package repository

import (
	"os"
	"regexp"
	"sync"
)

type RepositoryInterface interface {
	List(query, password string) []Entry
	Get(name, password string) (Entry, error)
	Add(entry Entry, password string) error
}

type Entry struct {
	Kind      string
	Name      string
	UserName  string
	Password  string
	TotpToken string
	Value     string
}

type Repository struct {
	data_file_path string
	lock           sync.RWMutex
	entries        []Entry
}

func NewRepository(data_file_path string) *Repository {
	r := Repository{
		data_file_path: data_file_path,
	}

	return &r
}

func (r *Repository) load(password string) error {
	r.entries = []Entry{}

	content, err := os.ReadFile(r.data_file_path)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		r.entries = []Entry{}
		return nil
	}

	entries, err := decode(content, password)
	if err != nil {
		return err
	} else {
		r.entries = entries
	}
	return nil
}

func (r *Repository) dump(password string) error {
	content, err := encode(r.entries, password)
	if err != nil {
		return err
	}
	os.WriteFile(r.data_file_path, content, 0644)
	return nil
}

func (r *Repository) List(query, password string) []Entry {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.load(password)
	if query == "" {
		return r.entries
	}

	result := []Entry{}
	regexp := regexp.MustCompile(query)
	for _, entry := range r.entries {
		if regexp.MatchString(entry.Name) {
			result = append(result, entry)
		}
	}
	return result
}

func (r *Repository) Get(name, password string) (Entry, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.load(password)
	regexp := regexp.MustCompile(name)
	found := Entry{}
	for _, entry := range r.entries {
		if regexp.MatchString(entry.Name) {
			found = entry
			break
		}
	}
	return found, nil
}

func (r *Repository) Add(entry Entry, password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.load(password)
	r.entries = append(r.entries, entry)
	r.dump(password)
	return nil
}
