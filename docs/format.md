# The Recipe format

A Recipe is one file, `recipes/<name>/recipe.yaml`, describing a real HTTP API's wire
behaviour precisely enough that Cauldron can stand in for it. This is the reference for
every key that file may contain. The Go source is the authority — `internal/recipe` — and
a test in that package fails if this document names a key that does not exist, or misses
one that does.

Nothing here tells you what to write. The interesting part of a Recipe is the header
comment: what was probed, what the provider does that its own documentation denies, and
what was deliberately not modelled. Read a few before writing one. `recipes/rootly`,
`recipes/mezmo`, `recipes/firehydrant` and `recipes/census` are good ones to start from.

## Two rules that are not fields

**Say where every claim came from.** Every conformance case carries `source`. A case that
was actually run against the live provider also carries `verified` with the date. A case
read from documentation and never run carries `source` alone. The distinction is the whole
point: `cauldron verify` reports the split, and a Recipe that dates a case it never ran is
worse than one that admits it read the manual.

**Never put a real secret in a Recipe.** Keys under `auth.keys` are fixtures, published
in a public repository, and read by anyone. Invent them.

---

## Top level

| key | type | meaning |
|---|---|---|
| `recipe` | string | The Recipe's name, matching its directory. |
| `version` | string | This Recipe's own version, not the provider's. |
| `capability` | string | What kind of provider this is, so a hundred Recipes can be found by what they do rather than by whether you remember the company's name. One of: `payments`, `banking`, `accounting`, `tax`, `payroll`, `email`, `sms`, `chat`, `push`, `voice`, `auth`, `identity`, `crm`, `support`, `marketing`, `brokerage`, `commerce`, `shipping`, `storage`, `database`, `queue`, `search`, `cdn`, `hosting`, `observability`, `analytics`, `flags`, `ci`, `vcs`, `issues`, `docs`, `calendar`, `files`, `media`, `ai`, `signing`, `scheduling`, `hr`, `forms`, `cms`, `infrastructure`. |
| `upstream` | block | Which API this is, and where its description lives. |
| `auth` | block | The credential, and what refusing one looks like. |
| `resources` | map | The records the API serves, by name. |
| `routes` | list | The endpoints, in order. |
| `responses` | block | How lists, single records and failures are wrapped. |
| `errors` | map | The provider's failures, by a name this Recipe chooses. |
| `fixtures` | map | Named sets of seed records. |
| `required_headers` | map | Headers a request must carry. |
| `webhooks` | block | What this provider delivers, and how it signs it. |
| `conformance` | list | The evidence that this Recipe resembles the provider. |

### upstream

| key | type | meaning |
|---|---|---|
| `api` | string | The provider's own version label: `v1`, `v2`, `2026-06-30`. |
| `docs` | string | Where a person reads about it. |
| `provider` | string | Whose API this is, which `api` alone cannot say. |
| `spec` | string | The URL of the provider's own machine-readable description. |
| `spec_hash` | string | A fingerprint of the parts of that description this Recipe claims. `cauldron drift` compares it. |
| `spec_seen` | string | The date the fingerprint was taken, `YYYY-MM-DD`. |
| `supersedes` | list | Versions of the same API this one replaced. |

Each entry under `supersedes`:

| key | type | meaning |
|---|---|---|
| `version` | string | The provider's label for it, in the same words `api` uses. |
| `host` | string | What a client of that version talks to — the only thing detection can tell versions apart by. |
| `note` | string | Why it matters, for a person reading the Recipe. |

Record the fingerprint **last**. Editing a Recipe moves it, because it covers the
intersection of the published description and the Recipe's own claims.

### required_headers

A map of header name to:

| key | type | meaning |
|---|---|---|
| `error` | string | The error to raise when the header is missing. |
| `methods` | list | Limits the requirement to those HTTP methods. |

---

## auth

