package encrypt

func Encrypt(text string, key string) string {
	var res []byte
	k := 0
	for i := 0; i < len(text); i++ {
		if k == len(key) {
			k = 0
		}
		res = append(res, text[i]^key[k])
		k++
	}
	return string(res)
}
