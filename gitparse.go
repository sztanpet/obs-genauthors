package main

func (a *app) gitAuthors() []contributor {
	// TODO working directory to the directory of the executable
	// assume the obs git repo is there, and just parse the output of
	// git log --all --format='%aN|%cE'
}
