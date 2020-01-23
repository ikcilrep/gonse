package nse

import (
	"encoding/binary"
	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/internal/errors"
)

// Int64ToBytes converts integer to byte array. It ignores padding, result is as short as possible.
func Int64ToBytes(integer int64) []byte {
	bytes := make([]byte, 8)
	binary.PutVarint(bytes, integer)
	lastNonZeroIndex := 7
	for ; lastNonZeroIndex > 0 && bytes[lastNonZeroIndex] == 0; lastNonZeroIndex-- {
	}
	return bytes[:lastNonZeroIndex+1]
}

// Int64sToBytes converts []int64 into []byte.
// For each int64 in the slice there is one byte indicating how many bytes to read next and those bytes.
func Int64sToBytes(data []int64) []byte {
	dataLength := len(data)
	resultLength := dataLength * 9
	result := make([]byte, resultLength)
	resultIndex := 0
	for dataIndex := 0; dataIndex < dataLength; dataIndex++ {
		integerBytes := Int64ToBytes(data[dataIndex])
		result[resultIndex] = byte(len(integerBytes))
		resultIndex++
		copy(result[resultIndex:], integerBytes)
		resultIndex += len(integerBytes)
	}
	return result[:resultIndex]
}

// BytesToInt64s converts result of Int64sToBytes back into []int64.
// It returns errors.WrongDataFormatError as an error when data doesn't appear to be a result of Int64sToBytes.
func BytesToInt64s(data []byte) ([]int64, error) {
	dataLength := len(data)

	resultLength := dataLength
	result := make([]int64, resultLength)
	resultIndex := 0
	for dataIndex := 0; dataIndex < dataLength; resultIndex++ {
		newDataIndex := dataIndex + int(data[dataIndex]) + 1
		if newDataIndex > dataLength {
			return nil, errors.WrongDataFormatError
		}

		result[resultIndex], _ = binary.Varint(data[dataIndex+1 : newDataIndex])
		dataIndex = newDataIndex
	}
	return result[:resultIndex], nil
}

// Int8sToBytes converts []int8 into []byte.
// Every int8 in the slice is treated like it would be unsigned.
func Int8sToBytes(data []int8) []byte {
	dataLength := len(data)
	result := make([]byte, dataLength)
	for index, value := range data {
		result[index] = bits.AsUnsigned(value)
	}
	return result
}

// BytesToInt8s converts []byte into []int8.
// Every byte in the slice is treated like it would be signed.
func BytesToInt8s(data []byte) []int8 {
	dataLength := len(data)
	result := make([]int8, dataLength)
	for index, value := range data {
		result[index] = bits.AsSigned(value)
	}
	return result
}
