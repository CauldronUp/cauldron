# affinity

Emulates the Affinity CRM API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Affinity's reference at `api-docs.affinity.co` and struck live against `api.affinity.co` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Three sentences, all plain text, all declared `text/html`.**

```
GET /persons          (no credential)      401  text/html   Unauthorized API Key.
GET /persons          Basic :notreal       401  text/html   Invalid or missing API key
GET /cauldron-nope    Basic :notreal       404  text/html   Unknown API endpoint
```

No JSON anywhere. No object, no field, no envelope — three bare sentences under a header claiming markup that is not there.

This collection now has four providers getting one header wrong, in four distinct ways:

| | Header says | Body is |
|---|---|---|
| [statsig](../statsig) | *nothing* | text |
| [wrike](../wrike) | `text/plain` | JSON |
| [maxio](../maxio) | `application/json` | text |
| **affinity** | `text/html` | text |

Only Maxio's is actively dangerous — a client trusting the header parses and throws. Affinity's is the one most likely to be **rendered into a page** by anything that trusts a content type.

**The three sentences are punctuated differently.** "Unauthorized API Key." has a full stop; the other two do not. Two capitalise "Key" and one does not. Three failures written at three different times, and prose that is not consistent with itself is the only way to tell them apart.

**"Invalid or missing" is the answer to sending a key.** The fifth provider here to describe the wrong failure — and, unusually, it names *both* cases, so it is also the least wrong of the five.

**The credential is the Basic password with an empty username.** `:key`, not `key:`. A credential helper filling in the username sends the key where Affinity does not look, and no refusal mentions which half.

**The primary address is also in the list.** Counting addresses by summing `emails` and adding `primary_email` counts one twice, and the list gives no clue which entry is primary.

**A person can have no address at all** — `primary_email: null` and `emails: []` on the same record. A CRM's whole job is not dropping people, and a client keying on either field drops this one.

**A person's kind is an unlabelled integer.** `type: 0` and `type: 1`, with no enumeration on the record, so what a person *is* has to be looked up in prose.

## Modelling limits

- **One route.** Persons. Organizations, opportunities, lists, list entries, field values, notes and the whole relationship-intelligence surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Affinity publishes a rendered reference site.
