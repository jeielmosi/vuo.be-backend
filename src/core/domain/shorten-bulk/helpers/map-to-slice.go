package helpers

type KeyValue[T any] struct {
	Key   string
	Value T
}

func MapToSlice[T any](mp map[string]T) []KeyValue[T] {
	size := len(mp)
	ans := make([]KeyValue[T], size)

	i := 0
	for key, val := range mp {
		ans[i].Key = key
		ans[i].Value = val
		i++
	}

	return ans
}
