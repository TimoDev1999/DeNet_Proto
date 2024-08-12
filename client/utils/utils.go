package utils

import "fmt"

var Addresses = GenerateTestAddresses(10000)

func GenerateTestAddresses(n int) []string {
	addresses := make([]string, n)
	for i := 0; i < n; i++ {
		addresses[i] = fmt.Sprintf("0x%040d", i+1)
	}
	return addresses
}

func ChunkAddresses(addresses []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(addresses); i += chunkSize {
		end := i + chunkSize
		if end > len(addresses) {
			end = len(addresses)
		}
		chunks = append(chunks, addresses[i:end])
	}
	return chunks
}
