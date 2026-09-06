# authorizenet

Emulates the Authorize.Net API (v1), for local development and tests.

**6 conformance cases, 4 checked against the live API on 2026-09-06.**

The live cases were struck against `apitest.authorize.net` with credentials nobody issued. No account was created and no transaction was attempted — `authenticateTestRequest` with invalid credentials is the one call that asks the gateway to do nothing at all.

## What this Recipe found

**A rejected credential comes back HTTP 200.** Recorded live:

```
POST https://apitest.authorize.net/xml/v1/request.api
{"authenticateTestRequest":{"merchantAuthentication":
  {"name":"cauldron-not-a-real-login","transactionKey":"cauldronNotAKey0"}}}

HTTP 200
{"messages":{"resultCode":"Error","message":[
  {"code":"E00007","text":"User authentication failed due to invalid authentication values."}]}}
```

`response.ok` is true. `status < 400` is true. `raise_for_status()` does nothing. The request was refused. Every failure on this API arrives that way, so the only thing separating success from failure is `messages.resultCode`, a string reading `"Ok"` or `"Error"`.

**The body begins with a UTF-8 byte-order mark, and this Recipe serves it.** The first three bytes are `EF BB BF`, before the `{`. Python's `json.loads` on the decoded text raises `Unexpected UTF-8 BOM (decode using utf-8-sig)`; Go's `encoding/json` refuses it too. A client that works against a fake without one throws on its first real response, so `responses.bom` exists because of this API — and Cauldron's own conformance runner had to learn to strip a mark before decoding, exactly as a surviving client does, the moment this Recipe started sending one.

**One URL for every operation.** There is a single path, `/xml/v1/request.api`, and what you are asking for is the top-level key of the body — `authenticateTestRequest`, `createTransactionRequest`, and so on. An HTTP-level view of an integration shows one endpoint doing everything, per-route rate limits are meaningless, and a proxy cannot tell a refund from a lookup without parsing the body.

**The credential's path depends on the operation.** `merchantAuthentication` is nested *beneath* whichever request you are making, so the secret lives at `authenticateTestRequest.merchantAuthentication.transactionKey` on one call and `createTransactionRequest.merchantAuthentication.transactionKey` on the next. There is no fixed location for it: it cannot be set once as a header, every call site carries it, it lands in anything that logs request bodies, and a scrubber trying to redact it has to know every operation name there is. For a payment gateway, that is a worse place for a secret to live than most.

**A message code carries its severity in its first character.** `I00001` is informational, `E00007` is an error, so `code` is a string to be parsed rather than a number to be compared — and `message` is an array even when it holds one entry.

## Modelling limits

- **`messages.resultCode` is not served.** It sits beside the message array rather than inside it, and this format expresses one envelope shape, not a scalar next to an array. The array is served, with its `code` and `text`.
- **One operation.** `authenticateTestRequest` is served because it is the one observable without an account. The rest of the gateway is a much larger surface and each part wants evidence rather than a shape copied from this.
- **The XML half is not served.** This API answers XML or JSON depending on the Content-Type, and the JSON is a transliteration of the XML — which is why the envelope looks the way it does.
