package offering

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestOfferingCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
