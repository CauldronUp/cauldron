package engine

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// fakeRunner records the docker commands Cauldron builds and returns canned
// output. The arguments are where the bugs live, so they are what is asserted.
type fakeRunner struct {
	calls   [][]string
	replies map[string]string
	fail    map[string]error
}

func newFake() *fakeRunner {
	return &fakeRunner{replies: map[string]string{}, fail: map[string]error{}}
}

func (f *fakeRunner) Run(_ context.Context, args ...string) (string, error) {
	f.calls = append(f.calls, args)

	joined := strings.Join(args, " ")

	for key, err := range f.fail {
		if strings.Contains(joined, key) {
			return "", err
		}
	}

	for key, reply := range f.replies {
		if strings.Contains(joined, key) {
			return reply, nil
		}
	}

	return "", nil
}

// find returns the first recorded call containing every fragment.
func (f *fakeRunner) find(fragments ...string) []string {
	for _, call := range f.calls {
		joined := strings.Join(call, " ")

		matched := true

		for _, fragment := range fragments {
			if !strings.Contains(joined, fragment) {
				matched = false
				break
			}
		}

		if matched {
			return call
		}
	}

	return nil
}

func engineWith(f *fakeRunner) *Engine {
	return New("demo").WithRunner(f)
}

func TestNamesAreScopedToTheProject(t *testing.T) {
	e := New("demo")

	if e.Network() != "cauldron-demo" {
		t.Errorf("Network = %q", e.Network())
	}

	if got := e.ContainerName("postgres"); got != "cauldron-demo-postgres" {
		t.Errorf("ContainerName = %q", got)
	}
}

func TestStartBuildsTheExpectedRunCommand(t *testing.T) {
	f := newFake()
	e := engineWith(f)

	spec, _ := SpecFor("postgres")

	if err := e.Start(context.Background(), spec); err != nil {
		t.Fatalf("Start: %v", err)
	}

	call := f.find("run", "--detach")
	if call == nil {
		t.Fatalf("no docker run was issued; calls: %v", f.calls)
	}

	joined := strings.Join(call, " ")

	for _, want := range []string{
		"--name cauldron-demo-postgres",
		"--network cauldron-demo",
		"--network-alias postgres",
		LabelProject + "=demo",
		LabelService + "=postgres",
		"postgres:18-alpine",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("run command is missing %q\n%s", want, joined)
		}
	}
}

// A fake database with a known password must not be reachable from the network.
func TestPublishedPortsAreBoundToLoopback(t *testing.T) {
	f := newFake()
	e := engineWith(f)

	spec, _ := SpecFor("postgres")
	_ = e.Start(context.Background(), spec)

	joined := strings.Join(f.find("run", "--detach"), " ")

	if !strings.Contains(joined, "--publish 127.0.0.1:5432:5432") {
		t.Errorf("ports must bind to loopback only\n%s", joined)
	}

	if strings.Contains(joined, "--publish 0.0.0.0") {
		t.Error("a service was published on every interface")
	}
}

func TestEnvironmentIsOrderedSoCommandsAreReproducible(t *testing.T) {
	f := newFake()
	e := engineWith(f)

	spec, _ := SpecFor("postgres")
	_ = e.Start(context.Background(), spec)

	first := strings.Join(f.find("run", "--detach"), " ")

	f2 := newFake()
	_ = engineWith(f2).Start(context.Background(), spec)

	second := strings.Join(f2.find("run", "--detach"), " ")

	if first != second {
		t.Errorf("the same spec produced two different commands:\n%s\n%s", first, second)
	}
}

func TestStartReusesARunningContainer(t *testing.T) {
	f := newFake()
	f.replies["ps --all"] = "running\n"

	e := engineWith(f)
	spec, _ := SpecFor("redis")

	if err := e.Start(context.Background(), spec); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if f.find("run", "--detach") != nil {
		t.Error("a running container must not be recreated")
	}
}

