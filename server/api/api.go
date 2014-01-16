package api

import (
	"fmt"
	"github.com/felixge/log"
	"github.com/felixge/quantastic/model"
	"github.com/felixge/quantastic/services/mite"
	"time"
)

type Config struct {
	Log  log.Interface
	Mite *mite.Client
}

func NewApi(config Config) *Api {
	return &Api{
		mite: config.Mite,
		log:  config.Log,
	}
}

func (api *Api) ReadTimeRange(timeRange *model.TimeRange) error {
	var (
		timeout       = 5 * time.Second // @TODO make configurable
		miteEntries   []mite.TimeEntry
		miteProjects  []*mite.Project
		miteCustomers []mite.Customer
		queries       = make(chan error, 3)
	)

	go func() {
		query := mite.NewTimeEntriesQuery()
		query.SetFrom(timeRange.Start)
		query.SetTo(timeRange.End)
		var err error
		miteEntries, err = api.mite.TimeEntries(query)
		queries <- err
	}()

	go func() {
		query := mite.NewProjectsQuery()
		var err error
		miteProjects, err = api.mite.Projects(query)
		queries <- err
	}()

	go func() {
		query := mite.NewCustomersQuery()
		var err error
		miteCustomers, err = api.mite.Customers(query)
		queries <- err
	}()

	if err := waitQueries(queries, timeout); err != nil {
		return err
	}

	timeRange.Root = timeCategoryRoot(miteCustomers, miteProjects)
	timeRange.Entries = timeEntries(miteEntries)

	for _, category := range timeRange.Root.All() {
		for _, entry := range timeRange.Entries {
			if category.Id == entry.Category.Id {
				entry.Category = category
				category.Entries = append(category.Entries, entry)
			}
		}
	}
	return nil
}

func timeCategoryRoot(miteCustomers []mite.Customer, miteProjects []*mite.Project) *model.TimeCategory {
	root := &model.TimeCategory{Name: "Root"}
	root.Children = make([]*model.TimeCategory, 0, len(miteCustomers))
	for _, miteCustomer := range miteCustomers {
		customerCategory := &model.TimeCategory{
			Name:   miteCustomer.Name,
			Parent: root,
		}
		for _, miteProject := range miteProjects {
			if miteCustomer.Id == miteProject.CustomerId {
				projectCategory := &model.TimeCategory{
					Id:     intToStr(miteProject.Id),
					Name:   miteProject.Name,
					Parent: customerCategory,
				}
				customerCategory.Children = append(customerCategory.Children, projectCategory)
			}
		}
		root.Children = append(root.Children, customerCategory)
	}
	return root
}

func timeEntries(miteEntries []mite.TimeEntry) []*model.TimeEntry {
	timeEntries := make([]*model.TimeEntry, 0, len(miteEntries))
	for _, miteEntry := range miteEntries {
		duration := time.Duration(miteEntry.Minutes) * time.Minute
		timeEntries = append(timeEntries, &model.TimeEntry{
			From:     miteEntry.CreatedAt,
			To:       miteEntry.CreatedAt.Add(duration),
			Category: &model.TimeCategory{Id: intToStr(miteEntry.ProjectId)},
		})
	}
	return timeEntries
}

func waitQueries(queries chan error, timeout time.Duration) error {
	results := 0
	for results < cap(queries) {
		select {
		case err := <-queries:
			if err != nil {
				return err
			}
			results++
		case <-time.After(timeout):
			return NewTimeoutError(timeout)
		}
	}
	return nil
}

func intToStr(val int) string {
	return fmt.Sprintf("%d", val)
}

//func (api *Api) GetTimeEntries(from, to time.Time) ([]*model.TimeEntry, error) {
//query := mite.NewTimeEntriesQuery()
//query.SetFrom(from)
//query.SetTo(to)
//miteEntries, err := api.mite.TimeEntries(query)
//if err != nil {
//return nil, err
//}
//timeEntries := make([]*model.TimeEntry, 0, len(miteEntries))
//for _, miteEntry := range miteEntries {
//duration := time.Duration(miteEntry.Minutes) * time.Minute
//timeEntries = append(timeEntries, &model.TimeEntry{
//From:     miteEntry.CreatedAt,
//To:       miteEntry.CreatedAt.Add(duration),
//Duration: duration,
//Category: &model.TimeCategory{
//Name:   miteEntry.ProjectName,
//Parent: &model.TimeCategory{Name: miteEntry.CustomerName},
//},
//})
//}
//return timeEntries, nil
//}

//func (api *Api) GetTimeCategories(from, to time.Time) ([]*model.TimeCategory, error) {
//}

type Api struct {
	log  log.Interface
	mite *mite.Client
}
