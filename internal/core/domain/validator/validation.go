package domain

func IsValidDocument(document string) bool {
	length := len(document)

	if length != 11 && length != 14 {
		return false
	}

	for i := 0; i < lenght; i++ {
		if document[i] < '0' || document[i] > '9' {
			return false
		}
	}

	return true
}