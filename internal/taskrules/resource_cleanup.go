package taskrules

type CloseRecorder struct {
	Closed []string
}

func (r *CloseRecorder) Close(name string) func() {
	return func() {
		for i, c := range r.Closed {
			if c == name {
				r.Closed = append(r.Closed[:i], r.Closed[i+1:]...)
				return
			}
		}
	}
}

func CleanupResource(recorder *CloseRecorder, name string) {
	recorder.Closed = append(recorder.Closed, name)
	closeFn := recorder.Close(name)
	defer closeFn()
}
