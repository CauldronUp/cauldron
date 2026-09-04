# Google Calendar

Emulates the Google Calendar API (v3), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A listing returns the recurring series, not its occurrences. The master event's start date is the day the series began, so a query for "this week" hands back an event dated years ago with an RRULE attached, and rendering it directly puts a daily standup on one day, in the past, exactly once -- expansion is opt-in, and this is what the API does when nobody opts in. Per the Recipe's own header, this may be the most confidently misread API in common use.

A few more defaults catch people the same way. An all-day event has no time at all -- `start.date` instead of `start.dateTime`, and the two never coexist -- so code that always reads `start.dateTime` finds nothing on every holiday and leave day. An all-day event's end date is exclusive, so a one-day event on the 18th ends on the 19th, and naive duration math is off by one in the direction that looks plausible. And a cancelled occurrence still shows up in the response, almost empty -- no summary, no start, no end -- so a loop that reads `.summary` off every item throws on the one meeting somebody cancelled.

A sync token also arrives only on the last page of a listing; every earlier page carries a page token and no sync token, so code that grabs whichever token shows up first stores the wrong one and the next incremental sync either replays everything or fails outright.

## Sources

- Documentation: https://developers.google.com/calendar/api/v3/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googlecalendar     # run it
cauldron verify googlecalendar -v # check every claim
```
