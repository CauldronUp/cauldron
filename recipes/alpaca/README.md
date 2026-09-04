# alpaca

Emulates the alpaca API (v2), for local development and tests.

**19 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An Alpaca order is not a fill. An order for ten shares can sit at partially_filled indefinitely -- three filled at one price, four at another, three never -- and that's a real position with real money spent, so code that waits for status == "filled" waits forever while holding shares it never accounted for. Worse, the order listing shows only open orders by default, so the common sequence is: submit an order, it fills, list your orders, see nothing, and conclude the order never existed -- nothing errored, the list was just answering a different question.

Every number in the API is a string (qty: "10", filled_avg_price: "185.42"), which is a live bug in any language that compares strings lexically -- "10" > "9" is false. filled_avg_price is null until something fills and then is only the average of what has filled so far, so it changes under a caller who read it once and stored it. And legs is null on an ordinary order but an array on a bracket order, so code that iterates it works on the orders someone tested with and throws on the rest.

This models the paper trading API only -- no order here ever progresses past the state it was created with, and nothing about margin, buying power, or market data is covered.

## Sources

- Documentation: https://docs.alpaca.markets/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve alpaca     # run it
cauldron verify alpaca -v # check every claim
```
