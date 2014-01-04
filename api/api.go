package api

import (
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
	}
}

func (api *Api) GetTimeEntries(from, to time.Time) ([]*model.TimeEntry, error) {
	query := mite.NewTimeEntriesQuery()
	query.SetFrom(from)
	query.SetTo(to)
	miteEntries, err := api.mite.TimeEntries(query)
	if err != nil {
		return nil, err
	}
	timeEntries := make([]*model.TimeEntry, 0, len(miteEntries))
	for _, miteEntry := range miteEntries {
		duration := time.Duration(miteEntry.Minutes) * time.Minute
		timeEntries = append(timeEntries, &model.TimeEntry{
			From:     miteEntry.CreatedAt,
			To:       miteEntry.CreatedAt.Add(duration),
			Duration: duration,
			Category: &model.TimeCategory{
				Name:   miteEntry.ProjectName,
				Parent: &model.TimeCategory{Name: miteEntry.CustomerName},
			},
		})
	}
	return timeEntries, nil
}

func (api *Api) GetTimeCategories(from, to time.Time) ([]*model.TimeCategory, error) {
	pQuery := mite.NewProjectsQuery()
	miteProjects, err := api.mite.Projects(pQuery)
	if err != nil {
		return nil, err
	}

	cQuery := mite.NewCustomersQuery()
	miteCustomers, err := api.mite.Customers(cQuery)
	if err != nil {
		return nil, err
	}

	categories := make([]*model.TimeCategory, 0, len(miteCustomers))
	for _, miteCustomer := range miteCustomers {
		category := &model.TimeCategory{Name: miteCustomer.Name}
		for _, miteProject := range miteProjects {
			if miteCustomer.Id == miteProject.CustomerId {
				category.Children = append(category.Children, &model.TimeCategory{
					Parent: category,
					Name:   miteProject.Name,
				})
			}
		}
		categories = append(categories, category)
	}
	return categories, nil
}

type Api struct {
	mite *mite.Client
}
