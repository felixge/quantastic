package main

type Context struct {
	Cmd *command
	Args []string
	Db   *Db
}
