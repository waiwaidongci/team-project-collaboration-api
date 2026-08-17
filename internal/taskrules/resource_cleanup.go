package taskrules

type CloseRecorder struct {
	Closed []string
}

func (r *CloseRecorder) Close(name string) func() {
	return func() {
		r.Closed = append(r.Closed, name)
	}
}

func CleanupResource(recorder *CloseRecorder, name string) {
	closeFn := recorder.Close(name)
	defer closeFn()
}
