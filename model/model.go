// Package model contains the model subsystem of quantastic. The model is an
// in-memory representation of all data in the system. Because there could
// potentially be a lot of data in the system, the model usually only gets
// partially populated with enough data from the data stores to perform a
// certain task. The model also contains all rules applying to the data it
// holds.
package model

import (
	"time"
)

type TimeEntry struct {
	Id       string
	Category *TimeCategory
	Note     string
	Created  time.Time
	From     time.Time
	To       time.Time
	Duration time.Duration
}

type TimeCategory struct {
	Id       string
	Parent   *TimeCategory
	Children []*TimeCategory
	Entries  []*TimeEntry
	Name     string
}

// Categories returns a list of categories including itself, as well as all of
// its children (recursively). This is useful when one needs to work with
// categories as a list, rather than tree.
func (c *TimeCategory) Categories() []*TimeCategory {
	results := make([]*TimeCategory, 0, 1+len(c.Children))
	results = append(results, c)
	for _, child := range c.Children {
		for _, childCategory := range child.Categories() {
			results = append(results, childCategory)
		}
	}
	return results
}

// Hierarchy a list starting with the root category and containing all
// categories leading up to the current category. The list ends with the
// category itself.
func (c *TimeCategory) Hierarchy() []*TimeCategory {
	hierarchy := make([]*TimeCategory, 0, 1)
	parent := c
	for parent != nil {
		hierarchy = append(hierarchy, parent)
		parent = parent.Parent
	}
	reverse := make([]*TimeCategory, len(hierarchy))
	for i, category := range hierarchy {
		reverse[len(hierarchy)-i-1] = category
	}
	return reverse
}
