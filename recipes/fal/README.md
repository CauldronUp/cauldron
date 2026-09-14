# fal

Emulates the fal queue status endpoint and the run endpoint beside it, for local development and tests.

**7 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from fal's own published JavaScript client ([`fal-ai/fal-js`](https://github.com/fal-ai/fal-js)), and struck live against `fal.run`, `queue.fal.run` and `rest.alpha.fal.ai` on 2026-09-13 with no credential, with a wrong one, on applications with two, three and four path segments, and on one that does not exist.

## What this Recipe found

**The refusal names an application you did not ask for.**

```
POST /fal-ai/flux/dev                  401  Cannot access application "fal-ai/flux"
POST /fal-ai/flux                      401  Cannot access application "fal-ai/flux"
POST /fal-ai/flux/dev/image-to-image   401  Cannot access application "fal-ai/flux"
```

Three different addresses, one name quoted back. The message keeps the first two path segments and drops the rest, so a caller debugging a 401 on `fal-ai/flux/dev` is told about `fal-ai/flux`, which is a different application. And "authentication is required" is what a wrong key is told too — the sentence does not change when a credential is present and refused.

**A GET on a POST-only endpoint is a 401.** The credential is checked before the method, so the failure a browser produces is about a token.

**Three hosts, three failure shapes.**

```
fal.run             {"detail":"Cannot access application \"fal-ai/flux\"…"}
queue.fal.run       {"status":"NOT_FOUND"}
rest.alpha.fal.ai   {"detail":"Authorization header is required"}
```

The queue answers with `status` and an uppercase enum where the run host answers with `detail` and a sentence — and the third host has `alpha` in its name and serves production token management.

**A doc comment describes a different field from the one it is on.** In the published client, `UrlOptions.subdomain` is typed `string` and documented "If `true`, the function will use the queue to run the function asynchronously and return the result in a separate call."

**The queue status is three shapes with no discriminator a type can use.** `queue_position` appears only when `IN_QUEUE`, `logs` only when it is not, and `metrics` only when `COMPLETED` and optionally even then. The client's own guard is `obj && obj.status && obj.response_url` — a truthiness check, so a status of the empty string would fail it.

**A log line's `source` is a type of one.** `source: "USER"`, and nothing else. Its `level` is `"STDERR" | "STDOUT" | "ERROR" | "INFO" | "WARN" | "DEBUG"` — two stream names in a severity enum.

**And the timestamp carries a note to the maintainer.** The shipped type says:

```ts
timestamp: string; // Using string to represent date-time format, but you could
                   // also use 'Date' type if you're going to construct Date objects.
```

**`Metrics` has one member and it is nullable** — `inference_time: number | null`; and the client's own `Result<T>` is `{data, requestId}`, camelCase, where every field on the wire is snake_case.

## Sources

- [`fal-ai/fal-js`](https://github.com/fal-ai/fal-js) — `libs/client/src/types/common.ts`: the queue status union, the log line, the metrics and the URL options.
- Live: `fal.run`, `queue.fal.run` and `rest.alpha.fal.ai`, struck 2026-09-13 with no credential, a wrong key, three path depths, and an application that does not exist.

## Modelling limits

- **Two routes, on one host.** Live these are two hostnames — `fal.run` for running and `queue.fal.run` for status — and this Recipe serves both paths from one sandbox. The finding is that the two answer failures in different shapes, which the cases assert.
- **The run endpoint answers only its failure.** A real key returns a model's output, whose shape depends on the model; this Recipe declines to invent one, so the route is here for the 401 it gives everyone else.
- **The application-not-found message is fixed.** Live it quotes the two-segment name from the request; the case that asserts it sends the name in the message.
- **No `spec:`.** fal documents this API as prose and publishes a JavaScript client; no OpenAPI description is served at any address this Recipe could find, so there is nothing for `cauldron drift` to record.
- **Nothing is mapped in detection.** fal is reached through `@fal-ai/client` on npm or `fal-client` on PyPI, and neither resolves to this host through a dependency file. Checked 2026-09-13.
