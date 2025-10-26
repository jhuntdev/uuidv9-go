package uuidv9

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// "UUIDv9" // Adjust to your actual module path
)

var (
	// uuidRegex     = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	uuidV1Regex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-1[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	uuidV4Regex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	uuidV9Regex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-9[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

func Test_UUIDv9(t *testing.T) {
	t.Run("should validate as a UUID", func(t *testing.T) {
		id1, _ := UUIDv9(UUIDv9Options{})
		id2, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4"})
		id3, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Timestamp: false})
		id4, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Checksum: true})
		id5, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Checksum: true, Version: true})
		id6, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Checksum: true, Legacy: true})

		assert.True(t, uuidRegex.MatchString(id1))
		assert.True(t, uuidRegex.MatchString(id2))
		assert.True(t, uuidRegex.MatchString(id3))
		assert.True(t, uuidRegex.MatchString(id4))
		assert.True(t, uuidRegex.MatchString(id5))
		assert.True(t, uuidRegex.MatchString(id6))
	})

	t.Run("should generate sequential UUIDs", func(t *testing.T) {
		id1, _ := UUIDv9(UUIDv9Options{})
		time.Sleep(2 * time.Millisecond)
		id2, _ := UUIDv9(UUIDv9Options{})
		time.Sleep(2 * time.Millisecond)
		id3, _ := UUIDv9(UUIDv9Options{})

		assert.True(t, id1 < id2)
		assert.True(t, id2 < id3)
	})

	t.Run("should generate sequential UUIDs with a prefix", func(t *testing.T) {
		id1, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4"})
		time.Sleep(2 * time.Millisecond)
		id2, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4"})
		time.Sleep(2 * time.Millisecond)
		id3, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4"})

		assert.True(t, id1 < id2)
		assert.True(t, id2 < id3)
		assert.Equal(t, "a1b2c3d4", id1[:8])
		assert.Equal(t, "a1b2c3d4", id2[:8])
		assert.Equal(t, "a1b2c3d4", id3[:8])
		assert.Equal(t, id1[14:18], id2[14:18])
		assert.Equal(t, id2[14:18], id3[14:18])
	})

	t.Run("should generate non-sequential UUIDs", func(t *testing.T) {
		idS, _ := UUIDv9(UUIDv9Options{Timestamp: false})
		time.Sleep(2 * time.Millisecond)
		idNs, _ := UUIDv9(UUIDv9Options{Timestamp: false})

		assert.NotEqual(t, idS[:4], idNs[:4])
	})

	t.Run("should generate non-sequential UUIDs with a prefix", func(t *testing.T) {
		idS, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Timestamp: false})
		time.Sleep(2 * time.Millisecond)
		idNs, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Timestamp: false})

		assert.Equal(t, "a1b2c3d4", idS[:8])
		assert.Equal(t, "a1b2c3d4", idNs[:8])
		assert.NotEqual(t, idS[14:18], idNs[14:18])
	})

	t.Run("should generate UUIDs with a checksum", func(t *testing.T) {
		id1, _ := UUIDv9(UUIDv9Options{Checksum: true})
		id2, _ := UUIDv9(UUIDv9Options{Timestamp: false, Checksum: true})

		assert.True(t, uuidRegex.MatchString(id1))
		assert.True(t, uuidRegex.MatchString(id2))
		assert.True(t, verifyChecksum(id1))
		assert.True(t, verifyChecksum(id2))
	})

	t.Run("should generate UUIDs with a version", func(t *testing.T) {
		id1, _ := UUIDv9(UUIDv9Options{Version: true})
		id2, _ := UUIDv9(UUIDv9Options{Timestamp: false, Version: true})

		assert.True(t, uuidRegex.MatchString(id1))
		assert.True(t, uuidRegex.MatchString(id2))
		assert.True(t, uuidV9Regex.MatchString(id1))
		assert.True(t, uuidV9Regex.MatchString(id2))
	})

	t.Run("should generate backward compatible UUIDs", func(t *testing.T) {
		id1, _ := UUIDv9(UUIDv9Options{Checksum: true, Legacy: true})
		id2, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Legacy: true})
		id3, _ := UUIDv9(UUIDv9Options{Timestamp: false, Legacy: true})
		id4, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Timestamp: false, Legacy: true})

		assert.True(t, uuidRegex.MatchString(id1))
		assert.True(t, uuidRegex.MatchString(id2))
		assert.True(t, uuidRegex.MatchString(id3))
		assert.True(t, uuidRegex.MatchString(id4))
		assert.True(t, uuidV1Regex.MatchString(id1))
		assert.True(t, uuidV1Regex.MatchString(id2))
		assert.True(t, uuidV4Regex.MatchString(id3))
		assert.True(t, uuidV4Regex.MatchString(id4))
	})

	t.Run("should correctly validate and verify checksum", func(t *testing.T) {
		id1, _ := UUIDv9(UUIDv9Options{Checksum: true})
		id2, _ := UUIDv9(UUIDv9Options{Timestamp: false, Checksum: true})
		id3, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Checksum: true})
		id4, _ := UUIDv9(UUIDv9Options{Prefix: "a1b2c3d4", Timestamp: false, Checksum: true})
		id5, _ := UUIDv9(UUIDv9Options{Checksum: true, Version: true})
		id6, _ := UUIDv9(UUIDv9Options{Checksum: true, Legacy: true})
		id7, _ := UUIDv9(UUIDv9Options{Timestamp: false, Checksum: true, Legacy: true})

		assert.True(t, isUUID(id1))
		assert.False(t, isUUID("not-a-real-uuid"))
		assert.True(t, isValidUUIDv9(id1, isValidUUIDv9Options{Checksum: true}))
		assert.True(t, isValidUUIDv9(id2, isValidUUIDv9Options{Checksum: true}))
		assert.True(t, isValidUUIDv9(id3, isValidUUIDv9Options{Checksum: true}))
		assert.True(t, isValidUUIDv9(id4, isValidUUIDv9Options{Checksum: true}))
		assert.True(t, isValidUUIDv9(id5, isValidUUIDv9Options{Checksum: true, Version: true}))
		assert.True(t, isValidUUIDv9(id6, isValidUUIDv9Options{Checksum: true, Version: true}))
		assert.True(t, isValidUUIDv9(id7, isValidUUIDv9Options{Checksum: true, Version: true}))
		assert.True(t, verifyChecksum(id1))
		assert.True(t, verifyChecksum(id2))
		assert.True(t, verifyChecksum(id3))
		assert.True(t, verifyChecksum(id4))
		assert.True(t, verifyChecksum(id5))
		assert.True(t, verifyChecksum(id6))
		assert.True(t, verifyChecksum(id7))
	})
}
