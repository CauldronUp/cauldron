# aptible

Emulates the Aptible accounts API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Struck live against `api.aptible.com` on 2026-09-13 with no credential and with a deliberately invalid one. Record shapes come from Aptible's own generated Go client, which is open.

## What this Recipe found

**A request with no credential at all answers 200.**

```
GET /accounts            200 {"_embedded":{"accounts":[]},"total_count":0,
                              "per_page":40,"current_page":1,
                              "_links":{"self":{"href":"…"}}}
Authorization: Bearer x  400 {"code":400,"error":"invalid_token","message":"Invalid Token"}
```

Not 401. The collection is scoped to whoever is asking; an unauthenticated caller owns nothing; the answer is a successful, well-formed, entirely empty page.

So a client that checks `response.ok` and reads `total_count` concludes the user has no environments — and a monitoring check written that way reports healthy for a token that was never sent.

A *wrong* token is refused, and refused with **400**. The two credential states that both mean "you are not signed in" answer a success and a bad-request, and 401 never appears at all.

**The root document is public too.** `GET /` answers a HAL index naming every collection — accounts, apps, databases, stacks, active_plans, database_images, external_aws_accounts — to a caller who has presented nothing.

**It is HAL, and it does not link its pages.** The envelope carries `total_count`, `per_page` and `current_page`, and `_links` holds exactly one key: `self`. No `next`, no `prev`, no `last`.

Hypermedia's whole argument is that a client follows links instead of building URLs. The one place this API could have made that true is the one place it did not.

**Every record carries `type` and `_type`.** Two discriminators, one of them underscored, side by side on the same object, saying the same word.

**And the field names are the engineering team's.** From Aptible's own generated client, on the public account record:

| | |
| --- | --- |
| `sweetness_stack` | `gentlemanjerry_endpoint` |
| `sweetness_stack_version` | `gentlemanjerry_certificate` |
| | `gentlemanjerry_docker_name` |
| | `gentlemanjerry_instance_id` |

Six fields named after two internal projects. A caller reading an account has to know that "sweetness" is the stack and "gentlemanjerry" is the log drain, and neither word appears anywhere in the product.

**`ca_private_key` is a property of the account model.** It is the one field in the whole struct marked `omitempty` — every other field, all thirty of them, is required by the generated constructor. So the model describing an environment carries a key named for a private key, and it is the only key the type treats as optional.

**A name and a handle are both the name.** `handle` is a plain string and `name` is nullable, so the human label that may be absent sits beside the one that may not, with nothing saying which a client should show.

**And four different container counts do not add up.** `container_count`, `app_container_count`, `database_container_count`, `total_app_count` and `total_database_count` all sit on one record, and the first is not the sum of the second and third. Nothing on the record says what the difference is.

## Sources

- Live: `api.aptible.com`, struck 2026-09-13.
- [`model_account.go`](https://github.com/aptible/aptible-api-go/blob/main/aptibleapi/model_account.go) — the account record, in Aptible's own generated client.

## Modelling limits

- **One route.** Accounts, which Aptible calls environments. Apps, databases, stacks, endpoints, services, operations, backups, plans and the rest are named in the root document and each want their own evidence.
- **The sandbox cannot scope a collection to its caller.** Aptible answers an anonymous request with the records that caller owns, which is none. Here the anonymous case is pinned to the empty fixture, which is the shape the live API returns; seeding records and asking anonymously would serve them, which Aptible would not.
- **It is served as `application/json`, not `application/hal+json`.** Aptible labels both its successes and its failures with the HAL media type. A Recipe cannot declare a success content type, so the shape is faithful and the label is not.
- **`_links` is not served.** The envelope's one link is `self`, whose value is the request URL; nothing in this Recipe reconstructs it, and its absence is the finding rather than a gap.
- **No `spec:`.** Aptible publishes no OpenAPI document at any address this Recipe could find. The record shape comes from the generated client, which is the artefact a description would have produced.
- **Nothing is mapped in detection.** Aptible is driven by its own CLI and by Terraform; the Go and Ruby clients are generated for internal use and a project holding one is rare enough that a name match would mostly be wrong. Checked 2026-09-13.
