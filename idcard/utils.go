package idcard

import (
	"strings"
)

func Clean(idnumbers string) (string, error) {
	if len(idnumbers) <= 0 {
		return "", ErrInvalidLength
	}
	if len(idnumbers) >= 32 {
		// return "", fmt.Errorf("%w length:%d", ErrInvalidLength, len(idnumbers))
		return "", ErrInvalidLength
	}
	idnumbers = strings.ReplaceAll(idnumbers, " ", "")
	chars := []rune(idnumbers)
	if len(chars) != 15 && len(chars) != 18 {
		return "", ErrInvalidLength
	}

	var cleaned = ""
	lastChar := chars[len(chars)-1]
	switch lastChar {
	case 'x', 'ｘ', 'Ｘ':
		chars[len(chars)-1] = 'X'
		cleaned = string(chars)
	default:
		cleaned = idnumbers
	}
	return cleaned, nil
}
