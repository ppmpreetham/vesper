package sites

var DefaultSites = []Site{
	{
		Name:           "Twitter",
		UrlTemplate:    "https://x.com/%s",
		UsernameRegex:  `^[A-Za-z0-9_]{1,15}$`,
		ErrorMsg:       "Sorry, that page does not exist!",
		Method:         "GET",
		DisplayURL:     "https://twitter.com/%s",
		ErrorRegex:     `Sorry, that page does not exist!`,
		TimeoutSeconds: 10,
	},
	{
		Name:          "GitHub",
		UrlTemplate:   "https://github.com/%s",
		UsernameRegex: `^[A-Za-z0-9][A-Za-z0-9-]{0,38}$`,
		ErrorMsg:      "This account doesn't exist",
		Method:        "GET",
		ErrorRegex:    `This account doesn't exist`,
	},
}
