package Auction

import (
	"time"
)

type Auction struct {
	Nom      string
	IdName   string
	Temps    time.Time
	Price    int
	Provider string
	Winner   string
}

type Subscription struct {
	Auction Auction
}



