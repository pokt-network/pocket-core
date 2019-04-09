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
// Do clients need to sign requests?
// How are we doing the devid / signature architecture ?
// Signature capture in the relay response objects? What does the object look like that is going to the validators?
// ************************************************************************************************************

var _ = Describe("Servicing", func() {
	
	Describe("Service Configuration", func() {
		It("should configure a third party blockchain endpoint list", func() {
			// add code
		})
		It("should test connection to each third party blockchain", func() {
			// add code
		})
	})
	
	Describe("Initialize servicing", func() {
		Context("Receives a message of a relay request to service from a client", func() {
			Describe("Message validation", func() {
				Describe("Unmarshal from bytes to fbs", func() {
					It("(the byte array) should be able to be unmarshalled into a relay flatbuffer", func() {
					
					})
				})
				Describe("A relay must contain: a data payload, blockchainhash, devid, signature, and an http method (url param optional)", func() {
					It("should contain a data payload", func() {
					
					})
					It("should contain a blockchainhash", func() {
					
					})
					It("should contain a devid", func() {
					
					})
					It("should contain a signature", func() {
						// TODO signatures
					})
					It("should contain an http method", func() {
					
					})
				})
				Describe("Proper Formatting", func() {
					It("should contain a properly formatted signature", func() {
					
					})
					It("should contain a properly formatted devid", func() {
					
					})
				})
				Describe("Field validation", func() {
					It("should contain a blockchain hash that is supported by the node", func() {
					
					})
					It("should contain a signature that corresponds with the devid", func() {
					
					})
					It("should contain a devid that generates a session that corresponds to the service node", func() {
					
					})
				})
			})
		})
	})
	
	Describe("Service preparation", func() {
		Context("After the request is validated, the service node must prepare for relay execution", func() {
			It("(the service node) should identify the proper endpoint to submit the relay request payload to", func() {
			
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
					// TODO websockets
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
