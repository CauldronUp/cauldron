package detect

import (
	"encoding/json"
)

type npmDetector struct{}

func (npmDetector) file() string { return "package.json" }

type npmManifest struct {
	Engines         map[string]string `json:"engines"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	PackageManager  string            `json:"packageManager"`
}

// npmServices maps an npm package to a backing service.
var npmServices = map[string]string{
	"ioredis":            "redis",
	"redis":              "redis",
	"pg":                 "postgres",
	"mysql2":             "mysql",
	"meilisearch":        "meilisearch",
	"bullmq":             "redis",
	"@aws-sdk/client-s3": "minio",
}

func (d npmDetector) detect(contents []byte, p *Project) error {
	var manifest npmManifest

	if err := json.Unmarshal(contents, &manifest); err != nil {
		return err
	}

	if version, ok := manifest.Engines["node"]; ok {
		p.add(Requirement{
			Kind:     KindRuntime,
			Name:     "node",
			Version:  normaliseVersion(version),
			Source:   d.file(),
			Evidence: "engines.node",
		})
	}

	all := map[string]string{}
	for pkg, constraint := range manifest.Dependencies {
		all[pkg] = constraint
	}
	for pkg, constraint := range manifest.DevDependencies {
		all[pkg] = constraint
	}

	// A project with a package.json needs Node even if engines is absent.
	if len(all) > 0 {
		p.add(Requirement{
			Kind:     KindRuntime,
			Name:     "node",
			Source:   d.file(),
			Evidence: d.file(),
		})
	}

	for pkg, constraint := range all {
		if service, ok := npmServices[pkg]; ok {
			p.add(Requirement{
				Kind:     KindService,
				Name:     service,
				Source:   d.file(),
				Evidence: pkg,
			})
		}

		if pkg == "vite" {
			p.add(Requirement{
				Kind:     KindService,
				Name:     "vite",
				Version:  normaliseVersion(constraint),
				Source:   d.file(),
				Evidence: pkg,
			})
		}

		if recipe, ok := recipeForNPM(pkg); ok {
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
