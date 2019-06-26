# Header Validation Toolkit

This toolkit is designed to ensure your headers are working as intended through a proxy. It's also designed as an acceptance test endpoint to ensure your expectations are met. Everyone loves acceptance tests.

## Expectations and Intentions

This toolkit has two core concepts: _intentions_ and _expectations_, both of which solve different problems.

### Exceptions

When working with complicated, non-centralized proxy configurations, it's often hard to ensure configuration expectations are met, especially in high-compliance scenarios. To ensure expectations are met, this toolkit has an endpoint to provide header validations.

| Endpoint | Use |
| --- |--- |
| `/expect` | Used for viewing the current headers. |
| `/expect/diff` | See the difference between your expect file and what is deployed. |

#### Expectation Header File

The tookit support an expectation file, which is a reference file for how headers should look after passing through your proxy. This file is designed to be based off your acceptance criteria tests, so whatever you want your target headers to look like, this file should reflect that. Here is the current file format:

```json
{
  "header-field-name": ["header-field-value-0"]
}
```

#### Expected Behaviour with Expectation Diffs

By default, the `header-reference.json` file only contains registered HTTP/1.1 header fields, but vendors love to put custom fields in their requests.

This is what an untracked (non-standard, custom) header looks like:

```json
{
	"name": "X-Amzn-Trace-Id",
	"have": {},
	"want": null
}
```

The reason this shows up is because there is no matching header in in our expectation header file, so there's no way to compare it against anything. If you want to start tracking the header, just add the header to the expectations file and restart the app, and it will look more normal:

```json
{
	"name": "X-Amzn-Trace-Id.[0]",
	"have": "Root=1-5cd1e3a3-45b120fc9d944ffc16de3c4c",
	"want": ""
}
```

If you're wondering about why header names generally look like `Header.[0]`, it's because of how Go parses the header frame. Go sets the header frames to be `map[string][]string`, which means it looks a bit like this:

```json
{
  "Header": [
    "value",
    "another-value"
  ]
}
```

According to RFC2616: 

> Multiple message-header fields with the same field-name MAY be present in a message if and only if the entire field-value for that header field is defined as a comma-separated list [i.e., #(values)]. It MUST be possible to combine the multiple header fields into one "field-name: field-value" pair, without changing the semantics of the message, by appending each subsequent field-value to the first, each separated by a comma. The order in which header fields with the same field-name are received is therefore significant to the interpretation of the combined field value, and thus a proxy MUST NOT change the order of these field values when a message is forwarded.

So what this means is that each header can use whatever separator it wants, but it has to use the same separator every single time. The problem is that in practice, no one can agree on a separator, which means valid separators can be `\s`, `,`, `-`, `;`, or just about any other character.

Go's header parsing implementation is RFC-compliant, but in reality it doesn't actually separate any of the strings, which means all your index values will be `0`. If Go started splitting the header values into separate string fields as the RFC originally intended, you would end up seeing something like this:

```json
[
    {
        "name": "X-Amzn-Trace-Id.[0]",
        "have": "Root=1-5cd1e3a3-45b120fc9d944ffc16de3c4c",
        "want": ""
    },
    {
        "name": "X-Amzn-Trace-Id.[1]",
        "have": "Root=1-5cd1e3a3-45b120f1523243fc16de3c4c",
        "want": ""
    }
]
```

So this is expected behaviour. If there is no diff between your expected headers and the diff, you'll see `HTTP/1.1 418 I'm a teapot`. Why? It was difficult to find a way to express when there is no diff, so being a teapot is so ridiculous that it can't possibly be a proxy that sets it.

### Intentions

While ensuring expectations are met, it's important to validate intentions are also met. This toolkit also provides a way to set headers and then generate a diff against reality.

| Endpoint | Use |
| --- |--- |
| `/intent` | Used for viewing the headers you intended to set. |
| `/intent/diff` | See the difference between your intent file and what is actually seen. |

#### Expected Behaviour for Intention Diffs

The behaviour should be identical to the expectations functionality.

## Use

1. `go build -v -o hs header-server/header-server.go`
1. `./hs -header-ref target-file.json`

Optionally,

1. `cf push`
