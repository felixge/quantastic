package model

import (
	"github.com/felixge/quantastic/util"
	"testing"
)

func TestTimeCategory_All(t *testing.T) {
	root := &TimeCategory{
		Name: "A",
		Children: []*TimeCategory{
			{
				Name: "B",
				Children: []*TimeCategory{
					{Name: "C"},
					{Name: "D"},
				},
			},
			{
				Name: "E",
				Children: []*TimeCategory{
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
	for _, category := range root.All() {
		got = append(got, category.Name)
	}
	if err := util.DeepEqual(got, expected); err != nil {
		t.Error(err)
	}
}

func TestTimeCategory_Chain(t *testing.T) {
	leaf := &TimeCategory{
		Name: "C",
		Parent: &TimeCategory{
			Name: "B",
			Parent: &TimeCategory{
				Name: "A",
			},
		},
	}

	expected := []string{"A", "B", "C"}
	got := []string{}
	for _, category := range leaf.Chain() {
		got = append(got, category.Name)
	}
	if err := util.DeepEqual(got, expected); err != nil {
		t.Error(err)
	}
}
