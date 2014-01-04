package model

import (
	"github.com/felixge/quantastic/testutil"
	"testing"
)

func TestCategory_Categories(t *testing.T) {
	root := &Category{
		Name: "A",
		Children: []*Category{
			{
				Name: "B",
				Children: []*Category{
					{Name: "C"},
					{Name: "D"},
				},
			},
			{
				Name: "E",
				Children: []*Category{
					{Name: "F"},
				},
			},
			{
				Name: "G",
			},
		},
	}

	expected := []string{"A", "B", "C", "D", "E", "F", "G"}
	got := []string{}
	for _, category := range root.Categories() {
		got = append(got, category.Name)
	}
	if err := testutil.DeepEqual(got, expected); err != nil {
		t.Error(err)
	}
}

func TestCategory_Hierarchy(t *testing.T) {
	leaf := &Category{
		Name: "C",
		Parent: &Category{
			Name: "B",
			Parent: &Category{
				Name: "A",
			},
		},
	}

	expected := []string{"A", "B", "C"}
	got := []string{}
	for _, category := range leaf.Hierarchy() {
		got = append(got, category.Name)
	}
	if err := testutil.DeepEqual(got, expected); err != nil {
		t.Error(err)
	}
}
