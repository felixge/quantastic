package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const timeFile = "time.json"

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

func OpenDb(dir string) (*Db, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}
	db := &Db{dir: dir, timeEntries: make(map[string]*TimeEntry)}
	if err := db.loadTimeEntries(); err != nil {
		return nil, err
	}
	return db, nil
}

type Db struct {
	dir         string
	timeEntries map[string]*TimeEntry
}

func (d *Db) TimeEntry(id string) (*TimeEntry, error) {
	entry, ok := d.timeEntries[id]
	if !ok {
		return nil, NotFoundError(id+" not found.")
	}
	return entry, nil
}

func (d *Db) SaveTimeEntry(e *TimeEntry) error {
	if e.Id == "" {
		e.Id = mustUUID()
	}
	d.timeEntries[e.Id] = e
	return d.saveTimeEntries()
}

func (d *Db) DeleteTimeEntry(id string) (error) {
	if _, ok := d.timeEntries[id]; !ok {
		return NotFoundError(id+" not found.")
	}
	delete(d.timeEntries, id)
	return d.saveTimeEntries()
}

func (d *Db) TimeEntriesByStart() ([]*TimeEntry, error) {
	sorted := make(timeEntriesByStart, 0, len(d.timeEntries))
	for _, entry := range d.timeEntries {
		sorted = append(sorted, entry)
	}
	sort.Sort(sorted)
	return []*TimeEntry(sorted), nil
}

func (d *Db) LatestTimeEntry() (*TimeEntry, error) {
	entries, err := d.TimeEntriesByStart()
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, NotFoundError("No time entries found.")
	}
	return entries[0], nil
}

func (d *Db) ActiveTimeEntry() (*TimeEntry, error) {
	entry, err := d.LatestTimeEntry()
	if err != nil {
		return nil, err
	}
	if !entry.Active() {
		return nil, NotFoundError("No active entry found.")
	}
	return entry, nil
}

func (d *Db) loadTimeEntries() error {
	data, err := ioutil.ReadFile(filepath.Join(d.dir, timeFile))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if err := json.Unmarshal(data, &d.timeEntries); err != nil {
			return err
		}
	}
	return nil
}

func (d *Db) saveTimeEntries() error {
	file, err := ioutil.TempFile("", timeFile)
	if err != nil {
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)
	if err := e.Encode(d.timeEntries); err != nil {
		return err
	}
	return os.Rename(file.Name(), filepath.Join(d.dir, timeFile))
}

func mustUUID() string {
	b := make([]byte, 16)
	if _, err := io.ReadAtLeast(rand.Reader, b, len(b)); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}

type TimeEntry struct {
	Id       string
	Category []string
	Start    time.Time
	End      *time.Time
	Note     string
}

func (t *TimeEntry) Valid() error {
	if len(t.Category) == 0 {
		return fmt.Errorf("Category must not be empty.")
	}
	for _, part := range t.Category {
		if part == "" {
			return fmt.Errorf("Category part must not be empty.")
		}
	}
	if t.Start.IsZero() {
		return fmt.Errorf("Start must not be empty.")
	}
	return nil
}

func (t *TimeEntry) Duration() time.Duration{
	if t.End == nil {
		return UtcNow().Sub(t.Start)
	} else {
		return t.End.Sub(t.Start)
	}
}

func (t *TimeEntry) Active() bool {
	return t.End == nil
}

type timeEntriesByStart []*TimeEntry

func (entries timeEntriesByStart) Len() int {
	return len(entries)
}

func (entries timeEntriesByStart) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

func (entries timeEntriesByStart) Less(i, j int) bool {
	return entries[i].Start.After(entries[j].Start)
}
