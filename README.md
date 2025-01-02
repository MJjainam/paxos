# Paxos (Ongoing)
Paxos is a bedrock algorithm in distributed systems. This project attempts to simulate Paxos in Go.

## Broad Points ensured when writing the code
1. The code must be easy to understand. 
2. Well documented.
3. Highly configurable that allows simulation of failures, adding network delays, etc.
4. Metrics and Observability.



## Glossary
| Term          | Definition                                                                 |
|---------------|----------------------------------------------------------------------------|
|Paxos          | A consensus algorithm for distributed systems                              |
|Actor          | In this repo, the actor means a go routine that can be either Proposer or Acceptor. In real world this would be a node|
|Proposer       | An actor in the algorithm that proposes a value                            | 
|Acceptor       | An actor in the algorithm that potentially accepts the value proposed by proposer | 
|Actor Number   | If there are n actors participating in the paxos, each actor will be identitified from 1..N|
|PRN| Proposal Request Number. If the actor "n" sends its first proposal request, the req number would be 1.n |



## Pseudo Code
### Happy Path 
1. A random actor assumes the role of a proposer and initiates proposal request to all other actors. The proposer will choose itself as a proposer with the probability of x (to be decided)
2. The acceptors return with positive ack if what they have seen is the highest PRN until then. If the PRN is lower then it will return neg ack.
3. If any of the acceptor returns a negative ack, then the proposer relinquishes the proposer role. Otherwise the proposer proposes a value to the actors that have accepted the proposals. 
4. The acceptors return with positive ack if the value is accepted. The acceptors sends this ack to all the actors. 
5. Every actors notes the positve ack by other actors. 
6. The actors remember who accepted what request and if they exceed the majority, then it is considered to be the agreed request number. 
7. The algorithm ends. 

### Failure Scenarios and Mitigations
...


## References
1. Paxos Simplified: https://www.youtube.com/watch?v=SRsK-ZXTeZ0
2. 
