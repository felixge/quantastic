package mite

import (
	"encoding/xml"
	"net/url"
	"time"
)

func (c *Client) Projects(query *ProjectsQuery) ([]*Project, error) {
	var q url.Values
	if query != nil {
		q = query.query
	}

	root := projectsRoot{}
	if err := c.get("/projects.xml", q, &root); err != nil {
		return nil, err
	}
	return root.Projects, nil
}

type projectsRoot struct {
	XMLName  xml.Name   `xml:"projects"`
	Projects []*Project `xml:"project"`
}

type Project struct {
	Id           int       `xml:"id"`
	Name         string    `xml:"name"`
	Note         string    `xml:"note"`
	Budget       int       `xml:"budget"`
	BudgetType   string    `xml:"budget-type"`
	Archived     bool      `xml:"archived"`
	CustomerId   int       `xml:"customer-id"`
	CustomerName string    `xml:"customer-name"`
	CreatedAt    time.Time `xml:"created-at"`
	UpdatedAt    time.Time `xml:"updated-at"`

	// Note implemented yet:
	//<hourly-rate type="integer">0</hourly-rate>
	//<active-hourly-rate nil="true"></active-hourly-rate>
	//<hourly-rates-per-service type="array"></hourly-rates-per-service>
}

func NewProjectsQuery() *ProjectsQuery {
	return &ProjectsQuery{query: make(url.Values)}
}

type ProjectsQuery struct {
	query url.Values
}

func (q *ProjectsQuery) SetName(val string) {
	q.query.Set("name", val)
}

func (q *ProjectsQuery) SetLimit(val int) {
	q.query.Set("limit", intToStr(val))
}

func (q *ProjectsQuery) SetPage(val int) {
	q.query.Set("page", intToStr(val))
}
