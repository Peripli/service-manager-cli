package broker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBrokerCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
