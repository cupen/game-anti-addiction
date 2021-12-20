package idcard

import "fmt"

func CheckBefore(idnumbers string) error {
	if len(idnumbers) != 15 && len(idnumbers) != 18 {
		return fmt.Errorf("invalid idnumber: length(%d): maybe you should `Clean` it before", len(idnumbers))
	}
	return nil
}

func Clean(idnumbers string) (cleaned string, isChanged bool, err error) {
	if len(idnumbers) <= 0 {
		return
	}
	if len(idnumbers) >= 32 {
		err = fmt.Errorf("invalid idnumbers: too long %d", len(idnumbers))
		return
	}

	chars := []rune(idnumbers)
	if len(chars) <= 0 {
		err = fmt.Errorf("empty []rune(idnumbers) ???")
		return
	}

	lastChar := chars[len(chars)-1]
	switch lastChar {
	case 'x', 'ｘ', 'Ｘ':
		chars[len(chars)-1] = 'X'
		isChanged = true
		cleaned = string(chars)
	default:
		cleaned = idnumbers
	}
	return
}
