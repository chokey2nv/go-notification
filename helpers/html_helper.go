package helpers

func StripHTML(input string) string {
	var out []rune
	inTag := false

	for _, r := range input {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				out = append(out, r)
			}
		}
	}

	return string(out)
}
