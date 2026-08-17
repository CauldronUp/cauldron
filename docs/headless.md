# Headless mode

Cauldron's job is the emulated providers. Your application, its database and
its web server are yours, and headless mode is how you say so: no containers,
no plan describing an environment Cauldron is not going to set up, and one line
of JSON instead of a banner addressed to a person.

```bash
cauldron serve --headless stripe woocommerce
```

```json
{"address":"http://127.0.0.1:4600","bind":"127.0.0.1","port":4600,"control":"http://127.0.0.1:4600/_cauldron/status","mounted":[{"recipe":"stripe","url":"http://127.0.0.1:4600/stripe"},{"recipe":"woocommerce","url":"http://127.0.0.1:4600/woocommerce"}],"missing":[],"unseeded":[]}
```

Everything below is exercised by the `headless` job in `.github/workflows/ci.yml`
on every push. Documentation nobody executes drifts, and the shape of this
output is something other people's environments depend on.

## The contract

| Field | What it is |
|---|---|
| `address` | A base URL that works **from this machine**. Always dialable, even when the bind is a wildcard. |
| `bind` | The interface actually bound. `0.0.0.0` has no usable URL of its own. |
| `port` | The port, separately, for callers arriving from somewhere else. |
| `control` | The control plane: seed, reset, fault, clock, request log. |
| `mounted` | One entry per emulated provider, with its base URL. |
| `missing` | Providers detected in the project that Cauldron **cannot** emulate. These still reach the real network. |
| `unseeded` | Providers a global `--fixture` did not fit. |

Four things worth knowing:

- **The line is printed after the socket is bound.** Reading it means a
  connection will be accepted, so there is no readiness race to poll around.
- **`missing` is not decoration.** A provider without a Recipe still reaches the
  real network. It is in the machine-readable output for exactly the callers
  that cannot see a banner.
- **`missing` and `unseeded` are `[]`, never `null`**, so you can iterate
  without checking.
- **JSON is on stdout, diagnostics on stderr.** `cauldron serve --headless … 2>/dev/null | head -1`
  gives you one parseable line and nothing else.

`address` is the field to be careful with. It is correct from the machine
Cauldron is running on, and that is all it promises. A caller in a container has
a different route to the same server, so use `port` and your own hostname for
the host — `[::]:4600` is not a URL anybody can use, which is why it is not
what gets reported.

## Cauldron in your compose file

There is no published image yet, so build it from the `Dockerfile` in this
repository:

```bash
docker build -t cauldron:latest https://github.com/CauldronUp/cauldron.git
```

It binds `0.0.0.0` by default, because a container that binds loopback is
reachable from nothing but itself.

```yaml
services:
  cauldron:
    build:
      context: https://github.com/CauldronUp/cauldron.git
    command: ["serve", "--headless", "--host", "0.0.0.0", "stripe", "woocommerce"]
    ports:
      # Only needed if something outside the compose network talks to it.
      - "4600:4600"
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "-", "http://127.0.0.1:4600/_cauldron/status"]
      interval: 2s
      timeout: 2s
      retries: 10

  app:
    build: .
    depends_on:
      cauldron:
        condition: service_healthy
    environment:
      STRIPE_API_BASE: http://cauldron:4600/stripe
      WOOCOMMERCE_URL: http://cauldron:4600/woocommerce
```

Your application reaches Cauldron at `http://cauldron:4600/<recipe>` — the
service name, on the compose network. Nothing needs publishing for that.

The healthcheck uses `wget` because it is in the image already; there is no
`curl`. The port is written out rather than derived, so change both if you pass
`--port`.

## Cauldron on the host, application in a container

Cauldron binds loopback by default, and **a container cannot reach the host's
loopback**. Bind the wildcard instead:

```bash
cauldron serve --headless --host 0.0.0.0 stripe
```

Then from the container:

```yaml
services:
  app:
    build: .
    extra_hosts:
      # Docker Desktop provides this already. On Linux this line adds it.
      - "host.docker.internal:host-gateway"
    environment:
      STRIPE_API_BASE: http://host.docker.internal:4600/stripe
```

Binding past loopback puts the control plane on the network, where anything
that can reach the machine can seed, reset and fault the providers. Cauldron
says so on stderr when you do it. In a container that message is expected —
there is no other useful thing to bind — and it stays on stderr so it never
lands in the JSON.

## Cauldron in a container, application on the host

Publish the port and point at localhost:

```bash
docker run --rm -p 4600:4600 cauldron:latest \
  serve --headless --host 0.0.0.0 stripe wordpress
```

```bash
export STRIPE_API_BASE=http://127.0.0.1:4600/stripe
```

