package util

func NoErrOrDie(err error) {
	if err != nil {
		panic(err)
	}
}
