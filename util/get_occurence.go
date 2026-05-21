package util

type Occurence[T any] struct {
	Count int
	Data  T
}

func GetOccurences[T any](arr []T, isEqual func(a, b T) bool) []*Occurence[T] {
	occurences := make([]*Occurence[T], 0, len(arr))

	i := 0
	j := i

	count := 0
	for j < len(arr) {
		if isEqual(arr[i], arr[j]) {
			count++
			j++
			continue
		}

		o := &Occurence[T]{
			Count: count,
			Data:  arr[i],
		}
		count = 0 // reset count
		occurences = append(occurences, o)
		i = j
	}

	// Clean up
	o := &Occurence[T]{
		Count: count,
		Data:  arr[i],
	}
	count = 0 // reset count
	occurences = append(occurences, o)

	return occurences
}
