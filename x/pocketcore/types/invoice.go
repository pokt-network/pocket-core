package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"sync"
)

var (
	globalAllInvoices *Invoices // holds every RelayProof of the node
	allInvoicesOnce   sync.Once // ensure only made once
)

// RelayProof of relay per application
type Invoice struct {
	SessionHeader `json:"invoice_header"` // the session invoiceHeader serves as an identifier for the invoice
	TotalRelays   int64                   `json:"total_relays"` // the total number of relays completed
	Proofs        []RelayProof            `json:"proofs"`       // a slice of RelayProof objects (RelayProof per relay)
}

// generate the merkle root of an invoice
func (i *Invoice) GenerateMerkleRoot() (root HashSum) {
		root, sortedProofs := GenerateRoot(i.Proofs)
	i.Proofs = sortedProofs
	return
}

// generate the merkle RelayProof for an invoice
func (i *Invoice) GenerateMerkleProof(index int) (proofs MerkleProofs, cousinIndex int) {
	return GenerateProofs(i.Proofs, index)
}

// every `invoice` the node holds in memory
type Invoices struct {
	M map[string]Invoice `json:"invoices"` // map[invoiceKey] -> Invoice
	l sync.Mutex         // a lock in the case of concurrent calls
}

// get all invoices the node holds
func GetAllInvoices() *Invoices {
	// only do once
	allInvoicesOnce.Do(func() {
		// if the all proofs object is nil
		if globalAllInvoices == nil {
			// initialize
			globalAllInvoices = &Invoices{M: make(map[string]Invoice)}
		}
	})
	return globalAllInvoices
}

func (i Invoices) GetInvoice(invoiceHeader SessionHeader) (invoice Invoice, found bool) {
	key := invoiceHeader.HashString()
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	invoice, found = i.M[key]
	return
}

func (i Invoices) IsUniqueProof(invoiceHeader SessionHeader, p RelayProof) bool {
	key := invoiceHeader.HashString()
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	if _, found := i.M[key]; found {
		// if RelayProof already stored in allProofs
		invoice := i.M[key]
		// iterate over invoices to see if unique // todo efficiency (store hashes in map)
		for _, proof := range invoice.Proofs {
			if proof.HashStringWithSignature() == p.HashStringWithSignature() {
				return false
			}
		}
	}
	return true
}

// add the RelayProof to the Invoices object
func (i Invoices) AddToInvoice(invoiceHeader SessionHeader, p RelayProof) sdk.Error {
	var invoice Invoice
	// generate the key for this specific RelayProof
	key := invoiceHeader.HashString()
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	if _, found := i.M[key]; found {
		// if RelayProof already stored in allProofs
		invoice = i.M[key]
	} else {
		// if RelayProof is not already stored, initialize all
		invoice.SessionHeader = invoiceHeader
		invoice.Proofs = make([]RelayProof, 0)
		invoice.TotalRelays = 0
	}
	// add RelayProof to the proofs object
	invoice.Proofs = append(invoice.Proofs, p)
	// increment total relay count
	invoice.TotalRelays = invoice.TotalRelays + 1
	// update POR
	i.M[key] = invoice
	return nil
}

func (i Invoices) GetTotalRelays(invoiceHeader SessionHeader) int64 {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// return the proofs object, corresponding to the invoiceHeader
	return i.M[invoiceHeader.HashString()].TotalRelays
}

// retrieve the single RelayProof from the all proofs object
func (i Invoices) GetProof(invoiceHeader SessionHeader, index int) RelayProof {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// return the proofs object, corresponding to the invoiceHeader
	invoice := i.M[invoiceHeader.HashString()].Proofs
	// do a nil check before indexing
	if invoice == nil {
		return RelayProof{}
	}
	// return the RelayProof at specific index
	return invoice[index]
}

// retrieve the proofs from the all proofs object
func (i Invoices) GetProofs(invoiceHeader SessionHeader) []RelayProof {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// return the proofs object, corresponding to the invoiceHeader
	return i.M[invoiceHeader.HashString()].Proofs
}

// delete invoice
func (i Invoices) DeleteInvoice(invoiceHeader SessionHeader) {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// delete the value corresponding to the invoiceHeader
	delete(i.M, invoiceHeader.HashString())
}

// structure used to store the RelayProof after verification
type StoredInvoice struct {
	SessionHeader   `json:"header"`
	ServicerAddress string `json:"address"`
	TotalRelays     int64  `json:"relays"`
}
