package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type POKT sdk.Coin

func RegisterPOKT(){
	err := sdk.RegisterDenom("pokt", sdk.NewDec(0)) // todo handle error
	if err != nil {
		fmt.Println(err.Error())
	}
}

func NewPOKT(numberOf int64) POKT {
	test := POKT(sdk.NewCoin("pokt", sdk.NewInt(numberOf)))
	return test
}
