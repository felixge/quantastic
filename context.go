package main

import (
	"github.com/felixge/quantastic/db"
)

type Context struct {
	Cmd *command
	Args []string
	Db   *db.Db
}
