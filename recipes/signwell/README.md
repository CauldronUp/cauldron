# signwell

Emulates the SignWell credentials API for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document SignWell embeds in each documentation page, and struck live against `www.signwell.com` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The message a client reads first is the one that cannot tell the two failures apart.**

```
(no key)     401 {"message":"Missing or invalid authorization key",
                  "meta":{"error":"missing_authorization_key_error",
                          "message":"Missing authorization key",
                          "messages":["Missing authorization key"]}}
wrong key    401 {"message":"Missing or invalid authorization key",
                  "meta":{"error":"api_key_unauthorized_error",
                          "message":"Not valid authorization token",
                          "messages":["Not valid authorization token"]}}
```

`body.message` is byte-identical on both — "Missing **or** invalid" — and the sentence that says which it is sits two levels down in `meta.message`. The top-level field, the one a client reaches for, is the one written to cover both cases.

**And the real message is in the body three times.** `meta.error` carries it as a code, `meta.message` as a string, and `meta.messages` as an array holding that same string and nothing else. Three shapes, one sentence, one level below a fourth copy that disagrees with all of them.

**Both codes end in `_error`, in a field called `error`.** `missing_authorization_key_error` and `api_key_unauthorized_error`, under `meta.error`.

**An unrouted path answers the marketing site.** `www.signwell.com` is the API host and the website, so `/api/v1/cauldron-nope` answers a full `<!DOCTYPE html>` page titled "SignWell" — a Turbo-powered web app, complete with viewport meta tags, where a JSON client expected an error.

**`account` and `workspace` are the same object.** In the reference's own example they carry the same `id`, the same `name`, the same `plan_tier` and the same four booleans. The response contains it twice under two names, and nothing says which a caller should read or when they might differ.

**Four identifiers, and the record's own is none of the obvious ones.** The top-level `id` is neither the user nor the account: `user.id`, `account.id` and `workspace.id` are three more, two of them equal. A caller asking "who am I" gets four answers.

**And a person has a name and a first name and no last name.** `user.name` is "Jerry Smith" and `user.first_name` is "Jerry", with nothing carrying "Smith" — so splitting the full name is the only way to get the other half.

## Sources

- Live: `www.signwell.com`, struck 2026-09-13.
- [Get credentials](https://developers.signwell.com/reference/getme) — the page, and the OpenAPI document embedded in it.

## Modelling limits

- **One route.** The credentials endpoint. Documents, templates, bulk sends, webhooks, API applications and the completed-PDF download each want their own evidence.
- **The unrouted 404 serves the page's opening, not the page.** Live it is the whole SignWell web application; here it is the first four lines, which is enough to reproduce what breaks — an HTML body where a JSON client expected an error.
- **The example's own names are kept, the addresses are not.** SignWell's reference uses generated fixture names and an `@pollich.example` address; the names are its and the addresses here are `example.com`, because an address in a shipped fixture should not resolve anywhere.
- **No `spec:`.** SignWell embeds a complete OpenAPI document per endpoint inside that endpoint's page, the same arrangement [mailtrap](../mailtrap), [lemlist](../lemlist) and [fillout](../fillout) use. There is nothing single to fingerprint.
- **Nothing is mapped in detection.** SignWell publishes no first-party client library, and a project using it holds an embed script or a webhook handler rather than a client of this API. Checked 2026-09-13.
