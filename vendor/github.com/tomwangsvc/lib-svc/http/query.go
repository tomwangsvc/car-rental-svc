package http

func UniqueQueryParameters(queryParameters map[string][]string) map[string]string {
	m := make(map[string]string, len(queryParameters))
	for k, v := range queryParameters {
		m[k] = v[0]
	}
	return m
}
