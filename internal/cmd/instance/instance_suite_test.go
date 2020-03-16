package instance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestInstanceCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
