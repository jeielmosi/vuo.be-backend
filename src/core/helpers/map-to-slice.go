package helpers

import helpers "github.com/jei-el/vuo.be-backend/src/core/helpers/types"

func MapToSlice[T any](mp map[string]T) []helpers.KeyValue[T] {
	size := len(mp)
	ans := make([]helpers.KeyValue[T], size)

	i := 0
	for key, val := range mp {
		ans[i].Key = key
		ans[i].Value = val
		i++
	}

	return ans
}
