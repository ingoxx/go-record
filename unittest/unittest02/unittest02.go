package unittest02

func add(n int) int {
	res := 0
	for i := 1; i <= n; i++ {
		res += i
	}
	return res
}

func sub(n1, n2 int) int {
	return n1 - n2
}

func mod(m1, m2 int) int {
	return m1 % m2
}
