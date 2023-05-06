# Developer Notes

## How to version APIs

See https://stackoverflow.com/questions/389169/best-practices-for-api-versioning

## Encoding, Marshaling and Serialization, huh?

In Go, "encoding", "marshaling", and "serialization" are often used interchangeably to refer to the process of converting a data structure into a format that can be stored or transmitted and later reconstructed. However, there are some subtle differences in their usage.

Encoding refers specifically to the process of converting data into a specific format, such as JSON, XML, or binary.

Marshaling is a term often used in the context of Go and is used to describe the process of converting Go data structures into a format that can be stored or transmitted.

Serialization is a more general term that is used to describe the process of converting data structures into a format that can be stored or transmitted, and can refer to any language or platform.

In general, all three terms refer to the same process, but the specific term used can depend on the context and the programming language being used.


### Why is `json.Decoder` preferred over `json.Unmarshal` to read the HTTP request body?

When using `json.Decoder`, the entire payload is not decoded in one action. Instead, `json.Decoder` reads the JSON data from the input stream (e.g., an HTTP request body) and produces a stream of tokens. Each token represents a JSON value, such as a string, number, object, or array.

The `json.Decoder` provides a Token() method that can be used to retrieve the next token from the input stream. This allows you to process the JSON data token by token, instead of having to read the entire payload into memory and then decode it all at once.

This can be particularly useful when dealing with large JSON payloads or in situations where you don't know the exact structure of the JSON data ahead of time. By processing the JSON data as a stream of tokens, you can extract only the data you need, and avoid consuming unnecessary memory or processing time.

`json.Unmarshal` reads the entire payload into memory before decoding.

When working with a message queue, you may still use `json.Decoder` to decode JSON messages received from the queue. However, in this case, you would typically need to wrap the message body in an `io.Reader` interface so that the `json.Decoder` can read the JSON data in a streaming fashion. This is because messages on a queue can be of arbitrary length, and the entire message may not be available in memory at once.

## Design Patterns

1. **Triage errors**

	It is beneficial to the caller if errors are triaged to make it more sensible in cases where third party calls are made and multiple forms of errors or difficult to decipher errors are returned.

2. **Consistency by helpers**

	Create helpers to ensure the same endpoint is being used so that if the app becomes more complex it can be changed in one place. This decision does not have to be made pre-emptively but when you see a pattern emerging move the
	repeated functionality away from the primary responsiblity of the called function.

3. **Unexpected vs. expected errors**

	Panic for developer errors (*unexpected errors*) and return for *expected errors*

4. **Do not trust input**

	The input should be considered to be in its raw format and sanitized or converted to the expected format.

5. **Validation is for business logic**

	Error handling is for application logic. So keep these separate. Typically, you would perform the necessary error handling first, and then proceed to peform business logic validation.
## Security

1. Information hiding in returned errors
	- Public APIs must not expose implementation details
2. Set a maximum request body size to limit DoS attacks
