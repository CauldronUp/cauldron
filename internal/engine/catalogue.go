package engine

import (
	"fmt"
	"net"
	"sort"
)

// Catalogue maps a detected service onto a container.
//
// Images are pinned to a major version rather than :latest. A local
// environment that silently changes database major version between two
// developers is precisely the class of "works on my machine" this tool exists
// to remove.
var catalogue = map[string]Spec{
	"postgres": {
		Service:    "postgres",
		Image:      "postgres:18-alpine",
		Ports:      []Port{{Host: 5432, Container: 5432}},
		Env:        map[string]string{"POSTGRES_USER": "cauldron", "POSTGRES_PASSWORD": "cauldron", "POSTGRES_DB": "cauldron"},
		Health:     []string{"pg_isready", "-U", "cauldron"},
		Volume:     "postgres",
		VolumePath: "/var/lib/postgresql/data",
	},
	"mysql": {
		Service:    "mysql",
		Image:      "mysql:8.4",
		Ports:      []Port{{Host: 3306, Container: 3306}},
		Env:        map[string]string{"MYSQL_ROOT_PASSWORD": "cauldron", "MYSQL_DATABASE": "cauldron", "MYSQL_USER": "cauldron", "MYSQL_PASSWORD": "cauldron"},
		Health:     []string{"mysqladmin", "ping", "-h", "127.0.0.1", "-pcauldron"},
		Volume:     "mysql",
		VolumePath: "/var/lib/mysql",
	},
	"redis": {
		Service:    "redis",
		Image:      "redis:8-alpine",
		Ports:      []Port{{Host: 6379, Container: 6379}},
		Health:     []string{"redis-cli", "ping"},
		Volume:     "redis",
		VolumePath: "/data",
	},
	"mailpit": {
		Service: "mailpit",
		Image:   "axllent/mailpit:latest",
		Ports:   []Port{{Host: 1025, Container: 1025}, {Host: 8025, Container: 8025}},
	},
	"meilisearch": {
		Service:    "meilisearch",
		Image:      "getmeili/meilisearch:v1.11",
		Ports:      []Port{{Host: 7700, Container: 7700}},
		Env:        map[string]string{"MEILI_MASTER_KEY": "cauldron", "MEILI_NO_ANALYTICS": "true"},
		Volume:     "meilisearch",
		VolumePath: "/meili_data",
	},
	"minio": {
		Service:    "minio",
		Image:      "minio/minio:latest",
		Ports:      []Port{{Host: 9000, Container: 9000}, {Host: 9001, Container: 9001}},
		Env:        map[string]string{"MINIO_ROOT_USER": "cauldron", "MINIO_ROOT_PASSWORD": "cauldron123"},
		Command:    []string{"server", "/data", "--console-address", ":9001"},
		Volume:     "minio",
		VolumePath: "/data",
	},
}

// Known reports whether a detected service can be booted as a container.
func Known(service string) bool {
	_, ok := catalogue[service]

	return ok
}

// SpecFor returns the container spec for a detected service.
func SpecFor(service string) (Spec, bool) {
	spec, ok := catalogue[service]
	if !ok {
		return Spec{}, false
	}

	// Copy the maps and slices so a caller adjusting one project's ports cannot
	// mutate the catalogue for the next.
	clone := spec
	clone.Ports = append([]Port(nil), spec.Ports...)
	clone.Health = append([]string(nil), spec.Health...)
	clone.Command = append([]string(nil), spec.Command...)

	if spec.Env != nil {
		clone.Env = make(map[string]string, len(spec.Env))
		for key, value := range spec.Env {
			clone.Env[key] = value
		}
	}

	return clone, true
}

// Catalogued returns every service Cauldron knows how to boot.
func Catalogued() []string {
	out := make([]string, 0, len(catalogue))

	for name := range catalogue {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}

// ErrPortInUse reports a host port that something else already holds.
type ErrPortInUse struct {
	Service string
	Port    int
}

func (e *ErrPortInUse) Error() string {
	return fmt.Sprintf("port %d is already in use, which %s needs", e.Port, e.Service)
}

// PortFree reports whether a host port can be bound.
//
// Both the wildcard address and loopback are tested, and this is not
// belt-and-braces. A container from another project publishing on 0.0.0.0
// still leaves 127.0.0.1 bindable on some platforms, so a loopback-only check
// reports "free" for a port Docker will then refuse. That failure surfaces as
// a wall of daemon output at the worst possible moment.
func PortFree(port int) bool {
	for _, address := range []string{fmt.Sprintf("0.0.0.0:%d", port), fmt.Sprintf("127.0.0.1:%d", port)} {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return false
		}

		_ = listener.Close()
	}

	return true
}

// CheckPorts returns the first conflict among a set of specs.
func CheckPorts(specs []Spec) error {
	for _, spec := range specs {
		for _, port := range spec.Ports {
			if !PortFree(port.Host) {
				return &ErrPortInUse{Service: spec.Service, Port: port.Host}
			}
		}
	}

	return nil
}
