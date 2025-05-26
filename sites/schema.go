package sites

type Site struct {
	Name               string // e.g., Twitter
	UrlTemplate        string // e.g., https://twitter.com/%s
	UsernameRegex      string // e.g., ^[A-Za-z0-9_]15
	ErrorMsg           string // What error message means not found
	ExistStatusCode    int    // HTTP status code for a found profile (e.g., 200)
	NotFoundStatusCode int    // HTTP status code for not found (e.g., 404)
	Method             string // GET or POST (some sites may require POST)
	TimeoutSeconds     int    // Custom timeout for this site
	CheckRedirect      bool   // Whether to check for HTTP redirects as a clue
	Disabled           bool   // Skip this site if true
	DisplayURL         string // Optional: the public-facing URL to show in results
	ErrorRegex         string // Regex for matching not found text in the response
	RateLimitSeconds   int    // Optional: how many seconds to wait between requests
	Notes              string // Any relevant notes or caveats about the site
}
