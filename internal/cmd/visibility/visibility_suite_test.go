package visibility

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestVisibilityCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
