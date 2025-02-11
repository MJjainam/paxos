package main

import (
	"fmt"
	"time"
)

const InactivityTimePeriod time.Duration = 2 * time.Second

/*
The actors will have four stages. Refer stage.go
*/

type Actor struct {
	N int
	Network
	/*
		Data is like a DB store

		The latest active cycle running is stored in "p_latest_proposal". The value is stored as <cycle>.<propser number>
		The previous cycle's

	*/
	Stage
	LatestPRN   PRN
	PValRespMap map[PRN]int //PRN to Count
	// Data      map[string]interface{}
}

func (a Actor) getRecvChan() chan Message {
	return a.Network.getRecvChan(a.N)
}

func (a Actor) Run() {

	//Begin reading from all the from channels
	messageChan := a.getRecvChan()

	//Run a continuous for loop
	for {
		//All actors will accept message via their from channels
		timer := time.NewTimer(InactivityTimePeriod)
	loop:
		for {
			select {
			case msg := <-messageChan:
				timer.Reset(InactivityTimePeriod)
				a.ProcessMessage(msg)

			case <-timer.C:
				fmt.Println("Timer expired")
				break loop
			}
		}

		//Proposal initiation after 2 seconds of inactivity
		//Send a proposal
		a.sendLeaderProposal()

	}

}

func (a Actor) sendLeaderProposal() {

	for i := 1; i <= TOTAL_ACTORS; i++ {
		//Get the latest cycle number, if none, start with 1
		latestCycle := a.getLatestPRN().getCycle()
		prn := newPRN(latestCycle+1, a.N)
		msg := Message{
			Sender:      a.N,
			Receiver:    i,
			MessageType: P_INIT,
			PInitMessage: PInitMessage{
				prn: prn,
			},
		}

		a.setLatestPRN(prn)
		a.sendTo(msg)

		msg = a.recvFrom(i) //Todo: handle the case when the message is not received for long time. Maybe the actor is down. Maybe add a timeout to the recvFrom function?
		if msg.MessageType != ACK {
			//Expected to receive ACK message, but received some other message.
			//stop P_INIT cycle, and process this message instead.
			a.ProcessMessage(msg)
		} else if msg.AckMessage.Success == false {
			//Some actor sent the N Ack. Stop Proposer intiation
			return
		}
	}

	//P_INIT phase completed successfully, propose itself to become the new leader.

}

func (a Actor) ParallelRead() chan Message {
	var messageChannel chan Message
	for _, ch := range a.From {
		go func(c chan Message) {
			for {
				msg := <-c
				messageChannel <- msg
			}
		}(ch)
	}
	return messageChannel
}

// There are two types of message this can process.
// 1. P_INIT message (Proposal Request message)
// 2. P_VAL message (Proposal Value message)
func (a Actor) ProcessMessage(msg Message) Message {
	switch msg.MessageType {

	case P_INIT:
		currentCycle := msg.PInitMessage.prn.getCycle()

		if prnObject, ok := a.Data["p_latest_proposal"]; ok {
			prn := prnObject.(PRN)
			cycle, _ := prn.Parse()
			if currentCycle >= cycle {
				//return negative ACK. There is a cycle number already running.
				return Message{
					Sender:      a.N,
					Receiver:    msg.Sender,
					MessageType: P_INIT_RESP,
					PInitRespMessage: PInitRespMessage{
						Success: false,
					},
				}
			}

		}
		//The cycle number is greater, hence the initiation has to begin.
		a.Data["p_latest_proposal"] = msg.PInitMessage.prn

	case P_VAL:
		//Check the stage. If the stage is FINAL, then reject.

		//Accept the message and return positive ack. Also, broadcast the acceptance to all the actors.

		/* Broadcast messages to all the actors of the acceptance */
		a.broadcastAcceptance(msg.prn)
		return Message{
			Sender:      a.N,
			Receiver:    msg.Sender,
			MessageType: P_VAL_RESP,
			PValRespMessage: PValRespMessage{
				Success: true,
				PRN:     msg.prn,
			},
		}

	}

}

/*Broadcast message to all the actors that a PRN value is accepted*/
func (a Actor) broadcastAcceptance(prn PRN) {
	// for i := 1; i <= TOTAL_ACTORS; i++ {
	// 	a.To[i] <- Message{
	// 		Sender:      a.N,
	// 		Receiver:    i,
	// 		MessageType: P_VAL_ACCEPT,
	// 		PValAcceptMessage: PValAcceptMessage{
	// 			PRN: prn,
	// 		},
	// 	}
	// }
}

/*
Return the current PRN if the P_INIT has already happened on the Actor. Otherwise returns empty string ""
*/
func (a Actor) getLatestPRN() PRN {
	if prnObject, ok := a.Data["p_latest_proposal"]; ok {
		prn := prnObject.(PRN)
		return prn
	} else {
		return "0.0"
	}
}

/*
There is P_INIT message received or started by the actor. Update it
*/
func (a Actor) setLatestPRN(prn PRN) {
	a.Data["p_latest_proposal"] = prn
}

/*
This function will set the value for the PRN. This information usually comes from P_VAL request
*/
func (a Actor) setPRNValue(prn PRN, value int) {

	var values map[PRN]int

	if valuesObj, ok := a.Data["values"]; ok {
		values = valuesObj.(map[PRN]int)
	} else {
		values = make(map[PRN]int)
	}
	values[prn] = value
	a.Data["values"] = values

}

// func (a Actor)
