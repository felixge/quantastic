package main

import (
	"fmt"
	"os"
	"sort"
	"time"
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
	options := [][]string{}
	if sinceStr, ok := flag("since", c.Args); ok {
		since, err := StringToTime(sinceStr)
		if err != nil {
			fatal("Invalid since value: %s", err)
		}
		options = append(options, []string{"SINCE:", TimeToString(&since)})
		for i, entry := range entries {
			if entry.Start.Before(since) {
				entries = entries[0:i]
				break
			}
		}
	}
	if untilStr, ok := flag("until", c.Args); ok {
		until, err := StringToTime(untilStr)
		if err != nil {
			fatal("Invalid until value: %s", err)
		}
		options = append(options, []string{"UNTIL:", TimeToString(&until)})
		for len(entries) > 0 {
			if entries[0].Start.After(until) {
				entries = entries[1:]
			} else {
				break
			}
		}
	}
	if len(options) > 0 {
		mustWriteTable(os.Stdout, options)
		fmt.Fprintf(os.Stdout, "\n")
	}
	if ok := boolFlag("group", c.Args); ok {
		m := map[string]time.Duration{}
		for _, entry := range entries {
			category := CategoryToString(entry.Category)
			m[category] += entry.Duration()
		}
		groups := groupTimeEntries{}
		for category, duration := range m {
			groups = append(groups, &groupTimeEntry{Category: category, Duration: duration})
		}
		sort.Sort(groups)
		data := make([][]string, 0, len(groups)+1)
		data = append(data, []string{"CATEGORY", "DURATION"})
		for _, group := range groups {
			data = append(data, []string{
				group.Category,
				DurationToString(group.Duration),
			})
		}
		mustWriteTable(os.Stdout, data)
		return
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

type groupTimeEntries []*groupTimeEntry

func (entries groupTimeEntries) Len() int {
	return len(entries)
}

func (entries groupTimeEntries) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

func (entries groupTimeEntries) Less(i, j int) bool {
	return entries[i].Duration > entries[j].Duration
}

type groupTimeEntry struct{
	Category string
	Duration time.Duration
}
