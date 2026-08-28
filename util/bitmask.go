package util

const byteBitSize = 8

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func SetBit[T Integer](data T, offset int) T {
	return data | (1 << offset)
}

func ResetBit[T Integer](data T, offset int) T {
	return data &^ (1 << offset)
}

func HasBit[T Integer](data T, offset int) bool {
	return data&(1<<offset) != 0
}

func SetBitSlice(data []byte, offset int) []byte {
	if offset < 0 {
		return data
	}
	index, pos := offset/byteBitSize, offset%byteBitSize
	if index >= len(data) {
		newSlice := make([]byte, index+1)
		copy(newSlice, data)
		data = newSlice
	}
	data[index] |= (0x01 << pos)
	return data
}

func ResetBitSlice(data []byte, offset int) []byte {
	if offset < 0 {
		return data
	}
	index, pos := offset/byteBitSize, offset%byteBitSize
	if index >= len(data) {
		return data
	}
	data[index] &^= (0x01 << pos)
	return data
}

func HasBitSlice(data []byte, offset int) bool {
	if offset < 0 {
		return false
	}
	index, pos := offset/byteBitSize, offset%byteBitSize
	if index >= len(data) {
		return false
	}
	return data[index]&(0x01<<pos) != 0
}

// Deprecated: Use SetBit[T] instead.
func SetBit32(data int32, offset int) int32 { return SetBit(data, offset) }

// Deprecated: Use ResetBit[T] instead.
func ResetBit32(data int32, offset int) int32 { return ResetBit(data, offset) }

// Deprecated: Use HasBit[T] instead.
func HasBit32(data int32, offset int) bool { return HasBit(data, offset) }

// Deprecated: Use SetBit[T] instead.
func SetBit64(data int64, offset int) int64 { return SetBit(data, offset) }

// Deprecated: Use ResetBit[T] instead.
func ResetBit64(data int64, offset int) int64 { return ResetBit(data, offset) }

// Deprecated: Use HasBit[T] instead.
func HasBit64(data int64, offset int) bool { return HasBit(data, offset) }
