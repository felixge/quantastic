package main

import (
	"os"
)

var cmdTimeLog = &command{
	name:        "time log",
	description: "Show entries.",
	usage:       "",
	fn:          cmdTimeLogFn,
}

func cmdTimeLogFn(c *Context) {
	entries, err := c.Db.TimeEntriesByStart()
	if err != nil {
		fatal("Could not get entries: %s", err)
	}
	data := make([][]string, 0, len(entries)+1)
	data = append(data, []string{"CATEGORY", "START", "END", "DURATION", "NOTE", "ID"})
	for _, entry := range entries {
		data = append(data, []string{
			CategoryToString(entry.Category),
			TimeToString(&entry.Start),
			TimeToString(entry.End),
			DurationToString(entry.Duration()),
			NoteExcerpt(entry.Note),
			entry.Id,
		})
	}
	mustWriteTable(os.Stdout, data)
}
