package server

// Converts little endian bytes to an int64.
func bytesToInt64(b []byte) int64 {
	var val int64

	val |= int64(b[0])
	val |= int64(b[1]) << 8
	val |= int64(b[2]) << 16
	val |= int64(b[3]) << 24

	return val
}
