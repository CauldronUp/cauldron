# sns

Emulates the AWS SNS API (2010-03-31), for local development and tests.

**5 conformance cases, 1 checked against the live API on 2026-09-06.**

The live case was struck against `sns.us-east-1.amazonaws.com` with no account and no credential — being refused is the thing being recorded, so nothing needed to be created.

## What this Recipe found

**It answers XML, and it is the only AWS Recipe here that does.** SQS, DynamoDB and Secrets Manager all speak the newer JSON protocol — a POST to `/` with `X-Amz-Target` naming the operation. SNS is the older Query protocol: parameters in the query string, `Action=ListTopics`, and `text/xml` coming back. A client written against the SQS Recipe next door and pointed at SNS calls `.json()` on an XML document and throws.

**A missing credential is 403, not 401.** Recorded live:

```
GET /?Action=ListTopics&Version=2010-03-31
HTTP 403  text/xml
<ErrorResponse xmlns="http://sns.amazonaws.com/doc/2010-03-31/">
  <Error>
    <Type>Sender</Type>
    <Code>MissingAuthenticationToken</Code>
    <Message>Request is missing Authentication Token</Message>
  </Error>
  <RequestId>5ba6480b-2695-5107-b1c6-2b1c4235b89a</RequestId>
</ErrorResponse>
```

No 401, no `WWW-Authenticate`, and nothing for an HTTP client's built-in retry-with-credentials path to catch. Code branching on 401 to refresh a token never refreshes; code reading 403 as "this key lacks permission" goes hunting through IAM when the problem is that no key was sent.

**`Type` is the retry signal, and it is not the code.** Every failure carries `Type`, either `Sender` or `Receiver`. Sender means the request was wrong and resending fails again; Receiver means AWS broke and retrying is right. That is the distinction a backoff loop needs, and it sits beside the code rather than in it.

**A malformed signature is answered with a hash of your own header.** Also struck live — an `Authorization` header with no date returns `IncompleteSignature` whose message ends `Authorization=y6ryXS/RhT3/69GgeQOKwb+lkSz/foW3fACGahjN+Ig=`. That is a digest, not the secret, so it is not a key leak — but it is a value derived from a credential arriving in an error string that goes straight to logs, and worth knowing before someone pastes one into a ticket.

**Listing topics gives you ARNs and nothing else.** `Topics.member.N` is an array of Topic objects and a Topic has one field: the ARN. No name, no display name, no subscription count. Anything a UI wants needs a second call per topic, and there is no batch form — so a hundred topics is a hundred and one requests.

## Modelling limits

- **The success bodies are recorded XML, not generated.** The Query protocol's `Topics.member.1.TopicArn` flattening is a shape this format builds no records into, so the listing is a fixed document and seeding a different fixture does not change it.
- **Topics only, and only listing them.** Publish, subscribe, the confirmation handshake and delivery-failure behaviour are the interesting parts of SNS, and each wants evidence of its own.
- **The signature is not verified.** Reproducing SigV4 is not what this project is for; what is reproduced is the difference between sending nothing and sending something malformed, which is where both live observations came from.
