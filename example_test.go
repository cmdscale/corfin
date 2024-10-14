// SPDX-FileCopyrightText: 2024 CmdScale GmbH
//
// SPDX-License-Identifier: CC0-1.0

package corfi_test

import (
	"fmt"

	"cmdscale.com/corfi"
)

func Example() {
	isin, err := corfi.NewISIN("DE000BAY0017") // non alphanumeric characters will be removed
	if err != nil {
		panic(err)
	}
	fmt.Println("ISIN:", isin.String())
	// Output:
	// ISIN: DE000BAY0017
}
