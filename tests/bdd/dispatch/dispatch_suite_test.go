package dispatch_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDispatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dispatch Suite")
}
