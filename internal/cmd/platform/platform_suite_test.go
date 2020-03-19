package platform

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPlatformCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
