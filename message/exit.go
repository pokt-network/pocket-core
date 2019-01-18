package message

import "github.com/pokt-network/pocket-core/const"

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

// "SendExitMSG" sends a message to the centralized dispatcher to unregister as service node.
func SendExitMSG() {
	m := NewExitMSG()
	SendMessage(RELAY, m, _const.DISPATCHIP, ExitPL{})
}
