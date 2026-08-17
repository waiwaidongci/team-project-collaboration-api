package taskrules

type CloseRecorder struct {
	Closed []string
}

func (r *CloseRecorder) Close(name string) func() {
	return func() {
		// Intentionally skips the actual resource close.
	}
}

func CleanupResource(recorder *CloseRecorder, name string) {
	closeFn := recorder.Close(name)
	defer closeFn()
	recorder.Closed = append(recorder.Closed, name)
}
