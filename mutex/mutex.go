package mutex

import (
	"../client"
	"../config"
	"../utils"
	"fmt"
	"strconv"
	"strings"
)

type Mutex struct {
	n		 	 uint32
	h 		 	 uint32
	moi 	 	 uint32
	tabM 	 	 []string
	accordSC 	 bool

	//channels
	chanFromClient	 chan string
	chanToClient 	 chan string
	chanFromNetwork  chan string
	chanToNetwork    chan string
	client  		 client.Client
}

const partName = "mut"

//initialisation d'un nouveau processus Mutex à partir de son id, 4 channel pour
//la commuinication avec client et avec network et une référence sur le client.
func NewMutex(id uint32, chanFromClient chan string, chanToClient chan string, chanFroNetwork chan string,
chanToNetwork chan string, client client.Client) *Mutex  {
	mut := new(Mutex)
	mut.n = config.DefaultNbrProc
	mut.moi = id
	mut.h = 0
	mut.accordSC = false
	mut.tabM = make([]string, 0)
	mut.client = client

	mut.chanToNetwork = chanToNetwork
	mut.chanFromNetwork = chanFroNetwork
	mut.chanToClient = chanToClient
	mut.chanFromClient = chanFromClient

	return mut
}

func (m *Mutex)Exec()  {
	go m.ManageMessage()
}

func (mess *Mutex)ManageMessage()  {

	for {
		select {
		//Message provided by network REQ, ACK, RELEASE
		case clientMsg := <-mess.chanFromNetwork:
			utils.PrintMessage(mess.moi, partName, clientMsg)

			//test si le message commence par le mot cle REQ
			if strings.HasPrefix(clientMsg, config.REQ) {
				utils.PrintMessage(mess.moi, partName, "REQ received")
				msgSplit := strings.Split(clientMsg, ",")
				Hi, err1 := strconv.ParseUint(msgSplit[0], 10, 32)
				i, err2 := strconv.ParseUint(msgSplit[1], 10, 32)
				if err1 == nil && err2 == nil {
					mess.req(uint32(Hi), uint32(i))
				}
			} else if strings.HasPrefix(clientMsg, config.ACK) {
				utils.PrintMessage(mess.moi, partName, "ACK received")
				msgSplit := strings.Split(clientMsg, ",")
				Hi, err1 := strconv.ParseUint(msgSplit[0], 10, 32)
				i, err2 := strconv.ParseUint(msgSplit[1], 10, 32)
				if err1 == nil && err2 == nil {
					mess.ack(uint32(Hi), uint32(i))
				}
			} else if strings.HasPrefix(clientMsg, config.REL) {
				utils.PrintMessage(mess.moi, partName, "REL received")
				msgSplit := strings.Split(clientMsg, ",")
				Hi, err1 := strconv.ParseUint(msgSplit[0], 10, 32)
				i, err2 := strconv.ParseUint(msgSplit[1], 10, 32)
				if err1 == nil && err2 == nil {
					mess.rel(uint32(Hi), uint32(i))
				}
			}

			//Message provide by client DEMANDE, ATTENTE, RELEASE
			case mutexMsg := <-mess.chanFromClient:
				utils.PrintMessage(mess.moi, partName, mutexMsg)
				if strings.HasPrefix(mutexMsg, config.DEMANDE){
					utils.PrintMessage(mess.moi, partName, "Demande received from client")
					mess.demande()
				} else if strings.HasPrefix(mutexMsg, config.ATTENTE){
					utils.PrintMessage(mess.moi, partName, "ATTENTE acoord for SC")
					mess.attente()
				}else if strings.HasPrefix(mutexMsg, config.FIN){
					utils.PrintMessage(mess.moi, partName, "END received from client")
					mess.fin()

				}
		}
	}

}

//demade() traitement demande, fait partie de l'olgorithme de Lamport
func (m *Mutex) demande()  {
	m.h++
	message := []string{config.REQ, fmt.Sprint(m.h), fmt.Sprint(m.moi)}
	m.tabM[m.moi] = strings.Join(message, ",")

	for i:= uint32(0); i<m.n; i++ {
			m.sendMessage(config.REQ, m.h, i)
	}

}

func (m*Mutex)fin()  {
	m.h++
	m.accordSC = false
	message := []string{config.REL, fmt.Sprint(m.h), fmt.Sprint(m.moi)}
	m.tabM[m.moi] = strings.Join(message, ",")

	for i:= uint32(0); i<m.n; i++{
			m.sendMessage(config.REL, m.h, i)
	}
}

func (m*Mutex)attente()  {
	if m.accordSC{
		// le traittement débloque le client qui passe en SC
	}
}

func (m*Mutex)req(hi uint32, i uint32)  {
	m.h = utils.Max(hi, m.h) + 1
	message := []string{config.REQ, fmt.Sprint(hi), fmt.Sprint(i)}
	m.tabM[i] = strings.Join(message, ",")
	m.sendMessage(config.ACK, m.h, i)
	m.verifieSC()

}

func (m*Mutex)ack(hi uint32, i uint32)  {
	m.h = utils.Max(m.h, hi)
	msgType := strings.Split(m.tabM[i], ",")[0]
	if msgType != config.REQ{
		message := []string{config.ACK, fmt.Sprint(hi), fmt.Sprint(i)}
		m.tabM[i] = strings.Join(message, ",")
	}
	m.verifieSC()
}

func (m*Mutex)rel(hi uint32, i uint32)  {
	m.h = utils.Max(hi, m.h) + 1
	message := []string{config.REL, fmt.Sprint(hi), fmt.Sprint(i)}
	m.tabM[i] = strings.Join(message, ",")
	m.verifieSC()
}

func (m*Mutex) verifieSC()  {
	 msgType := strings.Split(m.tabM[m.moi], ",")[0]
	if msgType != config.REQ{
		return
	}
	 plusAcienne := true
	 for i:= uint32(0); i<m.n; i++{
	 	if m.tabM[m.moi] >m.tabM[i] {
	 		plusAcienne = false
			break
		}else if (m.tabM[m.moi] == m.tabM[i]) && (m.moi>i){
			plusAcienne = false
			break
		}
	 }
	 if plusAcienne{
	 	m.accordSC = true
	 }
}


func (m*Mutex) sendMessage(methode string, hi uint32, i uint32)  {
	if i != m.moi {
		message := []string{methode, fmt.Sprint(hi),fmt.Sprint(i) }
		m.sendMessageToNetwork(strings.Join(message, ","))
	}
}

func (m*Mutex)sendMessageToNetwork(message string)  {
	utils.PrintMessage(m.moi, partName, "Send to network : "+message)
	m.chanToNetwork <- message
}


