package main

func assert(b bool, s string) {
	if !b {
		panic(s)
	}
}
