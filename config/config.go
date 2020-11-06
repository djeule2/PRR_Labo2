package config

import (
	"time"
	_ "time"
)

const DefaultNbrProc = 2
const DefaulPort = 8000
const DefaultHost = "Localhost"
const DEBUG = true

const TEMPS_TRANSMITION time.Duration = 2
const Temps_SC time.Duration=5
var nbrPrc  uint32
var Ports = make(map[uint32]uint32)
var Hosts   = make(map[uint32]string)
//var firstPort uint32
var Host string