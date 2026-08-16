package detect

import (
	"bufio"
	"bytes"
	"strings"
)

type goModDetector struct{}

func (goModDetector) file() string { return "go.mod" }

// goServices maps a Go module path prefix to a backing service.
var goServices = map[string]string{
	"github.com/redis/go-redis":      "redis",
	"github.com/lib/pq":              "postgres",
	"github.com/jackc/pgx":           "postgres",
	"github.com/go-sql-driver/mysql": "mysql",
}

func (d goModDetector) detect(contents []byte, p *Project) error {
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	inRequireBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "go ") {
			p.add(Requirement{
				Kind:     KindRuntime,
				Name:     "go",
				Version:  normaliseVersion(strings.TrimPrefix(line, "go ")),
				Source:   d.file(),
				Evidence: "go directive",
			})

			continue
		}

		switch {
		case strings.HasPrefix(line, "require ("):
			inRequireBlock = true
			continue
		case inRequireBlock && line == ")":
			inRequireBlock = false
			continue
		case strings.HasPrefix(line, "require "):
			d.requirement(strings.TrimPrefix(line, "require "), p)
			continue
		}

		if inRequireBlock {
			d.requirement(line, p)
		}
	}

	return scanner.Err()
}

// requirement handles a single line from a require block:
//
//	github.com/stripe/stripe-go/v76 v76.25.0 // indirect
func (d goModDetector) requirement(line string, p *Project) {
	// Indirect dependencies are transitive; booting a fake for something the
	// project never calls directly would be misleading.
	if strings.Contains(line, "// indirect") {
		return
	}

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return
	}

	path := fields[0]

	version := ""
	if len(fields) > 1 {
		version = normaliseVersion(fields[1])
	}

	for prefix, service := range goServices {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			p.add(Requirement{
				Kind:     KindService,
				Name:     service,
				Source:   d.file(),
				Evidence: path,
			})
		}
	}

	if recipe, ok := recipeForGoModule(path); ok {
		p.add(Requirement{
			Kind:     KindRecipe,
			Name:     recipe,
			Version:  version,
			Source:   d.file(),
			Evidence: path,
		})

		return
	}

	if looksLikeAPIClient(path) {
		p.addUnmatched(Requirement{
			Kind:     KindRecipe,
			Name:     path,
			Source:   d.file(),
			Evidence: path,
		})
	}
}
