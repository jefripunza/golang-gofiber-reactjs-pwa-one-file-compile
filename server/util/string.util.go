package util

import (
	"fmt"
	"regexp"
	"strings"
)

type String struct{}

func (ref String) ReplaceMessageWithDynamicData(message string, data map[string]interface{}) string {
	re := regexp.MustCompile(`#([^#]+)#`)
	return re.ReplaceAllStringFunc(message, func(m string) string {
		key := strings.Trim(m, "#") // Mendapatkan kunci tanpa karakter '#'
		if val, ok := data[key]; ok {
			return fmt.Sprintf("%v", val) // Mengganti dengan nilai dari data
		}
		return m // Jika kunci tidak ditemukan, kembalikan placeholder asli
	})
}