## Reading the line from a script

```bash
#!/usr/bin/env bash
set -euo pipefail

cauldron serve --headless --port 4600 stripe github > cauldron.json 2>cauldron.log &
cauldron_pid=$!
trap 'kill $cauldron_pid' EXIT

# The line appears once the socket is bound, so this waits for readiness
# rather than for a fixed number of seconds.
until [ -s cauldron.json ]; do sleep 0.1; done

base=$(jq -r '.mounted[] | select(.recipe == "stripe") | .url' cauldron.json)

# Anything Cauldron cannot emulate would reach the real network. Refuse rather
# than find out later.
if [ "$(jq '.missing | length' cauldron.json)" -gt 0 ]; then
  echo "Not emulated, would hit the real API: $(jq -c '.missing' cauldron.json)" >&2
  exit 1
fi

STRIPE_API_BASE="$base" ./run-tests
```

That `missing` check is the part worth copying. It is the difference between a
test suite that is offline and one that is offline except for the provider
nobody wrote a Recipe for yet.

## In CI

Cauldron is a single static binary with the Recipes compiled in, so the
simplest thing is usually not to involve a container at all:

```yaml
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - run: go install github.com/CauldronUp/cauldron/cmd/cauldron@latest
      - name: Start the emulated providers
        run: |
          cauldron serve --headless stripe wordpress > cauldron.json &
          until [ -s cauldron.json ]; do sleep 0.1; done
      - run: ./run-tests
        env:
          STRIPE_API_BASE: http://127.0.0.1:4600/stripe
```

If your job already runs in a container, build the image and start it as a
step. A GitHub Actions *service* container cannot take a custom command, and
the image's default detects nothing without recipes named, so a step is the
form that works:

```yaml
      - name: Start the emulated providers
        run: |
          docker build -t cauldron:ci https://github.com/CauldronUp/cauldron.git
          docker run -d --name cauldron -p 4600:4600 cauldron:ci serve --headless --host 0.0.0.0 stripe wordpress
          until curl -sf http://127.0.0.1:4600/_cauldron/status > /dev/null; do sleep 0.2; done
```

This repository's own `headless` CI job does exactly that, which is what keeps
this page honest.

## Pointing an SDK at it

Every provider names this differently, and several make it awkward on purpose.

```bash
# Stripe
STRIPE_API_BASE=http://127.0.0.1:4600/stripe

# WooCommerce and WordPress: the site URL, and the REST path follows it
WOOCOMMERCE_URL=http://127.0.0.1:4600/woocommerce
WP_URL=http://127.0.0.1:4600/wordpress

# GitHub
GITHUB_API_URL=http://127.0.0.1:4600/github
```

In code, it is whatever the SDK calls a base URL:

```php
// stripe-php
$stripe = new \Stripe\StripeClient([
    'api_key'  => 'sk_test_cauldron',
    'api_base' => getenv('STRIPE_API_BASE'),
]);
```

```js
// stripe-node
const stripe = new Stripe('sk_test_cauldron', {
  host: '127.0.0.1',
  port: 4600,
  protocol: 'http',
  basePath: '/stripe/v1',
});
```

```go
// go-github
client := github.NewClient(nil)
client.BaseURL, _ = url.Parse("http://127.0.0.1:4600/github/")
```

The credentials each Recipe accepts are in its `auth.keys`. They are fixtures,
not secrets, and they are in the repository on purpose.

## Driving it from your tests

The control plane is the same whether Cauldron is headless or not:

```bash
# Load seed data before a test
curl -X POST http://127.0.0.1:4600/_cauldron/stripe/seed -d '{"fixture":"small-shop"}'

# Break the next two calls, with the provider's real rate-limit shape
curl -X POST http://127.0.0.1:4600/_cauldron/stripe/fault -d '{"error":"rate_limit","count":2}'

# Age everything a month, so a subscription falls into dunning
curl -X POST http://127.0.0.1:4600/_cauldron/clock/advance -d '{"duration":"30d"}'

# Back to clean, between tests
curl -X POST http://127.0.0.1:4600/_cauldron/reset
```

A reset between tests is cheaper than a fresh process, and the identifier
generator is seeded, so the same fixture produces the same ids on every machine
and every run.

## What headless mode does not do

- It does not start your database, queue or mail catcher. `cauldron up` does
  that; headless mode is the opposite choice.
- It does not run your application. Nothing in Cauldron ever has.
- It does not proxy traffic or install certificates. Providers are reached by
  pointing a base URL at Cauldron, which means the swap is visible in your
  configuration rather than hidden in your network.
