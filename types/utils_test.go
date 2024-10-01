package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortJSON(t *testing.T) {
	cases := []struct {
		unsortedJSON string
		want         string
		wantErr      bool
	}{
		// simple case
		{unsortedJSON: `{"cosmos":"foo", "atom":"bar",  "tendermint":"foobar"}`,
			want: `{"atom":"bar","cosmos":"foo","tendermint":"foobar"}`, wantErr: false},
		// failing case (invalid JSON):
		{unsortedJSON: `"cosmos":"foo",,,, "atom":"bar",  "tendermint":"foobar"}`,
			want:    "",
			wantErr: true},
		// genesis.json
		{unsortedJSON: `{"consensus_params":{"block_size_params":{"max_bytes":22020096,"max_txs":100000,"max_gas":-1},"tx_size_params":{"max_bytes":10240,"max_gas":-1},"block_gossip_params":{"block_part_size_bytes":65536},"evidence_params":{"max_age":100000}},"validators":[{"pub_key":{"type":"AC26791624DE60","value":"c7UMMAbjFuc5GhGPy0E5q5tefy12p9Tq0imXqdrKXwo="},"power":100,"name":""}],"app_hash":"","genesis_time":"2018-05-11T15:52:25.424795506Z","chain_id":"test-chain-Q6VeoW","app_state":{"accounts":[{"address":"718C9C23F98C9642569742ADDD9F9AB9743FBD5D","coins":[{"denom":"Token","amount":1000},{"denom":"upokt","amount":50}]}],"upokt":{"pool":{"total_supply":50,"bonded_shares":"0","unbonded_shares":"0","bonded_pool":0,"unbonded_pool":0,"inflation_last_time":0,"inflation":"7/100"},"params":{"inflation_rate_change":"13/100","inflation_max":"1/5","inflation_min":"7/100","goal_bonded":"67/100","max_validators":100,"bond_denom":"upokt"},"candidates":null,"bonds":null}}}`,
			want:    `{"app_hash":"","app_state":{"accounts":[{"address":"718C9C23F98C9642569742ADDD9F9AB9743FBD5D","coins":[{"amount":1000,"denom":"Token"},{"amount":50,"denom":"upokt"}]}],"upokt":{"bonds":null,"candidates":null,"params":{"bond_denom":"upokt","goal_bonded":"67/100","inflation_max":"1/5","inflation_min":"7/100","inflation_rate_change":"13/100","max_validators":100},"pool":{"bonded_pool":0,"bonded_shares":"0","inflation":"7/100","inflation_last_time":0,"total_supply":50,"unbonded_pool":0,"unbonded_shares":"0"}}},"chain_id":"test-chain-Q6VeoW","consensus_params":{"block_gossip_params":{"block_part_size_bytes":65536},"block_size_params":{"max_bytes":22020096,"max_gas":-1,"max_txs":100000},"evidence_params":{"max_age":100000},"tx_size_params":{"max_bytes":10240,"max_gas":-1}},"genesis_time":"2018-05-11T15:52:25.424795506Z","validators":[{"name":"","power":100,"pub_key":{"type":"AC26791624DE60","value":"c7UMMAbjFuc5GhGPy0E5q5tefy12p9Tq0imXqdrKXwo="}}]}`,
			wantErr: false},
		// from the TXSpec:
		{unsortedJSON: `{"chain_id":"test-chain-1","sequence":1,"fee_bytes":{"amount":[{"amount":5,"denom":"photon"}],"gas":10000},"msg_bytes":{"inputs":[{"address":"696E707574","coins":[{"amount":10,"denom":"atom"}]}],"outputs":[{"address":"6F7574707574","coins":[{"amount":10,"denom":"atom"}]}]},"alt_bytes":null}`,
			want:    `{"alt_bytes":null,"chain_id":"test-chain-1","fee_bytes":{"amount":[{"amount":5,"denom":"photon"}],"gas":10000},"msg_bytes":{"inputs":[{"address":"696E707574","coins":[{"amount":10,"denom":"atom"}]}],"outputs":[{"address":"6F7574707574","coins":[{"amount":10,"denom":"atom"}]}]},"sequence":1}`,
			wantErr: false},

		// non-object values
		{unsortedJSON: "null", want: "null"},
		{unsortedJSON: "42", want: "42"},
		{unsortedJSON: "-3.14", want: "-3.14"},
		{unsortedJSON: "\"str\"", want: "\"str\""},
		{unsortedJSON: "[\"string\",12345,null]", want: "[\"string\",12345,null]"},
		{unsortedJSON: "[[3,1,2],{},[1,2,3]]", want: "[[3,1,2],{},[1,2,3]]"},
	}

	for tcIndex, tc := range cases {
		got, err := SortJSON([]byte(tc.unsortedJSON))
		if tc.wantErr {
			require.NotNil(t, err, "tc #%d", tcIndex)
			require.Panics(t, func() { MustSortJSON([]byte(tc.unsortedJSON)) })
		} else {
			require.Nil(t, err, "tc #%d, err=%s", tcIndex, err)
			require.NotPanics(t, func() { MustSortJSON([]byte(tc.unsortedJSON)) })
			require.Equal(t, got, MustSortJSON([]byte(tc.unsortedJSON)))
		}

		require.Equal(t, string(got), tc.want)
	}
}

