package mesh

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	log2 "log"
	"math/rand"
	"os"
	"time"
)

// isInvalidRelayCode - check if the error code is someone that block incoming relays for current session.
func isInvalidRelayCode(code sdk.CodeType) bool {
	for _, c := range invalidCodes {
		if c == code {
			return true
		}
	}

	return false
}

// getRandomNode - return a random servicer object from the list load at the start
func getRandomNode() *fullNode {
	mutex.Lock()
	address := servicerList[rand.Intn(len(servicerList))]
	mutex.Unlock()
	s, ok := servicerMap.Load(address)
	if !ok {
		return nil
	}
	return s.Node
}

// getAddressFromPubKeyAsString - return an address as string from a public key string
func getAddressFromPubKeyAsString(pubKey string) (string, error) {
	key, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	return sdk.GetAddress(key).String(), nil
}

// getNodeFromAddress - lookup a node from a servicer address
func getNodeFromAddress(address string) *fullNode {
	s, ok := servicerMap.Load(address)

	if !ok {
		return nil
	}

	return s.Node
}

// fileExist - check if file exists or not.
func fileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// sortJSONResponse - sorts json from a relay response
func sortJSONResponse(response string) string {
	var rawJSON map[string]interface{}
	// unmarshal into json
	if err := json.Unmarshal([]byte(response), &rawJSON); err != nil {
		return response
	}
	// marshal into json
	bz, err := json.Marshal(rawJSON)
	if err != nil {
		return response
	}
	return string(bz)
}

// newSdkErrorFromPocketSdkError - return a mesh node sdkErrorResponse from a pocketcore sdk.Error
func newSdkErrorFromPocketSdkError(e sdk.Error) *sdkErrorResponse {
	return &sdkErrorResponse{
		Code:      e.Code(),
		Codespace: e.Codespace(),
		Error:     e.Error(),
	}
}

// newPocketSdkErrorFromSdkError - return a pocketcore sdk.Error from a mesh node sdkErrorResponse
func newPocketSdkErrorFromSdkError(e *sdkErrorResponse) sdk.Error {
	return sdk.NewError(e.Codespace, e.Code, errors.New(e.Error).Error())
}

// servicerIsSupported - use on pocket node side to verify if the address is handled by the running process.
func servicerIsSupported(address string) error {
	if address == "" {
		return errors.New("missing query param address")
	} else {
		// if lean pocket enabled, grab the targeted servicer through the relay proof
		nodeAddress, err := sdk.AddressFromHex(address)
		if err != nil {
			return errors.New("could not convert servicer hex")
		}

		if _, ok := servicerMap.Load(nodeAddress.String()); !ok {
			return errors.New("failed to find correct servicer private key")
		}
	}

	return nil
}

// newWorkerPool - create pond.WorkerPool instance with the right params in place.
func newWorkerPool(name string, strategyName string, maxWorkers, maxCapacity, idleTimeout int) *pond.WorkerPool {
	panicHandler := func(p interface{}) {
		logger.Error(fmt.Sprintf("Worker %s task paniced: %v", name, p))
	}

	var strategy pond.ResizingStrategy

	switch app.GlobalMeshConfig.WorkerStrategy {
	case "lazy":
		strategy = pond.Lazy()
		break
	case "eager":
		strategy = pond.Eager()
		break
	case "balanced":
		strategy = pond.Balanced()
		break
	default:
		log2.Fatal(
			fmt.Sprintf(
				"strategy %s is not a valid option; allowed values are: lazy|eager|balanced",
				strategyName,
			),
		)
	}

	logger.Debug(
		fmt.Sprintf(
			"starting worker %s with MaxWorkers=%d and MaxCapacity=%d",
			name, maxWorkers, maxCapacity,
		),
	)

	return pond.New(
		// size worker dynamically based on amount of servicers.
		maxWorkers, maxCapacity,
		pond.IdleTimeout(time.Duration(idleTimeout)*time.Millisecond),
		pond.PanicHandler(panicHandler),
		pond.Strategy(strategy),
	)
}
