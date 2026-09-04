# Jira Service Management

Emulates the Jira Service Management API (servicedeskapi), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An SLA's `elapsedTime` is not elapsed time, it's service time. It only counts minutes during working hours if `withinCalendarHours` is true, and it stops entirely while `paused` is true -- so a ticket raised at 5pm Friday against a nine-to-five calendar can show an `elapsedTime` near zero on Monday morning, three days after it was actually opened. Every dashboard's obvious `now - startTime` computation quietly disagrees with the number the SLA actually breaches on, and the two diverge by exactly the hours nobody was working, which is most of them. On top of that, the `sla` field isn't even in a plain GET response unless explicitly expanded -- a request comes back with a status, a reporter, and timestamps, and nothing at all about whether it's about to breach.

There is also no single field that answers "is this breaching." A request can have zero or more SLAs, each with zero or more completed cycles and at most one ongoing cycle, so answering that question means walking a list where any entry might not have a cycle running at all. And `breachTime` gets populated even on cycles that never breached -- documented as the time it *would have* breached, on a ticket that was actually answered on time, which is a timestamp under an alarming name describing something that didn't happen.

## Sources

- Documentation: https://developer.atlassian.com/cloud/jira/service-desk/rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve jiraservicemanagement     # run it
cauldron verify jiraservicemanagement -v # check every claim
```
