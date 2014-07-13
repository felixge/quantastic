package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/felixge/asciitable"
	"github.com/felixge/pager"
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
	p, err := pager.Start("less")
	if err != nil {
		fatal("Failed to execute pager: %s", err)
	}
	defer p.Wait()
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

	if ok := boolFlag("table", c.Args); ok {
		table := asciitable.NewTable()
		table.AddRow("Category", "Start", "End", "Duration", "Note", "Id")
		var total time.Duration
		var prevDay string
		for _, entry := range entries {
			day := entry.Start.Format("2006-01-02")
			if day != prevDay {
				table.AddSeparator()
				table.AddRow("", day, "", "", entry.Start.Format("Monday"))
				table.AddSeparator()
			}
			prevDay = day
			total += entry.Duration()
			table.AddRow(
				CategoryToString(entry.Category),
				TimeString(&entry.Start),
				TimeString(entry.End),
				DurationToString(entry.Duration()),
				NoteExcerpt(entry.Note),
				entry.Id,
			)
		}
		table.AddSeparator()
		table.AddRow("Total", "", "", DurationToString(total), "", "")
		table.Fprint(os.Stdout)
	} else {
		stdout := bufio.NewWriter(os.Stdout)
		for i, entry := range entries {
			if i > 0 {
				fmt.Fprintf(stdout, "\n")
			}
			data := [][]string{
				[]string{"Id:", entry.Id},
				[]string{"Category:", CategoryToString(entry.Category)},
				[]string{"Start:", TimeToString(&entry.Start)},
				[]string{"End:", TimeToString(entry.End)},
				[]string{"Duration:", DurationToString(entry.Duration())},
			}
			mustWriteTable(stdout, data)
			if entry.Note != "" {
				fmt.Fprintf(stdout, "\n%s\n", prefixLines(entry.Note, "    "))
			}
		}
		stdout.Flush()
	}
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

type groupTimeEntry struct {
	Category string
	Duration time.Duration
}
