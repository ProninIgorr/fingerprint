package registrator

type KeySize struct {
	Key  interface{}
	Size int64
}

func GetKeySizes(ec map[interface{}]int) (unique, total int64) {
	for ks, count := range ec {
		s := ks.(KeySize).Size
		unique += s
		total += s * int64(count)
	}
	return
}
