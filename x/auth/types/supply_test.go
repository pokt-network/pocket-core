package types

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"

	sdk "github.com/pokt-network/pocket-core/types"

	"github.com/stretchr/testify/require"
)

func TestSupplyMarshalYAML(t *testing.T) {
	supply := DefaultSupply()
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.OneInt()))
	supply = supply.Inflate(coins)

	bz, err := yaml.Marshal(supply)
	require.NoError(t, err)
	bzCoins, err := yaml.Marshal(coins)
	require.NoError(t, err)

	want := fmt.Sprintf(`total:
%s`, string(bzCoins))

	require.Equal(t, want, string(bz))
	require.Equal(t, want, supply.String())
}
