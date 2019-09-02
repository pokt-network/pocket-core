package aat_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// ************************************************************************************************************
// Milestone: https://github.com/pokt-network/pocket-core/milestone/53
// Unanswered Questions?
// 1.- Do we need to define the behaviour of the service node on this test or on a separate test?
// ************************************************************************************************************

var _ = Describe("Application Auth Token", func() {

	PDescribe("Token Generation", func() {})

	PDescribe("Token Parsing", func() {

		PContext("Invalid Token", func() {
			PContext("Attributes are missing or null", func() {
				PContext("Missing version", func{
					PIt("should return missing version error", func() {})
				})

				PContext("Missing message", func{
					PIt("should return missing message error", func() {})
				})

				PContext("Missing appPubKey", func{
					PIt("should return missing appPubKey error", func() {})
				})

				PContext("Missing signature", func{
					PIt("should return missing signature error", func() {})
				})
			})

			PContext("version is incorrect", func() {
				PContext("Version not in semver format", func() {
					PIt("should return invalid version format error", func () {})
				})

				PContext("Version not supported", func() {
					PIt("should return version not supported error", func () {})
				})
			})

			PContext("message is incorrect", func() {
				PContext("Attributes are missing or null", func() {
					PContext("Missing applicationAddress", func{
						PIt("should return missing applicationAddress error", func() {})
					})
				})
			})

			PContext("appPubKey is incorrect", func() {
				PContext("appPubKey has an invalid format", func() {
					PIt("should return invalid appPubKey error", func () {})
				})
			})

			PContext("signature is incorrect", func() {
				PContext("appPubKey doesn't match signature", func() {
					PIt("should return invalid signature error", func () {})
				})

				PContext("signature doesn't match message", func() {
					PIt("should return invalid signature error", func () {})
				})
			})
		})
	})
})
