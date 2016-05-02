package core

type PCStatus int

type PlayerController struct {
	Address InboxAddress `json:"address"`
	ID      string       `json:"id"`
	MatchID uint         `json:"matchID"`
}

func NewPlayerController(addr InboxAddress) *PlayerController {
	c := PlayerController{}
	c.Address = addr
	c.ID = addr.String()
	c.MatchID = 0
	return &c
}
