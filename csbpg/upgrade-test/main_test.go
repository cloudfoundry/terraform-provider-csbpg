package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	oldVersion "github.com/cloudfoundry/terraform-provider-csbpg/120/csbpg"
	newVersion "github.com/cloudfoundry/terraform-provider-csbpg/dev/csbpg"
)

var _ = Describe("upgrade tests", Label("upgrade-test"), func() {
	Context("when running the old version of the code", func() {
		It("has 9 elements in the Provider Schema", func() {
			Expect(len(oldVersion.Provider().Schema)).To(Equal(9))
		})
	})

	Context("when running the local version of the code", func() {
		It("has 1 element in the Provider Schema", func() {
			Expect(len(newVersion.Provider().Schema)).To(Equal(1))
		})
	})
})