| key | type | meaning |
|---|---|---|
| `scheme` | string | One of `bearer`, `basic`, `header`, `query`, `body`, `none`. `body` is the credential as a field in the request body, which Canny does: it cannot be set once as a default header, so every call site carries it, and it lands in anything that logs bodies. |
| `header` | string | The header carrying the credential, when the scheme is `header`. |
| `param` | string | The parameter carrying it, when the scheme is `query` or `body`. |
| `prefix` | string | Stripped before comparison, e.g. `"Bearer "`. |
| `credential` | string | Which half of a basic credential carries the secret: `username` (the default, Twilio's account SID) or `password` (Mailgun, whose username is the constant `api`). |
| `keys` | list | The credentials the emulator accepts. Fixtures only. |
| `pattern` | string | Accepts any credential matching this regular expression, for schemes where the value is computed per request — AWS signs every call, so there is no fixed value to hold. **This is the acceptance test**, not a shape check. |
| `shape` | string | What a credential has to look like *before* its value is worth comparing. Failing it is its own verdict; passing it goes on to the key comparison. Every key must match it. |
| `shape_error` | string | The errors entry for a credential the shape ruled out. Falls back to `malformed_error`. |
| `absent_error` | string | The errors entry for a request carrying no credential at all. |
| `malformed_error` | string | The errors entry for a credential the carrier could not read: the prefix is missing, or the prefix was all there was. |
| `rejected_error` | string | The errors entry for a credential of the right shape that this Recipe does not hold. |
| `unentitled` | list | Keys that authenticate and are refused anyway. |
| `unentitled_error` | string | The errors entry those raise. |
| `after_routing` | bool | Check the credential only once the request has matched a route, so an unrouted path answers 404 and a wrong method answers 405 whatever was sent. |
| `also` | list | Other carriers the same credential may arrive in. |

### The verdicts

A credential check reaches one of six conclusions, and a Recipe naming no error at all
gets one message for every failing one — exactly as it did before the distinction existed.

- **accepted** — one this Recipe holds.
- **absent** — nothing presented. The header or query parameter is not on the request.
  `"Bearer "` with nothing after it is *not* absent; somebody sent that.
- **malformed** — present and unreadable: the declared prefix is missing, or was all there was.
- **misshapen** — read, and ruled out by `shape` before its value was compared.
- **rejected** — the right shape, and not one this Recipe holds.
- **unentitled** — one this Recipe holds, for something it may not have. Authentication
  succeeded and authorisation did not; the fallback status is 403 rather than 401.

`malformed` and `misshapen` are separate because providers separate them. Opsgenie answers
401 "Could not authenticate" to a header without its `GenieKey ` prefix and 422 "Key format
is not valid!" to a correctly prefixed value that is not a UUID.

### also — one secret, several carriers

An alternative names what differs, usually just the carrier, and inherits the rest. It
takes the same keys as `auth` itself, may not name alternatives of its own, and must name a
scheme, header or param — one naming nothing is the primary written twice.

Alternatives are tried only when the primary carrier holds nothing, and the most
informative verdict wins: somebody who put a credential in the primary and got it wrong has
made a different mistake from somebody who used another channel.

**The prefix is the one field where empty does not mean inherit.** A prefix belongs to a
carrier — `"Bearer "` means nothing to an `X-Api-Key` header — so write `prefix: "-"` for an
alternative taking the bare secret. Validation refuses an alternative that changes scheme
and says nothing about the prefix.

---

## resources

A map of resource name to:

| key | type | meaning |
|---|---|---|
| `collection` | string | The plural name the provider wraps lists in. |
| `id` | block | How identifiers are shaped. |
| `alias` | string | A second field a path may address this resource by. |
| `fields` | map | The record's fields. |
| `constants` | map | Fields the provider always sends with a fixed value, such as Stripe's `object` discriminator and `livemode`. |
| `version_field` | string | A field the provider keeps as an optimistic lock: a number that moves on every write, which a caller has to quote back. |
| `version_conflict` | string | The failure a stale write is refused with. |
| `version_missing` | string | The failure a write carrying no version at all is refused with. |

### id

| key | type | meaning |
|---|---|---|
| `style` | string | `prefixed` (`cus_abc123`), `numeric`, `timestamp` (Slack's `1767225600.000100`), `opaque` (a bare random string, SendGrid), `uuid`, `hex` (Intercom, and anything out of MongoDB) or `digits` (Discord snowflakes, and any provider whose ids are long numeric strings that must not be parsed as numbers). |
| `prefix` | string | The prefix, for the `prefixed` style. |
| `other_prefixes` | list | Prefixes a record may legitimately carry that this Recipe does not mint. |
| `length` | int | How long the generated part is. |
| `field` | string | The property the provider returns the identifier in. |
| `type` | string | The JSON type it travels as: `string` (the default) or `number`. |
| `pattern` | string | The shape an identifier must have for the provider to look it up at all. |
| `carried_by` | string | The field holding the identifier when the provider does not send it under a name of its own. |

### fields

A map of field name to:

| key | type | meaning |
|---|---|---|
| `type` | string | `string`, `integer`, `number`, `boolean`, `timestamp` (Unix seconds, Stripe and Twilio), `timestamp_ms` (milliseconds, Clerk and most JavaScript-first APIs), `timestamp_ms_string`, `datetime` (RFC 3339, GitHub and HubSpot), `msdate`, `list` or `map`. |
| `required` | bool | Whether a create must supply it. |
| `default` | any | The value applied when the caller does not. |
| `stamped` | bool | Whether a time field is filled in from the sandbox clock when unsupplied. |
| `null_when_unset` | bool | Send the field as `null` rather than leaving it out. |
| `in` | string | Nest this field under a sub-object on the wire. |
| `as` | string | The name this field takes on the wire, when it differs from the name it is stored under. |

---

## routes

A list, matched in order.

| key | type | meaning |
|---|---|---|
| `method` | string | The HTTP method. |
| `path` | string | The path, with `{name}` parameters. |
| `resource` | string | Which resource it serves. |
| `operation` | string | One of `create`, `get`, `list`, `update`, `delete`. |
| `status` | int | The success status, when it is not the default for the operation. |
| `scope` | list | Path parameters that partition this resource, so records do not leak between them. |
| `fields` | map | Constants this route adds to its response body. |
| `headers` | map | Response headers this route sets. |
| `returns` | list | Limits the response to the named fields, for routes answering with less than the record they touched. |
| `not_found` | string | The error for a record that is not there, instead of `resource_missing`. |
| `error` | string | A failure this route always answers with, whatever the request. |
| `empty_body` | bool | Send no body at all. |
| `raw` | block | A body sent verbatim, for a provider that does not send JSON. |
| `list` | block | Overrides the Recipe-wide list envelope for this route. |
| `envelope` | block | Overrides how this route wraps a single object. |
| `pagination` | block | How this listing pages. |
| `filters` | list | Query parameters that narrow a listing. |
| `beside` | list | Other resources whose records travel in the same response body, each under its own collection name. |
| `id_from` | string | Where the identifier comes from when it is not a path parameter: `query:channel` or `body:channel`. |
| `lookup_by` | string | The field that value is matched against, for routes addressing a record by something that is not its identifier. |
| `id_as` | string | Renames the identifier on this route alone. |
| `selects` | string | Disambiguates several routes sharing one path by what a GraphQL query asks for. |
| `selects_body` | string | The same job, looking anywhere in the request body. |
| `matches_header` | map | Request headers whose values pick this route. |
| `matches_query` | map | Query parameters whose values pick this route. |
| `deleted_body` | string | What a delete answers with: `receipt`, `record`, `flagged`, `id` or `empty`. |
| `deleted_key` | string | The key the identifier arrives under, for the `id` body. |
| `emits` | string | The webhook event this route fires, for providers whose event names are not `resource.created`. |
| `emits_when` | list | Events that fire only when a particular field changes. Each entry takes `event` and `field`. |
| `public` | — | Marks a route the provider answers without a credential. `true` or `always` never examines one; `when-absent` exempts only a request presenting nothing, so a wrong credential is still refused. |
| `auth` | block | This route's own credential. |

### A route's own auth

Takes every key `auth` takes, and **inherits field by field** anything it does not name —
so a route with a second surface's key writes only `keys`. The prefix is the exception,
exactly as it is for `also`: write `prefix: "-"` when moving to a carrier that takes the
bare secret. A route cannot be both `public` and credentialled.

`after_routing` on a route governs requests that reach it, and nothing else. Ordering is
mostly a claim about what happens when routing *fails*, and a request matching no route has
no route to take an ordering from — the Recipe's own decides there.

### pagination

| key | type | meaning |
|---|---|---|
| `style` | string | `cursor`, `offset`, `page`, or `none` -- the provider serves the whole collection, takes no page size and no position, and offers no next page. `none` is the whole declaration; naming any other key beside it is refused. |
| `limit` | int | The default page size. |
| `max_limit` | int | The largest page the provider will serve, for the ones that cap and trim. |
| `over_limit` | string | The failure for asking for a bigger page, for the ones that refuse instead. |
| `may_overshoot` | bool | The page size is advice and a page may carry **more** records than were asked for -- Missive: "A page may return more [items] than limit". Declared, a page that is not the last serves one extra, because `while len(page) == limit` terminates on the first page otherwise. |
| `may_undershoot` | bool | A page may carry **fewer** records than were asked for without being the last -- Onfleet "will return up to 64 tasks but may return fewer". Declared, a page that is not the last serves one fewer, never zero. The end is the cursor going absent, not a thin page. Not exclusive with `may_overshoot` -- Modern Treasury does both -- and when both are set the first page overshoots and the rest undershoot. |
| `limit_param` | string | The parameter carrying the page size, for providers not calling it `limit`. |
| `cursor_param` | string | The parameter carrying the position to resume from. |
| `first_page` | int | The number the provider gives its first page, for the `page` style. |
| `in` | string | Where the parameters travel: `query` (the default) or `body`. |

### filters

| key | type | meaning |
|---|---|---|
| `param` | string | The query parameter's name. |
| `field` | string | The record field it matches against. |
| `default` | string | The value applied when the parameter is absent — usually the reason the filter exists at all. |
| `all` | string | The value that turns the filter off, for providers that have one. |
| `values` | map | Expands a parameter value into the set of field values it covers, for filters whose vocabulary is not the field's. |

### raw

| key | type | meaning |
|---|---|---|
| `content_type` | string | The media type to send. |
| `text` | string | The body, byte for byte as it was recorded. |
| `empty` | bool | This route really does answer with nothing. |

---

## responses

How every response is wrapped, in four blocks.

| key | type | meaning |
|---|---|---|
| `list` | block | The envelope a collection travels in. |
| `resource` | block | The envelope a single record travels in. |
| `error` | block | The envelope every failure travels in. |
| `success` | block | Constants added to every successful response. |
| `bom` | bool | Put a UTF-8 byte-order mark in front of every body. Authorize.Net does, on everything: the first three bytes are `EF BB BF` and `json.loads` raises "Unexpected UTF-8 BOM" on it. Opt-in, because almost nobody sends one and a mark nobody sends would be its own infidelity. |

> **`responses.error` is not `errors`.** `responses.error` describes the *envelope* every
> failure is written into — where the message lives, what the code is called. `errors` is
> the table of individual failures. A key like `type_field` or `status_field` belongs to the
> envelope and is not valid inside an entry in the errors table; an entry has `type` and
> `status` instead. Two Recipes have been written with those keys in the wrong block.

### responses.list

| key | type | meaning |
|---|---|---|
| `style` | string | `envelope` (Stripe), `bare` (GitHub), `wrapped` (Shopify), `map` (Pusher, whose channels arrive as an object keyed by name, so looping over it as a list finds nothing) or `tuple`. |
| `key` | string | The wrapping property name, required when the style is `wrapped`. |
| `url` | bool | Echo the request path in the envelope, which Stripe does. |
| `cursor_field` | string | The property carrying the next cursor. |
| `cursor_null` | bool | Send it as `null` on the last page rather than leaving it out. |
| `cursor_url` | string | The cursor carries an address rather than a token: `absolute` or `path`. |
| `count_field` | string | How many records matched in total, which is not how many are on this page. |
| `count_means` | string | What that field counts, for providers where it is not how many matched. |
| `count_lookahead` | int | How many pages a lookahead count reaches, including the one being served. |
| `page_count_field` | string | How many records are on this page, for providers sending that beside the total. |
| `pages_field` | string | How many pages the whole set makes at this page size. |
| `page_field`, `limit_field` | string | Properties echoing the page number and page size the request asked for. |
| `count_as_string` | bool | Send the counts as strings. |
| `has_more_field` | string | A boolean saying whether more pages remain. |
| `complete_field` | string | A boolean saying the opposite. |
| `final_field` | string | A field sent only on the last page. |
| `link_header` | bool | Advertise the next page in an RFC 5988 `Link` header rather than in the body. |
| `prev_link` | bool | Add a `rel="prev"` beside it. |
| `entry_field` | string | Make each entry that one field's value rather than the whole record. |
| `entry_style` | string | `wrapped` wraps each item under the resource's own name. |
| `omit_when_empty` | bool | Leave the collection key out entirely when there is nothing to send. |
| `collapse_single` | bool | Send a collection of one as the object rather than a list of one. |
| `fields` | map | Constants added to a list response only. |

### responses.resource

| key | type | meaning |
|---|---|---|
| `style` | string | `bare` (the default) or `wrapped`. |
| `key` | string | The wrapping property name. |
| `array` | bool | Wrap the single object in a list. |

### responses.error

| key | type | meaning |
|---|---|---|
| `style` | string | `nested` (Stripe, the default), `flat` (GitHub), `list` (SendGrid, `{"errors": [{...}]}`, because one request can fail several ways at once), `string_list` (Datadog, the same array with bare strings in it), `string` or `text` (Trello, whose failures are not JSON at all, so a client calling `.json()` on one throws). |
| `key` | string | The property holding the array, when the style is `list`. |
| `message_field` | string | The property carrying the human-readable message in a flat envelope. |
| `code_field` | string | The property carrying the error code. |
| `code_type` | string | Whether the code is sent as a `string` or a `number`. |
| `status_field` | string | A property echoing the HTTP status inside the body, which Twilio does. |
| `type_field` | string | A property carrying the error category. |
| `fields` | map | Constants the provider adds to every error, such as GitHub's `documentation_url`. |

### responses.success

| key | type | meaning |
|---|---|---|
| `fields` | map | Constants added to every successful response. |

---

## errors

A map of a name this Recipe chooses to:

| key | type | meaning |
|---|---|---|
| `status` | int | The HTTP status. |
| `code` | string | The provider's code. |
| `type` | string | The provider's error category, often a much smaller set than the codes. |
| `message` | string | The sentence, in the provider's own words. |
| `headers` | map | Response headers this failure carries, such as `Retry-After`. |
| `fields` | map | Extra body properties, merged over the Recipe-wide ones. |
| `empty` | bool | A status line and no body at all. |
| `style` | string | Overrides the Recipe-wide envelope for this failure alone, because a provider can answer two shapes and the npm registry does. |
| `key` | string | Overrides the envelope key for this failure alone. |
| `message_field` | string | Overrides the field carrying the sentence, for this failure alone. |
| `code_field` | string | Overrides the field carrying the code; `-` removes it. |
| `type_field` | string | Overrides the field carrying the category, for this failure alone; `-` removes it. |

`{detail}` in a message or in a field is replaced per request with what the failure is
about — the credential that was refused, the path that was not found. NeverBounce answers
`Invalid API key 'whatever you sent'`, and New Relic echoes a wrong key into a body field.

Some names are reached by the runtime on its own, without anything pointing at them:
`unknown_route`, `method_not_allowed`, `resource_missing`, `authentication_error`,
`parameter_missing`, `conflict`, `invalid_request`, `id_malformed` and
`unsupported_operation`. Declaring one replaces the generic wording with the provider's.

---

## fixtures

A map of fixture name to a map of resource name to a list of records. `empty: {}` is a
fixture with nothing in it, which several Recipes use to prove a listing's empty shape.

```yaml
fixtures:
  small-team:
    alert:
      - id: 11111111-2222-4333-8444-555555555555
        message: Disk almost full
```

---

## webhooks

| key | type | meaning |
|---|---|---|
| `events` | list | The event types this provider delivers. |
| `payload` | any | The envelope the provider wraps the changed record in. |
| `signing` | block | How a delivery is signed. |

### webhooks.signing

| key | type | meaning |
|---|---|---|
| `scheme` | string | `hmac-sha256` or `none`. |
| `header` | string | The header the signature travels in. |
| `secret` | string | The signing secret. A fixture, like every other credential here. |
| `over` | string | What string the digest is taken over. |
| `encoding` | string | How the digest is written down. |
| `value` | string | How it is wrapped before it goes in the header. |
| `timestamp_header` | string | The header the signed timestamp travels in, for providers whose signature covers one the value does not carry. |

---

## conformance

A list of cases. Each is one request and what the provider answers.

| key | type | meaning |
|---|---|---|
| `name` | string | What the case claims, in a sentence. No two cases in one Recipe may share one. |
| `source` | string | The documentation or transcript the expectation came from. |
| `verified` | string | The date this case was last checked against the real API, `YYYY-MM-DD`. Only for a case that was actually run. |
| `fixture` | string | Seeded before the case runs. |
| `arm` | string | An entry in the errors table to install before this case's request, and only for it. |
| `request` | block | What to send. |
| `expect` | block | What must come back. |

### request

| key | type | meaning |
|---|---|---|
| `method` | string | The HTTP method. |
| `path` | string | The path. |
| `query` | map | Query parameters. |
| `headers` | map | Request headers. |
| `json` | any | A JSON body. An array is a body too. |
| `form` | map | `application/x-www-form-urlencoded`, which is what Stripe's own SDKs send. |
| `multipart` | list | `multipart/form-data`. The boundary is generated, because a fixture naming one would name something no client sends twice. |

Each part under `multipart`:

| key | type | meaning |
|---|---|---|
| `name` | string | The part's form field name. |
| `filename` | string | Marks the part as a file. |
| `content_type` | string | The part's own type, separate from the request's. |
| `text` | string | The part's content. |
| `json` | any | The part's content as a document, for a provider sending metadata beside a file. |

### expect

| key | type | meaning |
|---|---|---|
| `status` | int | The HTTP status. |
| `body` | map | Fields the response must carry, by dotted path or nested. |
| `matches` | map | Dotted field paths to regular expressions, for values correct in shape rather than exact — generated identifiers. |
| `absent` | list | Fields that must not appear. |
| `headers` | map | Response headers that must be present with those values. |
| `header_matches` | map | Response header names to regular expressions. |
| `absent_headers` | list | Response headers that must not appear. |
| `body_matches` | string | A regular expression applied to the raw body, without parsing it. |
| `no_body` | bool | The response body is empty. |
| `webhook` | block | What the request emitted. |

### expect.webhook

| key | type | meaning |
|---|---|---|
| `event` | string | The type the delivery must carry. |
| `body` | map | Dotted paths in the payload, envelope included. |
| `matches` | map | Regular expressions against payload paths. |
| `absent` | list | Paths the payload must not carry. |
| `none` | bool | The request emitted nothing at all — an event that fires when it should not is as wrong as one that does not fire. |
| `absent_events` | list | Events the request must not emit, for when it does emit something. |
| `signature` | string | A pattern the signature value must match. |
| `signature_header` | string | The header the signature travels in. |
| `header_matches` | map | Patterns against the delivery's other headers. |

---

## What validation refuses

Not an exhaustive list — `cauldron verify` will tell you the rest — but these are the rules
that exist because a Recipe once got them wrong quietly.

- **A scheme with nothing behind it.** `auth` naming a scheme and holding neither `keys` nor
  a `pattern` means "accept anything", so the emulator enforces no authentication at all.
  Four Recipes shipped that way. Use `scheme: none` to say a surface needs no credential.
- **A shape excluding the Recipe's own key.** A fake refusing the key its own README tells
  you to send is not a fake of anything.
- **`shape` and `pattern` together.** A pattern is itself the acceptance test, so the shape
  would decide nothing.
- **An override changing carrier without saying what happens to the prefix.** Empty means
  inherit, and inheriting `"Bearer "` into a Basic channel refuses credentials the provider
  accepts — a fake stricter than the thing it stands in for.
- **A paging block that says `none` and then describes paging.** `none` means the
  provider serves the whole collection; a page size or a parameter name beside it is
  a second claim the runtime ignores, so the file would describe paging that never
  happens.
- **A route both `public` and credentialled.** The exemption wins, so the credential is dead
  weight a reader takes for a claim.
- **A key in both `keys` and `unentitled`.** Nothing can say whether it works.
- **Naming an error that is not in the table**, from `auth`, a route, or a case's `arm`.
- **Two conformance cases sharing a name.**
- **Arming a failure and expecting a status it does not answer with**, so what the case
  installed changed nothing.
- **Claiming a field name nothing asserts.** A name a Recipe chooses is only a claim if a
  case asserts it where the value exists; renaming it otherwise would break nothing, which
  means the Recipe was not really saying anything.
