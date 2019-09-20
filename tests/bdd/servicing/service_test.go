package servicing

import (
	"fmt"
	"github.com/h2non/gock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/service"
)

var _ = Describe("Service", func() {
	Describe("Service Configuration", func() {
		// Contextualize this block on connection type (e.g. HTTP)
		Describe("Configuring a third party blockchain endpoint list", func() {

			Describe("Parsing a json file", func() {

				Context("Able to unmarshal the blockchain list into chain objects", func() {

					It("should return a nil error", func() {
						Expect(types.HostedChainsFromFile(chainsfile)).To(BeNil())
					})

					It("should have created a globally accessible list of blockchains", func() {
						Expect(types.GetHostedChains().Len()).ToNot(BeZero())
					})
				})

				Context("Unable to unmarshal the blockchain list into a slice of chain objects", func() {

					It("should return unparsable json error", func() {
						Expect(types.HostedChainsFromFile(brokenchainsfile)).ToNot(BeNil())
					})
				})
			})
		})
	})
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
				Expect(validEthRelay.Validate(hostedBlockchains)).To(BeNil())
			})
		})
	})

	Describe("RelayExecution", func() {
		Describe("HTTP", func() {

			Context("request happens successfully", func() {

				It("should return the result", func() {
					// two different endpoints for testing
					defer gock.Off()
					gock.New(GOODENDPOINT).Persist().Get("/").Reply(200).
						JSON(map[string]string{"id": "67", "jsonrpc": "2.0", "result": GOODRESULT})
					goodEndpointResponse, err := validEthRelay.Execute(hostedBlockchains)
					Expect(err).To(BeNil())
					Expect(goodEndpointResponse).To(ContainSubstring(GOODRESULT))
				})
			})

			Context("request happens unsuccessfully", func() {

				It("should return HTTP execute error", func() {
					gock.New(BADENDPOINT).Get("/").Reply(500).
						JSON(map[string]string{})
					resp, err := validBtcRelay.Execute(hostedBlockchains)
					fmt.Println(resp)
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Describe("Responding", func() {

			Context("The node is able to sign the relay response", func() {

				It("should return nil error", func() {
					goodEndpointResponse, err := validEthRelay.Execute(hostedBlockchains)
					rr := service.RelayResponse{
						Signature:   "",
						Response:    goodEndpointResponse, // todo amino encoding
						ServiceAuth: service.ServiceCertificate{ServiceCertificatePayload: service.ServiceCertificatePayload{Counter: 0}},
					}
					err = rr.Sign()
					Expect(err).To(BeNil())
				})
			})

			Context("The node did not sign the relay response", func() {

				It("should return a signature error", func() {
					goodEndpointResponse, err := validEthRelay.Execute(hostedBlockchains)
					rr := service.RelayResponse{
						Signature:   "",
						Response:    goodEndpointResponse, // todo amino encoding
						ServiceAuth: service.ServiceCertificate{ServiceCertificatePayload: service.ServiceCertificatePayload{Counter: 0}},
					}
					err = rr.Validate()
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})
	PDescribe("Relay batching", func() {})
})
