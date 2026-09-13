# shortpixel

Emulates the ShortPixel image reducer API for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from ShortPixel's published API reference at [`shortpixel.com/api-docs`](https://shortpixel.com/api-docs) and from its own PHP client, and struck live against `api.shortpixel.com` on 2026-09-13 with no key, with a wrong key, and on a path that does not exist.

## What this Recipe found

**Every refusal is an HTTP 200.** Struck live, headers and all:

```
HTTP/1.1 200 OK
content-type: text/html; charset=UTF-8

{
    "Status": {
        "Code": -401,
        "Message": "Missing API key. You need an API key to use this service."
    }
}
```

Two hundred OK for a request that was refused, a JSON body labelled `text/html`, and the real status in the body as a negative number. A client checking `response.ok` sees success; a client checking the content type before parsing sees a web page.

**The negative numbers are HTTP status codes with the meanings left behind.** `-401` is an invalid key, `-402` a wrong one, `-403` an exhausted quota, and `-404` "the maximum number of URLs in the optimization queue reached" — so Forbidden's number is a billing state and Not Found's is back-pressure. Beside them sit `-102` through `-117`, `-201` through `-207` and `-301` through `-306`, which resemble nothing at all. One number space, half borrowed and half invented, and no way to tell which half a code came from without the table.

**And the documented sentence for `-401` is not the sentence `-401` sends.** The reference gives "Invalid API key. Please check that the API key is the one provided to you." Live, an absent key answers "Missing API key. You need an API key to use this service." One code, two sentences, and the published one is not the one you get.

**Every size field spells "Lossless" wrong and every URL field spells it right.** From the reference, on one record:

```
LosslessURL       LoselessSize
WebPLosslessURL   WebPLoselessSize
AVIFLosslessURL   AVIFLoselessSize
```

Three families, and in all three the URL is correct and the size has lost an `s` and an `l`. A typo that consistent was copied rather than slipped, and a client cannot write one rule that spells both halves of the pair.

**The AVIF lossy size carries WebP's P.** The reference lists `WebPLossySize` and then `AVIFPLossySize` — an `AVIFP` that is not a format, sitting next to `AVIFLossyURL` and `AVIFLoselessSize`, which are not spelled that way.

**"NA" is the absent value.** A URL that was not requested and a size that was not produced both come back as the two-character string `NA`, so a size field holds a number or a word, and a client comparing it to zero has to check which it got.

**A record names the machine that served it.** `Server` is "The server in our network that did the file processing. Used for debugging purposes." — an internal hostname on the response, with a load balancer cookie (`ShortPixel-LB2=ShortPixel_spfe02`) beside it on an API with no session.

## Sources

- [`shortpixel.com/api-docs`](https://shortpixel.com/api-docs) — the published reference: the field table and the full list of codes.
- [`short-pixel-optimizer/shortpixel-php`](https://github.com/short-pixel-optimizer/shortpixel-php) — `lib/ShortPixel/Result.php`, whose own test for failure is that the number is negative.
- Live: `api.shortpixel.com`, struck 2026-09-13 with no key and with a wrong one, on `reducer.php` and `post-reducer.php`.

## Modelling limits

- **Two routes.** The reducer and the upload endpoint beside it. `reducer-sync.php` returns the optimized bytes with an `image/…` type rather than JSON, and `api-status.php` reads the key from the query string rather than the body, so neither fits one auth scheme with these two.
- **The HTML 404 is not served.** An unrouted path answers a full styled error page — `<title>404 - ShortPixel</title>`, Google Fonts preconnects and all — to an API client. This Recipe records that rather than emulating a browser page.
- **The success fixture is documentation-derived.** Optimizing an image needs a real key and real quota; the record here is the reference's own field list with plausible values, and every case that reads it is marked documentation-only.
- **Nothing is mapped in detection.** ShortPixel is reached through `shortpixel/shortpixel-php` on Packagist or its WordPress plugin, and neither name resolves to this host through a dependency file. Checked 2026-09-13.
