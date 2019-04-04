package servicing_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServicing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Servicing Suite")
}
