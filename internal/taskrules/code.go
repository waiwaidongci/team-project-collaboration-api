package taskrules

func CodeFor(err error) string {
	if err == nil {
		return "ok"
	}
	return "internal"
}
