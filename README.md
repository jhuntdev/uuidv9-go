# UUID v9

## Fast, lightweight, zero-dependency Go implementation of UUID version 9

The v9 UUID supports both sequential (time-based) and non-sequential (random) UUIDs with an optional prefix of up to four bytes, an optional checksum, and sufficient randomness to avoid collisions. It uses the UNIX timestamp for sequential UUIDs and CRC-8 for checksums. A version digit can be added if desired, but is omitted by default.

To learn more about UUID v9, please visit the website: https://uuidv9.jhunt.dev

## Installation

Install UUID v9 from Go Packages.

```bash
go get github.com/jhuntdev/uuid-v9
```

## Usage

```go
package main

import {
	"github.com/jhuntdev/uuid-v9"
}

func main() {
	orderedID := UUIDv9()
	prefixedOrderedID := UUIDv9(Option{Prefix: "a1b2c3d4"})
	unorderedID := UUIDv9(Option{Timestamp: false})
	prefixedUnorderedID := UUIDv9(Option{Prefix: "a1b2c3d4", Timestamp: false})
	orderedIDWithChecksum := UUIDv9(Option{Checksum: true})
	orderedIDWithVersion := UUIDv9(Option{Version: true})
	orderedIDWithLegacyMode := UUIDv9(Option{Legacy: true})

	isValid := isValidUUIDv9(orderedID)
	isValidWithChecksum := isValidUUIDv9(orderedIDWithChecksum, Option{Checksum: true})
	isValidWithVersion := isValidUUIDv9(orderedIDWithVersion, Option{Version: true})

	// orderedId = uuidv9()
	// prefixedOrderedId = uuidv9({ prefix: 'a1b2c3d4' })
	// unorderedId = uuidv9({ timestamp: false })
	// prefixedUnorderedId = uuidv9({ prefix: 'a1b2c3d4', timestamp: false })
	// orderedIdWithChecksum = uuidv9({ checksum: true })
	// orderedIdWithVersion = uuidv9({ version: true })
	// orderedIdWithLegacyMode = uuidv9({ legacy: true })

	// isValid = isValidUUIDv9(orderedId)
	// isValidWithChecksum = isValidUUIDv9(orderedIdWithChecksum, { checksum: true })
	// isValidWithVersion = isValidUUIDv9(orderedIdWithVersion, { version: true })
}
```

## Backward Compatibility

Some UUID validators check for specific features of v1 or v4 UUIDs. This causes some valid v9 UUIDs to appear invalid. Three possible workarounds are:

1) Use the built-in validator (recommended)
2) Use legacy mode*
3) Bypass the validator (not recommended)

_*Legacy mode adds version and variant digits to immitate v1 or v4 UUIDs depending on the presence of a timestamp._

## License

This project is licensed under the [MIT License](LICENSE).