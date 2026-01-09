package random

import (
	"math/rand"
)

func NewRandomString(Len int) string {
	rangArr := [][]int{
		{'0', '9'},
		{'A', 'Z'},
		{'a', 'z'},
	}

	var fullArr []byte
	for _, row := range rangArr {
		for i := row[0]; i < row[1]; i++ {
			fullArr = append(fullArr, byte(i))
		}
	}

	bytesArray := make([]byte, Len)
	for i := range bytesArray {
		randASCIICode := byte(fullArr[rand.Intn(len(fullArr))])
		bytesArray[i] = randASCIICode
	}

	return string(bytesArray)
}
