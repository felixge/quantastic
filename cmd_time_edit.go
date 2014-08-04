package main

import (
	"bufio"
	"bytes"
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
	entry.Category = StringToCategory(editedEntry.Category)
	entry.Note = editedEntry.Note
	if entry.Start, err = StringToTime(editedEntry.Start); err != nil {
		fatal("Bad Start: %s", err)
	}
	if editedEntry.End == "" {
		entry.End = nil
	} else if t, err := StringToTime(editedEntry.End); err != nil {
		fatal("Bad End: %s", err)
	} else {
		entry.End = &t
	}
	if err := entry.Valid(); err != nil {
		fatal("%s", err)
	}
	if err := c.Db.SaveTimeEntry(entry); err != nil {
		fatal("Failed to save entry: %s", err)
	}
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
	for _, entry := range editedEntries {
		fmt.Printf("entry: %#v\n", entry)
	}
}

type EditedEntry struct {
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
	type readState int
	const (
		argStart readState = iota + 1
		argVal
		argEnd
		argQuoteVal
		argQuoteEnd
		rowEnd
	)
	var (
		state = argStart
		br    = bufio.NewReader(r)
		arg   = bytes.NewBuffer(nil)
		debug = bytes.NewBuffer(nil)
		row   []string
		rows  [][]string
	)
	for {
		c, err := br.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		debug.WriteByte(c)
	statemachine:
		switch state {
		case argStart:
			if c == '"' {
				state = argQuoteVal
			} else {
				arg.WriteByte(c)
				state = argVal
			}
		case argVal:
			switch c {
			case ' ':
				state = argEnd
				goto statemachine
			case '\n':
				state = rowEnd
				goto statemachine
			default:
				arg.WriteByte(c)
			}
		case argEnd:
			row = append(row, strings.TrimSpace(arg.String()))
			arg.Reset()
			state = argStart
		case argQuoteVal:
			if c == '"' {
				state = argQuoteEnd
			} else {
				arg.WriteByte(c)
			}
		case argQuoteEnd:
			switch c {
			case ' ':
				state = argEnd
				goto statemachine
			case '\n':
				state = rowEnd
				goto statemachine
			default:
				return nil, fmt.Errorf("Invalid character after quote: %s", debug)
			}
		case rowEnd:
			fmt.Printf("row: %s\n", row)
			rows = append(rows, row)
			row = nil
			state = argStart
		}
	}
	return nil, nil
}
