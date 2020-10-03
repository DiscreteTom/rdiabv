package rdiabv

// return [v for v in itertools.product([0,1,2], repeat=x)]
func cartesianProduct(x int) (result [][]int) {
	result = [][]int{
		{0},
		{1},
		{2},
	}
	for i := 0; i < x-1; i++ {
		var tmp [][]int
		for _, tt := range result {
			tmp = append(tmp, append(tt, 0))
			tmp = append(tmp, append(tt, 1))
			tmp = append(tmp, append(tt, 2))
		}
		result = tmp
	}
	return
}
