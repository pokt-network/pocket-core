package session_test

import (
	"github.com/google/flatbuffers/go"
	"path/filepath"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/core"
	"github.com/pokt-network/pocket-core/types"
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
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
		
=======

>>>>>>> servicing milestone addition
		devid := []byte(core.SHA3FromString("foo"))
		blockhash := core.SHA3FromString("foo")
		requestedChain := core.Blockchain{Name: "eth", NetID: "1", Version: "1"}
		marshalBC, err := core.MarshalBlockchain(flatbuffers.NewBuilder(0), requestedChain)
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
		capacity := 100
=======
>>>>>>> servicing milestone addition
		if err != nil {
			Fail(err.Error())
		}
		requestedChainHash := core.SHA3FromBytes(marshalBC)
		absPath, _ := filepath.Abs("../fixtures/xsmallnodepool.json")
		nodelist, err := core.FileToNodes(absPath)
		if err != nil {
			Fail(err.Error())
		}
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
		
		Context("Invalid SessionSeed Data", func() {
			
=======

		Context("Invalid SessionSeed Data", func() {

>>>>>>> servicing milestone addition
			Context("Parameters are missing or null", func() {

				Context("Missing Devid", func() {
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
					NoDevIDSeed := core.SessionSeed{BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist, Capacity: capacity}
					It("should return missing devid error", func() {
						_, err := core.NewSession(NoDevIDSeed)
						Expect(err).To(Equal(core.NoDevIDError))
					})
				})
				
				Context("Missing Blockhash", func() {
					NoBlockhashSeed := core.SessionSeed{DevID: devid, RequestedChain: requestedChainHash, NodeList: nodelist, Capacity: capacity}
					It("should return missing blockhash error", func() {
						_, err := core.NewSession(NoBlockhashSeed)
						Expect(err).To(Equal(core.NoBlockHashError))
					})
				})
				
				Context("Missing Requested Chain", func() {
					NoRequestedChain := core.SessionSeed{DevID: devid, BlockHash: blockhash, NodeList: nodelist, Capacity: capacity}
					It("should return missing requested chain error", func() {
						_, err := core.NewSession(NoRequestedChain)
						Expect(err).To(Equal(core.NoReqChainError))
					})
				})
				
				Context("Missing Nodelist", func() {
					NoNodeListSeed := core.SessionSeed{DevID: devid, BlockHash: blockhash, RequestedChain: requestedChainHash, Capacity: capacity}
					It("should return missing nodelist error", func(){
						_, err := core.NewSession(NoNodeListSeed)
						Expect(err).To(Equal(core.NoNodeListError))
					})
				})
				
				Context("Missing Capacity", func() {
					NoCapacitySeed := core.SessionSeed{DevID: devid, BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return missing capacity error", func(){
						_, err := core.NewSession(NoCapacitySeed)
						Expect(err).To(Equal(core.NoCapacityError))
=======
					NoDevIDSeed := core.SessionSeed{BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return missing devid error", func() {
						_, err := core.NewSession(NoDevIDSeed)
						Expect(err).To(Equal(core.NoDevID))
>>>>>>> servicing milestone addition
					})
				})
			})

			Context("Devid is incorrect...", func() {

				Context("Devid is incorrect format", func() {
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
					invalidDevIDSeed := core.SessionSeed{DevID: []byte("invalidtest"), BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist, Capacity: capacity}
					It("should return `invalid developer id` error", func() {
						_, err := core.NewSession(invalidDevIDSeed)
						Expect(err).To(Equal(core.InvalidDevIDFormatError))
=======
					invalidDevIDSeed := core.SessionSeed{DevID: []byte("invalidtest"), BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return `invalid developer id` error", func() {
						_, err := core.NewSession(invalidDevIDSeed)
						Expect(err).To(Equal(core.InvalidDevIDFormat))
>>>>>>> servicing milestone addition
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
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
					invalidBlockHashFormatSeed := core.SessionSeed{DevID: devid, BlockHash: []byte("foo"), RequestedChain: requestedChainHash, NodeList: nodelist, Capacity:capacity}
					It("should return `invalid block hash` error", func() {
						_, err := core.NewSession(invalidBlockHashFormatSeed)
						Expect(err).To(Equal(core.InvalidBlockHashFormatError))
=======
					invalidBlockHashFormatSeed := core.SessionSeed{DevID: devid, BlockHash: []byte("foo"), RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return `invalid block hash` error", func() {
						_, err := core.NewSession(invalidBlockHashFormatSeed)
						Expect(err).To(Equal(core.InvalidBlockHashFormat))
>>>>>>> servicing milestone addition
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
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
					noNodesSeed := core.SessionSeed{DevID: devid, BlockHash: blockhash, RequestedChain: core.SHA3FromString("foo"), NodeList: nodelist, Capacity: capacity}
					It("should return `invalid blockchain` error", func() {
						_, err := core.NewSession(noNodesSeed)
						Expect(err).To(Equal(core.InsufficientNodesError))
=======
					noNodesSeed := core.SessionSeed{DevID: devid, BlockHash: blockhash, RequestedChain: core.SHA3FromString("foo"), NodeList: nodelist}
					It("should return `invalid blockchain` error", func() {
						_, err := core.NewSession(noNodesSeed)
						Expect(err).To(Equal(core.InsufficientNodes))
>>>>>>> servicing milestone addition
					})
				})
			})
		})
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
		
		Context("Valid SessionSeed Data", func() {
			absPath, _ := filepath.Abs("../fixtures/mediumnodepool.json")
			validSeed, _ := core.NewSessionSeed(devid, absPath, requestedChainHash, blockhash, capacity)
=======

		Context("Valid SessionSeed Data", func() {
			absPath, _ := filepath.Abs("../fixtures/mediumnodepool.json")
			validSeed, _ := core.NewSessionSeed(devid, absPath, requestedChainHash, blockhash)
>>>>>>> servicing milestone addition
			s, err := core.NewSession(validSeed)
			It("should not have returned any error", func() {
				Expect(err).To(BeNil())
			})
			Describe("Generating a valid session", func() {
				It("should generate a session key", func() {
					Expect(s.Key).ToNot(BeNil())
					Expect(len(s.Key)).ToNot(BeZero())
				})

				Describe("Node selection", func() {
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
					
					It("should find the core.NODECOUNT closest nodes to the session key", func() {
						Expect(len(s.Nodes)).To(Equal(core.NODECOUNT))
=======

					It("should find the core.NODECOUNT closest nodes to the session key", func() {
						Expect(len(s.Nodes.ValidatorNodes) + len(s.Nodes.ServiceNodes)).To(Equal(core.NODECOUNT))
>>>>>>> servicing milestone addition
					})

					It("should contain no duplicated nodes", func() {
						check := types.NewSet()
						for _, node := range s.Nodes {
							Expect(check.Contains(node.GID)).To(BeFalse())
							check.Add(node.GID)
						}
					})
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
					
					Describe("SessionNodes in an evenly distributed fashion", func() {
						
=======

					Describe("SessionNodes in an evenly distributed fashion", func() {

>>>>>>> servicing milestone addition
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
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
				
=======

				Describe("Role assignment", func() {

					It("should assign roles to each node", func() {
						Expect(s.Nodes.DelegatedMinter.GID).ToNot(BeEmpty())
						Expect(len(s.Nodes.ValidatorNodes)).To(Equal(core.MAXVALIDATORS))
						Expect(len(s.Nodes.ServiceNodes)).To(Equal(core.MAXSERVICERS))
					})

					PIt("should check the validity of the assigned roles", func() {
						// TODO need blockchain layer
					})

					It("should assign roles to nodes proportional to the protocol guidelines", func() {
						Expect(len(s.Nodes.ValidatorNodes) > len(s.Nodes.ServiceNodes))
					})
				})

>>>>>>> servicing milestone addition
				Describe("Deterministic from the seed data", func() {

					Context("2 sessions derived from valid same seed data", func() {
						It("should be = and valid", func() {
							s1, _ := core.NewSession(validSeed)
							s2, _ := core.NewSession(validSeed)
							s3, _ := core.NewSession(validSeed)
							s4, _ := core.NewSession(validSeed)
							Expect(s1).To(Equal(s2))
							Expect(s2).To(Equal(s3))
							Expect(s3).To(Equal(s4))
						})
					})

					Context("2 sessions derived from different valid seed data", func() {
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
						validSeed1, _ := core.NewSessionSeed(core.SHA3FromString("foo"), absPath, requestedChainHash, blockhash, capacity)
						validSeed2, _ := core.NewSessionSeed(core.SHA3FromString("bar"), absPath, requestedChainHash, blockhash, capacity)
=======
						validSeed1, _ := core.NewSessionSeed(core.SHA3FromString("foo"), absPath, requestedChainHash, blockhash)
						validSeed2, _ := core.NewSessionSeed(core.SHA3FromString("bar"), absPath, requestedChainHash, blockhash)
>>>>>>> servicing milestone addition
						It("should be != and valid", func() {
							s1, _ := core.NewSession(validSeed1)
							s2, _ := core.NewSession(validSeed2)
							Expect(s1).ToNot(Equal(s2))
						})
					})
				})
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
				
				Describe("Expose session data", func() {
					Context("Node data", func() {
						It("should not be nil or zero value", func(){
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
=======

				Describe("Expose node info", func() {

					Describe("For the developer", func() {

						It("should expose the devID", func() {
>>>>>>> servicing milestone addition
							Expect(s.DevID).ToNot(BeNil())
							Expect(s.DevID).ToNot(BeZero())
						})
					})
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
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
=======

					Describe("For the nodes", func() {

						It("should expose the nodes host and port", func() {
							for _, node := range s.Nodes.ValidatorNodes {
								Expect(node.IP).ToNot(BeEmpty())
								Expect(node.Port).ToNot(BeEmpty())
							}
							for _, node := range s.Nodes.ServiceNodes {
								Expect(node.IP).ToNot(BeEmpty())
								Expect(node.Port).ToNot(BeEmpty())
							}
						})

						It("should expose the unique identifier", func() {
							for _, node := range s.Nodes.ValidatorNodes {
								Expect(node.GID).ToNot(BeEmpty())
							}
							for _, node := range s.Nodes.ServiceNodes {
								Expect(node.GID).ToNot(BeEmpty())
							}
						})

						It("should expose the role", func() {
							for _, node := range s.Nodes.ValidatorNodes {
								Expect(node.Role).To(Or(Equal(core.DELEGATEDMINTER), Equal(core.VALIDATE)))
							}
							for _, node := range s.Nodes.ServiceNodes {
								Expect(node.Role).To(Equal(core.SERVICE))
							}
>>>>>>> servicing milestone addition
						})
					})
				})
			})
		})
	})
})