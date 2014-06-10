package main

var cmdTimeRm = &command{
	name:        "time rm",
	description: "Delete an entry.",
	usage:       "<id>",
	fn:          cmdTimeRmFn,
}

func cmdTimeRmFn(c *Context) {
	if len(c.Args) == 0 {
		fatal("%s", c.Cmd.Usage())
	}
	id := c.Args[0]
	if err := c.Db.DeleteTimeEntry(id); err != nil {
		fatal("Could not delete entry: %s", err)
	}
	mustPrintfStdout("Deleted entry %s.\n", id)
}
