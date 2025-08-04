package main

import "flag"

type Args struct {
	port       int
	replicaof  string
	dir        string
	dbfilename string
}

func GetArgs() Args {
	port := flag.Int("port", 6379, "The port on which the db listens")

	replicaof := flag.String("replicaof", "", "The host and port of master instance")

	dir := flag.String("dir", "tmp/redis", "The directory to store database files")

	dbfilename := flag.String("dbfilename", "dump.rdb", "The name of the database file")

	flag.Parse()

	return Args{
		port:       *port,
		replicaof:  *replicaof,
		dir:        *dir,
		dbfilename: *dbfilename,
	}
}
