package cmdio

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must1[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}
