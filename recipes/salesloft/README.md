# salesloft

Emulates the Salesloft v2 users API for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-11.**

Read from Salesloft's own API reference and struck live against `api.salesloft.com` on 2026-09-11 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The failure points at a hostname that is not a link.**

```
(no header)      401 {"error":"No Bearer Token attached to request, please see documentation at developers.salesloft.com"}
Bearer notreal   401 {"error":"Invalid Bearer token"}
```

`developers.salesloft.com` has no scheme, so it is not a URL. It will not open from a terminal, it does not linkify in a log viewer, and a client that wants to surface it has to guess `https://` in front. It is also the only machine-readable thing in the response — there is no code, no type, and no field beside `error`. The second failure drops even that.

**And the 401 names a scheme it will not send a challenge for.** No `WWW-Authenticate` on either failure, from an API whose own error message contains the words "Bearer Token". RFC 9110 makes that header mandatory on a 401; here the scheme is in the prose instead, where no client reads it.

That is the sixth challenge-header fault in this collection:

| Recipe | What it does with `WWW-Authenticate` |
| --- | --- |
| [cloudamqp](../cloudamqp) | omits it from the caller who needs it |
| [airship](../airship) | fills it with whichever scheme the caller already tried |
| [cloudsmith](../cloudsmith) | names a scheme that does not work |
| [fastmail](../fastmail) | points at an RFC 9728 metadata document |
| [zitadel](../zitadel) | puts a sentence where the scheme belongs |
| salesloft | names the scheme in the body and sends no header at all |

**Every path answers the credential failure, including the ones that do not exist.** `/v2/cauldron-nope` with no token gets the missing-token sentence, and with a wrong token gets the invalid-token sentence. A mistyped path is indistinguishable from a credential problem, and a client that retries its credentials on a 401 retries forever against a typo.

**A "List users" endpoint documents its success as one user.** The reference for `GET /v2/users` publishes a 200 schema of exactly four fields — `id`, `guid`, `name`, `email` — with no array around them and no envelope. It is the same schema the single-user endpoint `/v2/me` publishes, character for character. So a listing and a lookup are described identically, and neither description is a shape a listing can have.

**The envelope appears only inside a parameter description.** `include_paging_counts` is documented as "Whether to include `total_pages` and `total_count` in the **metadata**" — naming a `metadata` object that appears nowhere in the response schema. This Recipe does not invent one: sending a field the provider may not send is the more dangerous mistake, because code written against it works locally and fails in production.

**And the listing does not paginate by default.**

> `page` — The current page to fetch users from. **Defaults to returning all users**
>
> `per_page` — How many users to show per page in the range [1, 100]. Defaults to 25. **Results are only paginated if the page parameter is defined**

So a caller who sets a page size and no page number gets every user on the team in one response, and the parameter they set does nothing. `include_paging_counts` defaults to false, so the totals needed to page are off unless asked for; `visible_only` defaults to true, which hides the users the caller cannot act on and every deactivated user. The default behaviour of "list users" is an unpaginated, uncounted, filtered subset.

**Two identifiers, three spellings.** Every user carries an integer `id` and a `guid`, and the note on the second says "New endpoints will explicitly accept this over id" — without naming which endpoints those are. The filters take both, as `ids` (`integer[]`) and `guid` (`string[]`), and the manager filter spells the same concept a third way as `manager_user_guid`.

**Two filters are typed as strings so they can hold magic words.** `group_id` is `string[]` with the example `[1, _is_null]` — a number and a sentinel in one array — and `role_id` is `string[]` with `[1, User]`, taking an identifier or a role's name in the same position. `last_login` is a string whose documented example value is `never`.

**The declared credential is OAuth2 only.** The reference lists one security scheme: an authorization-code flow against `accounts.salesloft.com` with the scope `team:read`. The live error talks about bearer tokens, and Salesloft issues API keys as well. Neither appears in the description.

## Sources

- Live: `api.salesloft.com`, struck 2026-09-11.
- [List users](https://developers.salesloft.com/docs/api/users-index/) — the 21 filters and the four-field schema.
- [Fetch current user](https://developers.salesloft.com/docs/api/me-index/) — the same four-field schema.

## Modelling limits

- **Two routes.** Listing users and the current user. Accounts, people, cadences, emails, calls, meetings, opportunities, tasks and the rest of a reference with over a hundred resources each want their own evidence.
- **`per_page` works here and is inert live.** Salesloft paginates only when `page` is present, and Cauldron has no way to say that a declared page-size parameter applies only alongside a position. The sandbox honours `per_page` on its own where Salesloft ignores it. One provider is a note; if a second turns up, that is the field name.
- **No `metadata` envelope.** The reference names one in a parameter description and never describes it, so nothing is served beside the array. A Recipe that guessed at it would teach client code to read a field that may not exist.
- **No `spec:`.** The reference is a Docusaurus site built from an OpenAPI document at a build-time path (`api-spec/v2/v2.json`) that the bundle names and the site answers 404 for. The schemas here were read from the rendered reference.
- **Nothing is mapped in detection.** Salesloft publishes no first-party client library; integrations go through its App Directory or generic HTTP clients holding an OAuth token, so no package name says this API is in use. Checked 2026-09-11.
