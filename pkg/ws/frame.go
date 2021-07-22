package ws

import (
	"encoding/binary"
	"math"
	"bytes"
)

// Helper functions
func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}


// Frame
type frame struct {
	fin bool
	rsv []bool
	opcode byte
	mask bool

	payloadLen uint64
	maskKey []byte
	payload []byte

	raw []byte
	decoded string
}

// Constructor
func newFrame(data []byte) frame {
	raw := data

	// Header
	fin := data[0] & (1 << 7) > 0

	rsv := make([]bool, 3)
	for i := 0; i < 3; i++ {
		rsv[i] = data[0] & byte(1 << (6 - i)) > 0
	}

	opcode := data[0] & 0b00001111

	data = data[1:]

	// Body
	mask := data[0] & (1 << 7) > 0

	var payloadLen uint64
	payloadByte := data[0] & 0b01111111
	data = data[1:]

	if payloadByte := int(payloadByte); payloadByte <= 125 {
		// Normal payload length
		payloadLen = uint64(payloadByte)
	} else {
		// Extended payload length
		payloadByteSlice := make([]byte, 8)
		if payloadByte < 127 {
			payloadByteSlice = data[:2]
			data = data[2:]
		} else {
			payloadByteSlice = data[:8]
			data = data[8:]
		}

		// Binary to uint
		payloadLen = binary.BigEndian.Uint64(payloadByteSlice)
	}

	maskKey := make([]byte, 4)
	if mask {
		maskKey = data[:4]
		data = data[4:]
	}

	return frame{
		fin: fin,
		rsv: rsv,
		opcode: opcode,
		mask: mask,

		payloadLen: payloadLen,
		maskKey: maskKey,
		payload: data,

		raw: raw,
	}
}

// Methods
func (f *frame) decode() string {
	if f.decoded != "" {
		return f.decoded
	}

	decodedBytes := make([]byte, 0, f.payloadLen)

	var i uint64
	for i = 0; i < f.payloadLen; i++ {
		decodedBytes = append(decodedBytes, f.payload[i] ^ f.maskKey[i%4])
	}

	f.decoded = string(decodedBytes)

	return string(decodedBytes)
}

func (f *frame) encode() []byte {
	if f.raw != nil {
		return f.raw
	}

	final := make([]byte, 0)
	var current byte

	// Header
	current |= byte(boolToInt(f.fin) << 7)
	
	for i := 0; i < 3; i++ {
		current |= byte(boolToInt(f.rsv[i])<<(6-i))
	}

	current |= byte(f.opcode)

	final = append(final, current)
	current = 0

	// Body
	current |= byte(boolToInt(f.mask) << 7)

	if f.payloadLen <= 125 {
		// 1-byte length
		current |= byte(f.payloadLen)
		final = append(final, current)
	} else {
		// Multibyte length
		if (f.payloadLen <= math.MaxUint16) {
			// 2 bytes
			final = append(final, 126)

			mb := make([]byte, 0, 2)
			binary.BigEndian.PutUint16(mb, uint16(f.payloadLen))
			length := bytes.Trim(mb, "\x00")

			for i := 0; i < 2 - len(length); i++ {
				final = append(final, 0)
			}

			final = append(final, length...)
		} else {
			// 8 bytes
			final = append(final, 127)

			mb := make([]byte, 0, 2)
			binary.BigEndian.PutUint64(mb, f.payloadLen)
			length := bytes.Trim(mb, "\x00")

			for i := 0; i < 8 - len(length); i++ {
				final = append(final, 0)
			}

			final = append(final, length...)
		}
	}
	current = 0

	if f.mask {
		final = append(final, f.maskKey...)
	}

	if f.payloadLen > 0 {
		final = append(final, f.payload...)
	}


	f.raw = final
	return final
}
