package core

import (
	"log"
)

var _ = log.Println

type InboxClient struct {
	conn  InboxConnection
	id    int
	inbox *Inbox
}

func NewInboxClient(conn InboxConnection, inbox *Inbox, id int) *InboxClient {
	client := InboxClient{}
	client.conn = conn
	client.id = id
	client.inbox = inbox
	return &client
}

func (c *InboxClient) Listen() {
	c.listenRead()
	c.conn.Close()
	c.inbox.RemoveClient(c.id)
}

func (c *InboxClient) Accept(addr InboxAddress) bool {
	return c.conn.Accept(addr)
}

func (c *InboxClient) Write(msg *InboxMessage) {
	go func() {
		e := c.conn.WriteJSON(msg)
		log.Printf("write msg:%v, addr:%v\n", msg.Data, msg.Address)
		if e != nil {
			log.Printf("send message error:%v\n", e.Error())
		}
	}()
}

func (c *InboxClient) listenRead() {
	for {
		select {
		default:
			m := NewInboxMessage()
			e := c.conn.ReadJSON(m)
			if e != nil {
				log.Printf("read message error:%v\n", e.Error())
			}
			log.Printf("got message:%v", m)
			c.inbox.ReceiveMessage(m)
			if m.ShouldCloseConnection {
				return
			}
		}
	}

}
