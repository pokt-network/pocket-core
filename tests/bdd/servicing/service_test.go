package servicing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/x/service"
)

var _ = Describe("Service", func() {
	PDescribe("Service Configuration", func() {})
	Describe("Relay request validation", func() {
		Context("Invalid blockchain", func() {
			Context("Missing blockchain", func() {
				It("should return an empty blockchain error", func() {
					Expect(relayMissingBlockchain.Validate(hostedBlockchains)).To(Equal(service.EmptyBlockchainError))
				})
			})
			Context("Unsupported blockchain", func() {
				It("should return an unsupported blockchain error", func() {
					Expect(relayUnsupportedBlockchain.Validate(hostedBlockchains)).To(Equal(service.UnsupportedBlockchainError))
				})
			})
		})
		Context("Invalid payload", func() {
			Context("Missing payload data", func() {
				It("should return an empty payload data error", func() {
					Expect(relayMissingPayload.Validate(hostedBlockchains)).To(Equal(service.EmptyPayloadDataError))
				})
			})
			PContext("Invalid method", func() {})
			PContext("Invalid path", func() {})
		})
		Context("Invalid application authentication token", func() {
			Context("Invalid token version", func() {
				Context("Missing token version", func() {
					It("should return a missing token version error", func() {
						Expect(relayMissingTokenVersion.Validate(hostedBlockchains).Error()).To(ContainSubstring(service.MissingTokenVersionError.Error()))
					})
				})
				Context("Unsupported token version", func() {
					It("should return an unsupported token version error", func() {
						Expect(relayUnsupportedTokenVersion.Validate(hostedBlockchains).Error()).To(ContainSubstring(service.UnsupportedTokenVersionError.Error()))
					})
				})
			})
			Context("Invalid token message", func() {
				Context("Missing application public key in message body", func() {
					It("should return a missing app pub key error", func() {
						Expect(relayMissingTokenAppPubKey.Validate(hostedBlockchains).Error()).To(ContainSubstring(service.MissingApplicationPublicKeyError.Error()))
					})
				})
				Context("Missing client public key in message body", func() {
					It("should return a missing cli pub key error", func() {
						Expect(relayMissingTokenCliPubKey.Validate(hostedBlockchains).Error()).To(ContainSubstring(service.MissingClientPublicKeyError.Error()))
					})
				})
			})
			Context("Invalid token signature", func() {
				It("should return an invalid token signature error", func() {
					Expect(relayInvalidTokenSignature.Validate(hostedBlockchains).Error()).To(ContainSubstring(service.InvalidTokenSignatureErorr.Error()))
				})
			})
		})
		Context("Invalid Increment Counter", func() {
			Context("Increment count is negative", func() {
				It("should return a negative increment counter relay count error", func() {
					Expect(relayInvalidICCount.Validate(hostedBlockchains).Error()).To(ContainSubstring(service.NegativeICCounterError.Error()))
				})
			})
			Context("Invalid client increment counter signature", func() {
				It("should return an invalid incremnt counter signature error", func() {
					Expect(relayInvalidICSignature.Validate(hostedBlockchains).Error()).To(ContainSubstring(service.InvalidICSignatureError.Error()))
				})
			})
		})
		Context("Valid relay data", func() {
			It("should return no error", func() {
				Expect(validRelay.Validate(hostedBlockchains)).To(BeNil())
			})
		})
	})
	PDescribe("Relay execution", func() {})
	PDescribe("Relay batching", func() {})
})
