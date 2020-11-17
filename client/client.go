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
	newAuction 		bool
	auctionIds      int
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
	client.demandeEnCour = false
	client.newAuction = false

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

		c.demandeSectionCritique()
		if !c.checkIfUserIsNotUnique(c.user) {
			message := []string{config.ADD_USER, fmt.Sprint(c.user)}
			c.chanToMutex <- strings.Join(message, ",")
		} else {
			println("user not unique")
		}
		c.relacherSectionCritique()

	}
	print(menu())
	fmt.Scanln(&input)
	for input != "exit" {
		switch input {
		case "1":
			c.demandeSectionCritique()
			for _, lot := range c.auctions {
				println(toStringLot(lot))
			}
			c.relacherSectionCritique()
		case "2":
			idAuction := getValueFromUser(input, "enter the auction id")
			c.demandeSectionCritique()
			isFound, pos := c.search(idAuction)
			c.relacherSectionCritique()
			if isFound == true {
				newValue := getValueFromUser(input, "How much to bid?")
				intNewValue, newValueErr := strconv.Atoi(newValue)
				if newValueErr == nil {
					c.demandeSectionCritique()
					if intNewValue > c.auctions[pos].Price {
						println("You're winning the auction")
						c.subscriptions = append(c.subscriptions, Auction.Subscription{
							Auction: c.auctions[pos],
						})
						for _, item := range c.subscriptions {
							if item.Auction.IdName == c.auctions[pos].IdName {
								println("new winning for the auction : " + item.Auction.Nom + " id : " + item.Auction.IdName)
							}
							message := []string{config.NOTIFY_SUBSCRIBED, fmt.Sprint(idAuction)}
							c.chanToMutex <- strings.Join(message, ",")
						}
						c.auctions[pos].Winner = c.user
						message := []string{config.CHANGE_WINNER, fmt.Sprint(c.user), fmt.Sprint(c.auctions[pos].IdName)}
						c.chanToMutex <- strings.Join(message, ",")
						c.auctions[pos].Price = intNewValue
						message = []string{config.CHANGE_PRICE, fmt.Sprint(intNewValue), fmt.Sprint(c.auctions[pos].IdName)}
						c.chanToMutex <- strings.Join(message, ",")
					}
					c.relacherSectionCritique()
				}
			} else {
				println("id not found")
			}
		case "3":
			println("1 : subscribe to an auction by id\n" +
				"2 : subscribe to all new auctions\n" +
				"3 : unsubscribe to an auction by id\n" +
				"4 : unsubscribe to all new auctions\n")
			fmt.Scanln(&input)
			switch input {
			case "1":
				idAuction := getValueFromUser(input, "enter the auction id")
				c.demandeSectionCritique()
				isFound, pos := c.search(idAuction)
				c.relacherSectionCritique()
				if isFound == true {
					c.demandeSectionCritique()
					c.subscriptions = append(c.subscriptions, Auction.Subscription{
						Auction: c.auctions[pos],
					})
					c.relacherSectionCritique()
				} else {
					println("id not found")
				}
			case "2":
				c.newAuction = true
			case "3":
				idAuction := getValueFromUser(input,"enter the auction id")
				c.demandeSectionCritique()
				isFound, pos := c.search(idAuction)
				c.relacherSectionCritique()
				if isFound == true {
					c.demandeSectionCritique()
					for _, item := range c.subscriptions {
						if item.Auction.IdName == c.auctions[pos].IdName {
							c.removeSubscription(item)
						}
					}
					c.relacherSectionCritique()
				} else {
					println("id not found")
				}
			case "4":
				c.newAuction = false
			default:
			}
		case "4":
			auctionName := getValueFromUser(input, "enter auction name")
			valMin, errValMin := strconv.Atoi(getValueFromUser(input, "enter bid starting price"))
			temps, errTemps := strconv.Atoi(getValueFromUser(input, "enter time the auction may last in minutes"))

			time := time.Now().Add(time.Duration(temps*6e10))
			c.demandeSectionCritique()
			if errTemps == nil && errValMin == nil {
				toAdd := Auction.Auction{
					Nom:         auctionName,
					IdName:      strconv.Itoa(c.auctionIds),
					Temps:       time,
					Price:		 valMin,
					Provider:    c.user,
					Winner:		 "",
				}
				// subscription for the vendor
				c.subscriptions = append(c.subscriptions, Auction.Subscription{
					Auction: toAdd,
				})
				// subscriptions for the subscribed to all new auctions
				/*
				for _, item := range c.newAuction {
					println("new auction! : " + toAdd.Nom + " id : " + toAdd.IdName)
				}*/
				message := []string{config.ADD_AUCTION, fmt.Sprint(auctionName), fmt.Sprint(c.auctionIds),
					fmt.Sprint(temps), fmt.Sprint(valMin), fmt.Sprint(c.user),fmt.Sprint("")}
				c.chanToMutex <- strings.Join(message, ",")
				message = []string{config.INCREMENT_AUCTION_IDS}
				c.chanToMutex <- strings.Join(message, ",")
				message = []string{config.NOTIFY_NEWAUCTION,  fmt.Sprint(auctionName)}
				c.chanToMutex <- strings.Join(message, ",")
				c.relacherSectionCritique()
			}
		default :
			println("wrong entry")
		}
		println(menu())
		fmt.Scanln(&input)
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
			c.demandeEnCour = true
		} else if strings.HasPrefix(mutexMsg, config.VALUE) {

			utils.PrintMessage(c.id, partName, " New value received, updated \n")
			splitMsg := strings.Split(mutexMsg, ",")
			value, err := strconv.ParseUint(splitMsg[1], 10, 32)

			if err == nil{
				c.setValueSC(uint32(value))
			}
		} else if strings.HasPrefix(mutexMsg, config.ADD_USER) {
			splitMsg := strings.Split(mutexMsg, ",")
			c.users = append(c.users, splitMsg[1])
		}  else if strings.HasPrefix(mutexMsg, config.CHANGE_PRICE) {
			splitMsg := strings.Split(mutexMsg, ",")
			for i := range c.auctions {
				if c.auctions[i].IdName == splitMsg[2] {
					value, err := strconv.Atoi(splitMsg[1])

					if err == nil{
						c.auctions[i].Price = value
					}
				}
			}
		}  else if strings.HasPrefix(mutexMsg, config.CHANGE_WINNER) {
			splitMsg := strings.Split(mutexMsg, ",")
			for i := range c.auctions {
				if c.auctions[i].IdName == splitMsg[2] {
					c.auctions[i].Winner = splitMsg[1]
				}
			}
		} else if strings.HasPrefix(mutexMsg, config.ADD_AUCTION) {
			splitMsg := strings.Split(mutexMsg, ",")
			temps, err := strconv.Atoi(splitMsg[3])
			price, errPrice := strconv.Atoi(splitMsg[4])

			time := time.Now().Add(time.Duration(temps * 6e10))
			if err == nil && errPrice == nil {

				toAdd := Auction.Auction{
					Nom:      splitMsg[1],
					IdName:   splitMsg[2],
					Temps:    time,
					Price:    price,
					Provider: splitMsg[5],
					Winner:   splitMsg[6],
				}
				c.auctions = append(c.auctions, toAdd)
			}
		} else if strings.HasPrefix(mutexMsg, config.INCREMENT_AUCTION_IDS) {
			c.auctionIds++
		} else if strings.HasPrefix(mutexMsg, config.NOTIFY_SUBSCRIBED) {
			splitMsg := strings.Split(mutexMsg, ",")
			for _, item := range c.auctions {
				if item.IdName == splitMsg[1] {
					println("new winner for the auction : " + item.Nom + " id : " + item.IdName)
				}
			}
		} else if strings.HasPrefix(mutexMsg, config.NOTIFY_END) {
			splitMsg := strings.Split(mutexMsg, ",")
			for _, item := range c.auctions {
				if item.IdName == splitMsg[1] {
					if item.Winner == "" {
						println("The auction : " + item.IdName + " ended with no winner")
					} else {
						println("The auction : " + item.IdName + " ended, the winner is : " + item.Winner)
					}
				}
			}

		}else if strings.HasPrefix(mutexMsg, config.NOTIFY_NEWAUCTION) {
			splitMsg := strings.Split(mutexMsg, ",")
			if c.newAuction {
				println("New auction : " + splitMsg[1])
			}

		} else if strings.HasPrefix(mutexMsg, config.REMOVE_AUCTION) {
			splitMsg := strings.Split(mutexMsg, ",")
			for i, item := range c.auctions {
				if item.IdName == splitMsg[1] {
					c.auctions = append(c.auctions[:i], c.auctions[i+1:]...)
				}
			}
		}
	}
}

func (c *Client) demandeSectionCritique() {
	c.chanToMutex <- config.DEMANDE
	for !c.demandeEnCour {
	}
}

func (c *Client) relacherSectionCritique() {
	c.demandeEnCour = false
	c.chanToMutex <- config.FIN
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

func toStringLot(l Auction.Auction) string {
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
	for _, item := range c.users {
		if item == user {
			return true
		}
	}
	return false
}

func (c *Client) auctionGarbageCollector() {
	for {
		time.Sleep(time.Second*10)
		c.demandeSectionCritique()
		for _, item := range c.auctions {
			if item.Temps.Sub(time.Now()) < 0 {
				message := []string{config.NOTIFY_END, fmt.Sprint(item.IdName)}
				c.chanToMutex <- strings.Join(message, ",")
				message = []string{config.REMOVE_AUCTION, fmt.Sprint(item.IdName)}
				c.chanToMutex <- strings.Join(message, ",")
			}
		}
		c.relacherSectionCritique()
	}
}


