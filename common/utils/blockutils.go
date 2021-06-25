package utils

func CalculateBlocks(totalSize, blockSize int) (blocks int) {
	if totalSize%blockSize == 0 {
		return totalSize / blockSize
	}
	return totalSize/blockSize + 1
}
