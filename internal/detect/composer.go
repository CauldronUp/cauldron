package detect

import (
	"encoding/json"
	"regexp"
	"strings"
)

type composerDetector struct{}

func (composerDetector) file() string { return "composer.json" }

type composerManifest struct {
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

// composerServices maps a Composer package to a backing service Cauldron boots.
var composerServices = map[string]string{
	"predis/predis":               "redis",
	"laravel/horizon":             "horizon",
	"laravel/scout":               "meilisearch",
	"laravel/reverb":              "reverb",
	"meilisearch/meilisearch-php": "meilisearch",
	"aws/aws-sdk-php":             "minio",
}

func (d composerDetector) detect(contents []byte, p *Project) error {
	var manifest composerManifest

	if err := json.Unmarshal(contents, &manifest); err != nil {
		return err
	}

	if version, ok := manifest.Require["php"]; ok {
		p.add(Requirement{
			Kind:     KindRuntime,
			Name:     "php",
			Version:  normaliseVersion(version),
			Source:   d.file(),
			Evidence: "require.php",
		})
	}

	for pkg, constraint := range manifest.Require {
		if pkg == "php" || strings.HasPrefix(pkg, "ext-") {
			continue
		}

		if pkg == "laravel/framework" {
			p.Framework = "Laravel"
			p.add(Requirement{
				Kind:     KindRuntime,
				Name:     "laravel",
				Version:  normaliseVersion(constraint),
				Source:   d.file(),
				Evidence: pkg,
			})
		}

		if service, ok := composerServices[pkg]; ok {
			p.add(Requirement{
				Kind:     KindService,
				Name:     service,
				Source:   d.file(),
				Evidence: pkg,
			})
		}

		if recipe, ok := recipeForComposer(pkg); ok {
			p.add(Requirement{
				Kind:     KindRecipe,
				Name:     recipe,
				Version:  normaliseVersion(constraint),
				Source:   d.file(),
				Evidence: pkg,
			})

			continue
		}

		if looksLikeAPIClient(pkg) {
			p.addUnmatched(Requirement{
				Kind:     KindRecipe,
				Name:     pkg,
				Source:   d.file(),
				Evidence: pkg,
			})
		}
	}

	return nil
}

// versionPattern pulls the first sensible version number out of a Composer or
// npm constraint: "^8.5", ">=8.2 <9.0", "~24.0.1" all yield a usable value.
var versionPattern = regexp.MustCompile(`(\d+(?:\.\d+)*)`)

// normaliseVersion reduces a constraint to a concrete-looking version. It is
// intentionally lossy: Cauldron wants "which major to boot", not full semver
// resolution, and the lock file is consulted when precision matters.
func normaliseVersion(constraint string) string {
	constraint = strings.TrimSpace(constraint)
	if constraint == "" || constraint == "*" {
		return ""
	}

	match := versionPattern.FindString(constraint)

	return match
}
