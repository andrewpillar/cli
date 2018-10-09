// Simple library for building CLI applications in Go
package cli

type args []string

func (a args) Get(i int) string {
	if i >= len(a) {
		return ""
	}

	return a[i]
}

func (a *args) set(i int, s string) {
	if i >= len(*a) {
		return
	}

	(*a)[i] = s
}
