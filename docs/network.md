# Network conditions

Cauldron's fault injection covers what a provider does *deliberately* when
something is wrong: rate limits, validation errors, 5xx. `cauldron network`
covers what the network does to you regardless of what the provider intended.

The difference matters. A rate limit arrives as a well-formed 429 your client
library already understands and probably already retries. A connection that
hangs for ninety seconds and then dies arrives as nothing at all — no status, no
body, no error type your code has a branch for. That is the one that takes
production down at 3am, and it is the one no provider's sandbox will ever give
you.

```bash
cauldron network stripe --latency 800ms --jitter 200ms
cauldron network stripe --bandwidth 50
cauldron network stripe --timeout 30s
cauldron network stripe --reset --probability 0.1
cauldron network stripe --clear
```

## Why these words

The vocabulary is [Toxiproxy's](https://github.com/Shopify/toxiproxy) on
purpose: latency, jitter, bandwidth, timeout, reset, slice, limit.

Toxiproxy is the tool people already reach for when they want to break the
network between their application and its database. If you have run
`toxiproxy-cli toxic add -t latency` against Postgres, you already know what
`--latency` and `--jitter` mean here, and being made to learn a second set of
words for the same idea would buy nothing. Cauldron applies them to the
third-party APIs Toxiproxy cannot help you with, because those APIs do not exist
locally for it to sit in front of.

Working in PHP? [`mpge/toxiproxy-php`](https://github.com/mpge/toxiproxy-php)
gives the same vocabulary against real TCP services, and the two compose: point
Toxiproxy at your database and Cauldron at your payment provider, and every
dependency your application has can be made to misbehave from one test.

## What each condition does

| Flag | Effect |
|---|---|
| `--latency 800ms` | Delay every response by this long |
| `--jitter 200ms` | Vary that delay by ± this much |
| `--bandwidth 50` | Throttle the response body to 50 KB/s |
| `--timeout 30s` | Accept the request, answer nothing, close the connection after this long |
| `--reset` | Close the connection immediately, with no response at all |
| `--limit 1024` | Close after this many bytes of body — a response truncated mid-flight |
| `--slice 64` | Write the body in chunks of roughly this size |

And the modifiers, which work the same way they do on `cauldron fault`:

| Flag | Effect |
|---|---|
| `--probability 0.25` | Affect this share of requests (default: all of them) |
| `--count 3` | Affect only this many requests |
| `--for 30s` | Expire after this long of sandbox time |
| `--path /v1/charges` | Only affect paths containing this |
| `--clear` | Remove every armed condition |

`--reset` and `--timeout` end the request without a response. Everything else
degrades a response that still arrives.

## Where an HTTP emulator cannot be a TCP proxy

Two things are worth stating plainly, because a fake that quietly differs from
reality is worse than no fake at all.

**Only the body is throttled and sliced.** Response headers are written
normally. Throttling those would model a slow link more completely, but it would
mean breaking the `http.ResponseWriter` contract, and the body is where a
client's read loop actually lives. If you are testing header-timeout behaviour
specifically, put Toxiproxy in front of Cauldron.

**Slicing is HTTP chunking, not TCP segmentation.** `--slice` writes the body in
pieces and flushes each one, so a client that assumes a whole response arrives
in a single read will break — which is the bug being hunted. It does not control
packet boundaries below that.

Everything else is faithful. `--reset` hijacks the connection and closes it with
`SO_LINGER` at zero, so your client sees a genuine TCP reset rather than a clean
EOF; clients treat those differently and the difference is often the bug.

## Reproducibility

A probability below 1 draws from the sandbox seed, not the wall clock. The same
seed produces the same sequence of affected requests, so a failure you saw once
is a failure you can see again. Chaos that cannot be replayed is a bug report
nobody can act on.

## Seeing what happened

Armed conditions show up in `cauldron status`, and every affected request is
labelled in `cauldron requests`, so a baffling timing has a visible cause rather
than looking like a flake.

```
$ cauldron status
RECIPE   FIXTURE     REQUESTS  FAULTS  NETWORK
stripe   small-shop  14        0       latency 800ms ±200ms
```

## Ordering

Network conditions are applied before the provider's own behaviour, including
authentication and armed faults. That is the order reality uses: a connection
that never completes never reaches an application that could have rate limited
it.

## Over HTTP

```
POST /_cauldron/{recipe}/network
{"latency": "800ms", "jitter": "200ms", "probability": 0.5}

POST /_cauldron/{recipe}/network
{"clear": true}
```

Durations are strings so they carry their unit: `"800ms"` says what it means
where `800` would not.
