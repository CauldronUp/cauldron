# clevercloud

Emulates the Clever Cloud organisations API for local development and tests.

**8 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Clever Cloud serves without a credential at [`api.clever-cloud.com/v2/openapi.json`](https://api.clever-cloud.com/v2/openapi.json), and struck live on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**`id` is the error code.**

```
(no header)      401 {"id":2001,"message":"Not connected","type":"error"}
Bearer notreal   403 {"id":1504,"message":"This token's scope is invalid.","type":"error"}
/cauldron-nope   404 <!DOCTYPE html>… "Page not found"
```

Every resource in this API has an `id` too, and theirs is a string. So `body.id` is the identifier of a thing on success and a numeric error code on failure — same field, different type, with nothing but the status line to tell them apart.

**`type` is always `"error"`.** A discriminator whose only observed value repeats what the status code already said, on an envelope with no other category field in it.

**A token that does not exist is reported as a scope problem.** The wrong-credential answer is 403 "This token's scope is invalid" — not that the token is unknown, but that its permissions are wrong. A caller reading that goes looking at the scopes on a credential that was never valid at all.

**And the message carries a typographic apostrophe.** U+2019, not U+0027, in a machine-readable error string. A client matching on the text has to get the right one of two characters that look identical in most fonts.

**An unrouted path answers HTML.** A full `<!DOCTYPE html>` page titled "Page not found", from an API whose three other answers are JSON — a fourth shape, and the only one a `.json()` call throws on.

**217 paths and no security block.** The published document declares no `security`, no `securitySchemes`, and nothing anywhere saying how to authenticate — the same silence [kestra](../kestra)'s 194-path document keeps.

**And no operation declares a failure.** `/organisations` lists exactly one response, 200, on an API that answers 401 to every anonymous request. `cauldron drift` names all three:

```
clevercloud   not backed: no operation the Recipe routes to answers 401, which not_connected declares
              not backed: no operation the Recipe routes to answers 403, which invalid_scope declares
              not backed: no operation the Recipe routes to answers 404, which unknown_route declares
```

**`VAT` is the only upper-case field name on the record**, sitting beside `vatState`, which is camelCase. Two spellings of the same word, adjacent.

**An organisation carries a postal address and a phone number.** `address`, `city`, `zipcode`, `country`, `billingEmail`, `customerFullName` and `emergencyNumber` are all on the object a listing returns.

**And nothing on it is required.** Eighteen properties, no `required` array, `id` included.

## Sources

- Live: `api.clever-cloud.com`, struck 2026-09-13.
- [`openapi.json`](https://api.clever-cloud.com/v2/openapi.json) — served without a credential, 217 paths.

## Modelling limits

- **One route.** Listing organisations. Applications, add-ons, deployments, instances, logs, payments, the GitHub integration and the OAuth surface are 216 more paths.
- **The identifier prefix is this Recipe's.** `OrganisationView.id` is declared only as a string, with no pattern and no example in the document; `orga_` is the shape Clever Cloud's own console uses and is not a claim the description makes.
- **No credential scheme is modelled from the description**, because there is none in it. The bearer header here is what the live API reads, not what the document says.
- **Nothing is mapped in detection.** Clever Cloud is driven by its own `clever-tools` CLI, which a project holds as a dev dependency at most and usually installs globally — so a dependency list rarely names it, and never names which of the API's two auth schemes a project uses. Checked 2026-09-13.
