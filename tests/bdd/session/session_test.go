package session

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/tests/fixtures"
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/session"
)

var _ = Describe("Session", func() {
	// validate the incoming request
	PDescribe("Validation", func() {})
	// generate the session
	Describe("Generation", func() {
		validDeveloper := session.SessionDeveloper(fixtures.GenerateDeveloper())
		validNonNativeChain := fixtures.GenerateNonNativeBlockchain()
		validBlockID := session.SessionBlockID(fixtures.GenerateBlockHash())
		Describe("SessionKey Generation", func() {
			Context("Empty developer seed", func() {
				emptyDeveloper := session.SessionDeveloper(types.Developer{}) // empty developer
				It("should return empty developer public key error", func() {
					_, err := session.NewSessionKey(emptyDeveloper, validNonNativeChain, validBlockID)
					gomega.Expect(err).To(gomega.Equal(session.EmptyDevPubKeyError))
				})
			})
			Context("Empty NonNativeChain", func() {
				emptyNonNativeChain := session.SessionBlockchain{}
				It("should return empty nonNativeChain error", func() {
					_, err := session.NewSessionKey(validDeveloper, emptyNonNativeChain, validBlockID)
					gomega.Expect(err).To(gomega.Equal(session.EmptyNonNativeChainError))
				})
			})
			Context("Empty Block ID", func() {
				emptyBlockID := session.SessionBlockID{}
				It("should return empty blockID error", func() {
					_, err := session.NewSessionKey(validDeveloper, validNonNativeChain, emptyBlockID)
					gomega.Expect(err).To(gomega.Equal(session.EmptyBlockIDError))
				})
			})
			Context("All seed data is valid", func() {
				It("should return nil error", func() {
					_, err := session.NewSessionKey(validDeveloper, validNonNativeChain, validBlockID)
					gomega.Expect(err).To(gomega.BeNil())
				})
			})
		})
		Describe("Service Node Generation", func() {
			validSessionKey, _ := session.NewSessionKey(validDeveloper, validNonNativeChain, validBlockID)
			Context("Empty NonNativeChain", func() {
				emptyNonNativeChain := session.SessionBlockchain{}
				It("should return empty nonNativeChain error", func() {
					_, err := session.NewSessionNodes(emptyNonNativeChain, validSessionKey)
					gomega.Expect(err).To(gomega.Equal(session.EmptyNonNativeChainError))
				})
			})
			Context("Empty SessionKey", func() {
				invalidSessionKey := session.SessionKey{}
				It("should return invalid SessionKey error", func() {
					_, err := session.NewSessionNodes(validNonNativeChain, invalidSessionKey)
					gomega.Expect(err).To(gomega.Equal(session.EmptySessionKeyError))
				})
			})
			Context("Valid Seed data", func() {
				It("should return nil error", func() {
					_, err := session.NewSessionNodes(validNonNativeChain, validSessionKey)
					gomega.Expect(err).To(gomega.BeNil())
				})
			})
		})
		Context("Valid seed data for the session", func() {
			It("should return nil error", func() {
				_, err := session.NewSession(validDeveloper, validNonNativeChain, validBlockID)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
	// abide by the session 'rules'
	PDescribe("Rules", func() {})
	// gracefully terminate the session
	PDescribe("Termination", func() {})
})
