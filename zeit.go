package zeit

import (
	"fmt"
	"github.com/felixge/zeit/config"
	"github.com/felixge/zeit/mite"
)

func NewZeit(configPath string) (z *Zeit, err error) {
	z = &Zeit{}
	c := config.Config{}
	if err = config.Load("config.xml", &c); err != nil {
		return
	}
	fmt.Printf("config: %#v\n", c)

	z.Mite, err = mite.NewMite(c.Mite.Url, c.Mite.ApiKey)
	if err != nil {
		return
	}
	return
}

type Zeit struct {
	Mite   *mite.Mite
}

func (z *Zeit) Run() error {
	query := mite.NewTimeEntriesQuery()
	query.SetProjectId(1215321)
	entries, err := z.Mite.TimeEntries(query)
	if err != nil {
		return err
	}
	fmt.Printf("entries: %#v\n", entries)
	return nil
}
