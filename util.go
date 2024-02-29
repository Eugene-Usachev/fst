package fst

func getSizeForLen(len int) int {
	if len < 255 {
		return 1
	} else if len < 65535 {
		return 2
	} else {
		return 3
	}
}

func getLenAndSize(buf []byte) (int, int) {
	if buf[0] < 255 {
		return int(buf[0]), 1
	}

	length := int(buf[2])<<8 | int(buf[1])
	if length < 65535 {
		return length, 3
	}

	return int(buf[5])<<16 | int(buf[4])<<8 | int(buf[3]), 6
}

func getBytesFromLen(len int) []byte {
	if len < 255 {
		return []byte{byte(len)}
	} else if len < 65535 {
		return []byte{byte(255), byte(len), byte(len >> 8)}
	}

	return []byte{byte(255), byte(255), byte(255), byte(len), byte(len >> 8), byte(len >> 16)}
}

func getBytesForInt64(i int64) []byte {
	return []byte{
		byte(i),
		byte(i >> 8),
		byte(i >> 16),
		byte(i >> 24),
		byte(i >> 32),
		byte(i >> 40),
		byte(i >> 48),
		byte(i >> 56),
	}
}

func getInt64(buf []byte) int64 {
	return int64(buf[7]) << 56 | int64(buf[6]) << 48 | int64(buf[5]) << 40 | int64(buf[4]) << 32 |
		int64(buf[3]) << 24 | int64(buf[2]) << 16 | int64(buf[1]) << 8 | int64(buf[0])
}
