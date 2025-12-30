package utils

type anyError map[string]any

func (e anyError) Error() string {
	return ""
}
