package mite

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	pkgurl "net/url"
	"strings"
	"time"
)

const (
	maxResponseBody = 2 * 1024 * 1024
	dateFormat      = "2006-01-02"
)

func NewMite(url, apiKey string) (*Mite, error) {
	parsedUrl, err := pkgurl.Parse(url)
	if err != nil {
		return nil, err
	}
	return &Mite{url: *parsedUrl, apiKey: apiKey}, nil
}

type Mite struct {
	url    pkgurl.URL
	apiKey string
}

// http://mite.yo.lk/en/api/time-entries.html
func (m *Mite) TimeEntries(query *TimeEntriesQuery) ([]TimeEntry, error) {
	url := m.url
	url.Path = "/time_entries.xml"
	if query != nil {
		url.RawQuery = query.query.Encode()
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-MiteApiKey", m.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(res.Body, maxResponseBody))
	if err != nil {
		return nil, err
	}
	entries := timeEntries{}
	if err := xml.Unmarshal(body, &entries); err != nil {
		errs := errors{}
		if err2 := xml.Unmarshal(body, &errs); err2 == nil {
			return nil, fmt.Errorf(strings.Join(errs.Errors, ", "))
		}
		return nil, err
	}

	results := make([]TimeEntry, len(entries.Entries))
	for i, entry := range entries.Entries {
		results[i] = entry.TimeEntry
		results[i].DateAt = time.Time(entry.DateAt)
	}

	return results, nil
}

type errors struct {
	XMLName xml.Name `xml:"errors"`
	Errors  []string `xml:"error"`
}

type timeEntries struct {
	XMLName xml.Name    `xml:"time-entries"`
	Entries []timeEntry `xml:"time-entry"`
}

type timeEntry struct {
	TimeEntry
	DateAt date `xml:"date-at"`
}

type date time.Time

func (d *date) UnmarshalText(text []byte) error {
	parsed, err := time.Parse(dateFormat, string(text))
	if err != nil {
		return err
	}
	*d = date(parsed)
	return nil
}

type TimeEntry struct {
	Id           int       `xml:"id"`
	DateAt       time.Time `xml:"date-at"`
	Minutes      string    `xml:"minutes"`
	Revenue      float64   `xml:"revenue"`
	Billable     bool      `xml:"billable"`
	Note         string    `xml:"note"`
	UserId       int       `xml:"user-id"`
	UserName     string    `xml:"user-name"`
	ProjectId    int       `xml:"project-id"`
	ProjectName  string    `xml:"project-name"`
	ServiceId    int       `xml:"service-id"`
	ServiceName  string    `xml:"service-name"`
	CustomerId   int       `xml:"customer-id"`
	CustomerName string    `xml:"customer-name"`
	Locked       bool      `xml:"locked"`
	CreatedAt    time.Time `xml:"created-at"`
	UpdatedAt    time.Time `xml:"updated-at"`
}

func NewTimeEntriesQuery() *TimeEntriesQuery {
	return &TimeEntriesQuery{query: make(pkgurl.Values)}
}

type TimeEntriesQuery struct {
	query pkgurl.Values
}

func (q *TimeEntriesQuery) SetCustomerId(val int) {
	q.query.Set("customer_id", intToStr(val))
}

func (q *TimeEntriesQuery) SetProjectId(val int) {
	q.query.Set("project_id", intToStr(val))
}

func (q *TimeEntriesQuery) SetServiceId(val int) {
	q.query.Set("service_id", intToStr(val))
}

func (q *TimeEntriesQuery) SetUserId(val int) {
	q.query.Set("user_id", intToStr(val))
}

func (q *TimeEntriesQuery) SetBillable(val bool) {
	q.query.Set("billable", boolToStr(val))
}

func (q *TimeEntriesQuery) SetNote(val string) {
	q.query.Set("note", val)
}

func (q *TimeEntriesQuery) SetAt(val time.Time) {
	q.query.Set("at", dateToStr(val))
}

func (q *TimeEntriesQuery) SetFrom(val time.Time) {
	q.query.Set("from", dateToStr(val))
}

func (q *TimeEntriesQuery) SetTo(val time.Time) {
	q.query.Set("to", dateToStr(val))
}

func (q *TimeEntriesQuery) SetLocked(val bool) {
	q.query.Set("locked", boolToStr(val))
}

func intToStr(val int) string {
	return fmt.Sprintf("%d", val)
}

func boolToStr(val bool) string {
	return fmt.Sprintf("%t", val)
}

func dateToStr(val time.Time) string {
	return val.Format(dateFormat)
}
