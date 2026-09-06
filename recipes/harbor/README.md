# harbor

Emulates the Harbor container-registry API (v2.0), for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against Harbor's own Swagger 2.0 document, published in its own repository — `goharbor/harbor`, `api/v2.0/swagger.yaml`, 135 paths — read on 2026-09-06. Harbor is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**Every boolean in a project's metadata is a string.** `ProjectMetadata` declares eight settings and every one is `type: string`, with the document's own descriptions saying the valid values are `"true"` and `"false"`:

| setting | what it controls |
|---|---|
| `public` | whether anyone may pull |
| `enable_content_trust` | whether unsigned images may be pulled |
| `enable_content_trust_cosign` | the same, for cosign signatures |
| `prevent_vul` | whether vulnerable images may run |
| `auto_scan` | whether pushes are scanned |
| `auto_sbom_generation` | whether an SBOM is produced |
| `reuse_sys_cve_allowlist` | which allowlist applies |

So `project.metadata.public` is the string `"false"` on a private project — and `"false"` is truthy in JavaScript, Python, Ruby, PHP and Perl. Code that shows a badge when a project is public shows it on every project; code that refuses to push when `prevent_vul` is set refuses always.

These are the settings deciding whether unsigned images may be pulled and whether vulnerable images may run, and the type carrying them cannot be tested for truth.

**A header decides what a path segment means.** `X-Is-Resource-Name`, a boolean header defaulting to false. In the document's own words: "the parameter which supports both name and id in the path … When the `X-Is-Resource-Name` is false and the parameter can be converted to an integer, the parameter will be as an id, otherwise, it will be as a name."

So `/projects/123` is project id 123 — unless you send the header, in which case it is the project *named* "123". The URL does not say which resource it addresses; a header does. A project whose name is all digits is unreachable without it, and two callers sending the same URL can get two different projects.

**A repository name must be URL-encoded twice.** From the `repository_name` parameter's own description:

> If it contains slash, encode it twice over with URL encoding. e.g. `a/b` → `a%2Fb` → `a%252Fb`

Nearly every repository name contains a slash — `library/nginx` is the shape of the thing. And every HTTP client library encodes path parameters once, for you, correctly. Passing `library/nginx` to a generated client produces `library%2Fnginx` on the wire, which is one encoding, and Harbor wants two. The right call requires encoding *before* handing the value to the thing that will encode it again — and the wrong one is not an error, it is a lookup for a repository nobody has.

**One path segment holds two kinds of identifier.** `reference` is documented as "can be digest or tag", so a tag named `sha256…` and a digest occupy the same segment with nothing saying which was meant.

**The total is a header and so is the next page.** A listing answers `X-Total-Count` and `Link` together, with nothing in the body — Harbor is the third Recipe here to need `count_header`, after gitea and basecamp.

**Two fields one character apart, one deprecated.** `current_user_role_id` is "deprecated and will be removed in future versions"; `current_user_role_ids` is the array that replaced it. The singular carries only the highest role, so a caller with two roles reading it sees one.

**A zero repository count is sent rather than dropped.** `repo_count` carries `x-omitempty: false` — a vendor extension existing because Go would otherwise omit it, and an absent count reads as *unknown* where a zero reads as *empty*.

**The search parameter is a language.** `q` supports, in the document's own words, "exact match(k=v)", "fuzzy match(k=~v)", "range(k=[min~max])", "list with union releationship(k={v1 v2 v3})" and "list with intersetion relationship(k=(v1 v2 v3))" — braces, parentheses, brackets and tildes, all needing URL encoding, inside a single query-string value. The two spelling mistakes are the document's.

**A failure is an array even when there is one of them.** `{"errors": [{"code": …, "message": …}]}`, so a client reading `body.message` finds `undefined` and reports the reason as nothing.

## Modelling limits

- **Nothing here is verified against a live API.** Harbor is self-hosted and there is no public instance.
- **Projects, listed and fetched.** 135 paths is a registry: repositories, artifacts, tags, scanners, replication, retention, robot accounts, webhooks and the quota surface each want their own evidence.
- **`X-Is-Resource-Name` is recorded and not served.** Serving it would mean routing one path two ways from a header, and what this Recipe would then demonstrate is its own routing rather than Harbor's.
- **The double encoding is recorded and not served.** This Recipe has no route taking a repository name, and adding one to demonstrate an encoding rule would demonstrate Go's URL parsing instead.
- **Only basic auth is modelled**, which is all the document declares. Robot accounts are real and are not in it.
