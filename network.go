package main

type Network struct {
	channels []chan Message
}

func newNetwork(n int) Network {
	return Network{
		channels: make([]chan Message, n+1),
	}
}

func (nw Network) sendTo(m Message) {
	nw.channels[m.Sender] <- m
}

func (nw Network) recvFrom(i int) Message {
	msg := <-nw.channels[i]
	return msg
}

func (nw Network) getRecvChan(i int) chan Message {
	return nw.channels[i]
}
