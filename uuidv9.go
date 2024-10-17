package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
	"math/big"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func calcChecksum(hexString string) string {
	data := make([]byte, len(hexString)/2)
	_, _ = hex.Decode(data, []byte(hexString))

	const polynomial byte = 0x07
	var crc byte = 0x00

	for _, byteVal := range data {
		crc ^= byteVal
		for i := 0; i < 8; i++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc <<= 1
			}
		}
	}
	return fmt.Sprintf("%02x", crc)
}

func verifyChecksum(uuid string) bool {
	base16String := strings.ReplaceAll(uuid, "-", "")[:30]
	crc := calcChecksum(base16String)
	return crc == uuid[34:36]
}

func checkVersion(uuid string, version *int) bool {
	versionDigit := uuid[14:15]
	variantDigit := uuid[19:20]

	if version == nil || string(versionDigit) == fmt.Sprint(*version) {
		if string(versionDigit) == "9" || (strings.Contains("14", string(versionDigit)) && strings.Contains("89abAB", variantDigit)) {
			return true
		}
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
	if options.Checksum != false && options.Checksum && !verifyChecksum(uuid) {
		return false
	}
	if options.Version != false && options.Version && !checkVersion(uuid, nil) {
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
	return fmt.Sprintf("%s-%s-%s-%s-%s", str[:8], str[8:12], str[12:16], str[16:20], str[20:])
}

type UUIDv9Options struct {
	Prefix   string
	Timestamp interface{}
	Checksum bool
	Version  bool
	Legacy   bool
}

func uuidv9(options UUIDv9Options) (string, error) {
	var prefix string
	if options.Prefix != "" {
		if err := validatePrefix(options.Prefix); err != nil {
			return "", err
		}
		prefix = strings.ToLower(options.Prefix)
	}

	var center string
	switch t := options.Timestamp.(type) {
	case time.Time:
		center = fmt.Sprintf("%x", t.UnixNano())
	case int, string:
		var ts time.Time
		switch v := t.(type) {
		case int:
			ts = time.Unix(int64(v), 0)
		case string:
			parsedTime, err := time.Parse(time.RFC3339, v)
			if err == nil {
				ts = parsedTime
			}
		}
		center = fmt.Sprintf("%x", ts.UnixNano())
	default:
		center = fmt.Sprintf("%x", time.Now().UnixNano())
	}

	var checksum = false
	if options.Checksum != false {
		checksum = options.Checksum
	}

	var version = false
	if options.Version != false {
		version = options.Version
	}

	var legacy = false
	if options.Legacy != false {
		legacy = options.Legacy
	}

	var length = 32 - len(prefix) - len(center)
	if !!checksum {
		length -= 2
	}
	if !!legacy {
		length -= 2
	} else if !!version {
		length -= 1
	}
	suffix, err := randomBytes(length)
	if err != nil {
		return "", err
	}

	joined := prefix + center + suffix
	if legacy {
		joined = joined[:12] + "1" + joined[12:15] + randomChar("89ab") + joined[15:]
	} else if version {
		joined = joined[:12] + "9" + joined[12:]
	}

	if checksum {
		joined += calcChecksum(joined)
	}
	return addDashes(joined), nil
}