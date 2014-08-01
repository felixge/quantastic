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
	days := map[int]time.Duration{}
	for _, entry := range entries {
		if !CategoryEqual(entry.Category, category) {
			continue
		}
		if entry.Start.Year() != month.Year() || entry.Start.Month() != month.Month() {
			continue
		}
		days[entry.Start.Day()] += entry.Duration()
	}
	table := asciitable.NewTable()
	table.AddRow("Day", "Hours")
	table.AddSeparator()
	var total time.Duration
	for day := 1; day <= daysInMonth(month); day++ {
		dayTime := time.Date(month.Year(), month.Month(), day, 0, 0, 0, 0, time.Local)
		duration := days[day]
		//remainder := duration % (15 * time.Minute)
		//duration -= remainder
		//if remainder > 3*time.Minute {
		//duration += 15 * time.Minute
		//}
		total += duration
		table.AddRow(
			dayTime.Format("2006-01-02 (Mon)"),
			DurationToString(duration),
		)
	}
	table.AddSeparator()
	table.AddRow("Total", DurationToString(total))
	fmt.Printf("%s\n", table)
	//startStr := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month()-1, 1)
	//firstDay, err := time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), 1))
	//if err != nil {
	//fatal("Could not get first day of current month: %s", err)
	//}
	//num

	//if err != nil {
	//fatal("Could not get first day of next month: %s", err)
	//}

	//category := StringToCategory(c.Args[0])
	//month := time.Now().Add(-time.Month)
	//day := 1
	//days := map[int]time.Duration{}
	//daysInMonth := firstDayNextMonth.Add(-time.Hour).Day()
	//for _, entry := range entries {
	//if month != entry.Start.Month() {
	//continue
	//}
	//if !CategoryEqual(category, entry.Category) {
	//continue
	//}
	//days[entry.Start.Day()] += entry.Duration()
	//}
	//table := [][]string{{"DAY", "DURATION"}}
	//for day := 1; day < daysInMonth; day++ {
	//dayStr := fmt.Sprintf("%04d-%02d-%02d", mo)
	//table = append(table, []string{dayStr, DurationToString(days[day])})
	//}
	//mustWriteTable(os.Stdout, table)
}

func lastMonth() time.Time {
	now := time.Now()
	month := now.Month() - 1
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
