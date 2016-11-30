package cowbull_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCowbull(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cowbull Suite")
}
