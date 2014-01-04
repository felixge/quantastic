// Package model contains the model subsystem of quantastic. The model is an
// in-memory representation of all data in the system. Because there could
// potentially be a lot of data in the system, the model usually only gets
// populated with enough data to perform a certain task. The model MUST NOT be
// aware of storage or representation.
package model

// @TODO Figure out if ids should go into the model or the api/storage
// subsystem.
// @TODO Figure out how to deal with TimeEntries that fall into a range, but
// start before and/or end after it.
// @TODO TimeRange has TimeEntries which themselves have a time range. Meta!

import (
	"time"
)

// @TODO Write test
func Day(day time.Time) *TimeRange {
	return &TimeRange{
		Start: dayStart(day),
		End:   dayEnd(day),
	}
}

func dayStart(day time.Time) time.Time {
	h, m, s := day.Clock()
	return day.Add(-(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second))
}

func dayEnd(day time.Time) time.Time {
	h, m, s := day.Clock()
	return day.Add((23 - time.Duration(h)*time.Hour) + (59 - time.Duration(m)*time.Minute) + (59 - time.Duration(s)*time.Second))
}

type TimeRange struct {
	Start   time.Time
	End     time.Time
	Entries []*TimeEntry
	Root    *TimeCategory
}

func (t *TimeRange) Duration() time.Duration {
	return t.End.Sub(t.Start)
}

type TimeEntry struct {
	Id       string
	Note     string
	Created  time.Time
	From     time.Time
	To       time.Time
	Category *TimeCategory
}

func (t *TimeEntry) Duration() time.Duration {
	return t.To.Sub(t.From)
}

type TimeCategory struct {
	Id       string
	Name     string
	Parent   *TimeCategory
	Children []*TimeCategory
	Entries  []*TimeEntry
}

// Root returns if the category is the root category (i.e. it has no parent)
func (c *TimeCategory) Root() bool {
	return c.Parent == nil
}

// AllDuration returns the sum of durations of all time entries in the category
// as well as all of its children.
func (c *TimeCategory) AllDuration() (d time.Duration) {
	for _, category := range c.All() {
		d = d + category.Duration()
	}
	return
}

// Duration returns the sum of durations of all time entries in the category
// itself.
func (c *TimeCategory) Duration() (d time.Duration) {
	for _, entry := range c.Entries {
		d = d + entry.Duration()
	}
	return
}

// All returns a list of categories including itself, as well as all of its
// children (recursively). This is useful when one needs to work with
// categories as a list, rather than tree.
func (c *TimeCategory) All() []*TimeCategory {
	results := make([]*TimeCategory, 0, 1+len(c.Children))
	results = append(results, c)
	for _, child := range c.Children {
		for _, childCategory := range child.All() {
			results = append(results, childCategory)
		}
	}
	return results
}

// Chain a list starting with the root category and containing all categories
// leading up to the current category. The list ends with the category itself.
func (c *TimeCategory) Chain() []*TimeCategory {
	chain := make([]*TimeCategory, 0, 1)
	for c != nil {
		chain = append(chain, c)
		c = c.Parent
	}
	reverse := make([]*TimeCategory, len(chain))
	for i, category := range chain {
		reverse[len(chain)-i-1] = category
	}
	return reverse
}
