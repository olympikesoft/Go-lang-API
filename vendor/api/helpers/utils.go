package helpers

func containsInt(s []uint32, e uint32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Min(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}

func RoundPrice(num uint32) uint32 {
	realPrice := float64(num)/float64(100)
	realPrice += 0.49
	return uint32(realPrice) * 100
}