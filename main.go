package main

import (
	"./client"
	"./config"
	"./mutex"
	"./network"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var waitAllGroupe = &sync.WaitGroup{}

type demandesReqType []int

var demandesReq demandesReqType

//String() formate la valeur du flag.
func (d *demandesReqType) String() string {
	return fmt.Sprint(*d)
}

//on definit la valeur du flag de type demandesReqType
func (r*demandesReqType) Set(s string) error  {
	split := strings.Split(s, ",")

	for _, v := range split{
		i, err := strconv.ParseInt(v, 10, 64)

		if err == nil{
			*r = append(*r, int(i))
		} else {
			return err
		}
	}
	return nil
}

func main()  {
	var nbrProc, port int
	var host  string
	var defaultDemandesReq demandesReqType = []int{2, 8}


	flag.IntVar(&nbrProc, "n", config.DefaultNbrProc, "Number of process")
	flag.StringVar(&host, "host", config.DefaultHost, "Host")
	flag.IntVar(&port, "port", config.DefaulPort, "First port")
	flag.Var(&demandesReq, "REQ", "Times of the req demande")
	flag.Parse()
	config.FirsPort =uint32(port)
	config.NbrPrc = uint32(nbrProc)

	if demandesReq == nil{
		demandesReq = defaultDemandesReq
	}

	fmt.Println("First port: ", config.FirsPort)
	fmt.Println("Host: ", host)
	fmt.Println("Number of process: ", config.NbrPrc);

	for i:= uint32(0); i<config.NbrPrc; i++{
		config.Ports[i] = config.FirsPort +1
		config.Hosts[i] = host;
		waitAllGroupe.Add(1)
		exec(i);
	}
	waitAllGroupe.Wait()
}

func exec(id uint32)  {
	chanClientMutex := make(chan string)
	chanMutexClient := make(chan string)
	chanMutexNetwork := make(chan string)
	chanNetworkMutex := make(chan string)

	//on definit les différents temps d'attente pour les envois de REQ
	demReq := make([]int, 0)
	if int(id) < len(demandesReq){
		for k, v := range demandesReq {
			if(v != 0) && k%int(config.NbrPrc) == int(id){
				demReq = append(demReq, v)
			}

		}
	}
	fmt.Println("Id", id, ", REQ aux temps :", demReq)
	client := client.NewClient(id, chanMutexClient, chanClientMutex, demReq)
	mutex := mutex.NewMutex(id, chanClientMutex, chanMutexClient, chanNetworkMutex, chanMutexNetwork, *client)
	network := network.NewNetwork(id, chanMutexNetwork, chanNetworkMutex, *mutex)

	go network.Excecute();

}
