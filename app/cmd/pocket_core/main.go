package main

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app/cmd/cli"
	"time"
)

func main() {
	loc, _ := time.LoadLocation("EST")       // use other time zones such as MST, IST
	t := time.Date(2021, time.May, 31, 18, 0, 0, 0, loc)
	sleepDuration := time.Until(t)
	fmt.Println("Sleeping for ", sleepDuration)
	time.Sleep(sleepDuration)
	cli.Execute()
}