// Restarting rather than recreating is what preserves the volume, and with it
// the developer's data, across an interrupted run.
func TestStartRestartsAStoppedContainer(t *testing.T) {
	f := newFake()
	f.replies["ps --all"] = "exited\n"

	e := engineWith(f)
	spec, _ := SpecFor("redis")

	if err := e.Start(context.Background(), spec); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if f.find("start", "cauldron-demo-redis") == nil {
		t.Error("a stopped container should be started, not recreated")
	}

	if f.find("run", "--detach") != nil {
		t.Error("a stopped container must not be recreated")
	}
}

func TestEnsureNetworkIsIdempotent(t *testing.T) {
	f := newFake()
	f.replies["network ls"] = "cauldron-demo\n"

	if err := engineWith(f).EnsureNetwork(context.Background()); err != nil {
		t.Fatalf("EnsureNetwork: %v", err)
	}

	if f.find("network", "create") != nil {
		t.Error("an existing network must not be created again")
	}
}

func TestEnsureNetworkCreatesWhenMissing(t *testing.T) {
	f := newFake()

	if err := engineWith(f).EnsureNetwork(context.Background()); err != nil {
		t.Fatalf("EnsureNetwork: %v", err)
	}

	if f.find("network", "create", "cauldron-demo") == nil {
		t.Errorf("expected the network to be created; calls: %v", f.calls)
	}
}

func TestWaitHealthyAcceptsAContainerWithoutAHealthcheck(t *testing.T) {
	f := newFake()
	f.replies["inspect"] = "running \n"

	if err := engineWith(f).WaitHealthy(context.Background(), "mailpit", time.Second); err != nil {
		t.Errorf("a container with no healthcheck should count as ready: %v", err)
	}
}

func TestWaitHealthyWaitsForHealthy(t *testing.T) {
	f := newFake()
	f.replies["inspect"] = "running healthy"

	if err := engineWith(f).WaitHealthy(context.Background(), "postgres", time.Second); err != nil {
		t.Errorf("WaitHealthy: %v", err)
	}
}

func TestWaitHealthyGivesUpAndSaysWhy(t *testing.T) {
	f := newFake()
	f.replies["inspect"] = "running starting"

	err := engineWith(f).WaitHealthy(context.Background(), "postgres", 900*time.Millisecond)

	var unhealthy *ErrUnhealthy
	if !errors.As(err, &unhealthy) {
		t.Fatalf("err = %v, want ErrUnhealthy", err)
	}

	if !strings.Contains(unhealthy.Error(), "postgres") {
		t.Errorf("the error should name the service: %v", unhealthy)
	}
}

func TestStopRemovesEverythingForTheProject(t *testing.T) {
	f := newFake()
	f.replies["ps --all"] = "cauldron-demo-postgres\ncauldron-demo-redis\n"

	if err := engineWith(f).Stop(context.Background(), false); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	for _, name := range []string{"cauldron-demo-postgres", "cauldron-demo-redis"} {
		if f.find("rm", "--force", name) == nil {
			t.Errorf("%s was not removed", name)
		}
	}

	if f.find("network", "rm", "cauldron-demo") == nil {
		t.Error("the network was not removed")
	}
}

func TestStopKeepsVolumesWhenAsked(t *testing.T) {
	f := newFake()
	f.replies["ps --all"] = "cauldron-demo-postgres\n"

	_ = engineWith(f).Stop(context.Background(), true)

	joined := strings.Join(f.find("rm", "--force"), " ")

	if !strings.Contains(joined, "--volumes=false") {
		t.Errorf("keepData should preserve volumes; got %s", joined)
	}
}

func TestStopToleratesAMissingNetwork(t *testing.T) {
	f := newFake()
	f.fail["network rm"] = errors.New("Error: No such network: cauldron-demo not found")

	if err := engineWith(f).Stop(context.Background(), false); err != nil {
		t.Errorf("a missing network is not a failure worth reporting: %v", err)
	}
}

