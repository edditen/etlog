package utils

var (
	whiteBytes = map[byte]bool{' ': true, '\t': true, '\n': true}
)

func isWhiteByte(b byte) bool {
	if _, ok := whiteBytes[b]; ok {
		return true
	}
	return false
}

func TrimBytes(bs []byte) []byte {
	i := 0
	j := len(bs) - 1

	for ; i < len(bs) && isWhiteByte(bs[i]); i++ {
	}

	for ; j >= 0 && isWhiteByte(bs[j]); j-- {
	}
	return bs[i : j+1]
}
