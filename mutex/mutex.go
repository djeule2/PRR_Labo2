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
	mut.n = config.NbrPrc
	mut.moi = id
	mut.h = 0
	mut.accordSC = false
	mut.tabM = make([]string, mut.n)
	mut.client = client

	message := []string{config.REL, fmt.Sprint(mut.h)}

	for i:= uint32(0); i<mut.n; i++{
		mut.tabM[i] = strings.Join(message, ",")
	}

	mut.chanToNetwork = chanToNetwork
	mut.chanFromNetwork = chanFroNetwork
	mut.chanToClient = chanToClient
	mut.chanFromClient = chanFromClient

	return mut
}

func (m *Mutex)Exec()  {
	go m.client.Exec()
	go m.ManageMessage()
}

func (mess *Mutex)ManageMessage()  {

	for {
		select {
		//Message provided by network REQ, ACK, RELEASE
		case clientMsg := <-mess.chanFromNetwork:
			utils.PrintMessage(mess.moi, partName, clientMsg +"\n")

			//test si le message commence par le mot cle REQ
			if strings.HasPrefix(clientMsg, config.REQ) {
				utils.PrintMessage(mess.moi, partName, " REQ received \n")
				msgSplit := strings.Split(clientMsg, ",")
				Hi, err1 := strconv.ParseUint(msgSplit[1], 10, 32)
				i, err2 := strconv.ParseUint(msgSplit[2], 10, 32)

				if err1 == nil && err2 == nil {
					mess.req(uint32(Hi), uint32(i))
				}
			} else if strings.HasPrefix(clientMsg, config.ACK) {
				utils.PrintMessage(mess.moi, partName, " ACK received \n")
				msgSplit := strings.Split(clientMsg, ",")
				Hi, err1 := strconv.ParseUint(msgSplit[1], 10, 32)
				i, err2 := strconv.ParseUint(msgSplit[2], 10, 32)
				if err1 == nil && err2 == nil {
					mess.ack(uint32(Hi), uint32(i))

				}
			} else if strings.HasPrefix(clientMsg, config.VALUE){
				utils.PrintMessage(mess.moi, partName, config.VALUE)
				mess.chanToClient <- clientMsg

			} else if strings.HasPrefix(clientMsg, config.USERS){
			utils.PrintMessage(mess.moi, partName, config.USERS)
			mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.ADD_USER){
				utils.PrintMessage(mess.moi, partName, config.ADD_USER)
				mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.ADD_AUCTION){
				utils.PrintMessage(mess.moi, partName, config.ADD_AUCTION)
				mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.REMOVE_AUCTION){
				utils.PrintMessage(mess.moi, partName, config.REMOVE_AUCTION)
				mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.INCREMENT_AUCTION_IDS){
				utils.PrintMessage(mess.moi, partName, config.INCREMENT_AUCTION_IDS)
				mess.chanToClient <- clientMsg


			} else if strings.HasPrefix(clientMsg, config.NOTIFY_NEWAUCTION){
				utils.PrintMessage(mess.moi, partName, config.NOTIFY_NEWAUCTION)
				mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.NOTIFY_SUBSCRIBED){
				utils.PrintMessage(mess.moi, partName, config.NOTIFY_SUBSCRIBED)
				mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.NOTIFY_END){
				utils.PrintMessage(mess.moi, partName, config.NOTIFY_END)
				mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.CHANGE_WINNER){
				utils.PrintMessage(mess.moi, partName, config.CHANGE_WINNER)
				mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.CHANGE_PRICE) {
				utils.PrintMessage(mess.moi, partName, config.CHANGE_PRICE)
				mess.chanToClient <- clientMsg

			}else if strings.HasPrefix(clientMsg, config.USERS){
					utils.PrintMessage(mess.moi, partName, config.USERS)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.ADD_USER){
					utils.PrintMessage(mess.moi, partName, config.ADD_USER)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.ADD_AUCTION){
					utils.PrintMessage(mess.moi, partName, config.ADD_AUCTION)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.REMOVE_AUCTION){
					utils.PrintMessage(mess.moi, partName, config.REMOVE_AUCTION)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.INCREMENT_AUCTION_IDS){
					utils.PrintMessage(mess.moi, partName, config.INCREMENT_AUCTION_IDS)
					mess.chanToClient <- clientMsg


			} else if strings.HasPrefix(clientMsg, config.NOTIFY_NEWAUCTION){
					utils.PrintMessage(mess.moi, partName, config.NOTIFY_NEWAUCTION)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.NOTIFY_SUBSCRIBED){
					utils.PrintMessage(mess.moi, partName, config.NOTIFY_SUBSCRIBED)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.NOTIFY_END){
					utils.PrintMessage(mess.moi, partName, config.NOTIFY_END)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.CHANGE_WINNER){
					utils.PrintMessage(mess.moi, partName, config.CHANGE_WINNER)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.CHANGE_PRICE){
					utils.PrintMessage(mess.moi, partName, config.CHANGE_PRICE)
					mess.chanToClient <- clientMsg


			}else if strings.HasPrefix(clientMsg, config.REL) {
				utils.PrintMessage(mess.moi, partName, " REL received \n")
				msgSplit := strings.Split(clientMsg, ",")
				Hi, err1 := strconv.ParseUint(msgSplit[1], 10, 32)
				i, err2 := strconv.ParseUint(msgSplit[2], 10, 32)

				if err1 == nil && err2 == nil {

					mess.rel(uint32(Hi), uint32(i))

				}

			}

			//Message provide by client DEMANDE, ATTENTE, RELEASE
			case mutexMsg := <-mess.chanFromClient:
				utils.PrintMessage(mess.moi, partName, mutexMsg +"\n")
				if strings.HasPrefix(mutexMsg, config.DEMANDE){
					utils.PrintMessage(mess.moi, partName, " Demande received from client \n")
					mess.demande()
				} else if strings.HasPrefix(mutexMsg, config.VALUE){

					utils.PrintMessage(mess.moi, partName, " Value change received from client \n")
					splitMsg := strings.Split(mutexMsg, ",")
					v, error := strconv.ParseUint(splitMsg[1], 10, 32)
					if error == nil{
						mess.sendValue(uint32(v))
					}

				} else if strings.HasPrefix(mutexMsg, config.USERS){
					utils.PrintMessage(mess.moi, partName, config.USERS)
					// ....


				}else if strings.HasPrefix(mutexMsg, config.ADD_USER){
					utils.PrintMessage(mess.moi, partName, config.ADD_USER)
					//..


				}else if strings.HasPrefix(mutexMsg, config.ADD_AUCTION){
					utils.PrintMessage(mess.moi, partName, config.ADD_AUCTION)
					//...


				}else if strings.HasPrefix(mutexMsg, config.REMOVE_AUCTION){
					utils.PrintMessage(mess.moi, partName, config.REMOVE_AUCTION)



				}else if strings.HasPrefix(mutexMsg, config.INCREMENT_AUCTION_IDS){
					utils.PrintMessage(mess.moi, partName, config.INCREMENT_AUCTION_IDS)



				} else if strings.HasPrefix(mutexMsg, config.NOTIFY_NEWAUCTION){
					utils.PrintMessage(mess.moi, partName, config.NOTIFY_NEWAUCTION)



				}else if strings.HasPrefix(mutexMsg, config.NOTIFY_SUBSCRIBED){
					utils.PrintMessage(mess.moi, partName, config.NOTIFY_SUBSCRIBED)



				}else if strings.HasPrefix(mutexMsg, config.NOTIFY_END){
					utils.PrintMessage(mess.moi, partName, config.NOTIFY_END)



				}else if strings.HasPrefix(mutexMsg, config.CHANGE_WINNER){
					utils.PrintMessage(mess.moi, partName, config.CHANGE_WINNER)



				}else if strings.HasPrefix(mutexMsg, config.CHANGE_PRICE) {
					utils.PrintMessage(mess.moi, partName, config.CHANGE_PRICE)

				}else if strings.HasPrefix(mutexMsg, config.FIN){
					utils.PrintMessage(mess.moi, partName, " END received from client \n")
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
	message := []string{config.REL, fmt.Sprint(m.h), fmt.Sprint(m.moi)}
	m.tabM[m.moi] = strings.Join(message, ",")

	for i:= uint32(0); i<m.n; i++{
			m.sendMessage(config.REL, m.h, i)
	}
	m.accordSC = false
}

func (m*Mutex)attente()  {

	if m.accordSC {
		// le traittement débloque le client qui passe en SC
		m.chanToClient <- config.AUTORISATION
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
	m.h = utils.Max(m.h, hi) +1
	msgType := strings.Split(m.tabM[i], ",")[0]
	if msgType != config.REQ{
		message := []string{config.ACK, fmt.Sprint(hi),  fmt.Sprint(i)}
		m.tabM[i] = strings.Join(message, ",")
	}
	m.verifieSC()
	m.attente()
}

func (m*Mutex)rel(hi uint32, i uint32)  {
	m.h = utils.Max(hi, m.h) + 1
	message := []string{config.REL, fmt.Sprint(hi), fmt.Sprint(i)}
	m.tabM[i] = strings.Join(message, ",")
	m.verifieSC()
	m.attente()
}

func (m*Mutex) verifieSC()  {
	msgType := strings.Split(m.tabM[m.moi], ",")[0]
	if msgType != config.REQ {
		return
	}
	hm := strings.Split(m.tabM[m.moi], ",")[1]

	 plusAcienne := true
	 for i:= uint32(0); i<m.n; i++{
	 	if m.moi == i{
	 		continue
		}
		hi:=  strings.Split(m.tabM[i], ",")[1]

	 	if hm > hi {
	 		plusAcienne = false
			break
		}else if (hm == hi) && (m.moi>i){
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
		message := []string{methode, fmt.Sprint(hi), fmt.Sprint(m.moi), fmt.Sprint(i) }
		m.sendMessageToNetwork(strings.Join(message, ","))
	}
}

func (m*Mutex)sendValue(v uint32)  {
	for i:= uint32(0); i< m.n; i++{
		if i != m.moi{
			mess := []string{config.VALUE, fmt.Sprint(v), fmt.Sprint(m.moi), fmt.Sprint(i)}
			m.sendMessageToNetwork(strings.Join(mess, ","))
		}
	}

}

func (m*Mutex)sendMessageToNetwork(message string)  {
	utils.PrintMessage(m.moi, partName, " Send to network : "+message +"\n")
	m.chanToNetwork <- message
}


