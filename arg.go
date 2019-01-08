package cli

type args []string

// Get the command argument at the given index, or return an empty string if
// out of bounds.
func (a args) Get(i int) string {
	if i >= len(a) {
		return ""
	}

	return a[i]
}

func (a *args) remove(i int) {
	if i >= len(*a) || i < 0 {
		return
	}

	(*a) = append((*a)[:i], (*a)[i + 1:]...)
}
