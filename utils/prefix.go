package utils

func PrefixRange(prefix []byte) (start []byte, end []byte) {
	if len(prefix) == 0 {
		return nil, nil
	}

	start = prefix

	end = make([]byte, len(prefix))
	copy(end, prefix)

	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < 0xFF {
			end[i]++
			end = end[:i+1]
			return start, end
		}
	}

	return start, nil
}
