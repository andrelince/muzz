package slice

func Map[I, O any](in []I, mapFunc func(in I) O) []O {
	out := make([]O, len(in))
	for i := range in {
		out[i] = mapFunc(in[i])
	}
	return out
}
