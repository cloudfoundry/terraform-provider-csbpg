package main

import (
	"fmt"

	oldVersion "github.com/cloudfoundry/terraform-provider-csbpg/120/csbpg"
	newVersion "github.com/cloudfoundry/terraform-provider-csbpg/dev/csbpg"
)

func main() {
	fmt.Printf("%v\n", oldVersion.Provider().Schema)
	fmt.Printf("%v\n", newVersion.Provider().Schema)
}
