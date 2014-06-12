package main

import (
	"time"
)

var cmdTimeEnd = &command{
	name:        "time end",
	description: "Finish an entry.",
	usage:       "",
	fn:          cmdTimeEndFn,
}

func cmdTimeEndFn(c *Context) {
	entry, err := c.Db.LatestTimeEntry()
	if err != nil {
		fatal("Failed to get latest entry: %s", err)
	}
	if entry.End != nil {
		fatal("Latest entry is already finished.")
	}
	now := time.Now()
	entry.End = &now
	if err := c.Db.SaveTimeEntry(entry); err != nil {
		fatal("Could not save entry: %s", err)
	}
}
