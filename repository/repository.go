package repository

import (
	"os"
	"regexp"
	"sync"

	"github.com/google/uuid"
)

const (
	NoteType    = "Note"
	WebsiteType = "Website"
)

type RepositoryInterface interface {
	List(query, password string) []Entry
	Get(ID, password string) (Entry, error)
	Add(entry Entry, password string) error
	Delete(ID, password string) error
	Update(entry Entry, password string) error
	CheckPassword(password string) error
}

type Entry struct {
	ID        string
	Kind      string
	Name      string
	URL       string
	UserName  string
	Password  string
	TotpToken string
	Note      string
}

type Repository struct {
	dataFilePath string
	lock         sync.RWMutex
	entries      []Entry
}

func NewRepository(dataFilePath string) *Repository {
	r := Repository{
		dataFilePath: dataFilePath,
	}

	return &r
}

func (r *Repository) load(password string) error {
	r.entries = []Entry{}

	content, err := os.ReadFile(r.dataFilePath)
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
	os.WriteFile(r.dataFilePath, content, 0644)
	return nil
}

func (r *Repository) CheckPassword(password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	content, err := os.ReadFile(r.dataFilePath)
	if err != nil {
		return err
	}

	// if file is empty, we haven't saved any data yet
	if len(content) == 0 {
		return nil
	}

	saltForPasswordHash := content[:32]
	passwordHash := content[32:64]

	if verifyPassword(password, passwordHash, saltForPasswordHash) {
		return nil
	} else {
		return ErrInvalidPassword
	}
}

func (r *Repository) List(query, password string) []Entry {
	r.lock.Lock()
	defer r.lock.Unlock()

	err := r.load(password)
	if err != nil {
		return []Entry{}
	}
	if query == "" {
		return r.entries
	}

	var result []Entry
	searchPattern := regexp.MustCompile("(?i)" + query)
	for _, entry := range r.entries {
		if searchPattern.MatchString(entry.Name) {
			result = append(result, entry)
		}
	}
	return result
}

func (r *Repository) Get(ID, password string) (Entry, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.load(password)
	found := Entry{}
	for _, entry := range r.entries {
		if entry.ID == ID {
			found = entry
			break
		}
	}
	return found, nil
}

func (r *Repository) Add(entry Entry, password string) error {
	entry.ID = uuid.New().String()

	r.lock.Lock()
	defer r.lock.Unlock()
	r.load(password)
	r.entries = append(r.entries, entry)
	r.dump(password)
	return nil
}

func (r *Repository) Delete(ID, password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.load(password)
	for i, entry := range r.entries {
		if entry.ID == ID {
			r.entries = append(r.entries[:i], r.entries[i+1:]...)
			break
		}
	}
	r.dump(password)
	r.load(password)
	return nil
}

func (r *Repository) Update(entry Entry, password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.load(password)
	for i, e := range r.entries {
		if e.ID == entry.ID {
			r.entries[i] = entry
			break
		}
	}
	r.dump(password)
	return nil
}
