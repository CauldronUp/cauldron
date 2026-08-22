// Errors: how a provider says no.

package recipe

// Error is a named failure mode that `cauldron fault` can inject.
type Error struct {
	Status int    `yaml:"status"`
	Code   string `yaml:"code"`
	// Type is the provider's error category, which is often a much smaller set
	// than the codes. Stripe has four types and dozens of codes, and client
	// libraries switch on the type. Empty falls back to the code, which is
	// wrong often enough that every Recipe should set it.
	Type    string `yaml:"type"`
	Message string `yaml:"message"`
	// Style overrides the Recipe-wide error envelope for this failure alone,
	// because a provider can answer two shapes and the npm registry does.
	//
	// Checked against registry.npmjs.org on 2026-08-22: a package that does
	// not exist answers {"error":"Not found"}, and a version that does not
	// exist on a package that does answers the bare JSON string "version not
	// found: 99.99.99". Same status, same registry, one object and one
	// string. Code reading body.error off the second finds undefined, and
	// code that reports body.error as the reason reports "undefined".
	Style   string            `yaml:"style"`
	Headers map[string]string `yaml:"headers"`
	// Fields are extra body properties this failure carries, merged over the
	// Recipe-wide ones. Dropbox describes each failure with its own nested
	// union, so a single set of constants would make every error claim to be
	// the same one.
	Fields map[string]any `yaml:"fields"`
}
