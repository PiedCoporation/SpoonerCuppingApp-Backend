package arrayutils

// A \ B
func Difference[T comparable](a, b []T) []T {
	m := make(map[T]struct{})
	for _, x := range b {
		m[x] = struct{}{}
	}

	var diff []T
	for _, x := range a {
		if _, ok := m[x]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}
