package main

import (
	"fmt"
	"github.com/theplant/luhn"
	"regexp"
	"sort"
	"strconv"
)

var digitReg = regexp.MustCompile("[0-9]")

// returns map of integers positions (position:integer) in string
func createIntMap(input string) map[int]int {
	result := make(map[int]int)
	for pos, char := range input {
		if digitReg.MatchString(string(char)) {
			num, _ := strconv.Atoi(string(char))
			result[pos] = num
		}
	}

	return result
}

// returns positions of PAN number in position mapped integers
func findPanInIntMap(m map[int]int, from int) []int {
	var result []int

	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	if from > len(keys) {
		return result
	}

	sort.Ints(keys)
	keys = keys[from:]
	keysLen := len(keys)

	// min PAN number length
	if keysLen < 16 {
		return result
	}

OUTER:
	for i = 0; i < keysLen-15; i++ {
		var valueToCheck string
		for ii := 0; ii < (keysLen - i); ii++ {
			valueToCheck += fmt.Sprintf("%d", m[keys[i+ii]])
			if valueToCheck[0] == '0' {
				continue OUTER
			}
			num, err := strconv.Atoi(valueToCheck)
			if len(valueToCheck) < 16 {
				continue
			}
			if err != nil {
				continue OUTER
			}
			if luhn.Valid(num) {
				return keys[i : i+ii+1]
			}
		}
	}

	return result
}

// returns masked string with * in input positions
func maskWithPositions(input string, pos []int) string {
	sort.Ints(pos)
	for _, p := range pos {
		input = input[:p] + "*" + input[p+1:]
	}

	return input
}

// MaskPan returns string with masked PANs (Primary Account Number) in it
// Pan detection based on Luhn algorithm (check wiki)
func MaskPan(input string) string {
	intMap := createIntMap(input)
	result := input
	found := true
	panPositions := findPanInIntMap(intMap, 0)
	for found {
		if panPositions != nil {
			result = maskWithPositions(result, panPositions[4:len(panPositions)-4])
			panPositions = findPanInIntMap(intMap, panPositions[len(panPositions)-1])
		} else {
			found = false
		}
	}
	return result
}
