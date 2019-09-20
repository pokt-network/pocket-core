package dispatch

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/x/dispatch"
	"github.com/pokt-network/pocket-core/x/session"
)

var _ = Describe("Dispatch", func() {
	PDescribe("Dispatch Peers", func() {
		PContext("Insufficient tenermint peers", func() {
			PIt("should return insufficient tendermint peers error", func() {

			})
		})
		PContext("Insufficient alive peers", func() {
			PIt("should return insufficient alive peers error", func() {

			})
		})
		PContext("Sufficient peers", func() {
			PIt("should return nil error and a valid peers object", func() {

			})
		})
	})
	Describe("Dispatch Session Generation", func() {
		Context("Invalid application", func() {
			It("should return DispatchSessionGeneration error", func() {
				_, err := dispatch.DispatchSession(session.SessionAppPubKey(invalidApplication), validNonNativeChain)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring(dispatch.SessionGenerationError.Error()))
			})
		})
		Context("Invalid non-native chain", func() {
			It("should return DispatchSessionGeneration error", func() {
				_, err := dispatch.DispatchSession(validApplication, invalidNonNativeChain)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring(dispatch.SessionGenerationError.Error()))
			})
		})
		Context("Valid seed data", func() {
			It("should return nil error and a valid session", func() {
				sess, err := dispatch.DispatchSession(validApplication, validNonNativeChain)
				Expect(err).To(BeNil())
				Expect(sess).ToNot(BeNil())
				Expect(*sess).ToNot(Equal(&session.Session{}))
			})
		})
	})
})
