package mite

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_url "net/url"
	"strings"
	"time"
)

// @TODO make sure passing "nil" for any query argument works

const (
	maxResponseBody = 2 * 1024 * 1024
	dateFormat      = "2006-01-02"
)

func NewClient(url, apiKey string) (*Client, error) {
	parsedUrl, err := _url.Parse(url)
	if err != nil {
		return nil, err
	}
	return &Client{url: *parsedUrl, apiKey: apiKey}, nil
}

type Client struct {
	url    _url.URL
	apiKey string
}

func (c *Client) get(path string, query _url.Values, root interface{}) error {
	url := c.url
	url.Path = path
	if query != nil {
		url.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-MiteApiKey", c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(res.Body, maxResponseBody))
	if err := xml.Unmarshal(body, root); err != nil {
		root := errorsRoot{}
		if err := xml.Unmarshal(body, &root); err == nil {
			return fmt.Errorf(strings.Join(root.Errors, ", "))
		}
		return err
	}
	return nil
}

// http://mite.yo.lk/en/api/time-entries.html
func (c *Client) TimeEntries(query *TimeEntriesQuery) ([]TimeEntry, error) {
	var q _url.Values
	if query != nil {
		q = query.query
	}

	root := timeEntriesRoot{}
	if err := c.get("/time_entries.xml", q, &root); err != nil {
		return nil, err
	}

	results := make([]TimeEntry, len(root.Entries))
	for i, entry := range root.Entries {
		results[i] = entry.TimeEntry
		results[i].DateAt = time.Time(entry.DateAt)
	}

	return results, nil
}

type errorsRoot struct {
	XMLName xml.Name `xml:"errors"`
	Errors  []string `xml:"error"`
}

type timeEntriesRoot struct {
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
	Minutes      int       `xml:"minutes"`
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
	return &TimeEntriesQuery{query: make(_url.Values)}
}

type TimeEntriesQuery struct {
	query _url.Values
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

func (c *Client) Customers(query *CustomersQuery) ([]Customer, error) {
	var q _url.Values
	if query != nil {
		q = query.query
	}

	root := customersRoot{}
	if err := c.get("/customers.xml", q, &root); err != nil {
		return nil, err
	}
	return root.Customers, nil
}

type customersRoot struct {
	XMLName   xml.Name   `xml:"customers"`
	Customers []Customer `xml:"customer"`
}

type Customer struct {
	Id        int       `xml:"id"`
	Name      string    `xml:"name"`
	Note      string    `xml:"note"`
	Archived  bool      `xml:"archived"`
	CreatedAt time.Time `xml:"created-at"`
	UpdatedAt time.Time `xml:"updated-at"`

	// Not implemented yet
	//<hourly-rate type="integer">0</hourly-rate>
	//<active-hourly-rate nil="true"></active-hourly-rate>
}

func NewCustomersQuery() *CustomersQuery {
	return &CustomersQuery{query: make(_url.Values)}
}

type CustomersQuery struct {
	query _url.Values
}

func (q *CustomersQuery) SetName(val string) {
	q.query.Set("name", val)
}

func (q *CustomersQuery) SetLimit(val int) {
	q.query.Set("limit", intToStr(val))
}

func (q *CustomersQuery) SetPage(val int) {
	q.query.Set("page", intToStr(val))
}
