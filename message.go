package main

/*
Think of this as a library file
*/

/*Types of messages
1. Ack message
2. Proposal Request message
3. Proposal Value message
4. Value Acceptance Broadcast
*/

type MessageType int32

const (
	P_INIT      MessageType = 0 //Proposal Initiate
	P_INIT_RESP MessageType = 1
	P_VAL       MessageType = 2 //Proposal Value
	P_VAL_RESP  MessageType = 3 //Proposal Value Acceptance
)

type Message struct {
	Sender   int
	Receiver int
	MessageType
	PInitMessage
	PInitRespMessage
	PValMessage
	PValRespMessage
}

type PInitMessage struct {
	prn PRN
}

type PInitRespMessage struct {
	Success bool
}

type PValMessage struct {
	PRN
	Value int
}

type PValRespMessage struct {
	Success bool
	PRN
}
