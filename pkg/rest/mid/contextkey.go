package mid

// Great article about context keys
// https://medium.com/@matryer/context-keys-in-go-5312346a868d

// ContextKey represents the key for context value
type ContextKey string

func (c ContextKey) String() string {
	return "rest request " + string(c)
}
