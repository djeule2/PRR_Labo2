package client

import (
	"../Auction"
	"../config"
	"../utils"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	id				uint32
	valueSC 		uint32
	chanFromMutex 	chan string
	chanToMutex 	chan string
	demandeTime		[]int
	demandeEnCour	bool
	user			string
	users 			[]string
	auctions 		[]Auction.Auction
	subscriptions 	[]Auction.Subscription
	newAuction 		[]Auction.NewAuctionSub
}

var (
	auctionIds = 0
)

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

	var input string
	unique := false
	for unique==false {
		fmt.Scanln(&input)
		if runtime.GOOS == "windows" {
			c.user = strings.TrimRight(input, "\r\n")
		} else {
			c.user = strings.TrimRight(input, "\n")
		}
		// SECTION CRITIQUE
		c.chanToMutex <- config.DEMANDE

		if !c.checkIfUserIsNotUnique(c.user) {
			c.addUser()
			message := []string{config.ADD_USER, fmt.Sprint(c.user)}
			c.chanToMutex <- strings.Join(message, ",")
		} else {
			println("user not unique")
		}
		c.chanToMutex <- config.FIN
		// SECTION CRITIQUE
	}
	print(menu())
	fmt.Scanln(&input)
	for input != "exit" {
		switch input {
		case "1":
			// SECTION CRITIQUE
		case "2":
		case "3":
		case "4":
		}
	}

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
		utils.PrintMessage(c.id, partName, mutexMsg)
		if strings.HasPrefix(mutexMsg, config.AUTORISATION){
			utils.PrintMessage(c.id, partName, " Authorization to access to SC change the value \n")
			c.traitment()
		} else if strings.HasPrefix(mutexMsg, config.VALUE) {

			utils.PrintMessage(c.id, partName, " New value received, updated \n")
			splitMsg := strings.Split(mutexMsg, ",")
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
	message := []string{config.USERS, fmt.Sprint(c.users)}
	c.chanToMutex <- strings.Join(message, ",")
	c.chanToMutex <- config.FIN
}

// demandeReq est une goroutine qui effectuera les demandes de client
//pour leurs accès à la section critique
func (c*Client)demadeReq()  {
	for _, v :=range c.demandeTime {
		time.Sleep(time.Duration(v)*time.Second)
		utils.PrintMessage(c.id, partName, " Client request= "+strconv.Itoa(v) +"\n")
		c.chanToMutex <- config.DEMANDE

	}

}
//on affiche la valeur de la variable partager
func (c*Client) getValueSC() uint32  {
	return c.valueSC
}

//on modifie la valeur de la variable partager
func (c*Client)setValueSC(newValue uint32) {
	utils.PrintMessage(c.id, partName, "Value before change = " +fmt.Sprint(c.valueSC) +"\n")
	c.valueSC = newValue
	utils.PrintMessage(c.id, partName, "Value after change= " +fmt.Sprint(c.valueSC) +"\n")
}

func (c*Client)changeValue() {
	utils.PrintMessage(c.id, partName, " change value which random : \n")
	c.setValueSC(uint32(rand.Intn(200) +10))
}



func (c *Client) addUser() {
	c.users = append(c.users, c.user)
	message := []string{config.USERS, fmt.Sprint(c.users)}
	c.chanToMutex <- strings.Join(message, ",")
	c.chanToMutex <- config.FIN
}

func ToStringLot(l Auction.Auction) string {
	return fmt.Sprintf("name : %v\t price : %v tokens\t id : %v\t remaining time : %.3v minutes",
		l.Nom, l.Price, l.IdName, l.Temps.Sub(time.Now()).Minutes())
}

func menu() string {
	return "==== menu ====\n" +
		"1 : see auctions\n" +
		"2 : bid\n" +
		"3 : manage notifications\n" +
		"4 : create auction\n" +
		"exit : exit\n"
}

func (c *Client) removeUser(user string) {
	for i, item := range c.users {
		if item == user {
			c.users = append(c.users[:i], c.users[i+1:]...)
		}
	}
}

func (c *Client) removeAuction(auction Auction.Auction) {
	for i, item := range c.auctions {
		if item == auction {
			c.auctions = append(c.auctions[:i], c.auctions[i+1:]...)
		}
	}
}

func (c *Client) removeSubscription(sub Auction.Subscription) {
	for i, item := range c.subscriptions {
		if item == sub {
			c.subscriptions = append(c.subscriptions[:i], c.subscriptions[i+1:]...)
		}
	}
}

func (c *Client) search(id string) (bool, int) {
	for i, item := range c.auctions {
		if id == item.IdName {
			return true, i
		}
	}
	return false, -1
}

func getValueFromUser(input string, message string) string {
	println(message)
	fmt.Scanln(&input)
	return input
}

func (c *Client) checkIfUserIsNotUnique(user string) bool {
	// SECTION CRITIQUE DEBUT
	for _, item := range c.users {
		if item == user {
			return true
		}
	}
	return false
	// SECTION CRITIQUE FIN
}

