package session_test

import (
	. "github.com/onsi/ginkgo"
)

// ************************************************************************************************************
// Milestone: Sessions
//
// Tentative Timeline (4-6 weeks)
//
// Unanswered Questions?
// - Where is the seed data coming from? 99% sure -> *WORLD STATE*
// - How is the state of the session nodes maintained?
// - Is service request validation ddos safe?
// - Forking behavior
// ************************************************************************************************************

var _ = Describe("Session", func() {

	Describe("Session Creation \\ Computing", func() {

		Context("Invalid Seed Data", func() {
			// n-1 Block Hash (probably world state), blockchains requested [list] (probably  world state), node [list] (probably world state))

			Context("Parameters are missing or null", func() {

				It("should return `missing parameters` error", func() {
					// Code goes here ...
				})
			})

			Context("Devid1 is incorrect...", func() {

				Context("Devid is incorrect format", func() {

					It("should return `invalid developer id` error", func() {
						// Code goes here ...
					})
				})
			})

			Context("Block Hash is incorrect...", func() {
				// n-1 block
				Context("Not a valid block hash format", func() {

					It("should return `invalid block hash` error", func() {
						// Code goes here ...
					})
				})

				Context("Block hash is expired", func() {

					It("should error", func() {
						// Code goes here ...
					})
				})
			})

			Context("Blockchains list invalid...", func() {

				Context("No nodes are associated with a blockchain in the list", func() {

					It("should return `invalid blockchain list` error", func() {
						// Code goes here ...
					})
				})
			})

			Context("Node list is invalid...", func() {

				Context("Node structure is in the incorrect format", func() {

					It("should return `invalid node list` error", func() {
						// Code goes here ...
					})
				})
				// ... any other reason may be incorrect?
			})
		})

		Context("Valid Seed Data", func() {

			Describe("Generating a valid session", func() {

				It("should generate a session key", func() {
					// Code goes here ...
				})

				Describe("Node selection", func() {

					It("should find the 5 closest nodes to the session key", func() {
						// Code goes here ...
					})

					It("should contain no duplicated nodes", func() {
						// Code goes here ...
					})

					Describe("Nodes in an evenly distributed fashion", func() {

						Context("Small pool of nodes, small number of trials", func() {

							It("should result in evenly distributed nodes", func() {
								// Code goes here ...
							})
						})

						Context("Small pool of nodes, large number of trials", func() {

							It("should be evenly distributed", func() {
								// Code goes here ...
							})
						})

						Context("Large pool of nodes, small number of trials", func() {

							It("should be evenly distributed", func() {
								// Code goes here ...
							})
						})

						Context("Large pool of nodes, large number of trials", func() {

							It("should be evenly distributed", func() {
								// Code goes here ...
							})
						})
					})
				})

				Describe("Role assignment", func() {

					It("should assign roles to each node", func() {
						// including service nodes, validator nodes, and the delegated minter
					})

					It("should check the validity of the assigned roles", func() {
						// validators must meet the admission requirements
						// no duplicates
					})

					It("should assign roles to nodes proportional to the protocol guidelines", func() {
						// 1:2 ratio for servicers to validators currently
					})
				})

				Describe("Deterministic from the seed data", func() {

					Context("2 sessions derived from valid same seed data", func() {
						It("should be = and valid", func() {
							// code goes here...
						})
					})

					Context("2 sessions derived from different valid seed data", func() {
						It("should be != and valid", func() {
							// code goes here...
						})
					})
				})

				Describe("Expose node info", func() {

					Describe("For the developer", func() {

						It("should expose the devID", func() {
							// Code goes here ...
						})
					})

					Describe("For the nodes", func() {

						It("should expose the validator node host, port and net protocol", func() {
							// Code goes here ...
						})

						It("should expose the unique identifier", func() {
							// Code goes here ...
						})

						It("should expose the role", func() {
							// Code goes here ...
						})
					})
				})
			})
		})
	})
})
