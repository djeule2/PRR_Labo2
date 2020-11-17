package Auction

import (
	"time"
)

type Auction struct {
	Nom      string
	IdName   string
	Temps    time.Time
	Price    int
	provider string
	winner   string
}

type Subscription struct {
	ch chan string
	auction Auction
}
type NewAuctionSub struct {
	user string
	ch chan string
}