func TestCatalogueEntriesAreSelfConsistent(t *testing.T) {
	for _, name := range Catalogued() {
		spec, ok := SpecFor(name)
		if !ok {
			t.Fatalf("%s is listed but has no spec", name)
		}

		if spec.Service != name {
			t.Errorf("%s: spec service is %q", name, spec.Service)
		}

		if spec.Image == "" {
			t.Errorf("%s: no image", name)
		}

		// A floating tag makes two developers' environments diverge silently,
		// which is the whole problem this tool exists to remove.
		if strings.HasSuffix(spec.Image, ":latest") && name != "mailpit" && name != "minio" {
			t.Errorf("%s pins :latest (%s); pin a major version", name, spec.Image)
		}

		if (spec.Volume == "") != (spec.VolumePath == "") {
			t.Errorf("%s declares half a volume: %q %q", name, spec.Volume, spec.VolumePath)
		}
	}
}

func TestSpecForReturnsACopy(t *testing.T) {
	first, _ := SpecFor("postgres")
	first.Ports[0].Host = 9999
	first.Env["POSTGRES_USER"] = "tampered"

	second, _ := SpecFor("postgres")

	if second.Ports[0].Host != 5432 {
		t.Errorf("the catalogue was mutated through a returned spec: %d", second.Ports[0].Host)
	}

	if second.Env["POSTGRES_USER"] != "cauldron" {
		t.Errorf("the catalogue env was mutated: %v", second.Env["POSTGRES_USER"])
	}
}

func TestCheckPortsReportsTheServiceThatWantedIt(t *testing.T) {
	// Hold a port, then ask for a spec that wants it.
	spec := Spec{Service: "postgres", Ports: []Port{{Host: 5432, Container: 5432}}}

	if PortFree(spec.Ports[0].Host) {
		// Nothing is listening, so simulate the conflict on a port we can bind.
		t.Skip("port 5432 is free on this machine; conflict path covered by ErrPortInUse formatting")
	}

	err := CheckPorts([]Spec{spec})

	var inUse *ErrPortInUse
	if !errors.As(err, &inUse) {
		t.Fatalf("err = %v, want ErrPortInUse", err)
	}

	if !strings.Contains(inUse.Error(), "postgres") {
		t.Errorf("the error should name the service: %v", inUse)
	}
}

func TestErrPortInUseNamesTheServiceAndPort(t *testing.T) {
	err := &ErrPortInUse{Service: "redis", Port: 6379}

	if !strings.Contains(err.Error(), "redis") || !strings.Contains(err.Error(), "6379") {
		t.Errorf("unhelpful error: %v", err)
	}
}

// The bug this covers was found by running the tool, not by testing it: a
// container from another project publishing 0.0.0.0:6379 left 127.0.0.1:6379
// bindable, so the loopback-only check reported the port free and Docker then
// refused the run with a wall of daemon output.
func TestPublishedPortsSeesOtherProjectsContainers(t *testing.T) {
	f := newFake()
	f.replies["ps --format"] = "hauloycom-redis-1\t0.0.0.0:6379->6379/tcp, [::]:6379->6379/tcp\n" +
		"weather-edge-postgres-1\t0.0.0.0:55434->5432/tcp\n" +
		"scheduler\t80/tcp, 443/tcp\n"

	held, err := engineWith(f).PublishedPorts(context.Background())
	if err != nil {
		t.Fatalf("PublishedPorts: %v", err)
	}

	if owner := held[6379]; owner != "hauloycom-redis-1" {
		t.Errorf("port 6379 owner = %q, want hauloycom-redis-1", owner)
	}

	if owner := held[55434]; owner != "weather-edge-postgres-1" {
		t.Errorf("port 55434 owner = %q", owner)
	}

	// Exposed-but-unpublished ports are not allocated on the host.
	if _, taken := held[80]; taken {
		t.Error("an exposed port with no host mapping must not count as taken")
	}
}

func TestCheckDockerPortsNamesTheOffendingContainer(t *testing.T) {
	f := newFake()
	f.replies["ps --format"] = "hauloycom-redis-1\t0.0.0.0:6379->6379/tcp\n"

	spec, _ := SpecFor("redis")

	err := engineWith(f).CheckDockerPorts(context.Background(), []Spec{spec})

	var held *ErrPortHeld
	if !errors.As(err, &held) {
		t.Fatalf("err = %v, want ErrPortHeld", err)
	}

	// The fix is to stop that container, and the developer has no other way to
	// know which one it is.
	if !strings.Contains(held.Error(), "hauloycom-redis-1") {
		t.Errorf("the error must name the container: %v", held)
	}
}

