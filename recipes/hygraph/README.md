# hygraph

Emulates the Hygraph Content API for local development and tests.

**8 conformance cases, 1 checked against the live API on 2026-09-13.**

Read from Hygraph's own API reference, and struck live against `api-eu-central-1.hygraph.com` on 2026-09-13 with a project identifier that is malformed and with one that is well formed and does not exist.

## What this Recipe found

**Getting the project id wrong two ways gives two content types.**

```
/v2/cauldron-nope/master        400 text/plain        {"errors":[{"message":"project identifier \"cauldron-nope\" is invalid"}],"data":null}
/v2/<a valid-shaped id>/master  404 application/json  {"errors":[{"message":"project or environment not found"}],"data":null}
```

Same envelope, same GraphQL shape, two statuses — and the JSON body on the 400 is labelled `text/plain`. The mislabelled one is the one a caller reaches by mistyping the id, which is the more common of the two mistakes.

**Both carry `"data": null` beside `errors`.** GraphQL's specification says that if an error is encountered before execution begins, the `data` entry should not be present. Neither of these reached a resolver, and both have it, present and null.

**And neither status is GraphQL's.** A GraphQL endpoint answers 200 for a query that fails, with the reason in `errors`. These answer 400 and 404 with the reason in `errors`. One endpoint produces failures at three statuses in one envelope, so a client has to read `errors` whatever the status line says.

**A GET is answered.** Both were struck with no query, no body and no `Content-Type` — and the endpoint replied with its error rather than refusing the method.

**The second message merges two things.** "project **or** environment not found" — the path carries both, and the failure does not say which of the two was wrong.

**`publishedAt` is non-null on every content entry.** From the system-field table:

| field | type | description |
| --- | --- | --- |
| `publishedAt` | `DateTime!` | Timestamp for when the content entry was published |

The exclamation mark is GraphQL for "never null", on a timestamp recording an event that may not have happened. An entry in the draft stage has not been published, and the schema says the field is always there.

**Every entry carries three references to people** — `createdBy`, `updatedBy` and `publishedBy`, each a `User`. And the `User` model is documented as having the system fields "but not `createdBy`, `updatedBy`, and `publishedBy`", so the model that exists to record who did things is the one thing in the schema that does not record who did things to it.

**A token is a kind of user.** `kind: UserKind!` is `MEMBER`, `PAT` or `PUBLIC`. "Who published this" can answer with a personal access token, or with the public, and a client rendering an author's name has to handle both.

**And every record can query itself in other stages.** `documentInStages` is a field on the entry that returns the same entry as it exists elsewhere, so one object is a handle to all of its own versions.

## Sources

- Live: `api-eu-central-1.hygraph.com`, struck 2026-09-13.
- [System fields](https://hygraph.com/docs/api-reference/schema/system-fields) — the fields every model gets, and the `User` model that does not get all of them.
- [API Reference](https://hygraph.com/docs/api-reference) — the endpoint shape and the regional hosts.

## Modelling limits

- **One collection.** The `User` model, which Hygraph ships in every project and documents itself. Every other type in a project's schema is defined by the customer, so there is no shape this Recipe could claim for them.
- **The 404 variant is recorded, not served.** A well-formed but unknown project id answers `404 application/json`; a malformed one answers `400 text/plain`. Cauldron declares one failure per route, and this Recipe serves the mislabelled 400 because that is the one a typo produces.
- **`createdBy`, `updatedBy`, `publishedBy` and `documentInStages` are not modelled.** Each returns a `User` or a list of entries and would need its own resolver; they are quoted above because their existence is the finding.
- **One region.** Hygraph serves each project from a regional host, and a project id says nothing about which. This Recipe models one endpoint shape; the others differ only in host.
- **No `spec:`.** The Content API's schema is generated per project from the models a customer defines, so there is no single description to fetch.
- **Nothing is mapped in detection.** A project using Hygraph holds a generic GraphQL client pointed at a per-project regional URL, so the dependency names the protocol and never the endpoint. Checked 2026-09-13.
