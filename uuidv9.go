package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

var uuidRegex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func calcChecksum(hexString string) string {
	// This function must match the Python implementation exactly
	data := make([]byte, len(hexString)/2)
	_, err := hex.Decode(data, []byte(hexString))
	if err != nil {
		fmt.Printf("Error decoding hex in calcChecksum: %v\n", err)
		return "00" // Return a default in case of error
	}

	const polynomial byte = 0x07
	var crc byte = 0x00

	for i, byteVal := range data {
		crc ^= byteVal
		for j := 0; j < 8; j++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc <<= 1
			}
		}
		fmt.Printf("After byte %d (0x%02x): crc=0x%02x\n", i, byteVal, crc)
	}

	result := fmt.Sprintf("%02x", crc&0xFF)
	fmt.Printf("Final checksum for '%s': %s\n", hexString, result)
	return result
}

func verifyChecksum(uuid string) bool {
	// This function needs to exactly match Python's verify_checksum behavior
	// Python: def verify_checksum(uuid):
	//    clean_uuid = uuid.replace('-', '')[0:30]
	//    checksum = calc_checksum(clean_uuid)
	//    return checksum == uuid[34:36]

	// Only work with properly formatted UUIDs
	if !uuidRegex.MatchString(uuid) {
		fmt.Printf("UUID doesn't match regex: %s\n", uuid)
		return false
	}

	// Remove dashes and extract the first 30 chars for checksum calculation
	cleanUuid := strings.ReplaceAll(uuid, "-", "")
	if len(cleanUuid) < 32 {
		fmt.Printf("Clean UUID too short: %s\n", cleanUuid)
		return false
	}

	// Calculate checksum on first 30 characters
	base16String := cleanUuid[:30]
	calculated := calcChecksum(base16String)
	actual := uuid[34:36]

	fmt.Printf("Verifying UUID: %s\n", uuid)
	fmt.Printf("Clean UUID (first 30): %s\n", base16String)
	fmt.Printf("Calculated checksum: %s, Actual checksum at position 34-36: %s\n", calculated, actual)

	return calculated == actual
}

func checkVersion(uuid string, version *int) bool {
	versionDigit := uuid[14:15]
	variantDigit := uuid[19:20]

	if version == nil {
		return versionDigit == "9" || ((versionDigit == "1" || versionDigit == "4") && strings.Contains("89abAB", variantDigit))
	}

	versionStr := fmt.Sprint(*version)
	if versionDigit == versionStr {
		if versionStr == "1" || versionStr == "4" {
			return strings.Contains("89abAB", variantDigit)
		}
		return true
	}

	return false
}

func isUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

type validateUUIDv9Options struct {
	Checksum bool
	Version  bool
}

func isValidUUIDv9(uuid string, options validateUUIDv9Options) bool {
	if !isUUID(uuid) {
		return false
	}
	if options.Checksum && !verifyChecksum(uuid) {
		return false
	}
	if options.Version && !checkVersion(uuid, nil) {
		return false
	}
	return true
}

func randomBytes(count int) (string, error) {
	bytes := make([]byte, count)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}

// func randomChar(chars string) string {
// 	return string(chars[rand.Int(0, big.NewInt(int64(len(chars))))])
// }

func randomChar(chars string) string {
	n := len(chars)
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
	return string(chars[index.Int64()])
}

var base16Regex = regexp.MustCompile(`^[0-9a-fA-F]+$`)

func isBase16(str string) bool {
	return base16Regex.MatchString(str)
}

func validatePrefix(prefix string) error {
	if len(prefix) > 8 {
		return fmt.Errorf("prefix must be no more than 8 characters")
	}
	if !isBase16(prefix) {
		return fmt.Errorf("prefix must be only hexadecimal characters")
	}
	return nil
}

