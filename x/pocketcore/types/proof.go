package types

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
	"strconv"
	"sync"
	"time"
)

// Proof per relay
type Proof struct {
	Index              int
	SessionBlockHeight int64
	ServicerPubKey     string
	Token              AAT
	Signature          string
}

type PORHeader struct {
	ApplicationPubKey  string
	Chain              string
	SessionBlockHeight int64
}

// ProofOfRelay per application
type ProofOfRelay struct {
	PORHeader
	TotalRelays int64
	Proofs      []Proof // map[index] -> Proofs
}

// structure to map out all proofs
type AllProofs struct {
	M map[string]ProofOfRelay // map[appPubKey+chain+blockheight] -> ProofOfRelay
	l sync.Mutex
}

var (
	globalAllProofs *AllProofs
	apOnce          sync.Once
)

func GetAllProofs() *AllProofs {
	apOnce.Do(func() {
		if globalAllProofs == nil {
			*globalAllProofs = AllProofs{M: make(map[string]ProofOfRelay)}
		}
	})
	return globalAllProofs
}

func (ap AllProofs) AddProof(header PORHeader, p Proof, maxRelays int) sdk.Error {
	var por = ProofOfRelay{}
	porKey := KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(p.SessionBlockHeight)))
	// first check to see if all proofs contain a proof of relay for that specific application and session
	ap.l.Lock()
	defer ap.l.Unlock()
	if _, found := ap.M[porKey]; found {
		por = ap.M[porKey]
	} else {
		// if not found fill in the header info
		por.SessionBlockHeight = header.SessionBlockHeight
		por.ApplicationPubKey = header.ApplicationPubKey
		por.Chain = header.Chain
		por.Proofs = make([]Proof, maxRelays)
		por.TotalRelays = 0
	}
	// check to see if evidence was already stored
	if pf := por.Proofs[p.Index]; pf.Signature != "" {
		return NewDuplicateProofError(ModuleName)
	}
	// else add the proof to the slice
	por.Proofs[p.Index] = p
	// increment total relay count
	por.TotalRelays = por.TotalRelays + 1
	// update POR
	ap.M[porKey] = por
	// punch their ticket
	err := GetAllTix().PunchTicket(header, p.Index)
	if err != nil {
		return err
	}
	return nil
}

func (ap AllProofs) GetProof(header PORHeader, index int) *Proof {
	ap.l.Lock()
	defer ap.l.Unlock()
	por := ap.M[KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(header.SessionBlockHeight)))].Proofs
	if por == nil {
		return nil
	}
	return &por[index]
}

func (ap AllProofs) DeleteProofs(header PORHeader) {
	ap.l.Lock()
	defer ap.l.Unlock()
	porKey := KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(header.SessionBlockHeight)))
	delete(ap.M, porKey)
}

func (p Proof) Validate(maxRelays int64, servicerVerifyPubKey string) sdk.Error {
	// check for negative counter
	if p.Index < 0 {
		return NewNegativeICCounterError(ModuleName)
	}
	if int64(p.Index) > maxRelays {
		return NewMaximumIncrementCounterError(ModuleName)
	}
	// validate the service token
	if err := p.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	// validate the public key correctness
	if p.ServicerPubKey != servicerVerifyPubKey {
		return NewInvalidNodePubKeyError(ModuleName) // the public key is not this nodes, so they would not get paid
	}
	return SignatureVerification(p.Token.ClientPublicKey, p.HashString(), p.Signature)
}
func (p Proof) HashString() string {
	return hex.EncodeToString(SHA3FromString(p.ServicerPubKey + p.Token.HashString() + strconv.Itoa(p.Index) + strconv.Itoa(int(p.SessionBlockHeight)))) // todo standardize
}
func (p Proof) Hash() []byte {
	return SHA3FromString(p.ServicerPubKey + p.Token.HashString() + strconv.Itoa(p.Index) + strconv.Itoa(int(p.SessionBlockHeight))) // todo standardize
}

func (ph PORHeader) String() string {
	return strconv.Itoa(int(ph.SessionBlockHeight)) + ph.ApplicationPubKey + ph.Chain // todo standardize
}

// Tickets are used to order the relays
type Ticket struct {
	IsTaken    bool
	lastAccess int64
}

type Tickets struct {
	T []Ticket
	l sync.Mutex
}

const Timeout = 10000 // ms todo parameterize

type AllTickets map[string]Tickets // map[keyforPORHeaeer] -> Tickets

var (
	globalAllTickets *AllTickets
	ticketsonce      sync.Once
)

func GetAllTix() *AllTickets {
	ticketsonce.Do(func() {
		if globalAllTickets == nil {
			*globalAllTickets = make(map[string]Tickets)
		}
	})
	return globalAllTickets
}

func (at *AllTickets) GetNextTicket(header PORHeader, maxRelays int) int {
	key := KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(header.SessionBlockHeight)))
	tix, found := (*at)[key]
	if !found {
		// create a new tickets
		tix = *NewTickets(maxRelays, Timeout)
	}
	res := tix.GetNextTicket()
	(*at)[key] = tix // persist it to the allTix structure
	return res
}

func (at *AllTickets) PunchTicket(header PORHeader, ticketNumber int) sdk.Error {
	key := KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(header.SessionBlockHeight)))
	tix, found := (*at)[key]
	if !found {
		return NewTicketsNotFoundError(ModuleName) // should not happen
	}
	err := tix.PunchTicket(ticketNumber)
	if err != nil {
		return err
	}
	(*at)[key] = tix
	return nil
}

func NewTickets(maxRelays, timeout int) (tix *Tickets) {
	tix = &Tickets{T: make([]Ticket, maxRelays)}
	go func() { // todo are tickets closed upon deletion of tickets structure?
		for now := range time.Tick(time.Second) {
			tix.l.Lock()
			for _, ticket := range tix.T {
				if now.Unix()-ticket.lastAccess > int64(timeout) && ticket.lastAccess != -1 {
					ticket.IsTaken = false
				}
			}
			tix.l.Unlock()
		}
	}()
	return
}

func (tix *Tickets) GetNextTicket() int {
	tix.l.Lock()
	defer tix.l.Unlock()
	for i, ticket := range tix.T {
		if !ticket.IsTaken && ticket.lastAccess != -1 {
			ticket.lastAccess = time.Now().Unix()
			ticket.IsTaken = true
			tix.T[i] = ticket
			return i
		}
	}
	return -1 // out of tickets
}

func (tix *Tickets) PunchTicket(ticketNumber int) sdk.Error { // todo collisions possible
	tix.l.Lock()
	defer tix.l.Unlock()
	ticket := tix.T[ticketNumber]
	if ticket.lastAccess == -1 {
		return NewDuplicateTicketError(ModuleName)
	}
	ticket.lastAccess = -1
	tix.T[ticketNumber] = ticket
	return nil
}
