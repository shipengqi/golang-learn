package iteration

func Repeat(character string, times int) string {
	str := ""
	for i := 0; i < times; i ++ {
		str += character
	}
	return str
}