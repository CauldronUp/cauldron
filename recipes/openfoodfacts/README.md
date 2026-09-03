# Open Food Facts

Emulates the Open Food Facts API (v2), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

Food Facts's, where a product that is not there is a 200 with
`status: 0`. The HTTP status says it worked and the truth is a number in the
body, so `response.ok` is true and there is **no `product` key at all** --
`body.product.product_name` throws a TypeError on the answer a barcode scan is
most likely to get. And `status: 0` has more than one cause: `product not found`
and `no code or invalid code` share the number, so a client branching on it
cannot tell a barcode that does not exist from one it should not have sent.
**And the code you sent is not always the code that comes back**: asking for
`0000000000000`, thirteen zeros, answers `"code": "00000000"` -- eight -- with no
field anywhere recording what was asked for.

Then the nutrients. `energy` is **kilojoules**: the same object carries
`energy: 2252` and `energy-kcal: 539`, and the unqualified name is the kJ figure,
so reading `energy` and labelling it calories is out by a factor of 4.18. Every
nutriment is repeated -- `fat`, `fat_100g`, `fat_value` and `fat_unit` are four
keys for one number and its unit, and energy has twelve. The key names mix
separators inside one key: `added-sugars_100g`. `brands` is a comma-joined string
rather than a list, crowd-sourced, and one of the three on the reference product
is not a brand. And `nutriments_estimated` is a parallel object of the same
nutrients, computed rather than declared, keyed identically, with nothing inside
either saying which is which.

## Sources

- Documentation: https://openfoodfacts.github.io/openfoodfacts-server/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openfoodfacts     # run it
cauldron verify openfoodfacts -v # check every claim
```
