package main

func sort(values []int, l, r int) {
	if l >= r {
		return
	}

	pivot := values[l]
	i := l + 1

	for j := l + 1; j <= r; j++ {
		if pivot > values[j] {
			values[i], values[j] = values[j], values[i]
			i++
		}
	}

	values[l], values[i-1] = values[i-1], pivot

	sort(values, l, i-2)
	sort(values, i, r)
}
