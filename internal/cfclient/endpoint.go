package cfclient

const (
	getMeta   = "https://codeforces.com"
	referer   = "https://codeforces.com/problemset/status"
	getSubmit = "https://codeforces.com/data/submitSource"

	csrfHeader    = "X-Csrf-Token"
	refererHeader = "Referer"

	submissionId = "submissionId"
	csrfToken    = "csrf_token"

	xQuery = "//meta[@name='X-Csrf-Token']"
)
