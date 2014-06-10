package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	categorySeparator = ":"
	timeLayout        = "2006-01-02 15:04:05"
)

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func StringToCategory(category string) []string {
	parts := strings.Split(category, categorySeparator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		result = append(result, strings.TrimSpace(part))
	}
	return result
}

func CategoryToString(category []string) string {
	return strings.Join(category, categorySeparator)
}

func TimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Local().Format(timeLayout)
}

func StringToTime(s string) (time.Time, error) {
	return time.ParseInLocation(timeLayout, s, time.Local)
}

func UtcNow() time.Time {
	return time.Now().UTC()
}

func NoteExcerpt(note string) string {
	note = strings.Replace(note, "\n", "", -1)
	// iterate the string rune wise to avoid splitting a multibyte char
  var c int
	for i, _ := range note {
		c++
		if c == 20 {
			return note[0:i]
		}
	}
	return note
}

func DurationToString(d time.Duration) string {
	hours := d / time.Hour
	d -= hours * time.Hour
	min := d / time.Minute
	d -= min * time.Minute
	sec := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", hours, min, sec)
}

func StringToDuration(s string) (time.Duration, error) {
	var errInvalid = fmt.Errorf("Invalid duration format: %s", s)
	var d time.Duration
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return d, errInvalid
	}
	dParts := []time.Duration{}
	for _, part := range parts {
		i64, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			return d, errInvalid
		}
		dParts = append(dParts, time.Duration(i64))
	}
	fmt.Printf("dParts: %#v\n", dParts)
	d = time.Hour * dParts[0] + time.Minute * dParts[1] * time.Second * dParts[2]
	return d, nil
}

func writeTable(w io.Writer, data [][]string) (err error) {
	lengths := make([]int, len(data[0]))
	for _, row := range data {
		for i, val := range row {
			if l := len(val); l > lengths[i] {
				lengths[i] = l
			}
		}
	}
	bw := bufio.NewWriter(w)
	defer func() {
		if err == nil {
			err = bw.Flush()
		} else {
			bw.Flush()
		}
	}()
	for _, row := range data {
		for i, val := range row {
			last := i+1 == len(row)
			if pad := lengths[i] - len(val); pad > 0 && !last {
				val += strings.Repeat(" ", pad)
			}
			_, err = io.WriteString(bw, val)
			if err != nil {
				return
			}
			if last {
				_, err = io.WriteString(bw, "\n")
			} else {
				_, err = io.WriteString(bw, " ")
			}
			if err != nil {
				return
			}
		}
	}
	return
}

func mustPrintfStdout(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(os.Stdout, format, args...); err != nil {
		panic(err)
	}
}

func mustWriteTable(w io.Writer, data [][]string) {
	if err := writeTable(w, data); err != nil {
		panic(err)
	}
}
