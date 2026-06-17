package util

func ArrToMap[T any, K comparable, V any](items []T, itemToKV func(item T) (K, V)) map[K]V {
	m := map[K]V{}

	for _, item := range items {
		k, v := itemToKV(item)
		m[k] = v
	}

	return m
}