func addDashes(str string) string {
	// Ensure the string is exactly 32 characters
	if len(str) > 32 {
		str = str[:32]
	} else if len(str) < 32 {
		str = str + strings.Repeat("0", 32-len(str))
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", str[:8], str[8:12], str[12:16], str[16:20], str[20:])
}

type UUIDv9Options struct {
	Prefix    string
	Timestamp interface{}
	Checksum  bool
	Version   bool
	Legacy    bool
}

// UUIDv9 generates a UUID version 9 string
// This is compatible with the Python implementation
// Options:
//   - Prefix: Optional prefix for the UUID (up to 8 hexadecimal characters)
//   - Timestamp: If true or nil, includes current timestamp; can be custom int or time.Time
//   - Checksum: If true, includes a checksum in the last 2 characters
//   - Version: If true, sets the version character to '9'
//   - Legacy: If true, makes the UUID compatible with v1 or v4 format
func UUIDv9(optionalOptions ...UUIDv9Options) (string, error) {
	// Get config from options
	var options UUIDv9Options
	if len(optionalOptions) > 0 {
		options = optionalOptions[0]
	} else {
		options = UUIDv9Options{} // Default options
	}
	prefix := options.Prefix
	timestamp := options.Timestamp
	checksum := options.Checksum
	version := options.Version
	legacy := options.Legacy

	// Apply default to timestamp if not specified
	if timestamp == nil {
		timestamp = true
	}

	// Validate the prefix is base16 (lowercase hex only)
	if prefix != "" {
		if err := validatePrefix(prefix); err != nil {
			return "", err
		}
		prefix = strings.ToLower(prefix)
	}

	// Generate timestamp component if requested
	center := ""
	if timestamp == true {
		// Convert nanoseconds to milliseconds to match Python behavior
		timeMs := time.Now().UnixNano() / 1000000
		center = fmt.Sprintf("%x", timeMs)
	}

	// Calculate how many random bytes are needed based on options
	// Base UUID is 32 hex chars (16 bytes)
	length := 32 - len(prefix) - len(center)

	// Adjust for optional components
	if checksum {
		length -= 2 // reserve 2 chars (1 byte) for checksum
	}

	if legacy {
		length -= 2 // reserve 2 chars (1 byte) for legacy UUID v1/v4 marking
	} else if version {
		length -= 1 // reserve 1 char (half-byte) for version marking
	}

	// Ensure we're generating a full UUID (32 hex chars)
	if length < 0 {
		length = 0 // Safety check for very long prefixes
	}

	// Each byte produces 2 hex chars, so divide by 2
	suffix, err := randomBytes(length / 2)
	if err != nil {
		return "", err
	}

	// Join all components
	joined := prefix + center + suffix

	// Add version and variant if requested
	if legacy {
		// Match Python implementation: Place a '1' or '4' at position 12, and variant at position 16
		pos12 := 12
		pos16 := 16

		if len(joined) < pos12+1 {
			// Add padding if needed
			joined = fmt.Sprintf("%0*s", pos12+1, joined)
		}

		// Take parts before and after position 12
		part1 := joined[:pos12]
		part2 := ""

		if len(joined) > pos12 {
			part2 = joined[pos12+1:]
		}

		// Insert '1' or '4' at position 12
		timeChar := "1"
		if timestamp == false {
			timeChar = "4"
		}

		joined = part1 + timeChar + part2

		// Add variant at position 16
		if len(joined) < pos16+1 {
			// Add padding if needed
			joined = fmt.Sprintf("%0*s", pos16+1, joined)
		}

		// Take parts before and after position 16
		part1 = joined[:pos16]
		part2 = ""

		if len(joined) > pos16 {
			part2 = joined[pos16+1:]
		}

		// Add a random variant character ('8', '9', 'a', or 'b')
		variant := randomChar("89ab")

		joined = part1 + variant + part2
	} else if version {
		// Add a '9' at position 12 to indicate UUIDv9
		pos12 := 12

		if len(joined) < pos12+1 {
			// Add padding if needed
			joined = fmt.Sprintf("%0*s", pos12+1, joined)
		}

		// Take parts before and after position 12
		part1 := joined[:pos12]
		part2 := ""

		if len(joined) > pos12 {
			part2 = joined[pos12+1:]
		}

		// Insert '9' at position 12
		joined = part1 + "9" + part2
	}

	// Create a dashed UUID without the checksum first
	uuidWithoutChecksum := addDashes(joined)

	// Add checksum if requested - Must be added AFTER version is set
	if checksum {
		// Calculate checksum on the first 30 chars of the UUID without dashes
		cleanUuid := strings.ReplaceAll(uuidWithoutChecksum, "-", "")
		base16String := cleanUuid[:30]
		checksum := calcChecksum(base16String)

		// Replace the last two characters of the UUID with the checksum
		// so that it appears at positions 34-36 in the final dashed format
		return uuidWithoutChecksum[:34] + checksum, nil
	}

	return uuidWithoutChecksum, nil
}
