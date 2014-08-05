package main

import (
	"strings"

	"github.com/kylelemons/godebug/pretty"
)
import "testing"

func Test_readEntries(t *testing.T) {
	tests := []struct {
		Input    string
		Expected []*EditedEntry
	}{
		{
			Input:    "",
			Expected: nil,
		},
		{
			Input: "2014-08-05 08:35:18 , 2014-08-05 09:22:51 , Work:Quantastic , 8ff6d7c02427e9a4e5e1e9931ef15a55\n",
			Expected: []*EditedEntry{
				{
					Id:       "8ff6d7c02427e9a4e5e1e9931ef15a55",
					Start:    "2014-08-05 08:35:18",
					End:      "2014-08-05 09:22:51",
					Category: "Work:Quantastic",
				},
			},
		},
		{
			Input: "2014-08-05 08:35:18 , 2014-08-05 09:22:51 , Work:Quantastic , 8ff6d7c02427e9a4e5e1e9931ef15a55\n" +
				"2014-09-05 03:21:92 , 2014-09-05 04:10:05 , Sports:Beach Volleyball , 4e0a6d93a0b589cac27d4e5bc0eb90bb\n",
			Expected: []*EditedEntry{
				{
					Id:       "8ff6d7c02427e9a4e5e1e9931ef15a55",
					Start:    "2014-08-05 08:35:18",
					End:      "2014-08-05 09:22:51",
					Category: "Work:Quantastic",
				},
				{
					Id:       "4e0a6d93a0b589cac27d4e5bc0eb90bb",
					Start:    "2014-09-05 03:21:92",
					End:      "2014-09-05 04:10:05",
					Category: "Sports:Beach Volleyball",
				},
			},
		},
	}
	for i, test := range tests {
		entries, err := readEntries(strings.NewReader(test.Input))
		if err != nil {
			t.Errorf("test %d: %s", i, err)
			continue
		}
		if diff := pretty.Compare(entries, test.Expected); diff != "" {
			t.Errorf("test %d: %s", i, diff)
			continue
		}
	}
}
