package session

func either(l, r string) string {
	return firstNonEmpty(l, r)
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
