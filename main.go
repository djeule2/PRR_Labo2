package main

import (
	"./client"
	"./config"
	"./mutex"
	"./network"
	"flag"
	"sync"
)

var waitAllGoroutine = &sync.WaitGroup{}

func main()  {
	var nbrProc, port int
	var host  string

	flag.IntVar(&nbrProc, "n", config.DefaultNbrProc, "Number of process")
	flag.StringVar(&host, "host", config.DefaultHost, "Host")
	flag.IntVar(&port, "port", config.DefaulPort, "First port")

}