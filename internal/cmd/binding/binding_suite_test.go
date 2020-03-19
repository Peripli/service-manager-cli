package binding

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBindingCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
