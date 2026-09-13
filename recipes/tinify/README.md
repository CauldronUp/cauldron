# tinify

Emulates the Tinify compression endpoint for local development and tests.

**7 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from Tinify's published HTTP reference at [`tinify.com/developers/reference`](https://tinify.com/developers/reference), and struck live against `api.tinify.com` on 2026-09-13 with no credential, with a wrong one, with a malformed body, and on a path that does not exist.

## What this Recipe found

**One endpoint, two documented responses, and they share no keys.** The reference prints `POST /shrink` twice on the same page. Uploading a file:

```
HTTP/1.1 201 Created
{"input": {"size": 207565, "type": "image/jpeg"}}
```

Pointing at a URL:

```
HTTP/1.1 201 Created
{"output": {"size": 30734, "type": "image/png"}}
```

Same method, same path, same status — and a client reading `body.output.size` gets nothing from the first, while one reading `body.input.size` gets nothing from the second. Whichever half of the page a reader implements against is the half that works.

**The result is not in the body.** The address of the compressed image is the `Location` header. A client that parses JSON and ignores headers has the sizes and no way to fetch the thing it just paid to make.

**And the bill is a header too.** `Compression-Count` carries the account's running monthly usage on every response, so the number the plan is priced on is only ever visible to code that reads response headers. The download adds `Image-Width` and `Image-Height` the same way, because the body there is the image and there is nowhere else to put them.

**The body is validated before the credential.** Struck live, with no `Authorization` header at all:

```
{"x":1}                     400  Request is invalid: request body requires key 'source';
                                 request body has unknown key 'x'.
{"source":{"url":"..."}}    401  Credentials are invalid.
```

An anonymous caller can enumerate the request schema one key at a time — the validator names the required key and names the unknown one — and only a well-formed request gets as far as being refused.

**"Credentials are invalid" is what you are told for sending none.** That 401 answered a request carrying no credential at all, and a wrong key gets the same sentence.

**A real path with the wrong method does not exist.** `GET /shrink` and `GET /cauldron-nope` answer the identical body: `{"error":"Not found","message":"This endpoint does not exist."}`. The endpoint that exists reports that it does not, and a typo in the verb is indistinguishable from a typo in the path.

**And the `error` field is a reason phrase, mis-cased.** `"Unauthorized"`, `"Bad request"`, `"Not found"` — HTTP's own titles, two of the three in a capitalisation the standard does not use, in the field a client would switch on.

## Sources

- [`tinify.com/developers/reference`](https://tinify.com/developers/reference) — the published HTTP reference: both worked `POST /shrink` exchanges, the header list, and the Basic Auth instructions.
- Live: `api.tinify.com`, struck 2026-09-13 with no credential, a wrong credential, a malformed body, a wrong method, and an unrouted path.

## Modelling limits

- **One endpoint.** `POST /shrink`, in both of its documented forms. The resize, convert, preserve and store options and the `GET /output/{id}` download are the rest — the download answers image bytes, which this Recipe does not serve.
- **The two forms are told apart by `Content-Type`.** The reference distinguishes them by what is in the body — binary or JSON — and this Recipe routes on the header that accompanies each, which is what the page's own examples send.
- **The credential is checked before the body here, and live it is the other way round.** That ordering is the finding above; a declarative Recipe authenticates first, so the sandbox refuses an anonymous malformed request where the real API describes it. The live transcript is recorded rather than emulated. The authenticated 400 is exact.
- **`Location` and `Compression-Count` are constants.** Live, one is the address of the object just created and the other is an account's running total.
- **No `spec:`.** Tinify documents this API as a reference page of worked HTTP exchanges; there is no machine-readable description at any address this Recipe could find, so there is nothing for `cauldron drift` to record.
- **Nothing is mapped in detection.** Tinify is reached through the `tinify` package on npm, PyPI, RubyGems or Packagist, and none of those names resolve to this host through a dependency file. Checked 2026-09-13.
