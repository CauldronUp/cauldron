# acuity

Emulates the Acuity Scheduling API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Acuity's reference at `developers.acuityscheduling.com` and struck live against `acuityscheduling.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Three fields, one fact, three renderings.**

```
401 {"status_code":401,"message":"Unauthorized","error":"unauthorized"}
```

`status_code` is the number. `message` is the reason phrase, capitalised. `error` is the reason phrase, lower-cased. **Every one is derivable from the status line**, and the body carries nothing the status did not already say.

It is the most compact example of a pattern several providers here have. [beehiiv](../beehiiv) sends the status three times and adds a real `errors` array beside it. [terraformcloud](../terraformcloud) sends it twice and stops. Acuity sends it three times and stops — so a client that parses the body learns the status, and a client that reads the status line learns the same thing sooner.

**The two phrase fields disagree about case.** `Unauthorized` and `unauthorized`, in one object, both meaning the reason phrase, and neither field name says which is the one to branch on.

**An unknown path is a 401 too**, so the route is judged after the credential.

**The start is a timestamp and the end is a clock time.** `datetime` is a full offset timestamp and `endTime` is `"10:45"` — a wall clock with no date and no zone. Computing a duration means combining a date from one field with a time from another, and **an appointment crossing midnight cannot be represented at all**.

**The offset has no colon.** `-0400` rather than `-04:00` — ISO 8601 basic format, not RFC 3339, so a strict parser refuses it. The same shape [dlocal](../dlocal) has.

**A cancelled appointment is in the listing** with no default filter, so a client counting bookings counts it. (The spelling is the American one, which is worth knowing before grepping for it.)

**`max` is the page size.** Not `limit` — and `max` usually means the largest permitted value rather than how many to send.

## Modelling limits

- **One route.** Appointments. Availability, appointment types, calendars, forms, products and the whole blocked-time surface each want their own evidence.
- **Nothing is mapped in detection.** Acuity is a Squarespace product now and publishes no client of its own on any registry.
- **No `spec:`.** Acuity publishes a rendered reference site.
