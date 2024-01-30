package lib

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func GetMultiWhois(ctx context.Context, domains []string) ([]whoisparser.WhoisInfo, error) {
	allWhois := make([]whoisparser.WhoisInfo, 0, len(domains))

	whoisCh := make(chan whoisparser.WhoisInfo, len(domains))
	errorCh := make(chan error, len(domains))

	for _, domain := range domains {
		go getChanWhois(ctx, domain, whoisCh, errorCh)
	}

	for i := 0; i < len(domains); i++ {
		select {
		case whois := <-whoisCh:
			allWhois = append(allWhois, whois)
		case err := <-errorCh:
			log.Printf("Error fetching WHOIS for domain %s: %v", domains[i], err)
			allWhois = append(allWhois, whoisparser.WhoisInfo{}) // Append an empty struct on error
		case <-ctx.Done():
			log.Printf("Context done for domain %s: %v", domains[i], ctx.Err())
			return nil, ctx.Err()
		}
	}

	close(whoisCh)
	close(errorCh)

	return allWhois, nil
}

// getChanWhois sends Whois data to a channel
func getChanWhois(ctx context.Context, domain string, whoisCh chan<- whoisparser.WhoisInfo, errorCh chan<- error) {
	select {
	case <-ctx.Done():
		errorCh <- ctx.Err()
		return
	case <-time.After(2 * time.Second): // Timeout for each WHOIS lookup
		errorCh <- fmt.Errorf("Lookup timed out after 2 seconds")
		return
	default:
		raw, err := whois.Whois(domain)
		if err != nil {
			errorCh <- err
			return
		}

		result, err := whoisparser.Parse(raw)
		if err != nil {
			errorCh <- err
			return
		}

		whoisCh <- result
	}
}
