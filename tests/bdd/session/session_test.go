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
		validApplication := session.SessionApplication(fixtures.GenerateApplication())
		validNonNativeChain := fixtures.GenerateNonNativeBlockchain()
		validBlockID := session.SessionBlockID(fixtures.GenerateBlockHash())
		Describe("SessionKey Generation", func() {
			Context("Empty Application seed", func() {
				emptyApplication := session.SessionApplication(types.Application{}) // empty Application
				It("should return empty Application public key error", func() {
					_, err := session.NewSessionKey(emptyApplication, validNonNativeChain, validBlockID)
					gomega.Expect(err).To(gomega.Equal(session.EmptyAppPubKeyError))
				})
			})
			Context("Empty NonNativeChain", func() {
				emptyNonNativeChain := session.SessionBlockchain{}
				It("should return empty nonNativeChain error", func() {
					_, err := session.NewSessionKey(validApplication, emptyNonNativeChain, validBlockID)
					gomega.Expect(err).To(gomega.Equal(session.EmptyNonNativeChainError))
				})
			})
			Context("Empty Block ID", func() {
				emptyBlockID := session.SessionBlockID{}
				It("should return empty blockID error", func() {
					_, err := session.NewSessionKey(validApplication, validNonNativeChain, emptyBlockID)
					gomega.Expect(err).To(gomega.Equal(session.EmptyBlockIDError))
				})
			})
			Context("All seed data is valid", func() {
				sessionkey, err := session.NewSessionKey(validApplication, validNonNativeChain, validBlockID)
				It("should return nil error", func() {
					gomega.Expect(err).To(gomega.BeNil())
				})
				It("should return non-nil error", func(){
					gomega.Expect(sessionkey).ToNot(gomega.BeNil())
				})
			})
		})
		Describe("Service Node Generation", func() {
			validSessionKey, _ := session.NewSessionKey(validApplication, validNonNativeChain, validBlockID)
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
				sessNodes, err := session.NewSessionNodes(validNonNativeChain, validSessionKey)
				It("should return nil error", func() {
					gomega.Expect(err).To(gomega.BeNil())
				})
				It("should return non-nil sessionNodes", func(){
					gomega.Expect(sessNodes).ToNot(gomega.BeNil())
				})
			})
		})
		Context("Valid seed data for the session", func() {
			sess, err := session.NewSession(validApplication, validNonNativeChain, validBlockID)
			It("should return nil error", func() {
				gomega.Expect(err).To(gomega.BeNil())
			})
			It("should return non-nil session", func(){
				gomega.Expect(sess).ToNot(gomega.BeNil())
			})
		})
	})
	// abide by the session 'rules'
	PDescribe("Rules", func() {})
	// gracefully terminate the session
	PDescribe("Termination", func() {})
})
