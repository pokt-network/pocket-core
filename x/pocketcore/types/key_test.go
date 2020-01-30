package types

//
//func TestKeyForClaim(t *testing.T) {
//	ctx := newContext(t, false)
//	appPubKey := getRandomPubKey().RawString()
//	ethereum, err := NonNativeChain{
//		Ticker:  "eth",
//		Netid:   "4",
//		Version: "v1.9.9",
//		Client:  "geth",
//		Inter:   "",
//	}.HashString()
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	sh := SessionHeader{
//		ApplicationPubKey:  appPubKey,
//		Chain:              ethereum,
//		SessionBlockHeight: 1,
//	}
//	// invalid session header
//	invalidSessHeader := SessionHeader{
//		ApplicationPubKey:  appPubKey,
//		Chain:              ethereum,
//		SessionBlockHeight: -1,
//	}
//	assert.Panics(t, func() { KeyForClaim(ctx, getRandomValidatorAddress(), invalidSessHeader) })
//	// invalid address
//	invalidAddr := types.Address{}
//	assert.Panics(t, func() { KeyForClaim(ctx, invalidAddr, sh) })
//	key := KeyForClaim(ctx, getRandomValidatorAddress(), sh)
//	assert.NotNil(t, key)
//	assert.NotEmpty(t, key)
//}
//
//func TestKeyForClaims(t *testing.T) {
//	// invalid address
//	invalidAddr := types.Address{}
//	assert.Panics(t, func() { KeyForClaims(invalidAddr) })
//	key := KeyForClaims(getRandomValidatorAddress())
//	assert.NotNil(t, key)
//	assert.NotEmpty(t, key)
//	assert.Len(t, key, 21)
//}
//
//func TestKeyForProof(t *testing.T) {
//	ctx := newContext(t, false)
//	appPubKey := getRandomPubKey().RawString()
//	ethereum, err := NonNativeChain{
//		Ticker:  "eth",
//		Netid:   "4",
//		Version: "v1.9.9",
//		Client:  "geth",
//		Inter:   "",
//	}.HashString()
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	sh := SessionHeader{
//		ApplicationPubKey:  appPubKey,
//		Chain:              ethereum,
//		SessionBlockHeight: 1,
//	}
//	// invalid session header
//	invalidSessHeader := SessionHeader{
//		ApplicationPubKey:  appPubKey,
//		Chain:              ethereum,
//		SessionBlockHeight: -1,
//	}
//	assert.Panics(t, func() { KeyForInvoice(ctx, getRandomValidatorAddress(), invalidSessHeader) })
//	// invalid address
//	invalidAddr := types.Address{}
//	assert.Panics(t, func() { KeyForInvoice(ctx, invalidAddr, sh) })
//	key := KeyForInvoice(ctx, getRandomValidatorAddress(), sh)
//	assert.NotNil(t, key)
//	assert.NotEmpty(t, key)
//}
//
//func TestKeyForProofs(t *testing.T) {
//	// invalid address
//	invalidAddr := types.Address{}
//	assert.Panics(t, func() { KeyForInvoices(invalidAddr) })
//	key := KeyForInvoices(getRandomValidatorAddress())
//	assert.NotNil(t, key)
//	assert.NotEmpty(t, key)
//	assert.Len(t, key, 21)
//}
