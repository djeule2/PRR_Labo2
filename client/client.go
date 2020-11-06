package client

type Client struct {
	id				uint32
	chanFromMutex 	chan string
	chanToMutex 	chan string
	
}

/*
* Initialise un nouveau processus Client à partir de son id, deux channnel pour
* communiquer avec mutex et un tableau de message auquel il pourra demander la section critique
*/
func NewClient(id uint32, chanFronMutex chan string, chanToMutex chan string, messageTab []int) *Client  {
	client := new(Client)
	client.id = id
	client.chanFromMutex = chanFronMutex
	client.chanToMutex = chanToMutex
	return client
}

func (c *Client)traitmentMessage()  {

	
}