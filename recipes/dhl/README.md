# DHL

Emulates the DHL API (v2), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A DHL event timestamp carries no time zone at all -- "2026-02-01T09:15:00" with the offset sitting in a separate sibling field, GMTOffset, that most code never reads -- and because a parcel moves between countries, the offset differs between events on the very same shipment, so the error from reading the timestamp alone isn't even constant across one shipment's history. A tracking number that matches nothing is also a 404 with an RFC 7807 problem document rather than a 200 with an empty array, and inside that document status is a string ("404") while the actual HTTP status is a number -- comparing the two without converting fails silently.

Events carry a stable two-letter typeCode alongside a description that's prose translated to whatever language the request asked for, so matching on the description works until someone sets a language header. A piece within a multi-piece shipment also has its own events that aren't summarized by the shipment's overall status -- one box can be delivered while another is still in transit, and the shipment-level field describes neither exactly.

This models DHL Express's MyDHL tracking shape only; rating, shipment creation and label generation aren't covered, and Cauldron checks only the API key half of the credential, accepting any secret.

## Sources

- Documentation: https://developer.dhl.com/api-reference/shipment-tracking
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dhl     # run it
cauldron verify dhl -v # check every claim
```
