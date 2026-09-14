# clicksend

Emulates the ClickSend SMS delivery-receipt listing for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-13.**

Read from the API documentation at [`developers.clicksend.com`](https://developers.clicksend.com/docs/messaging/sms) and from ClickSend's own PHP SDK, and struck live on 2026-09-13 with no credential, with a wrong basic credential, with an empty password, with a bearer token, on a path that does not exist, and with methods the paths do not take.

## What this Recipe found

**A Basic-auth API answers 401 with no `WWW-Authenticate` header.** RFC 9110 says a server generating a 401 **MUST** send one. The live response carries `Content-Type`, `Content-Length`, three rate-limit headers, `x-content-type-options`, `x-frame-options`, a Cloudflare ray id and a TLS version — and no challenge. A browser never prompts, and a client cannot discover from the response that Basic is what the server wants.

**The rate-limit family is spelled two ways in one response.**

```
ratelimit-reset: 47
x-ratelimit-limit: 6000
x-ratelimit-remaining: 5993
```

The reset uses the unprefixed spelling from the IETF's rate-limit draft; the other two use the de-facto `x-` one. A client reading `x-ratelimit-reset` finds nothing, and a client reading `ratelimit-limit` finds nothing either.

**The body repeats the status and names the failure a third time.** Every response carries `http_code`, `response_code` and `response_msg`, so a 401 is:

```json
{"http_code":401,"response_code":"UNAUTHORIZED","response_msg":"Authorization failed.","data":null}
```

The number, a word and a sentence, for one fact — and `data`, present and null.

**A wrong method is reported as a wrong endpoint, and told to try again.** `PUT /v3/sms/history` and `GET /v3/cauldron-nope` both answer:

```json
{"http_code":404,"response_code":"NOT_FOUND","response_msg":"Invalid endpoint. Please try again.","data":null}
```

Trying again produces the same answer.

**The success message is "Here are your data."** — a sentence, in the field a client switches on, on every successful read.

**The paginator's records live under `data.data`.** The envelope's `data` is a Laravel paginator — `total`, `per_page`, `current_page`, `last_page`, `next_page_url`, `prev_page_url`, `from`, `to` — whose own `data` is the array of records. So the path to a receipt is `data.data[0]`, and the two `data` keys, one level apart, mean different things.

**A delivery status is the string `"201"`.** On a receipt:

```json
{"status_code": "201", "status_text": "Success: Message received on handset."}
```

An HTTP created-code, as a string, for a handset delivery. Beside it `error_code` and `error_text` are null, and so are `custom_string` and `digits`.

**And the message id is an uppercase UUID.** `D6D16B28-46AC-484A-AB0A-A08CD08EF75C`. RFC 4122 says generators should output lowercase and that comparison should be case-insensitive — which code doing a string compare against a lowercased copy will not be.

**The official PHP SDK types every response as a string.** In `lib/Api/SMSApi.php`, eighteen methods are annotated `@return string`; the eighteen `WithHttpInfo` variants return `array`, being that same string with the status and headers beside it. No model, no parsing — the caller gets the JSON text and does the rest. And `smsHistoryGet` declares `@param int $page` with the default `'1'` and `@param int $limit` with the default `'10'`: string defaults on integer parameters.

Also pinned: `timestamp_send` and `timestamp` are two Unix seconds that are equal in the documentation's own example; `message_type` is `"sms"` on the SMS receipts endpoint; and the documented `q` parameter takes a query language of its own — "Custom query Example: `from:{number},status_code:201`." — inside a query parameter's value.

## Sources

- [ClickSend SMS documentation](https://developers.clicksend.com/docs/messaging/sms) — the base URL, Basic authentication, and the receipt record.
- [View MMS History](https://developers.clicksend.com/docs/messaging/mms/other/view-mms-history) — the envelope and the paginator, quoted in full there.
- [`ClickSend/clicksend-php`](https://github.com/ClickSend/clicksend-php) `lib/Api/SMSApi.php` — the `@return string` annotations and the `'1'` / `'10'` defaults.
- Live: `rest.clicksend.com`, struck 2026-09-13 with no credential, a wrong Basic credential, an empty password, a bearer token, an unrouted path, and wrong methods on two paths.

## Modelling limits

- **No description is published.** ClickSend serves no OpenAPI document, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Reading receipts needs a real account; the record here is the documentation's own example plus a second record of the same shape, and every case reading them is marked documentation-only.
- **`from` and `to` are not served.** The documented paginator carries them beside `total` and `per_page`; this format has no name for a pair of one-based row offsets, and inventing constants for them would be worse than leaving them out.
- **The page size is the documented example's, not a measured default.** The envelope in the documentation shows `per_page: 15`; the PHP SDK's own default for `limit` is the string `'10'`. This Recipe serves 15 and records the discrepancy above.
- **One route of many.** The SMS receipts listing. SMS history, send, price, receipts-by-message, MMS, voice, fax, email, contacts and the account surface are the rest.
- **Nothing is mapped in detection.** ClickSend is reached through `clicksend/clicksend-php`, `clicksend` on npm, or a plain Basic call, and none of them resolves to this host through a dependency file. Checked 2026-09-13.
