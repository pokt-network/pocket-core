package servicing

import (
	. "github.com/onsi/ginkgo"
)

// ************************************************************************************************************
// Milestone: Servicing
//
// Tentative Timeline (8-12 weeks)
//
// Unanswered Questions?
// What DSA?
// Who do we send the result to? -> The client and validators
// How are we doing writes?
// ************************************************************************************************************

var _ = Describe("Servicing", func() {
	
	Describe("Service Configuration", func() {
		It("should compute/retrieve blockchain hashes", func() {
			// add code
		})
		It("should configure a third pary blockchain endpoint list", func() {
			// add code
		})
		It("should test connection to each third party blockchain", func() {
			// add code
		})
	})
	
	Describe("Initialize servicing", func() {
		Context("Receives a message of a relay request to service from a client", func() {
			Describe("message validation", func() {
				It("should be able to be unmarshalled", func() {
					// add code
				})
				It("should not contain any missing fields", func() {
					// add code
				})
				It("should contain a valid signature", func() {
					// add code
				})
				Describe("Valid session id", func() {
					It("should be properly formatted", func() {
						// add code
					})
					It("should be valid for the specific service node", func() {
						// add code
					})
				})
				It("should be for a supported blockchain", func() {
					// add code
				})
			})
			It("should identify the proper endpoint to submit the relay request payload to", func() {
			
			})
		})
	})
	
	Describe("Execute the relay", func() {
		Describe("JSON RPC", func() {
			Context("request happens successfully", func() {
				It("should return the result", func() {
					// add code
				})
			})
			Context("request happens unsuccessfully", func() {
				It("should return JSON RPC execute error", func() {
					// add code
				})
			})
		})
		Describe("REST", func() {
			Context("request happens successfully", func() {
				It("should return the result", func() {
					// add code
				})
			})
			Context("request happens unsuccessfully", func() {
				It("should return REST execute error", func() {
					// add code
				})
			})
		})
		Describe("Websockets", func() {
			Context("request happens successfully", func() {
				It("should return the result", func() {
					// add code
				})
			})
			Context("request happens unsuccessfully", func() {
				It("should return WS execute error", func() {
					// add code
				})
			})
		})
		
		Describe("Responding", func() {
			It("should sign the relay response", func() {
				// add code
			})
			It("should send the relay response back to the original requester", func() {
				// add code
			})
			It("should send the relay response to the validators", func() {
				// add code
			})
		})
	})
})
