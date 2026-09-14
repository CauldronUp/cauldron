# tailscale

Emulates the Tailscale device listing for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from Tailscale's own Go client ([`tailscale-client-go`](https://github.com/tailscale/tailscale-client-go)), and struck live against `api.tailscale.com` on 2026-09-13 with no header, with a wrong bearer, with a Basic credential, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**One sentence for four different causes.** All four answer `{"message":"API token invalid"}` with a 401:

```
(no header)                    an absent credential
Bearer tskey-api-notarealkey   a wrong one
-u notarealkey:                a Basic header instead of a bearer
PUT on a real path             a method the route does not take
```

"Invalid", for a token that was never sent — and the same words for a request whose only fault was the verb.

**An unrouted path is Go's plain-text default, wrapped in JSON.** `{"message":"404 page not found"}` — the exact string `net/http` writes as `text/plain`, lifted into a message field. The sentence carries its own status code, so a client printing `message` prints "404" twice.

**A timestamp can be the empty string.** The vendor's own client defines a `Time` type whose unmarshaller says: "a blank string will keep the time at its zero value". That code exists because the API sends `""` where a date belongs — not `null`, not an absent key — on `created`, `expires` and `lastSeen`. A strict JSON date parser in any other language throws on it.

**The listing envelope has no type at all.** The client reads it as `map[string][]Device` and returns `resp["devices"]`, so the one key wrapping the collection is a string literal in the middle of a function rather than a declared shape. There is no count, no cursor and no paging of any kind beside it.

**A failure can be attributed to a user.** `APIError` is `{message, data}` where `data` is `[{user, errors}]` — a list of per-user objects, each with its own list of strings. One HTTP failure can carry several people's problems, and the live 401 carries none of it: `data` is simply absent.

**A device carries two public keys.** `machineKey` and `nodeKey` are on every record in the listing, so enumerating devices hands out both.

**And two fields say what it is called.** `name` and `hostname`, side by side, with nothing saying which one a human picked — and `addresses` beside them, a list, for the two or three ways to reach the same machine.

**Six booleans and no enum anywhere.** `authorized`, `keyExpiryDisabled`, `blocksIncomingConnections`, `isExternal`, `connectedToControl`, `updateAvailable` — every state on the record is its own flag.

## Sources

- [`tailscale/tailscale-client-go`](https://github.com/tailscale/tailscale-client-go) — `tailscale/client.go`: the `Device`, `Time` and `APIError` types.
- Live: `api.tailscale.com`, struck 2026-09-13 with no header, a wrong bearer, a Basic credential, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The device listing. ACLs, DNS, keys, devices-by-id, routes, webhooks, users, contacts and the logging surface are the rest.
- **The tailnet is a path parameter and is not partitioned.** Live, `-` means "the tailnet of the credential"; this Recipe serves the same devices for any tailnet name.
- **No `spec:`.** Tailscale's `api.md` now redirects to a documentation site, and no OpenAPI description is served at any address this Recipe could find — so there is nothing for `cauldron drift` to record.
- **The success fixture is client-derived.** Listing devices needs a real key; the records here are `Device`'s own fields with values of the declared shapes, including the empty-string timestamps the client's `Time` type exists to survive, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Tailscale is reached through `tailscale-client-go`, the `tailscale` CLI, or a plain HTTP call carrying a bearer token, and none of those resolve to this host through a dependency file. Checked 2026-09-13.
