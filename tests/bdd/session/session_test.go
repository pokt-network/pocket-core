package session_test

import (
	"github.com/google/flatbuffers/go"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/legacy"
)

// ************************************************************************************************************
// Milestone: Sessions
//
// Tentative Timeline (4-6 weeks)
//
// Unanswered Questions?
// - Where is the seed data coming from? 99% sure -> *WORLD STATE*
// - Is service request validation ddos safe?
// - Forking behavior
// TODO deciding on GID format
// ************************************************************************************************************

var _ = Describe("Session", func() {

	Describe("Session Creation \\ Computing", func() {

		devid := []byte(legacy.SHA3FromString("foo"))
		blockhash := legacy.SHA3FromString("foo")
		requestedChain := legacy.Blockchain{Name: "eth", NetID: "1", Version: "1"}
		marshalBC, err := legacy.MarshalBlockchain(flatbuffers.NewBuilder(0), requestedChain)
		capacity := 100
		if err != nil {
			Fail(err.Error())
		}
		requestedChainHash := legacy.SHA3FromBytes(marshalBC)
		absPath, _ := filepath.Abs("../fixtures/xsmallnodepool.json")
		nodelist, err := legacy.FileToNodes(absPath)
		if err != nil {
			Fail(err.Error())
		}

		Context("Invalid SessionSeed Data", func() {

			Context("Parameters are missing or null", func() {

				Context("Missing Devid", func() {
					NoDevIDSeed := legacy.SessionSeed{BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist, Capacity: capacity}
					It("should return missing devid error", func() {
						_, err := legacy.NewSession(NoDevIDSeed)
						Expect(err).To(Equal(legacy.NoDevIDError))
					})
				})

				Context("Missing Blockhash", func() {
					NoBlockhashSeed := legacy.SessionSeed{DevID: devid, RequestedChain: requestedChainHash, NodeList: nodelist, Capacity: capacity}
					It("should return missing blockhash error", func() {
						_, err := legacy.NewSession(NoBlockhashSeed)
						Expect(err).To(Equal(legacy.NoBlockHashError))
					})
				})

				Context("Missing Requested Chain", func() {
					NoRequestedChain := legacy.SessionSeed{DevID: devid, BlockHash: blockhash, NodeList: nodelist, Capacity: capacity}
					It("should return missing requested chain error", func() {
						_, err := legacy.NewSession(NoRequestedChain)
						Expect(err).To(Equal(legacy.NoReqChainError))
					})
				})

				Context("Missing Nodelist", func() {
					NoNodeListSeed := legacy.SessionSeed{DevID: devid, BlockHash: blockhash, RequestedChain: requestedChainHash, Capacity: capacity}
					It("should return missing nodelist error", func() {
						_, err := legacy.NewSession(NoNodeListSeed)
						Expect(err).To(Equal(legacy.NoNodeListError))
					})
				})

				Context("Missing Capacity", func() {
					NoCapacitySeed := legacy.SessionSeed{DevID: devid, BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return missing capacity error", func() {
						_, err := legacy.NewSession(NoCapacitySeed)
						Expect(err).To(Equal(legacy.NoCapacityError))
					})
				})
			})

			Context("Devid is incorrect...", func() {

				Context("Devid is incorrect format", func() {
					invalidDevIDSeed := legacy.SessionSeed{DevID: []byte("invalidtest"), BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist, Capacity: capacity}
					It("should return `invalid developer id` error", func() {
						_, err := legacy.NewSession(invalidDevIDSeed)
						Expect(err).To(Equal(legacy.InvalidDevIDFormatError))
					})
				})

				Context("Devid is not found in world state", func() {

					PIt("should error", func() {
						// TODO need a world state
					})
				})
			})

			Context("Block Hash is incorrect...", func() {

				Context("Not a valid block hash format", func() {
					invalidBlockHashFormatSeed := legacy.SessionSeed{DevID: devid, BlockHash: []byte("foo"), RequestedChain: requestedChainHash, NodeList: nodelist, Capacity: capacity}
					It("should return `invalid block hash` error", func() {
						_, err := legacy.NewSession(invalidBlockHashFormatSeed)
						Expect(err).To(Equal(legacy.InvalidBlockHashFormatError))
					})
				})

				PContext("Block hash is expired", func() {

					PIt("should error", func() {
						// TODO need a world state
					})
				})
			})

			Context("Requested Blockchain is invalid...", func() {

				Context("No nodes are associated with a blockchain", func() {
					noNodesSeed := legacy.SessionSeed{DevID: devid, BlockHash: blockhash, RequestedChain: legacy.SHA3FromString("foo"), NodeList: nodelist, Capacity: capacity}
					It("should return `invalid blockchain` error", func() {
						_, err := legacy.NewSession(noNodesSeed)
						Expect(err).To(Equal(legacy.InsufficientNodesError))
					})
				})
			})
		})

		Context("Valid SessionSeed Data", func() {
			absPath, _ := filepath.Abs("../fixtures/mediumnodepool.json")
			validSeed, _ := legacy.NewSessionSeed(devid, absPath, requestedChainHash, blockhash, capacity)
			s, err := legacy.NewSession(validSeed)
			It("should not have returned any error", func() {
				Expect(err).To(BeNil())
			})
			Describe("Generating a valid session", func() {
				It("should generate a session key", func() {
					Expect(s.Key).ToNot(BeNil())
					Expect(len(s.Key)).ToNot(BeZero())
				})

				Describe("Node selection", func() {

					It("should find the core.NODECOUNT closest nodes to the session key", func() {
						Expect(len(s.Nodes)).To(Equal(legacy.NODECOUNT))
					})

					It("should contain no duplicated nodes", func() {
						check := legacy.NewSet()
						for _, node := range s.Nodes {
							Expect(check.Contains(node.GID)).To(BeFalse())
							check.Add(node.GID)
						}
					})

					Describe("SessionNodes in an evenly distributed fashion", func() {

						Context("Small pool of nodes, small number of trials", func() {

							PIt("should result in evenly distributed nodes", func() {
								// TODO using golangs built in random
								// TODO need crypto consideration to make truly random
							})
						})

						Context("Small pool of nodes, large number of trials", func() {

							PIt("should be evenly distributed", func() {
								// TODO using golangs built in random
								// TODO need crypto consideration to make truly random
							})
						})

						Context("Large pool of nodes, small number of trials", func() {

							PIt("should be evenly distributed", func() {
								// TODO using golangs built in random
								// TODO need crypto consideration to make truly random
							})
						})

						Context("Large pool of nodes, large number of trials", func() {

							PIt("should be evenly distributed", func() {
								// TODO using golangs built in random
								// TODO need crypto consideration to make truly random
							})
						})
					})
				})

				Describe("Deterministic from the seed data", func() {

					Context("2 sessions derived from valid same seed data", func() {
						It("should be = and valid", func() {
							s1, _ := legacy.NewSession(validSeed)
							s2, _ := legacy.NewSession(validSeed)
							s3, _ := legacy.NewSession(validSeed)
							s4, _ := legacy.NewSession(validSeed)
							Expect(s1).To(Equal(s2))
							Expect(s2).To(Equal(s3))
							Expect(s3).To(Equal(s4))
						})
					})

					Context("2 sessions derived from different valid seed data", func() {
						validSeed1, _ := legacy.NewSessionSeed(legacy.SHA3FromString("foo"), absPath, requestedChainHash, blockhash, capacity)
						validSeed2, _ := legacy.NewSessionSeed(legacy.SHA3FromString("bar"), absPath, requestedChainHash, blockhash, capacity)
						It("should be != and valid", func() {
							s1, _ := legacy.NewSession(validSeed1)
							s2, _ := legacy.NewSession(validSeed2)
							Expect(s1).ToNot(Equal(s2))
						})
					})
				})

				Describe("Expose session data", func() {
					Context("Node data", func() {
						It("should not be nil or zero value", func() {
							Expect(s.Nodes).ToNot(BeNil())
							Expect(len(s.Nodes)).ToNot(BeZero())
						})
					})
					Context("Blockhash", func() {
						It("should not be nil or zero value", func() {
							Expect(s.BlockHash).ToNot(BeNil())
							Expect(len(s.BlockHash)).ToNot(BeZero())
						})
					})
					Context("Devid", func() {
						It("should not be nil or zero value", func() {
							Expect(s.DevID).ToNot(BeNil())
							Expect(s.DevID).ToNot(BeZero())
						})
					})
					Context("Chain", func() {
						It("should not be nil or zero value", func() {
							Expect(s.Chain).ToNot(BeNil())
							Expect(s.Chain).ToNot(BeZero())
						})
					})
					Context("Key", func() {
						It("should not be nil or zero value", func() {
							Expect(s.Key).ToNot(BeNil())
							Expect(len(s.Key)).ToNot(BeZero())
						})
					})
					Context("Capacity", func() {
						It("should not be nil or zero value", func() {
							Expect(s.Capacity).ToNot(BeZero())
						})
					})
				})
			})
		})
	})
})
