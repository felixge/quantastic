package main

import (
	"fmt"
	"time"

	"github.com/felixge/asciitable"
)

var cmdTimeInvoice = &command{
	name:        "time invoice",
	description: "Invoice client.",
	usage:       "<category>",
	fn:          cmdTimeInvoiceFn,
}

func cmdTimeInvoiceFn(c *Context) {
	entries, err := c.Db.TimeEntriesByStart()
	if err != nil {
		fatal("Could not get entries: %s", err)
	}
	if len(c.Args) < 1 {
		fatal("Missing category argument.")
	}
	category := StringToCategory(c.Args[0])
	if err != nil {
		fatal("Invalid category: %s", err)
	}
	month := lastMonth()
	if monthStr, ok := flag("month", c.Args); ok {
		month, err = time.Parse("2006-01-02", monthStr)
		if err != nil {
			fatal("Invalid month date: %s", err)
		}
	}
	type InvoiceDay struct {
		Duration time.Duration
		Notes    string
	}
	days := map[int]InvoiceDay{}
	for _, entry := range entries {
		if !CategoryEqual(entry.Category, category) {
			continue
		}
		if entry.Start.Year() != month.Year() || entry.Start.Month() != month.Month() {
			continue
		}
		invoiceDay := days[entry.Start.Day()]
		invoiceDay.Duration += entry.Duration()
		if entry.Note != "" {
			invoiceDay.Notes += entry.Note + "\n"
		}
		days[entry.Start.Day()] = invoiceDay
	}
	table := asciitable.NewTable()
	table.AddRow("Day", "Hours", "Total")
	table.AddSeparator()
	var total time.Duration
	for day := 1; day <= daysInMonth(month); day++ {
		dayTime := time.Date(month.Year(), month.Month(), day, 0, 0, 0, 0, time.Local)
		invoiceDay := days[day]
		total += invoiceDay.Duration
		table.AddRow(
			dayTime.Format("2006-01-02 (Mon)"),
			DurationToString(invoiceDay.Duration),
			DurationToString(total),
		)
	}
	fmt.Printf("%s\n\n", table)
	for day := 1; day <= daysInMonth(month); day++ {
		dayTime := time.Date(month.Year(), month.Month(), day, 0, 0, 0, 0, time.Local)
		invoiceDay := days[day]
		if invoiceDay.Notes != "" {
			fmt.Printf(
				"%s - %s\n\n%s\n",
				dayTime.Format("2006-01-02 (Mon)"),
				DurationToString(invoiceDay.Duration),
				prefixLines(invoiceDay.Notes, "    "),
			)
		}
	}
}

func lastMonth() time.Time {
	now := time.Now()
	month := now.Month()
	if month == 0 {
		month = 12
	}
	return time.Date(now.Year(), month, 1, 0, 0, 0, 0, time.Local)
}

func daysInMonth(t time.Time) int {
	month := t.Month() + 1
	if month == 13 {
		month = 1
	}
	nextMonth := time.Date(t.Year(), month, 1, 0, 0, 0, 0, time.Local)
	return nextMonth.Add(-1 * time.Hour).Day()
}
