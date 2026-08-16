// Package recipes bundles the first-party Cauldron Recipes into the binary.
//
// The embed directive has to live alongside the files it embeds, which is why
// this package exists rather than the Recipes being read from disk at runtime.
// Shipping them inside the binary keeps the single-file promise: `cauldron
// recipe info stripe` works on a fresh machine with no network and nothing to
// install.
package recipes

import "embed"

//go:embed */recipe.yaml
var FS embed.FS
