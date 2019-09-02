package aat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Auth Token Suite")
}
