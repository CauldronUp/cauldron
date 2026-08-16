// Package engine starts and stops the containers a project needs.
//
// It drives the docker command line rather than the Docker SDK. The SDK pulls
// in a very large dependency tree for what amounts to a handful of calls, and
// the CLI is present on every machine that has Docker at all, already knows
// the user's context and credentials, and behaves identically on macOS, Linux
// and Windows. Keeping Cauldron at one dependency is worth more here than
// avoiding a subprocess.
package engine

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Label keys used to find and clean up everything Cauldron owns. Without them,
// a crashed run would leave orphans that nothing knows how to remove.
const (
	LabelManaged = "dev.cauldron.managed"
	LabelProject = "dev.cauldron.project"
	LabelService = "dev.cauldron.service"
)

// Runner executes a docker command. It exists so the engine can be tested
// without a daemon: the tests assert on the arguments Cauldron builds, which
// is where the bugs actually live.
type Runner interface {
	Run(ctx context.Context, args ...string) (string, error)
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}

	return string(out), nil
}

// Engine manages containers for one project.
type Engine struct {
	runner  Runner
	project string
}

// New returns an engine for a project name.
func New(project string) *Engine {
	return &Engine{runner: execRunner{}, project: project}
}

// WithRunner swaps the command runner, for tests.
func (e *Engine) WithRunner(r Runner) *Engine {
	e.runner = r

	return e
}

// ErrUnavailable means Docker is not usable.
type ErrUnavailable struct{ Reason string }

func (e *ErrUnavailable) Error() string {
	return "Docker is not available: " + e.Reason
}

// Available reports whether a Docker daemon is reachable. It is checked before
// anything else so the failure is one clear sentence rather than a wall of
// container errors.
func (e *Engine) Available(ctx context.Context) error {
	if _, err := exec.LookPath("docker"); err != nil {
		return &ErrUnavailable{Reason: "the docker command is not on PATH"}
	}

	if _, err := e.runner.Run(ctx, "info", "--format", "{{.ServerVersion}}"); err != nil {
		return &ErrUnavailable{Reason: "the daemon is not responding, is Docker running?"}
	}

	return nil
}

// Network returns the project's network name.
func (e *Engine) Network() string { return "cauldron-" + e.project }

// ContainerName returns the container name for a service.
func (e *Engine) ContainerName(service string) string {
	return "cauldron-" + e.project + "-" + service
}

// EnsureNetwork creates the project network if it is missing.
func (e *Engine) EnsureNetwork(ctx context.Context) error {
	out, err := e.runner.Run(ctx, "network", "ls", "--filter", "name=^"+e.Network()+"$", "--format", "{{.Name}}")
	if err != nil {
		return err
	}

	if strings.TrimSpace(out) == e.Network() {
		return nil
	}

	_, err = e.runner.Run(ctx, "network", "create",
		"--label", LabelManaged+"=true",
		"--label", LabelProject+"="+e.project,
		e.Network())

	return err
}

// publishedPortPattern pulls the host port out of a docker ps ports column,
// e.g. "0.0.0.0:6379->6379/tcp".
var publishedPortPattern = regexp.MustCompile(`:(\d+)->`)

// PublishedPorts returns every host port Docker currently has allocated,
// across all projects.
//
// The operating system check alone is not enough: Docker refuses a port its
// own allocator already holds even when the OS would allow the bind.
func (e *Engine) PublishedPorts(ctx context.Context) (map[int]string, error) {
	out, err := e.runner.Run(ctx, "ps", "--format", "{{.Names}}\t{{.Ports}}")
	if err != nil {
		return nil, err
	}

	held := map[int]string{}

	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		name, ports, found := strings.Cut(strings.TrimSpace(line), "\t")
		if !found {
			continue
		}

		for _, match := range publishedPortPattern.FindAllStringSubmatch(ports, -1) {
			port, err := strconv.Atoi(match[1])
			if err != nil {
				continue
			}

			if _, seen := held[port]; !seen {
				held[port] = name
			}
		}
	}

	return held, nil
}

// ErrPortHeld reports a port another container already publishes. Naming the
// container matters: the fix is to stop that one, and the developer has no
// other way to know which it is.
type ErrPortHeld struct {
	Service   string
	Port      int
	Container string
}

func (e *ErrPortHeld) Error() string {
	return fmt.Sprintf("port %d is already published by container %s, which %s needs", e.Port, e.Container, e.Service)
}

// CheckDockerPorts reports the first spec whose port another container holds.
func (e *Engine) CheckDockerPorts(ctx context.Context, specs []Spec) error {
	held, err := e.PublishedPorts(ctx)
	if err != nil {
		return nil // Not worth failing the run over a diagnostic query.
	}

	for _, spec := range specs {
		for _, port := range spec.Ports {
			if owner, taken := held[port.Host]; taken {
				return &ErrPortHeld{Service: spec.Service, Port: port.Host, Container: owner}
			}
		}
	}

	return nil
}

// Port maps a container port onto the host.
type Port struct {
	Host      int
	Container int
}

// Spec describes a container to run.
type Spec struct {
	Service string
	Image   string
	Ports   []Port
	Env     map[string]string
	// Health is a command run inside the container to decide readiness. Without
	// one, "started" is not the same as "usable" and the app races the database.
	Health []string
	// Command overrides the image's default command.
	Command []string
	// Volume is a named volume mounted for persistence, if any.
	Volume string
	// VolumePath is where the volume mounts inside the container.
	VolumePath string
}

