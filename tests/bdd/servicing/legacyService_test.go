package servicing

import (
	"github.com/pokt-network/pocket-core/legacy"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// ************************************************************************************************************
// Milestone: Servicing
//
// Tentative Timeline (8-12 weeks)
//
// Unanswered Questions?
// How are we structuring the token
// ************************************************************************************************************

var (
	chainsfp, _                  = filepath.Abs("../../fixtures/legacy/chains.json")
	brokenchainsfp, _            = filepath.Abs("../../fixtures/legacy/brokenChains.json")
	chainHash, _                 = legacy.GenerateChainHash(legacy.Blockchain{Name: "eth", NetID: "1", Version: "1"})
	unsupportedChainHash, _      = legacy.GenerateChainHash(legacy.Blockchain{Name: "foo", NetID: "1", Version: "1"})
	payload                      = []byte("{\"jsonrpc\":\"2.0\",\"method\":\"net_version\",\"params\":[],\"id\":67}")
	httpMethod                   = []byte("POST")
	path                         = []byte("/testpath")
	token                        = legacy.Token{ExpDate: []byte("foo")}
	privateKey, _                = legacy.NewPrivateKey()
	privateKeyBytes              = legacy.FromECDSA(privateKey)
	publicKey                    = legacy.NewPublicKey(privateKey)
	CompressedPublicKey          = legacy.CompressPublicKey(publicKey.X, publicKey.Y)
	nonce                        = 1
	Relay                        = legacy.Relay{Blockchain: chainHash, Payload: payload, DevID: CompressedPublicKey, Token: token, Method: httpMethod, Path: path, Nonce: nonce}
	devid                        = []byte(legacy.SHA3FromString("foo"))
	gid                          = "foo"
	blockhash                    = legacy.SHA3FromString("foo")
	nodePoolfp, _                = filepath.Abs("../../fixtures/legacy/mediumnodepool.json")
	capacity                     = 100
	validSeed, _                 = legacy.NewSessionSeed(devid, nodePoolfp, chainHash, blockhash, capacity)
	RelayMissingDevID            = legacy.Relay{Blockchain: chainHash, Payload: payload, Token: token, Method: httpMethod, Path: path, Nonce: nonce}
	RelayMissingPayload          = legacy.Relay{Blockchain: chainHash, DevID: CompressedPublicKey, Token: token, Method: httpMethod, Path: path, Nonce: nonce}
	RelayMissingBlockchain       = legacy.Relay{Payload: payload, DevID: CompressedPublicKey, Token: token, Method: httpMethod, Path: path, Nonce: nonce}
	RelayMissingToken            = legacy.Relay{Blockchain: chainHash, Payload: payload, DevID: CompressedPublicKey, Method: httpMethod, Path: path, Nonce: nonce}
	RelayMissingMethod           = legacy.Relay{Blockchain: chainHash, Payload: payload, DevID: CompressedPublicKey, Token: token, Method: httpMethod, Path: path, Nonce: nonce}
	RelayInvalidDevID            = legacy.Relay{Blockchain: chainHash, Payload: payload, DevID: []byte("foo"), Token: token, Method: httpMethod, Path: path, Nonce: nonce}
	RelayMissingNonce            = legacy.Relay{Blockchain: chainHash, Payload: payload, DevID: CompressedPublicKey, Token: token, Method: httpMethod, Path: path}
	RelayUnsupportedChain        = legacy.Relay{Blockchain: unsupportedChainHash, Payload: payload, DevID: CompressedPublicKey, Token: token, Method: httpMethod, Path: path, Nonce: nonce}
	RelayMessageMissingSignature = legacy.RelayMessage{Relay: Relay}
)

var _ = Describe("Servicing", func() {

	Describe("Service Configuration", func() {

		// Contextualize this block on connection type (e.g. HTTP)
		Describe("Configuring a third party blockchain endpoint list", func() {

			Describe("Parsing a json file", func() {

				Context("Able to unmarshal the blockchain list into chain objects", func() {

					It("should return a nil error", func() {
						Expect(legacy.HostedChainsFromFile(chainsfp)).To(BeNil())
					})

					It("should have created a globally accessible list of blockchains", func() {
						Expect(legacy.GetHostedChains().Len()).ToNot(BeZero())
					})
				})

				Context("Unable to unmarshal the blockchain list into a slice of chain objects", func() {

					It("should return unparsable json error", func() {
						Expect(legacy.HostedChainsFromFile(brokenchainsfp)).ToNot(BeNil())
					})
				})
			})
		})

		Describe("Testing connection to each third party blockchain via HTTP", func() {

			Context("Failed connection to a blockchain", func() {

				It("should return unreachable chain error", func() {
					Expect(legacy.HostedChainsFromFile(chainsfp)).To(BeNil())
					Expect(legacy.TestChains().Error()).To(ContainSubstring(legacy.UnreachableAt))
				})
			})

			Context("Every connection succeeded", func() {

				PIt("should return a HostedChains object", func() {
					// clear the chains
					legacy.GetHostedChains().Clear()
					// assuming google is accessible by http
					legacy.GetHostedChains().AddChain(legacy.Chain{Hash: "test", URL: "https://google.com"})
					Expect(legacy.TestChains()).To(BeNil())
					// put the chains.json back as it was
					legacy.GetHostedChains().Clear()
					Expect(legacy.HostedChainsFromFile(chainsfp)).To(BeNil())
				})
			})
		})
	})

	Describe("Initialize servicing", func() {

		Context("Receives a message of a relay request to service from a client", func() {

			Describe("Message validation", func() {

				//Describe("Unmarshal from bytes to fbs", func() {
				//
				//	Context("(the byte array) is able to be unmarshalled into a relay", func() {
				//
				//		It("should return a relay object", func() {
				//			// marshal and unmarshal a relay message object
				//			rm := types.RelayMessage{Relay: Relay, ApplicationSignature: signature}
				//			b, err := core.MarshalRelayMessage(flatbuffers.NewBuilder(0), rm)
				//			Expect(err).To(BeNil())
				//			Expect(b).ToNot(BeNil())
				//			res := core.UnmarshalRelayMessage(b)
				//			Expect(res).ToNot(BeNil())
				//			Expect(res.Relay).To(Equal(Relay))
				//		})
				//	})
				//})

				Describe("Message contents", func() {

					Context("Contains all fields", func() {

						PIt("should return nil error", func() {
							Expect(Relay.ErrorCheck()).To(BeNil())
						})
					})

					Context("Doesn't contain a data payload", func() {

						It("should return missing data payload error", func() {
							Expect(RelayMissingPayload.ErrorCheck()).To(Equal(legacy.MissingPayloadError))
						})
					})

					Context("Doesn't contains a blockchainhash", func() {

						It("should return nil error", func() {
							Expect(RelayMissingBlockchain.ErrorCheck()).To(Equal(legacy.MissingBlockchainError))
						})
					})

					Context("Doesn't contain a devid", func() {

						It("should return missing devid error", func() {
							Expect(RelayMissingDevID.ErrorCheck()).To(Equal(legacy.MissingDevidError))
						})
					})

					Context("Doesn't contain a token", func() {

						It("should return missing token", func() {
							Expect(RelayMissingToken.ErrorCheck()).To(Equal(legacy.InvalidTokenError))
						})
					})

					Context("Doesn't contain a client signature", func() {

						It("should return missing signature error", func() {
							Expect(RelayMessageMissingSignature.ErrorCheck()).To(Equal(legacy.MissingSignatureError))
						})
					})

					Context("Doesn't contain an http method", func() {

						It("should replace the http method with POST", func() {
							RelayMissingMethod.ErrorCheck()
							Expect(RelayMissingMethod.Method).To(Equal([]byte(legacy.DefaultHTTPMethod)))
						})
					})

					Context("Nonce is zero", func() {
						It("should return a zero nonce error ", func() {
							Expect(RelayMissingNonce.ErrorCheck()).To(Equal(legacy.ZeroNonceError))
						})
					})
				})

				Describe("Proper Formatting", func() {

					Context("Doesn't contain a properly formatted devid", func() {

						It("should return improper devid format error", func() {
							Expect(RelayInvalidDevID.ErrorCheck()).To(Equal(legacy.InvalidDevIDError))
						})
					})
				})

				Describe("Field validation", func() {

					Context("Contains a blockchain hash that is not supported by the node", func() {

						It("should return unsupported chain error", func() {
							Expect(RelayUnsupportedChain.ErrorCheck()).To(Equal(legacy.UnsupportedBlockchainError))
						})
					})

					PContext("Contains an invalid token", func() {

						PIt("should return an invalid token error", func() {
							// todo figure out what constitutes an invalid token
							// todo no signature
						})
					})

					Context("A devid/seed that generates a session that doesn't correspond to the service node", func() {
						It("should return an invalid session error", func() {
							s, err := legacy.NewSession(validSeed)
							Expect(err).To(BeNil())
							Expect(s.ValidityCheck(gid)).To(Equal(legacy.InvalidSessionError))
						})
					})
				})
			})
		})
	})

	Describe("Execute the relay", func() {
		resp, err := legacy.RouteRelay(Relay)

		Describe("HTTP", func() {

			Context("request happens successfully", func() {

				PIt("should return the result", func() {
					Expect(err).To(BeNil())
					Expect(resp).ToNot(BeNil())
				})
			})

			Context("request happens unsuccessfully", func() {

				PIt("should return HTTP execute error", func() {
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Describe("Responding", func() {

			Context("The node is able to sign the relay response", func() {

				It("should return nil error", func() {
					// signature, err := core.Sign([]byte(resp), privateKeyBytes)
					signature, err := legacy.Sign(legacy.SHA3FromString("test"), privateKeyBytes)
					Expect(err).To(BeNil())
					Expect(signature).ToNot(BeNil())
				})
			})

			Context("The node is unable to sign the relay response", func() {

				It("should return a signature error", func() {
					_, err := legacy.Sign([]byte(resp), []byte("foo"))
					Expect(err).ToNot(BeNil())
				})
			})

			PContext("The node is able to respond to the client over http", func() {

				PIt("should return nil error", func() {
					// TODO need network
				})
			})

			PContext("The node is unable to respond to the client over http", func() {

				PIt("should return 500 error", func() {
					// TODO need network
				})
			})
		})
	})

})
