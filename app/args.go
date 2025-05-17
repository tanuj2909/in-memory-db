package main

import "flag"

type Args struct {
	port int
}

func GetArgs() Args {
	port := flag.Int("port", 6379, "The port on which the db listens")

	flag.Parse()

	return Args{
		port: *port,
	}
}
