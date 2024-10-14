<!--
SPDX-FileCopyrightText: 2024 CmdScale GmbH
SPDX-License-Identifier: CC0-1.0
-->

# Corporate Finance Go module `cmdscale.com/corfi`

[![CmdScale Project](https://github.com/cmdscale/.github/raw/main/profile/assets/CmdShield.svg)](https://cmdscale.com/)
[![Go Reference](https://pkg.go.dev/badge/cmdscale.com/corfi.svg)](https://pkg.go.dev/cmdscale.com/corfi)

## ISIN Validation

Ensures that the check-digit of a given [International Securities Identification Numbers (ISIN)](https://en.wikipedia.org/wiki/International_Securities_Identification_Number) is correct:

```go
package main

import "fmt"
import "cmdscale.com/corfi"

func main() {
	isin, err := corfi.NewISIN("DE000BAY0017") // non alphanumeric characters will be ignored
	if err != nil {
		fmt.Println(err) // isin is invalid
	}
	fmt.Println("ISIN is well-formed", isin.String())
}
```

## Installation

```sh
go get cmdscale.com/corfi
```

## License

BSD 3-Clause "New" or "Revised" License
