package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTimespec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Test Suite")
}
