package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTerraformProviderCSBPG(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TerraformProviderCSBPG Suite")
}
