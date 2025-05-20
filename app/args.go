package main

import "flag"

type Args struct {
	port      int
	replicaof string
}

func GetArgs() Args {
	port := flag.Int("port", 6379, "The port on which the db listens")

	replicaof := flag.String("replicaof", "", "The host and port of master instance")

	flag.Parse()

	return Args{
		port:      *port,
		replicaof: *replicaof,
	}
}
