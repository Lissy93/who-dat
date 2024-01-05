package structs

// SingleBody defines the JSON body for
// getting Whois data of a single domain
type SingleBody struct {
	Domain string
}

// MultiBody defines the JSON body for
// getting Whois data of an array of domains
type MultiBody struct {
	Domains []string
}
