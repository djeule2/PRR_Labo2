package network

import (
	"../config"
	"../mutex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"../utils"
)

const partName = "net"

type Network struct {
	procId uint32
	ready bool

	//channel
	chanFromMutex 	chan string
	chanToMutex 	chan string
	in 				chan string   //tous les messages de processus entrants
	out 			chan string   //tous les messages de processus sortants
	entering		chan struct { // cette structure sert à résoudre la concurrence pour les connexions des processus
		uint32
		net.Conn
	}
	mutex mutex.Mutex
	site map[uint32]net.Conn

}

var dialTimeout = struct {
	timeOut time.Duration
	max time.Duration
}{
	20*time.Second,
	 60*time.Second,
}

//NewNet initialisation d'un nouveau network
func NewNetwork(prcId uint32, chanFromMutex chan string, chanToMutex chan string, mutex mutex.Mutex)*Network {
	network := new(Network)
	network.procId = prcId
	network.in = make(chan string)
	network.out = make(chan string)
	network.entering = make(chan struct{
		uint32
		net.Conn
	})
	network.mutex = mutex
	network.site = make(map[uint32]net.Conn)
	return network

}

/**
* execute() est la goroutine qui s'occupe d'initialiser la connection tcp,
* de lancer les autre goroutine afin de garantir la connection au autre goroutine.
* Elle attend egalement que tous les clients soivent connecter
*
*/
func (n *Network) Excecute()  {
	address := config.Hosts[n.procId]+ ": " +strconv.FormatUint(uint64(config.Ports[n.procId]), 10)
	listener, err := net.Listen("tcp", address)
	utils.PrintMessage(n.procId, partName, "processus is running and listening on port" + fmt.Sprint(config.Ports[n.procId]))
	if err != nil{
		log.Fatal(err)
	}
	go n.broadcaster()
	go n.connectToAll()
	go n.checkReady()

	for{
		conn, err := listener.Accept() //f attent une connection client
		if err!= nil{
			log.Print(err)
			continue
		}
		utils.PrintMessage(n.procId, partName, "processu adress"+conn.RemoteAddr().String())
		go n.handleConn(conn)
	}
}


/**
* broadcaster est une goroutine qui gère les message entre site ou de processus mutex
* pour chacune des message elle est traité et transmis au mutex ou aux autres site
*/

func (n*Network)broadcaster()  {
	for{
		select {
		//traitement des message provenant du mutex puis les transmet au autres site
		case message:= <-n.chanFromMutex:
			utils.PrintMessage(n.procId, partName, "message from mutex "+ message)
			utils.PrintMessage(n.procId, partName, "Sending message to other processus")
			n.out <- message
			//pour messages reçu d'un site  on fait une attente puis on transmet
		case message:=<-n.in:
			time.Sleep(config.TEMPS_TRANSMITION*time.Second)
			utils.PrintMessage(n.procId, partName, "Sent message to mutex")
			go func() {n.chanToMutex<-message}()

        // arriver d'un nouveau site
		case conn := <-n.entering:
			n.site[conn.uint32] = conn.Conn
		}
	}

}
func (n *Network) handleConn(conn net.Conn)  {
	who := conn.RemoteAddr().String() // adress of processu
	utils.PrintMessage(n.procId, partName, "Handleconn for site "+ who)

	go n.sendMessageToNetwork()

	buffer := make([]byte, 4096)
	for{
		utils.PrintMessage(n.procId, partName, "wait message from site "+ who)
		m, err := conn.Read(buffer)
		msg := string(buffer[0:m])
		utils.PrintMessage(n.procId, partName, "Message from other processus arrive [site =" + who+"], msg: "+msg)
		messsegeSplit := strings.Split(msg, ",")
		if strings.HasPrefix(msg, "i") {
			id, err := strconv.Atoi(messsegeSplit[1])

			if err == nil{
				utils.PrintMessage(n.procId, partName, "ConnecToOne =" +strconv.Itoa(id))
				n.connectToOne(uint32(id))
			}

		}else {
			n.in <- msg
		}
		if err!=nil || m == 0{
			conn.Close()
			break
		}

	}
}

//sendMessageToNetwork est la goroutine pour communiquer avec les autres sites
// envoyer les messages reçus par processus mutex

func (n*Network) sendMessageToNetwork() {
	for{
		select {
		case message:= <-n.out:
			msgSplit := strings.Split(message, ",")
			id, err := strconv.Atoi(msgSplit[3])

			if err == nil{
				newMsgSplit:= []string{msgSplit[0], msgSplit[1], msgSplit[2]}
				newMsg := strings.Join(newMsgSplit, ",")

				utils.PrintMessage(n.procId, partName, "send message for site i =" +strconv.Itoa(id))
				n.site[uint32(id)].Write([]byte(newMsg))
			}


		}
	}

}

func (n*Network) connectToOne(procId uint32) {
	adress := config.Hosts[procId] + ":" + strconv.FormatUint(uint64(config.Ports[procId]), 10)
	d := net.Dialer{Timeout: dialTimeout.max}
	utils.PrintMessage(n.procId, partName, "new connection with the site "+fmt.Sprint(procId))
	conn, err := d.Dial("tcp", adress)
	if err!=nil{
		log.Fatal(err)
	}
	n.entering<- struct {
		uint32
		net.Conn
	}{ procId , conn}

	utils.PrintMessage(n.procId, partName, "connection establishrd with site " +fmt.Sprint(procId))

}

//connectToAll permet de se connecter automatiquement aux autre site à par le premier

func (n*Network) connectToAll()  {
	for p, v := range config.Ports{
		if p<n.procId{
			adress := config.Hosts[p] + ":" + strconv.Itoa(int(v))
			d := net.Dialer{Timeout: dialTimeout.max}
			utils.PrintMessage(n.procId, partName, "new connection with the site "+adress)
			conn, err := d.Dial("tcp", adress)

			if err != nil{
				log.Fatal(err)
			}
			n.entering <- struct {
				uint32
				net.Conn
			}{p, conn }

			message := []string{"i", fmt.Sprint(n. procId)}
			conn.Write([]byte(strings.Join(message, ",")))

		}
	}

}

//checReady est la fonction d'attente que tous les processus soient prêt.
//une fois tous le monde prêt nous pouvons lancé le processus mutex
func (n*Network) checkReady()  {
	for n.ready == false{
		if len(n.site) == int(config.DefaultNbrProc)-1{
			n.ready = true
			utils.PrintMessage(n.procId, partName, "Ready, start mutex")
			n.mutex.Exec()
		}
		time.Sleep(120*time.Millisecond)
	}

}