func TestTimeFormatAndParse(t *testing.T) {
	cases := []struct {
		RFC3339NanoStr     string
		SDKSortableTimeStr string
		Equal              bool
	}{
		{"2009-11-10T23:00:00Z", "2009-11-10T23:00:00.000000000", true},
		{"2011-01-10T23:10:05.758230235Z", "2011-01-10T23:10:05.758230235", true},
	}
	for _, tc := range cases {
		timeFromRFC, err := time.Parse(time.RFC3339Nano, tc.RFC3339NanoStr)
		require.Nil(t, err)
		timeFromSDKFormat, err := time.Parse(SortableTimeFormat, tc.SDKSortableTimeStr)
		require.Nil(t, err)

		require.True(t, timeFromRFC.Equal(timeFromSDKFormat))
		require.Equal(t, timeFromRFC.Format(SortableTimeFormat), tc.SDKSortableTimeStr)
	}
}
func Test_CompareVersionStrings(t *testing.T) {
	comp, err := CompareVersionStrings("0.9.1.1", "0.10.0")
	assert.Nil(t, err)
	assert.Equal(t, comp, -1)

	comp, err = CompareVersionStrings("0.10.0", "0.9.2")
	assert.Nil(t, err)
	assert.Equal(t, comp, 1)

	comp, err = CompareVersionStrings("1.0", "0.9.9")
	assert.Nil(t, err)
	assert.Equal(t, comp, 1)

	comp, err = CompareVersionStrings("0.0.0.1", "0.0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, comp, 0)

	comp, err = CompareVersionStrings("0.0.0.1.0", "0.0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, comp, 1)

	comp, err = CompareVersionStrings("v0.9.1.1", "0.10.0")
	assert.NotNil(t, err)

	comp, err = CompareVersionStrings("", "1")
	assert.NotNil(t, err)
}

func TestCompareStringMaps(t *testing.T) {
	m1 := map[string]int{}
	m2 := map[string]int{}
	assert.True(t, CompareStringMaps(m1, m2))

	// m1 is non-empty and m2 is empty
	m1["a"] = 10
	m1["b"] = 100
	assert.False(t, CompareStringMaps(m1, m2))

	// m1 and m2 are not empty and identical
	m2["b"] = 100
	m2["a"] = 10
	assert.True(t, CompareStringMaps(m2, m1))

	// m1 is non-empty and m2 is nil
	m2 = nil
	assert.False(t, CompareStringMaps(m1, m2))
	assert.False(t, CompareStringMaps(nil, m1))

	// m1 and m2 are both nil
	m1 = nil
	assert.True(t, CompareStringMaps(m1, m2))

	// Empty and nil maps are identical
	assert.True(t, CompareStringMaps(nil, map[string]int{}))
}
