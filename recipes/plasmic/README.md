# plasmic

Emulates the Plasmic loader code bundle for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-14.**

Read from the types in Plasmic's own `@plasmicapp/loader-fetcher`, and struck live on 2026-09-14 with no credential, with a wrong project token, on a path that does not exist, with a method the path does not take, on the studio API and on the root.

## What this Recipe found

**The status line says 499 and the body says what really happened.**

```
no credential   HTTP/1.1 499 Unknown    Upstream responded with 401
wrong token     HTTP/1.1 499 Unknown    Upstream responded with 404
```

499 is nginx's non-standard code for a client that closed the connection before the server replied. It is not in the IANA registry, nothing generated it here but the server, its reason phrase is the literal word "Unknown", and the real status arrives as English prose in a `text/plain` body. A client switching on `response.status` sees a number that means nothing; the number that means something is inside a sentence.

**And a wrong credential is reported as the upstream not finding something.** The only difference between sending no token and sending a wrong one is `401` against `404`, inside that sentence, under the same 499.

**Four requests, four error formats, on one deployment.**

| status | content type | body |
|---|---|---|
| `499` | `text/plain;charset=UTF-8` | `Upstream responded with 401` |
| `404` | `application/json` | `{"error":{"name":"NotFoundError","statusCode":404,"message":"Not Found"}}` |
| `405` | `text/plain` | `Error: 119` (with `Google-Edge-Cache: bad request.`) |
| `403` | `application/json` | `{"error":{"name":"ForbiddenError","statusCode":403,"message":"Must be a normal user"}}` |

The 405's entire body is a CDN's internal error number. The 403 — struck on `studio.plasmic.app` — says "Must be a normal user", which is the sentence an API gives a caller about what kind of account they need to be.

**The JSON errors carry a JavaScript class name and repeat the status.** `error.name` is `NotFoundError` or `ForbiddenError`, the constructor name, and `error.statusCode` is the status again, inside the body it was already on the status line of.

**The success is JavaScript source code, as data.** The loader bundle's `modules` is `{browser: [...], server: [...]}`, and each entry is a `CodeModule` — `{fileName, code, imports, type: "code"}`. The API hands back a string of executable code per file, split by where it is meant to run.

**A subtype declares a field's type as `never`.** `ComponentMeta.plumeType` is `string | undefined`. `PageMeta extends ComponentMeta` and redeclares it `plumeType: never` — TypeScript's bottom type. A page can never carry one, said in a way no wire format can express and no runtime can check.

**A field enumerates what is not there.** `filteredIds` is `Record<string, string[]>`, documented as "the list of component IDs that are **not** included in the bundle" — an absence, listed, per project.

Also pinned: two booleans are named for defaults — `deferChunksByDefault` and `disableRootLoadingBoundaryByDefault` — on a per-request response; `ProjectMeta.indirect` is a bare boolean with no explanation anywhere, sitting under a source comment that reads "Keep in sync with platform/wab ProjectMeta", which is a request to a human embedded in a shared type; and the client pins the wire format with `const VERSION = "10"` in its own source, rather than in any path or header this API exposes.

## Sources

- [`plasmicapp/plasmic`](https://github.com/plasmicapp/plasmic) — `packages/loader-fetcher/src/api.ts` for `LoaderBundleOutput`, `ComponentMeta`, `PageMeta`, `ProjectMeta` and the module types.
- [Plasmic loader API documentation](https://docs.plasmic.app/learn/loader-api/).
- Live: `codegen.plasmic.app` and `studio.plasmic.app`, struck 2026-09-14 with no credential, a wrong project token, an unrouted path, a wrong method, and the root.

## Modelling limits

- **No description is published.** Plasmic serves no OpenAPI document for the loader API, so there is no `spec` to fingerprint; the shapes here come from the client library's own types.
- **The code in the fixture is a stub.** A real bundle carries the compiled output of a Plasmic project — thousands of lines per component. The finding is that executable code is the payload, and one representative module says that without shipping anyone's site.
- **The 403 is recorded, not served.** It was struck on `studio.plasmic.app`, a different host with a different surface; this Recipe serves the codegen host.
- **One route of several.** The published code bundle. The HTML API, the preview bundle, the assets endpoint and the studio's project API are the rest.
- **`platform` and `browserOnly` are not modelled.** The bundle's shape depends on them; this Recipe serves one.
- **Nothing is mapped in detection.** Plasmic is reached through `@plasmicapp/loader-nextjs` or `@plasmicapp/loader-react`, and neither resolves to this host through a dependency file. Checked 2026-09-14.
