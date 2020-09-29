package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"net/url"
	"strconv"
	"strings"
)

// TODO shared code among modules below

const (
	httpsPrefix = "https://"
	httpPrefix  = "http://"
	colon       = ":"
	period      = "."
)

func ValidateServiceURL(u string) sdk.Error {
	u = strings.ToLower(u)
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return ErrInvalidServiceURL(ModuleName, err)
	}
	if u[:8] != httpsPrefix && u[:7] != httpPrefix {
		return ErrInvalidServiceURL(ModuleName, fmt.Errorf("invalid url prefix"))
	}
	temp := strings.Split(u, colon)
	if len(temp) != 3 {
		return ErrInvalidServiceURL(ModuleName, fmt.Errorf("needs :port"))
	}
	port, err := strconv.Atoi(temp[2])
	if err != nil {
		return ErrInvalidServiceURL(ModuleName, fmt.Errorf("invalid port, cant convert to integer"))
	}
	if port > 65535 || port < 0 {
		return ErrInvalidServiceURL(ModuleName, fmt.Errorf("invalid port, out of valid port range"))
	}
	if !strings.Contains(u, period) {
		return ErrInvalidServiceURL(ModuleName, fmt.Errorf("must contain one '.'"))
	}
	return nil
}

const (
	NetworkIdentifierLength = 2
)

func ValidateNetworkIdentifier(chain string) sdk.Error {
	// decode string into bz
	h, err := hex.DecodeString(chain)
	if err != nil {
		return ErrInvalidNetworkIdentifier(ModuleName, err)
	}
	// ensure length isn't 0
	if len(h) == 0 {
		return ErrInvalidNetworkIdentifier(ModuleName, fmt.Errorf("net id is empty"))
	}
	// ensure length
	if len(h) > NetworkIdentifierLength {
		return ErrInvalidNetworkIdentifier(ModuleName, fmt.Errorf("net id length is > %d", NetworkIdentifierLength))
	}
	return nil
}
