package fetch

// Option is a functional option for configuring a Fetch instance.
type Option func(*Fetch)

// WithLimit returns an Option that sets the maximum number of results to return.
func WithLimit(n int) Option {
	return func(f *Fetch) {
		f.limit = n
	}
}
