package util

func AddrToString(txt *string) string {
	if txt != nil {
		return *txt
	}

	return ""
}