package lib

import "fmt"

func GetMatchingOpeningParenIndex(str string, closingParenIndex int) (int, error) {
	var s = NewStack()
	for i := closingParenIndex; i >= 0; i-- {
		var char = str[i]
		if char == ')' {
			s.Push(char)
		} else if char == '(' {
			s.Pop()
			if s.Size() == 0 {
				return i, nil
			}
		}
	}
	return 0, fmt.Errorf("Could not find opening parenthesis for index %d", closingParenIndex)
}

func GetMatchingClosingParenIndex(str string, openParenIndex int) (int, error) {
	var s = NewStack()

	for i := openParenIndex; i < len(str); i++ {
		var char = str[i]
		if char == '(' {
			s.Push(char)
		} else if char == ')' {
			s.Pop()
			if s.Size() == 0 {
				return i, nil
			}
		}
	}
	return 0, fmt.Errorf("Could not find closing parenthesis for index %d", openParenIndex)
}

func ParenthesesAreBalanced(str string) bool {
	var s = NewStack()

	for i := range str {
		var char = str[i]
		if char == '(' {
			s.Push(char)
		} else if char == ')' {
			s.Pop()
		}
	}
	return s.Size() == 0
}

func StrIsComposedOfBytes(str string, allowedBytes map[byte]bool) bool {
	for i := range str {
		_, ok := allowedBytes[str[i]]
		if !ok {
			return false
		}
	}
	return true
}

func RemoveDuplicates(ints []int) []int {
	keys := make(map[int]bool)
	uniqueInts := []int{}
	for i := range ints {
		if _, value := keys[ints[i]]; !value {
			keys[ints[i]] = true
			uniqueInts = append(uniqueInts, ints[i])
		}
	}
	return uniqueInts
}