// Start launches a container, reusing one that is already running.
func (e *Engine) Start(ctx context.Context, spec Spec) error {
	name := e.ContainerName(spec.Service)

	state, err := e.state(ctx, name)
	if err != nil {
		return err
	}

	switch state {
	case "running":
		return nil
	case "exited", "created", "paused":
		// Restarting beats recreating: the volume, and therefore the data, is
		// preserved across an interrupted run.
		_, err := e.runner.Run(ctx, "start", name)

		return err
	}

	args := []string{
		"run", "--detach",
		"--name", name,
		"--network", e.Network(),
		"--network-alias", spec.Service,
		"--label", LabelManaged + "=true",
		"--label", LabelProject + "=" + e.project,
		"--label", LabelService + "=" + spec.Service,
		"--restart", "unless-stopped",
	}

	for _, port := range spec.Ports {
		// Bound to loopback deliberately. A fake database listening on every
		// interface is an invitation, and this one has a known password.
		args = append(args, "--publish", fmt.Sprintf("127.0.0.1:%d:%d", port.Host, port.Container))
	}

	for _, key := range sortedKeys(spec.Env) {
		args = append(args, "--env", key+"="+spec.Env[key])
	}

	if len(spec.Health) > 0 {
		args = append(args,
			"--health-cmd", strings.Join(spec.Health, " "),
			"--health-interval", "2s",
			"--health-retries", "15",
			"--health-start-period", "2s",
		)
	}

	if spec.Volume != "" && spec.VolumePath != "" {
		args = append(args, "--volume", e.volumeName(spec.Volume)+":"+spec.VolumePath)
	}

	args = append(args, spec.Image)
	args = append(args, spec.Command...)

	if _, err := e.runner.Run(ctx, args...); err != nil {
		// A failed run can still leave a created container behind. Left there,
		// the next attempt finds it in "created" state and tries to start it,
		// which fails for the same reason and hides the real cause. Failures
		// have to clean up after themselves.
		_, _ = e.runner.Run(ctx, "rm", "--force", name)

		// Translate the daemon's wall of output into the one fact that matters.
		if strings.Contains(err.Error(), "port is already allocated") {
			port := 0
			if len(spec.Ports) > 0 {
				port = spec.Ports[0].Host
			}

			return &ErrPortHeld{Service: spec.Service, Port: port, Container: "another container"}
		}

		return err
	}

	return nil
}

func (e *Engine) volumeName(volume string) string {
	return "cauldron-" + e.project + "-" + volume
}

// state returns a container's state, or "" when it does not exist.
func (e *Engine) state(ctx context.Context, name string) (string, error) {
	out, err := e.runner.Run(ctx, "ps", "--all", "--filter", "name=^"+name+"$", "--format", "{{.State}}")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

// ErrUnhealthy means a container started but never became usable.
type ErrUnhealthy struct {
	Service string
	Last    string
}

func (e *ErrUnhealthy) Error() string {
	return fmt.Sprintf("%s did not become ready (last state: %s)", e.Service, e.Last)
}

// WaitHealthy blocks until a container reports healthy, or the context ends.
//
// Containers without a healthcheck report no status at all, so those are
// treated as ready once running. Claiming readiness we cannot verify would be
// worse than saying nothing.
func (e *Engine) WaitHealthy(ctx context.Context, service string, timeout time.Duration) error {
	name := e.ContainerName(service)
	deadline := time.Now().Add(timeout)
	last := "unknown"

	for time.Now().Before(deadline) {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		out, err := e.runner.Run(ctx, "inspect", "--format", "{{.State.Status}} {{if .State.Health}}{{.State.Health.Status}}{{end}}", name)
		if err != nil {
			return err
		}

		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) == 0 {
			return &ErrUnhealthy{Service: service, Last: "missing"}
		}

		status := fields[0]
		health := ""

		if len(fields) > 1 {
			health = fields[1]
		}

		last = strings.TrimSpace(status + " " + health)

		switch {
		case status != "running":
			// Keep waiting: it may still be starting.
		case health == "" || health == "healthy":
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}

	return &ErrUnhealthy{Service: service, Last: last}
}

// Running lists the services this project currently has running.
func (e *Engine) Running(ctx context.Context) ([]string, error) {
	out, err := e.runner.Run(ctx, "ps",
		"--filter", "label="+LabelProject+"="+e.project,
		"--format", "{{.Label \""+LabelService+"\"}}")
	if err != nil {
		return nil, err
	}

	var services []string

	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			services = append(services, line)
		}
	}

	return services, nil
}

// Stop removes every container for this project, and the network with them.
func (e *Engine) Stop(ctx context.Context, keepData bool) error {
	out, err := e.runner.Run(ctx, "ps", "--all",
		"--filter", "label="+LabelProject+"="+e.project,
		"--format", "{{.Names}}")
	if err != nil {
		return err
	}

	var problems []string

	for _, name := range strings.Split(strings.TrimSpace(out), "\n") {
		if name = strings.TrimSpace(name); name == "" {
			continue
		}

		if _, err := e.runner.Run(ctx, "rm", "--force", "--volumes="+fmt.Sprint(!keepData), name); err != nil {
			problems = append(problems, err.Error())
		}
	}

	if _, err := e.runner.Run(ctx, "network", "rm", e.Network()); err != nil {
		// A network that was never created, or is already gone, is not a
		// failure worth reporting to someone running "cauldron down".
		if !strings.Contains(err.Error(), "not found") {
			problems = append(problems, err.Error())
		}
	}

	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}

	return nil
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))

	for key := range m {
		out = append(out, key)
	}

	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}

	return out
}
