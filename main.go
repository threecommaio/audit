// Copyright © 2018 ThreeComma.io <hello@threecomma.io>

package main

import "github.com/threecommaio/audit/cmd"
import "github.com/threecommaio/audit/pkg"

func main() {
	cmd.Execute(audit.Version)
}
