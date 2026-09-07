# terraformcloud

Emulates the HCP Terraform API for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against HashiCorp's reference at `developer.hashicorp.com` and struck live against `app.terraform.io` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Every failure is two fields and neither is a sentence.**

```
GET /api/v2/organizations   (no credential)      401 {"errors":[{"status":"401","title":"unauthorized"}]}
GET /api/v2/organizations   Bearer not-a-token   401 {"errors":[{"status":"401","title":"unauthorized"}]}
GET /api/v2/cauldron-nope   Bearer not-a-token   404 {"errors":[{"status":"404","title":"not found"}]}
```

`status` restates the HTTP status as a string. `title` is the status phrase, lower-cased. There is no `detail`, no `code`, no `id`, no `source`. **The whole body of a failure is the status line, twice** — a client that reads it learns nothing the status did not already say.

JSON:API defines `title` as a summary "that SHOULD NOT change from occurrence to occurrence of the problem" and `detail` as the explanation of *this* occurrence. Only the invariant half is here, which is the half that could have been a constant in the client.

**A missing credential and a wrong one are byte-identical**, so the two failures with different fixes are one response — and there is no field that could have separated them even in principle.

**The 404 has the same shape, and that is the good half.** One parser reads every failure this API produces, on every route, under `application/vnd.api+json` throughout. Uniformity is worth something. HCP Terraform buys it by never saying anything.

**Field names are hyphenated.** `cost-estimation-enabled`, `created-at`. JSON:API recommends it and almost nothing follows the recommendation, so the accessor here is `org["cost-estimation-enabled"]` rather than a property in any language.

**Paging lives under `meta` rather than `links`.** JSON:API defines `links.next`; HCP Terraform sends a `meta.pagination` object with hyphenated keys, so a client written against the specification finds no cursor and treats every page as the last.

**The identifier is the name somebody typed.** `id` and `attributes.name` are the same string, so renaming an organization changes every URL and every stored reference to it.

**The product was renamed and the host was not.** Terraform Cloud became HCP Terraform in 2024; the API is still at `app.terraform.io` and the paths still say `/api/v2` — the same split [maxio](../maxio) has with `chargify.com`.

## Modelling limits

- **One route.** Organizations. Workspaces, runs, plans, applies, state versions, variables, policy sets and the whole VCS-integration surface each want their own evidence.
- **Nothing is mapped in detection.** The Terraform CLI is not a client of this API in the sense detection means, and no library on npm, Packagist or the Go module proxy calls `app.terraform.io` under an obvious name.
- **No `spec:`.** HashiCorp publishes a rendered reference site.
