# messagebird

Emulates the MessageBird message listing for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from MessageBird's published API reference and its own PHP client, and struck live against `rest.messagebird.com` on 2026-09-13 with no header, with a wrong access key, with a Bearer header, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A phone number is a JSON integer and a network code beside it is a string.** From the reference's own message object:

```json
"recipient": 31612345678,
"recipientCountryPrefix": 31,
"mccmnc": "20408",
"mcc": "204",
"mnc": "08"
```

`recipient` is documented `type: integer`. `mnc` is a string, and it has to be — `"08"` loses its leading zero the moment it becomes a number. An MSISDN loses its leading zero exactly the same way, and is a number anyway. One object, two identifiers, opposite decisions, for the same reason.

**And the network code is on the record three times.** `mccmnc` is `"20408"`, `mcc` is `"204"`, `mnc` is `"08"` — the concatenation and both of its halves, as three fields on one recipient.

**A field called SentCount counts the ones that have not been sent.** From the reference: "`totalSentCount` — The count of recipients that have the message pending (status `sent`, and `buffered`)." Pending, under a name that says sent. And the example has `totalCount: 1`, `totalSentCount: 1` and `totalDeliveredCount: 1` for a single recipient, so three counts describe one person and adding them gives three.

**The status and the reason disagree in the vendor's own example.** `"status": "sent"` with `"statusReason": "successfully delivered"`, on the same recipient object.

**`mclass` has four valid values and two documented ones.** "Indicated the message type. `1` is a normal message, `0` is a flash message. (0-3 are valid values)" — two of four named, the special case is the zero, and the sentence has a typo in its first word.

**`gateway` is a bare integer with no list.** "The SMS route that is used to send the message. This is for advanced users." The reference prints `240` in one example and `10` in another and explains neither.

**The credential failure names a query parameter.**

```json
{"errors":[{"code":2,"description":"Request not allowed (incorrect access_key)","parameter":"access_key"}]}
```

The credential travels in `Authorization: AccessKey <key>`, and the failure blames a parameter called `access_key` — and says "incorrect" to a request that presented nothing. A `Bearer` header gets the identical body, so the scheme word is not what is being read.

**The prose field is called `description`.** Not `message`, not `detail` — with an integer `code` beside it, and a `parameter` that is present and `null` when there is nothing to blame.

**`links` carries `first`, `previous`, `next` and `last`**, with `first` and `last` the same URL on a one-page listing; `count` and `totalCount` are two different counts in one envelope; a record carries `href`, its own address; `price.amount` is a floating point number; and a `PUT` answers 411 Length Required with an HTML page whose text reads "POST requests require a `Content-length` header".

## Sources

- [`developers.messagebird.com/api/sms-messaging/`](https://developers.messagebird.com/api/sms-messaging/) — the published reference: the message object, the recipients table, and the list envelope.
- [`messagebird/php-rest-api`](https://github.com/messagebird/php-rest-api) — `Message` and `Recipient`, whose comments carry the same field descriptions.
- Live: `rest.messagebird.com`, struck 2026-09-13 with no header, a wrong access key, a Bearer header, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** Listing messages. Sending, balance, contacts, groups, voice, conversations, verify and the lookup surface are the rest.
- **No `spec:`.** MessageBird documents this API as a reference page with worked examples; there is no OpenAPI description at any address this Recipe could find, so there is nothing for `cauldron drift` to record.
- **`links.first` and `links.last` are constants.** Live they are addresses computed from the offset; the forward and backward links are served from the request.
- **The 411 is not served.** A `PUT` never reaches the API — a load balancer answers an HTML page first. That is recorded rather than emulated.
- **The success fixture is reference-derived.** Listing messages needs a real access key; the records here are the reference's own example objects with their field set intact, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** MessageBird is reached through `messagebird` on npm, PyPI or Packagist, and none of those names resolve to this host through a dependency file. Checked 2026-09-13.
