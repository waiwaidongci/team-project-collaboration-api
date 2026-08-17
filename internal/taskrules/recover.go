package taskrules

func RecoverPanic(fn func()) (recovered any) {
	fn()
	return nil
}
