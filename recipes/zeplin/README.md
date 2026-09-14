# zeplin

Emulates the Zeplin project listing for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the API reference at [`docs.zeplin.dev`](https://docs.zeplin.dev/reference/getprojects), and struck live on 2026-09-14 with no credential, with a wrong bearer token, with the token and no scheme, on a path that does not exist, with a method the path does not take, and on the root.

## What this Recipe found

**`message` holds the code and `detail` holds the sentence.** The two fields have swapped their conventional jobs:

```json
{"message":"invalid_token","detail":"Authorization header is missing"}
```

`message` — the field a client library will surface to a human — is a constant. `detail` — which reads as the machine-readable extra — is the prose.

**And on a routing failure the same field holds a sentence after all.** `{"message":"Not Found"}`. So `message` carries `invalid_token` on one response and `Not Found` on the next, and nothing in the body says which kind it is.

**The commonest credential mistake gets the least help.** Three ways to fail, struck live:

```
(no header)            {"message":"invalid_token","detail":"Authorization header is missing"}
Bearer <wrong token>   {"message":"invalid_token"}
<token, no prefix>     {"message":"invalid_token","detail":"Did you forget to add authorization type \"Bearer\" as a prefix?"}
```

A missing header is explained. A missing `Bearer ` prefix is explained, in a question, with escaped quotes inside the JSON string. A token that is simply *wrong* — what happens when a key is rotated or revoked, which is most of the time — gets no `detail` at all. The field is absent rather than empty, so code reading `body.detail` finds `undefined`.

**And the code is about a token on a request that has none.** `invalid_token`, for a request with no `Authorization` header at all.

**Seven counters ride on one project record.** `number_of_members`, `number_of_screens`, `number_of_components`, `number_of_connected_components`, `number_of_text_styles`, `number_of_colors`, `number_of_spacing_tokens` — seven aggregates computed over other collections, on the record that names the project.

**Two fields are called status and they are different types.** `status` is a string, `active` or `archived`. `workflow_status` is an object. One record, two statuses, and the second exists only on organisation projects.

**A parameter takes a magic word or an identifier in the same slot.** `workspace` is documented "it can be `personal` or the id of organization" — so a workspace named by a real id and the absence of one share a field, and an organisation whose id happened to be the string `personal` could not be addressed.

Also pinned: `created` and `updated` are Unix integers under names with no suffix saying so; `organization` and `workflow_status` are present only on organisation projects, so the record's shape depends on where it lives; and `limit` defaults to 30 over a documented range of 1 to 100.

## Sources

- [Zeplin `GET /v1/projects`](https://docs.zeplin.dev/reference/getprojects) — the project fields and the `limit` / `offset` / `workspace` / `status` parameters.
- Live: `api.zeplin.dev`, struck 2026-09-14 with no credential, a wrong bearer token, a token with no scheme, an unrouted path, a wrong method, and the root.

## Modelling limits

- **No description is published.** Zeplin serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Listing projects needs a real token; the records here are the reference's own field list with values of the documented types, and every case reading them is marked documentation-only.
- **`linked_styleguide`, `organization` and `workflow_status` are not served.** All three are objects the reference describes only by name, and the last two exist only on organisation projects — which is the finding, and is recorded above rather than guessed at.
- **`workspace` is not modelled.** Its value is either a magic word or an organisation id, and serving one of those two meanings would pick a side the reference does not.
- **One route of many.** The project listing. Screens, components, styleguides, notifications, webhooks and the organisation surface are the rest.
- **Nothing is mapped in detection.** Zeplin is reached through `@zeplin/cli` or a plain bearer call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
