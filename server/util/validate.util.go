package util

import (
	"fmt"
	"regexp"
	"strconv"
)

type Validate struct{}

func (ref Validate) NumberOnly(value interface{}) (int, error) {
	var number int
	switch v := value.(type) {
	case int:
		number = v
	case int32:
		number = int(v)
	case float64:
		number = int(v)
	case string:
		p, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		number = p
	default:
		return 0, fmt.Errorf("value harus berupa angka")
	}

	return number, nil
}

func (ref Validate) IsPhoneNumber(phoneNumber string) bool {
	// Pola regex untuk validasi nomor telepon internasional
	regex := `\+(9[976]\d|8[987530]\d|6[987]\d|5[90]\d|42\d|3[875]\d|2[98654321]\d|9[8543210]|8[6421]|6[6543210]|5[87654321]|4[987654310]|3[9643210]|2[70]|7|1)\d{1,14}$`

	// Membuat objek regex
	re := regexp.MustCompile(regex)

	// Mengecek apakah nomor telepon sesuai dengan pola regex
	return re.MatchString(phoneNumber)
}
