package main

/*
Actors will have 4 stages.

1. PINIT stage:
Proposer Initialization phase. In this phase the actor will asume the role of a proposer and send initiation message.
If all the remaining actors send back positive PInitResp message, then the actor is moved to PVal stage
2. PVAL stage:
In PVal stage the actor will send the message to all the actors (in parallel) and expect back acceptance message. If the number of acceptances are > majority, then it moves to the FINISH stage.
3. ACCEPT stage
Acceptors stage. This can move to PINIT stage if it crosses the inactivity period. And moves to FINISH stage if it has received > majority acceptances for a PRN
4. FINISH stage. Done.

Based on what stage the actor is in, appropriate message will be responded.
1. PINIT stage
	a. PInit message (This will be rejected, PInitRespMessage with false success will be responded)
	b. PInitResp message (Expected. Move forward)
	c. PVal message (This is invalid. Respond with Negative PValResp)
	d. PValResp message (This is received when some other actor has accepted the PRN.
		Increment the counter for the PRN. And if the value surpasses majority, assume that as the leader. Move to FINISH stage)

2. PVAL stage
	a. PInit message (This will be rejected by responding with false success)
	b. PInitResp message (This is invalid, as PVAL stage is reached only after accepting all the PInitResp message)
	c. PVal message (Someone else is sending proposal value, reject this)
	d. PValResp message (Same as 1.d )

3. ACCEPT stage
	a. PInit message (This can be accepted or rejected based on the value inside PInit message. If the cycle is lesser than or equal to the last seen cycle then reject it. Or accpet it)
	b. PInitResp message (This is invalid)
	c. PVal message (If the PVal's PRN is whatever is stored in last PRN, then accept. Otherwise reject. Upon accepting this broadcase the PValResp message to all actors)
	d. PValResp message (Same as 1.d)

4. FINISH stage
	a. PInit message. Invalid, reject
	b. PInitResp message. Invalid, reject
	c. PVal message. Invalid, reject
	d. PValResp message. Valid, increment.
*/

type Stage string

const PINIT_STAGE Stage = "PINIT"
const PVAL_STAGE Stage = "PVAL"
const ACCEPT_STAGE Stage = "ACCEPT"
const FINISH_STAGE Stage = "FINISH"
const INVALID_STAGE Stage = "INVALID"

type StageHandler interface {
	handle(a Actor, msg Message) (Stage, []Message)
}

type PInitStageHandler struct{}

func (PInitStageHandler) handle(a Actor, msg Message) (Stage, []Message) {
	switch msg.MessageType {
	case P_INIT:
		ret := make([]Message, 1)
		return PINIT_STAGE, append(ret, Message{
			Sender:      a.N,
			Receiver:    msg.Sender,
			MessageType: P_INIT_RESP,
			PInitRespMessage: PInitRespMessage{
				Success: false,
			},
		})
	case P_INIT_RESP:
		if msg.Sender == TOTAL_ACTORS /*last sender*/ || (msg.Sender == TOTAL_ACTORS-1 && a.N == TOTAL_ACTORS) {
			return PVAL_STAGE, nil //Move to PVAL_STAGE
		} else {
			ret := make([]Message, 1)
			return PINIT_STAGE, append(ret, Message{
				Sender:      a.N,
				Receiver:    msg.Sender + 1,
				MessageType: P_INIT,
				PInitMessage: PInitMessage{
					prn: a.LatestPRN,
				},
			})
		}

	case P_VAL:
		if msg.prn.GreaterThan(a.LatestPRN) {
			//Received greater PRN. Respect it and move to ACCEPT stage.
			a.LatestPRN = msg.prn

			//Broadcast to all the actors about the acceptance
			ret := make([]Message, TOTAL_ACTORS)
			for i := 1; i <= len(ret); i++ {
				ret = append(ret, Message{
					Sender:      a.N,
					Receiver:    i,
					MessageType: P_VAL_RESP,
					PValRespMessage: PValRespMessage{
						Success: true,
						PRN:     msg.PValMessage.PRN,
					},
				})
			}
			return ACCEPT_STAGE, ret
		}
		ret := make([]Message, 1)
		return PINIT_STAGE, append(ret, Message{
			Sender:      a.N,
			Receiver:    msg.Sender,
			MessageType: P_VAL_RESP,
			PValRespMessage: PValRespMessage{
				Success: false,
				PRN:     msg.prn,
			},
		})

	case P_VAL_RESP:
		if msg.PValRespMessage.Success == true {
			if _, ok := a.PValRespMap[msg.PValRespMessage.PRN]; ok {
				a.PValRespMap[msg.PValRespMessage.PRN] += 1
			} else {
				a.PValRespMap[msg.PValRespMessage.PRN] = 1
			}

			if a.PValRespMap[msg.PValRespMessage.PRN] > TOTAL_ACTORS/2 {
				return FINISH_STAGE, nil
			}
		}

	}

	return INVALID_STAGE, nil

}