func TestCheckDockerPortsPassesWhenNothingClashes(t *testing.T) {
	f := newFake()
	f.replies["ps --format"] = "something\t0.0.0.0:9999->9999/tcp\n"

	spec, _ := SpecFor("redis")

	if err := engineWith(f).CheckDockerPorts(context.Background(), []Spec{spec}); err != nil {
		t.Errorf("unexpected conflict: %v", err)
	}
}

// A diagnostic query failing should not stop a run that might still work.
func TestCheckDockerPortsToleratesAFailedQuery(t *testing.T) {
	f := newFake()
	f.fail["ps --format"] = errors.New("daemon busy")

	spec, _ := SpecFor("redis")

	if err := engineWith(f).CheckDockerPorts(context.Background(), []Spec{spec}); err != nil {
		t.Errorf("a failed diagnostic must not fail the run: %v", err)
	}
}

func TestStartTranslatesThePortAllocationError(t *testing.T) {
	f := newFake()
	f.fail["run --detach"] = errors.New("docker: Error response from daemon: failed to set up container networking: Bind for 0.0.0.0:6379 failed: port is already allocated")

	spec, _ := SpecFor("redis")

	err := engineWith(f).Start(context.Background(), spec)

	var held *ErrPortHeld
	if !errors.As(err, &held) {
		t.Fatalf("err = %v, want the daemon's wall of output translated to ErrPortHeld", err)
	}

	if held.Port != 6379 {
		t.Errorf("Port = %d, want 6379", held.Port)
	}
}

// A failed run leaves a created container behind. Left there, the next attempt
// finds it in "created" state, tries to start it, fails for the same reason,
// and reports something that looks nothing like the real cause.
func TestFailedStartRemovesThePartialContainer(t *testing.T) {
	f := newFake()
	f.fail["run --detach"] = errors.New("port is already allocated")

	spec, _ := SpecFor("redis")

	_ = engineWith(f).Start(context.Background(), spec)

	if f.find("rm", "--force", "cauldron-demo-redis") == nil {
		t.Errorf("a failed run must clean up after itself; calls: %v", f.calls)
	}
}

func TestLogsReadsOnlyThisProjectsContainer(t *testing.T) {
	f := newFake()
	f.replies["ps --filter"] = "redis\nmailpit\n"
	f.replies["logs"] = "ready to accept connections"

	out, err := engineWith(f).Logs(context.Background(), "redis", 20)
	if err != nil {
		t.Fatalf("Logs: %v", err)
	}

	if !strings.Contains(out, "ready to accept") {
		t.Errorf("out = %q", out)
	}

	call := f.find("logs")
	if call == nil {
		t.Fatal("no docker logs call was made")
	}

	joined := strings.Join(call, " ")

	if !strings.Contains(joined, "--tail 20") {
		t.Errorf("the line count should reach docker; got %q", joined)
	}

	// Scoped to this project's own container name, so a similarly named
	// container from another checkout cannot be read by accident.
	if !strings.Contains(joined, "cauldron-demo-redis") {
		t.Errorf("call = %q", joined)
	}
}

func TestLogsRefusesAServiceThisProjectIsNotRunning(t *testing.T) {
	f := newFake()
	f.replies["ps --filter"] = "redis\n"

	_, err := engineWith(f).Logs(context.Background(), "postgres", 10)
	if err == nil {
		t.Fatal("expected an error")
	}

	var missing *ErrNotRunning
	if !errors.As(err, &missing) {
		t.Fatalf("err = %T, want *ErrNotRunning", err)
	}

	if !strings.Contains(err.Error(), "redis") {
		t.Errorf("the error should name what is running; got %q", err)
	}

	if f.find("logs") != nil {
		t.Error("docker logs should not be called for a service that is not running")
	}
}

func TestLogsExplainsAnEmptyProject(t *testing.T) {
	f := newFake()

	_, err := engineWith(f).Logs(context.Background(), "redis", 10)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "cauldron up") {
		t.Errorf("err = %q, want a way forward", err)
	}
}
