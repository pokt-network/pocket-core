package session

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/tests/fixtures"
	"github.com/pokt-network/pocket-core/x/pocketcore/session"
)

var _ = Describe("Session", func() {
	// generate the session
	Describe("Generation", func() {
		validApplication := session.SessionAppPubKey(fixtures.GenerateApplication().PubKey)
		validNonNativeChain := fixtures.GenerateNonNativeBlockchain()
		validBlockID := session.SessionBlockID(fixtures.GenerateBlockHash())
		nodesPointer, err := fixtures.GetNodes()
		if err != nil {
			Fail("Could not retrieve mock nodesPointer for testing")
		}
		allActiveNodes := *nodesPointer
		Describe("SessionKey Generation", func() {
			Context("Empty Application seed", func() {
				emptyApplication := session.SessionAppPubKey("") // empty Application
				It("should return empty Application public key error", func() {
					_, err := session.NewSessionKey(emptyApplication, validNonNativeChain, validBlockID)
					Expect(err).To(Equal(session.EmptyAppPubKeyError))
				})
			})
			Context("Empty NonNativeChain", func() {
				emptyNonNativeChain := session.SessionBlockchain{}
				It("should return empty nonNativeChain error", func() {
					_, err := session.NewSessionKey(validApplication, emptyNonNativeChain, validBlockID)
					Expect(err).To(Equal(session.EmptyNonNativeChainError))
				})
			})
			Context("Empty Block ID", func() {
				emptyBlockID := session.SessionBlockID{}
				It("should return empty blockID error", func() {
					_, err := session.NewSessionKey(validApplication, validNonNativeChain, emptyBlockID)
					Expect(err).To(Equal(session.EmptyBlockIDError))
				})
			})
			Context("All seed data is valid", func() {
				sessionkey, err := session.NewSessionKey(validApplication, validNonNativeChain, validBlockID)
				It("should return nil error", func() {
					Expect(err).To(BeNil())
				})
				It("should return non-nil error", func() {
					Expect(sessionkey).ToNot(BeNil())
				})
			})
		})
		Describe("Service Node Generation", func() {
			validSessionKey, _ := session.NewSessionKey(validApplication, validNonNativeChain, validBlockID)
			Context("Empty NonNativeChain", func() {
				emptyNonNativeChain := session.SessionBlockchain{}
				It("should return empty nonNativeChain error", func() {
					_, err := session.NewSessionNodes(emptyNonNativeChain, validSessionKey, allActiveNodes)
					Expect(err).To(Equal(session.EmptyNonNativeChainError))
				})
			})
			Context("Empty SessionKey", func() {
				invalidSessionKey := session.SessionKey{}
				It("should return invalid SessionKey error", func() {
					_, err := session.NewSessionNodes(validNonNativeChain, invalidSessionKey, allActiveNodes)
					Expect(err).To(Equal(session.EmptySessionKeyError))
				})
			})
			Context("Valid Seed data", func() {
				sessNodes, err := session.NewSessionNodes(validNonNativeChain, validSessionKey, allActiveNodes)
				It("should return nil error", func() {
					Expect(err).To(BeNil())
				})
				It("should return non-nil sessionNodes", func() {
					Expect(sessNodes).ToNot(BeNil())
				})
			})
		})
		Context("Valid seed data for the session", func() {
			sess, err := session.NewSession(validApplication, validNonNativeChain, validBlockID, allActiveNodes)
			It("should return nil error", func() {
				Expect(err).To(BeNil())
			})
			It("should return non-nil session", func() {
				Expect(sess).ToNot(BeNil())
			})
		})
	})
})
