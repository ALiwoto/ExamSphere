package masterServer

var (
	// If you are behind a load-balancer such as cloudflare or
	// nginx, you should enable this option to get the real client IP,
	// for nginx, it's usually "X-Real-IP" and for cloudflare, it's
	// "CF-Connecting-IP".
	// Otherwise, just set it to an empty string.
	ProxyHeader   = ""
	CaseSensitive = false
)
