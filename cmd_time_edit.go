package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/felixge/quantastic/db"
)

var cmdTimeEdit = &command{
	name:        "time edit",
	description: "Edit an entry.",
	usage: `[<id>]

If <id> is not given, it defaults to the latest entry.`,
	fn: cmdTimeEditFn,
}

var fieldRegExp = regexp.MustCompile("^([^:]+):\\s*(.*?)\\s*$")

func cmdTimeEditFn(c *Context) {
	var entry *db.TimeEntry
	var err error
	if len(c.Args) == 0 {
		entry, err = c.Db.LatestTimeEntry()
	} else if strings.HasPrefix(c.Args[0], "~") {
		var num int
		_, err := fmt.Sscanf(c.Args[0], "~%d", &num)
		if err != nil {
			fatal("Bad argument: %s", c.Args[0])
		}
		entries, err := c.Db.TimeEntriesByStart()
		if err != nil {
			fatal("Could not load entries: %s", err)
		}
		if num > len(entries) {
			num = len(entries)
		}
		editEntries(c, entries[0:num])
		return
	} else {
		entry, err = c.Db.TimeEntry(c.Args[0])
	}
	if err != nil {
		fatal("Could not load entry: %s", err)
	}
	editEntry(c, entry)
}

func editEntry(c *Context, entry *db.TimeEntry) {
	table := [][]string{
		{"CATEGORY:", CategoryToString(entry.Category)},
		{"START:", TimeToString(&entry.Start)},
		{"END:", TimeToString(entry.End)},
	}
	editor := NewEditor()
	editor.TmpPrefix = "time-entry-" + entry.Id
	defer editor.Close()
	err := writeTable(editor, table)
	if err == nil {
		if entry.Note == "" {
			_, err = fmt.Fprintf(editor, "\n# Insert note here\n")
		} else {
			_, err = fmt.Fprintf(editor, "\n%s\n", entry.Note)
		}
	}
	if err != nil {
		fatal("Could not write to editor: %s", err)
	}
	if err := editor.Run(); err != nil {
		fatal("Failed to run editor %s: %s", editor.Command, err)
	}
	editedEntry, err := readEntry(editor)
	if err != nil {
		fatal("%s", err)
	}
	if err := applyEditedEntry(entry, editedEntry, true); err != nil {
		fatal("%s", err)
	}
	if err := entry.Valid(); err != nil {
		fatal("%s", err)
	}
	if err := c.Db.SaveTimeEntry(entry); err != nil {
		fatal("Failed to save entry: %s", err)
	}
}

func applyEditedEntry(entry *db.TimeEntry, edited *EditedEntry, note bool) error {
	start, err := StringToTime(edited.Start)
	if err != nil {
		return fmt.Errorf("Bad Start: %s", edited.Start)
	}
	entry.Start = start
	if edited.End != "" {
		end, err := StringToTime(edited.End)
		if err != nil {
			return fmt.Errorf("Bad End: %s", edited.End)
		}
		entry.End = &end
	} else {
		entry.End = nil
	}
	entry.Category = StringToCategory(edited.Category)
	if note {
		entry.Note = edited.Note
	}
	return nil
}

func editEntries(c *Context, entries []*db.TimeEntry) {
	editor := NewEditor()
	editor.TmpPrefix = "time-entries"
	defer editor.Close()
	table := [][]string{}
	for _, entry := range entries {
		row := []string{
			TimeToString(&entry.Start),
			", " + TimeToString(entry.End),
			", " + CategoryToString(entry.Category),
			", " + entry.Id,
		}
		table = append(table, row)
	}
	if err := writeTable(editor, table); err != nil {
		fatal("Could not write to editor: %s", err)
	}
	if err := editor.Run(); err != nil {
		fatal("Failed to run editor %s: %s", editor.Command, err)
	}
	editedEntries, err := readEntries(editor)
	if err != nil {
		fatal("Failed to read entries: %s", err)
	}
	for _, editedEntry := range editedEntries {
		found := false
		for _, entry := range entries {
			if entry.Id != editedEntry.Id {
				continue
			}
			found = true
			if err := applyEditedEntry(entry, editedEntry, false); err != nil {
				fatal("%s", err)
			}
			if err := entry.Valid(); err != nil {
				fatal("%s", err)
			}
			if err := c.Db.SaveTimeEntry(entry); err != nil {
				fatal("Failed to save entry: %s", err)
			}
			break
		}
		if !found {
			entry := &db.TimeEntry{}
			if err := applyEditedEntry(entry, editedEntry, false); err != nil {
				fatal("%s", err)
			}
			if err := entry.Valid(); err != nil {
				fatal("%s", err)
			}
			if err := c.Db.SaveTimeEntry(entry); err != nil {
				fatal("Failed to save entry: %s", err)
			}
		}
	}
	// Delete entries if they were removed
	for _, entry := range entries {
		found := false
		for _, edited := range editedEntries {
			if edited.Id == entry.Id {
				found = true
				break
			}
		}
		if !found {
			if err := c.Db.DeleteTimeEntry(entry.Id); err != nil {
				fatal("%s", err)
			}
		}
	}
}

type EditedEntry struct {
	Id       string
	Category string
	Start    string
	End      string
	Note     string
}

func readEntry(r io.Reader) (*EditedEntry, error) {
	entry := &EditedEntry{}
	scanner := bufio.NewScanner(r)
	noteSection := false
	for scanner.Scan() {
		line := scanner.Text()
		matches := fieldRegExp.FindStringSubmatch(line)
		if noteSection {
			if !strings.HasPrefix(line, "#") {
				entry.Note += line + "\n"
			}
			continue
		} else if len(matches) != 3 {
			if line != "" {
				return nil, fmt.Errorf("Bad line: %s", line)
			} else {
				noteSection = true
				continue
			}
		}
		field, val := matches[1], matches[2]
		switch strings.ToLower(field) {
		case "category":
			entry.Category = val
		case "start":
			entry.Start = val
		case "end":
			entry.End = val
		default:
			return nil, fmt.Errorf("Unknown field: %s", field)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Failed to scan: %s", err)
	}
	entry.Note = strings.TrimSpace(entry.Note)
	return entry, nil
}

func readEntries(r io.Reader) ([]*EditedEntry, error) {
	var (
		cr      = csv.NewReader(r)
		entries []*EditedEntry
	)
	for {
		row, err := cr.Read()
		if err == io.EOF {
			return entries, nil
		} else if err != nil {
			return nil, err
		}
		if len(row) != 4 {
			return nil, fmt.Errorf("Bad row: %s", strings.Join(row, ","))
		}
		for i, val := range row {
			row[i] = strings.TrimSpace(val)
		}
		entries = append(entries, &EditedEntry{
			Start:    row[0],
			End:      row[1],
			Category: row[2],
			Id:       row[3],
		})
	}
}
