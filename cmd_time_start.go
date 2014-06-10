package main

import (
	"time"
)

var cmdTimeStart = &command{
	name:        "time start",
	description: "Create a new entry.",
	usage:       "<category>",
	fn:          cmdTimeStartFn,
}

func cmdTimeStartFn(c *Context) {
	if len(c.Args) == 0 {
		fatal("%s", c.Cmd.Usage())
	}
	category := StringToCategory(c.Args[0])
	activeEntry, err := c.Db.ActiveTimeEntry()
	now := time.Now()
	if err != nil {
		if _, ok := err.(NotFoundError); !ok {
			fatal("Failed to get latest entry: %s", err)
		}
	} else {
		activeEntry.End = &now
		if err := c.Db.SaveTimeEntry(activeEntry); err != nil {
			fatal("Failed to set end for active entry: %s", err)
		}
	}
	entry := &TimeEntry{
		Category: category,
		Start:    now,
	}
	if err := entry.Valid(); err != nil {
		fatal("Invalid entry: %s", err)
	}
	if err := c.Db.SaveTimeEntry(entry); err != nil {
		fatal("Failed to save entry: %s", err)
	}
}
