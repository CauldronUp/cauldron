package detect

import "strings"

// provider maps a dependency to the Recipe that emulates it.
//
// The mapping is deliberately explicit rather than fuzzy. A wrong guess here is
// worse than no guess: booting the wrong fake sends a developer chasing a bug
// that doesn't exist.
type provider struct {
	// recipe is the canonical Recipe name.
	recipe string
	// packages are exact dependency names, per ecosystem.
	composer []string
	npm      []string
	gomod    []string
}

func providers() []provider {
	return []provider{
		{
			recipe:   "stripe",
			composer: []string{"stripe/stripe-php"},
			npm:      []string{"stripe"},
			gomod:    []string{"github.com/stripe/stripe-go"},
		},
		{
			recipe:   "shopify",
			composer: []string{"shopify/shopify-api"},
			npm:      []string{"@shopify/shopify-api", "@shopify/admin-api-client"},
			gomod:    []string{"github.com/bold-commerce/go-shopify"},
		},
		{
			recipe:   "github",
			composer: []string{"knplabs/github-api"},
			npm:      []string{"@octokit/rest", "@octokit/core", "octokit"},
			gomod:    []string{"github.com/google/go-github"},
		},
		{
			recipe:   "slack",
			composer: []string{"slack-php/slack-php-api"},
			npm:      []string{"@slack/web-api", "@slack/bolt"},
			gomod:    []string{"github.com/slack-go/slack"},
		},
		{
			recipe:   "twilio",
			composer: []string{"twilio/sdk"},
			npm:      []string{"twilio"},
			gomod:    []string{"github.com/twilio/twilio-go"},
		},
		{
			recipe:   "adyen",
			composer: []string{"adyen/php-api-library"},
			npm:      []string{"@adyen/api-library"},
		},
		{
			recipe:   "clerk",
			composer: []string{"clerk/clerk-sdk-php"},
			npm:      []string{"@clerk/clerk-sdk-node", "@clerk/backend"},
		},
		{
			recipe:   "hubspot",
			composer: []string{"hubspot/api-client"},
			npm:      []string{"@hubspot/api-client"},
		},
		{
			recipe:   "xero",
			composer: []string{"xeroapi/xero-php-oauth2"},
			npm:      []string{"xero-node"},
		},
		{
			recipe:   "mailgun",
			composer: []string{"mailgun/mailgun-php"},
			npm:      []string{"mailgun.js"},
		},
		{
			recipe:   "sendgrid",
			composer: []string{"sendgrid/sendgrid"},
			npm:      []string{"@sendgrid/mail"},
		},
		{
			recipe:   "openai",
			composer: []string{"openai-php/client", "openai-php/laravel"},
			npm:      []string{"openai"},
			gomod:    []string{"github.com/sashabaranov/go-openai"},
		},
	}
}

// suspectedAPIClient matches dependency names that look like third-party API
// clients we have no Recipe for. These are reported rather than ignored so the
// developer knows exactly which calls will still hit the real network.
var suspectedAPIClientHints = []string{
	"-sdk", "sdk-", "-api", "api-client", "-client",
}

// recipeForComposer resolves a Composer package to a Recipe name.
func recipeForComposer(pkg string) (string, bool) {
	for _, p := range providers() {
		for _, candidate := range p.composer {
			if strings.EqualFold(candidate, pkg) {
				return p.recipe, true
			}
		}
	}

	return "", false
}

// recipeForNPM resolves an npm package to a Recipe name.
func recipeForNPM(pkg string) (string, bool) {
	for _, p := range providers() {
		for _, candidate := range p.npm {
			if strings.EqualFold(candidate, pkg) {
				return p.recipe, true
			}
		}
	}

	return "", false
}

// recipeForGoModule resolves a Go module path to a Recipe name. Module paths
// carry major-version suffixes (/v76), so this matches on prefix.
func recipeForGoModule(path string) (string, bool) {
	for _, p := range providers() {
		for _, candidate := range p.gomod {
			if path == candidate || strings.HasPrefix(path, candidate+"/") {
				return p.recipe, true
			}
		}
	}

	return "", false
}

// looksLikeAPIClient is a heuristic used only to populate Unmatched. It never
// causes a fake to boot.
func looksLikeAPIClient(pkg string) bool {
	lower := strings.ToLower(pkg)

	for _, hint := range suspectedAPIClientHints {
		if strings.Contains(lower, hint) {
			return true
		}
	}

	return false
}
