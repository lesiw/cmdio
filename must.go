package cmdio

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustv[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
