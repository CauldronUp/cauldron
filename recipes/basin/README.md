# basin

Emulates the Basin forms API for local development and tests.

**10 conformance cases, 2 checked against the live API on 2026-09-13.**

Read from the Swagger document Basin serves without a credential at [`swagger.yaml`](https://usebasin.com/api_docs/v1/swagger.yaml), and struck live against `usebasin.com` on 2026-09-13 with no credential and with a deliberately invalid one.

## What this Recipe found

**The two credential states answer two different shapes at two different statuses.**

```
(no key)    400 {"error":"No API key was provided or no user session exists."}
wrong key   401 {"error":"invalid_token","error_description":"The access token is invalid","state":"unauthorized"}
```

One is a sentence in `error`. The other is OAuth 2.0's `error` / `error_description` pair, where `error` holds a *code*. So `body.error` is prose on the first and a machine-readable token on the second, and a client that prints it shows a sentence or the word `invalid_token` depending on which failure it met.

**The missing-key message conflates two kinds of caller.** "No API key was provided **or no user session exists**" — a browser session and an API key, in one sentence, to a caller who has no browser. And it is a 400, so the status for a missing credential is the status for a malformed request.

**`state` is not what `state` means.** OAuth 2.0 defines `state` as a value the *client* supplies and the server echoes back, for CSRF protection. Here the server invents one and puts the word `unauthorized` in it — a third field saying what the status code and the error code have both already said.

**The listing declares no response body.** `GET /api/v1/forms` has exactly one response, `200`, and its whole declaration is `description: Success`. No content, no schema, no example. Twenty-one responses elsewhere in the document do declare a body; this is not one of them, so a generated client for the most obvious endpoint in the API returns nothing typed at all.

**The same collection is declared at two paths.** `/api/v1/forms` and `/api/v1/forms/` are separate entries in the document, as are `/api/v1/form_webhooks` and `/api/v1/form_webhooks/`. One resource, two paths, differing by a trailing slash.

**A form carries a third party's secret.** `turnstile_secret` is a property of the form record, beside `turnstile_site_key` — and a site key is public by design while a secret key is the half that must never leave the server. Both are on the object a GET returns.

**Three CAPTCHA vendors, three booleans.** `force_recaptcha`, `force_hcaptcha` and `force_turnstile`, with nothing saying what happens if more than one is true.

**Two list encodings on one record.** `notification_emails`, `notification_cc_emails` and `notification_bcc_emails` are comma-separated strings; `allowed_domains`, `blocked_domains` and `content_blacklist` are arrays. Both are lists, and which encoding you get depends on the field.

**And the record uses both vocabularies at once.** `whitelist_source_domains` and `content_blacklist` sit beside `allowed_domains` and `blocked_domains` — the old words and the new ones, on the same object, for overlapping jobs.

**`notification_emails` has a three-way protocol encoded in one string.** From its own description:

> Omit this field to leave the current recipients unchanged, send an empty string to remove all recipients, or send a list to replace them.

Absent, empty and populated are three different commands, in a field a form encoder cannot leave out.

## Sources

- Live: `usebasin.com`, struck 2026-09-13.
- [`swagger.yaml`](https://usebasin.com/api_docs/v1/swagger.yaml) — served without a credential.

## Modelling limits

- **One route.** Listing forms. Submissions, webhooks, mail templates, domains, projects and form views each want their own evidence.
- **The record comes from the single-form endpoint.** The listing declares no response body at all, so the 63-property shape here is the one `GET /api/v1/forms/{id}` publishes — which is the only place in the document that says what a form looks like.
- **Twenty-four of those properties are modelled.** The rest are styling, mail-template and redirect settings that say nothing the ones here do not.
- **`turnstile_secret` is a fixture value.** It is shaped like Cloudflare's test key and is not a credential for anything; the finding is that the field is on the record at all.
- **There is no page-size parameter.** `page` and `query` are the operation's whole parameter list, so a caller cannot ask for a different page size and nothing in the response says what the size is.
- **Nothing is mapped in detection.** Basin is reached by pointing an HTML form's `action` at it, which is markup rather than a dependency, and it publishes no client library. Checked 2026-09-13.
