package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var cmdTimeEdit = &command{
	name:        "time edit",
	description: "Edit an entry.",
	usage:       "<id>",
	fn:          cmdTimeEditFn,
}

var fieldRegExp = regexp.MustCompile("^([^:]+):\\s*(.*?)\\s*$")

func cmdTimeEditFn(c *Context) {
	if len(c.Args) == 0 {
		fatal("%s", c.Cmd.Usage())
	}
	entry, err := c.Db.TimeEntry(c.Args[0])
	if err != nil {
		fatal("Could not load entry: %s", err)
	}
	table := [][]string{
		{"CATEGORY:", CategoryToString(entry.Category)},
		{"START:", TimeToString(&entry.Start)},
		{"END:", TimeToString(entry.End)},
	}
	file, err := ioutil.TempFile("", "time-entry-"+entry.Id)
	if err != nil {
		fatal("Could not open tmp file: %s", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	err = writeTable(file, table)
	if err == nil {
		if entry.Note == "" {
			_, err = fmt.Fprintf(file, "\n# Insert note here\n")
		} else {
			_, err = fmt.Fprintf(file, "\n%s\n", entry.Note)
		}
	}
	if err != nil {
		fatal("Could not write to tmp file: %s", err)
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		fatal("No EDITOR env var configured.")
	}

	cmd := exec.Command("/bin/sh", "-c", editor+" '"+file.Name()+"'")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("Failed to run editor %s: %s", editor, err)
	}
	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		fatal("Failed to seek tmp file: %s", err)
	}
	editedEntry, err := readEntry(file)
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
