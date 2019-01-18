package message

import "github.com/pokt-network/pocket-core/const"

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

// "SendEntryMSG" sends a message to the centralized dispatcher to register as service node.
func SendEntryMSG() {
	m := NewEnterMSG()
	SendMessage(RELAY, m, _const.DISPATCHIP, EnterPL{})
}
