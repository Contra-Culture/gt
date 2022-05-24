package gt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gt Suite")
}
