# AWS DynamoDB

Emulates the AWS DynamoDB API (20120810), for local development and tests.

**18 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Every DynamoDB attribute is wrapped in a type descriptor -- a name isn't "Ada", it's {"S": "Ada"}, and a number is {"N": "42"}, the digits as a string -- so code written against a document-style interface and pointed at the low-level API reads undefined everywhere, while code that unwraps by hand does arithmetic on strings. A query that matches nothing is still a success with Count 0 and no Items key at all, so iterating response.Items on an empty result throws rather than looping zero times.

Count and ScannedCount are also different numbers, and the gap between them is the real cost: a filter is applied after the read, so a query that scans a thousand items and returns one is billed for the thousand, and that gap is the only warning anyone gets. And a paginated result isn't finished when a page comes back short -- it's finished only when LastEvaluatedKey is absent, and a page can legitimately be empty while still carrying a key to continue from, so stopping on an empty page loses the rest of the table.

Cauldron doesn't verify SigV4 signatures (the shape is checked, not the signature) and doesn't evaluate key conditions or filter expressions -- a query returns everything in the table rather than what the expression would actually select.

## Sources

- Documentation: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dynamodb     # run it
cauldron verify dynamodb -v # check every claim
```
