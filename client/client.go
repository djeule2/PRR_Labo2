package client

import (
	"../utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"../config"
	"time"
)

type Client struct {
	id				uint32
	valueSC 		uint32
	chanFromMutex 	chan string
	chanToMutex 	chan string
	demandeTime		[]int
	demandeEnCour	bool

}

const partName  = "clt"

/*
* Initialise un nouveau processus Client à partir de son id, deux channnel pour
* communiquer avec mutex et un tableau de message auquel il pourra demander la section critique
 */
func NewClient(id uint32, chanFronMutex chan string, chanToMutex chan string, demandeTime []int) *Client  {
	client := new(Client)
	client.id = id
	client.chanFromMutex = chanFronMutex
	client.chanToMutex = chanToMutex
	client.demandeTime = demandeTime

	return client
}

//Exec  lance les deux goroutines du client, pour communiquer avec le mutex dans les deux sens
func (c *Client) Exec()  {
	go c.receiveMessage()
	go c.demadeReq()

}

/**
* receiveMessage est une goroutine qui écoute les messages provenant du mutex, si le message
* commence par OK le client entre en sc et modifie la vriable partagées. si le message reçu
* possède le prefixe UPDATE la variable dans la section critique à été modifier par un autre
*processus, client vas la mettre à jour simplement
 */
func (c *Client)receiveMessage()  {
	for{
		mutexMsg := <- c.chanFromMutex
		utils.PrintMessage(c.id, partName, mutexMsg);
		if strings.HasPrefix(mutexMsg, config.AUTORISATION){
			utils.PrintMessage(c.id, partName, "Authorization to access to SC change the value")
			c.traitment()
		} else if strings.HasPrefix(mutexMsg, config.UPDATE) {
			utils.PrintMessage(c.id, partName, "New value received, updated")
			splitMsg := strings.Split(mutexMsg, ", ");
			value, err := strconv.ParseUint(splitMsg[1], 10, 32)
			if err == nil{
				c.setValueSC(uint32(value))
			}
		}
	}
}

func (c *Client) traitment() {
	time.Sleep(config.Temps_SC*time.Second)
	c.changeValue()
	message := []string{config.UPDATE, fmt.Sprint(c.valueSC)}
	c.chanToMutex <- strings.Join(message, ",")
	c.chanToMutex <- config.FIN
}

// demandeReq est une goroutine qui effectuera les demandes de client
//pour leurs accès à la section critique
func (c*Client)demadeReq()  {
	for _, v :=range c.demandeTime {
		for c.demandeEnCour{
			time.Sleep(100*time.Second)
		}
		time.Sleep(time.Duration(v)*time.Second)
		utils.PrintMessage(c.id, partName, "Client request= "+strconv.Itoa(v))
		c.chanToMutex<-config.DEMANDE
		c.demandeEnCour = true
	}

}
//on affiche la valeur de la variable partager
func (c*Client) getValueSC() uint32  {
	return c.valueSC
}

//on modifie la valeur de la variable partager
func (c*Client)setValueSC(newValue uint32) {
	utils.PrintMessage(c.id, partName, "Value before change =" +fmt.Sprint(c.valueSC))
	c.valueSC = newValue
	utils.PrintMessage(c.id, partName, "Value after changen= " +fmt.Sprint(c.valueSC))
}

func (c*Client)changeValue() {
	utils.PrintMessage(c.id, partName, "change value which random")
	c.setValueSC(uint32(rand.Intn(200) +10))
}

