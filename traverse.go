package genie

func (l *Lamp) TraverseCommands(do func(command *Command)) {
	l.RootCommand.AnchorPaths()
	do(l.RootCommand)
	for _, sc := range l.RootCommand.SubCommands {
		traverse(sc, do)
	}
}

func traverse(c *Command, do func(command *Command)) {
	do(c)
	for _, sc := range c.SubCommands {
		traverse(sc, do)
	}
}
