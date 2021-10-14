package pcdmask

import (
	"strconv"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateIntMap(t *testing.T) {
	actual := createIntMap("12 4d6")
	expect := map[int]int{
		0: 1,
		1: 2,
		3: 4,
		5: 6,
	}
	assert.Equal(t, expect, actual)
	assert.Equal(t, 0, len(createIntMap("abcd")))
}

func TestFindIntInMap(t *testing.T) {
	assert.True(t, nil == findPanInIntMap(convertIntToMap(""), 0), "False positive result")
	assert.True(t, nil == findPanInIntMap(convertIntToMap("0"), 0), "False positive result")
	assert.Equal(
		t,
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		findPanInIntMap(convertIntToMap("4000160000000004"), 0),
		"Error in finding PAN (only PAN)",
	)
	assert.Equal(
		t,
		[]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		findPanInIntMap(convertIntToMap("00400016000000000400"), 0),
		"Error in finding PAN (PAN and surround zeroes)",
	)
}

func TestMaskWithPositions(t *testing.T) {
	assert.Equal(t, "****", maskWithPositions("1234", []int{0, 1, 2, 3}))
	assert.Equal(t, "12**", maskWithPositions("1234", []int{2, 3}))
	assert.Equal(t, "**34", maskWithPositions("1234", []int{0, 1}))
}

func TestMaskSinglePan(t *testing.T) {
	assert.Equal(t, "4000********0004", MaskPan("4000160000000004"))
	assert.Equal(t, "4000 **** **** 0004", MaskPan("4000 1600 0000 0004"))
	assert.Equal(t, `{"pan":"4000 **** **** 0004"}`, MaskPan(`{"pan":"4000 1600 0000 0004"}`))
	assert.Equal(t, `{"pan":"4000 1600 0000 0005"}`, MaskPan(`{"pan":"4000 1600 0000 0005"}`))
}

func TestMaskMultiplePans(t *testing.T) {
	assert.Equal(t, `4000********0004 ♥ 4000********0004`, MaskPan(`4000160000000004 ♥ 4000160000000004`))
	assert.Equal(t, "55 4000********0004 & 4000********0004", MaskPan("55 4000160000000004 & 4000160000000004"))
	// TODO interesting false positive here - actual "40001600******** & *000160000000004"
	// assert.Equal(t, "4000160000000005 & 4000********0004", MaskPan("4000160000000005 & 4000160000000004"))
	// TODO and here - actual "4000**********05 & 4111111111111111"
	// assert.Equal(t, "4000160000000005 & 4111********1111", MaskPan("4000160000000005 & 4111111111111111"))
	assert.Equal(t, "4000160000000005 & 2223********0010", MaskPan("4000160000000005 & 2223000048410010"))
	assert.Equal(t, "4000********0004 & 4000160000000005", MaskPan("4000160000000004 & 4000160000000005"))
}

func TestMaskMultipleDividedPans(t *testing.T) {
	actual := `4000160000000004 ♥ 4000 ♥ 1600 0000 ♥ 0004`
	expect := `4000********0004 ♥ 4000 ♥ **** **** ♥ 0004`
	assert.Equal(t, expect, MaskPan(actual))
}

func TestMaskWithMultilineString(t *testing.T) {
	actual := `400016
0000000004`
	expect := `4000**
******0004`
	assert.Equal(t, expect, MaskPan(actual))
}

func TestLongPan(t *testing.T) {
	assert.Equal(t, "6771***********0006", MaskPan("6771830000000000006"))
}

func BenchmarkMask500(b *testing.B) {
	symbols := "NYvMT1DUoalXFfOkBy0ZpqZ5rLKCOV2w1MjlaxaK9x9pVtkncLAarWCBm4mye8i3iIphrI1xnTjzINAPQ7fBuQMvAXGnZaJEYLnB7" +
		"nJPATcbyvK37qLgNBfuiSnRSsc23s0hvSfB2doczZegkOPml4vOLQVFb7K310SGBIrhaV7YogaCH1N8g1punNWylBetnR0gYrALkpkTB38zX" +
		"ksNRxSTjSucgVjD5e7QNGAHKCuCikvIuJeejKQe5DFFMFDKlP62zcR2pW5U2OghmBGAFkILXTySgCIr3v9jSLfuFeKcMfbDM7EbbJRBPcLNu" +
		"RYngi9wFzdICfhnXBS6uYCdPzP69JEzBZAHLqLMOXyydlOKvNVVlA6pNAi34rcnJzR2hpFrkDIleuI0BrzlyVo7dj3SSuLj9W3brzQFfJTKQ" +
		"kWbEsP9xl3cD7NdfdN4drJrqQu40001600000000045ENjdmT42VmDRrHzYota7YT98WwB6KVNZ"
	for i := 0; i < b.N; i++ {
		MaskPan(symbols)
	}
}

func convertIntToMap(input string) map[int]int {
	var r = make(map[int]int)
	for pos, char := range input {
		num, _ := strconv.Atoi(string(char))
		r[pos] = num
	}
	return r
}
