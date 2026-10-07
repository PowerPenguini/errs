# errs
Small toolkit for clean error handling with REST API.

`WriteError` preserves the first typed error boundary in the cause chain. An
`ErrorList` at that boundary produces an `errors` array, including field names.
A typed `Error` wrapping a list keeps its own HTTP status, code, and message;
the list in its cause does not override that response. Ordinary wrappers and
joined errors are supported.

Untyped errors return HTTP 500 with `unknown_error` and a generic message.
Infrastructure diagnostics stay in the original error for server-side logging.
