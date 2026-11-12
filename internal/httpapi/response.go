package httpapi

/*
PACKAGE: httpapi

FILE: response.go

PURPOSE:

Provide reusable HTTP request and response utilities.

This file centralises:
1. JSON request decoding.
2. JSON response encoding.
3. Structured error responses.
4. HTTP method error responses.
5. Application error to HTTP status mapping.

DESIGN PRINCIPLES:

- Keep response formatting consistent.
- Avoid duplicating JSON encoding logic.
- Reject malformed and unexpected request data.
- Use appropriate HTTP status codes.
- Keep HTTP-specific behaviour outside the service layer.

This file must not:
- Implement task business logic.
- Access TaskManager directly.
- Perform storage operations.
- Register HTTP routes.
*/

/*
SECTION 1: IMPORTS

Import the packages needed for:

- encoding/json:
  JSON encoding and decoding.

- errors:
  Comparing application errors using errors.Is().

- io:
  Detecting EOF and handling request bodies.

- net/http:
  HTTP request and response interfaces,
  status codes and request body size limits.

- strings:
  Processing HTTP header values if needed.

Only import packages that are actually used.
*/
import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

	"task_manager/internal/task"
)

/*
SECTION 2: REQUEST JSON DECODING

Function:
decodeJSON

Suggested signature:

decodeJSON(

	w http.ResponseWriter,
	r *http.Request,
	target any,

) bool

PURPOSE:

Decode the JSON request body into the supplied destination.

Example:

Incoming request:

	{
	    "id": 1,
	    "title": "Learn Go"
	}

Destination:
A pointer to the appropriate request DTO.

REQUIREMENTS:

A. REQUEST BODY LIFECYCLE

1. Ensure the request body is closed after processing.
2. Do not attempt to read the request body multiple times.
3. Handle empty request bodies as invalid input.

B. REQUEST SIZE LIMIT

1. Limit the maximum request body size.
2. Use http.MaxBytesReader.
3. A reasonable initial limit is 1 MB.
4. Reject requests that exceed the limit.

Reason:
Prevent unnecessarily large request bodies from consuming
excessive server resources.

C. CONTENT TYPE

1. Inspect the Content-Type request header.
2. Require application/json if that is your API contract.
3. Reject unsupported media types.
4. Return HTTP 415 Unsupported Media Type.

Consider whether to support optional JSON parameters
such as charset=utf-8.

D. JSON DECODING

1. Create a JSON decoder.
2. Configure it to reject unknown JSON fields.
3. Decode the request body into the supplied destination.
4. Handle malformed JSON.
5. Handle type mismatches.
6. Handle empty input.

E. SINGLE JSON VALUE VALIDATION

After successfully decoding the first JSON value:

1. Attempt to decode one additional value.
2. The second decode must return io.EOF.
3. If another JSON value is present, reject the request.
4. If malformed trailing data exists, reject the request.

Example of invalid input:

	{
	    "id": 1,
	    "title": "Learn Go"
	}

	{
	    "id": 2,
	    "title": "Learn Rust"
	}

The decoder must reject this request.

F. ERROR HANDLING

If decoding fails:

- Return HTTP 400 for malformed JSON.
- Return HTTP 400 for unknown fields.
- Return HTTP 400 for multiple JSON values.
- Return HTTP 400 for empty request bodies.
- Return HTTP 413 for oversized request bodies.
- Return HTTP 415 for unsupported content types.

Use a clear and consistent error response.

G. RETURN VALUE

Return true only when:

- The content type is acceptable.
- The request body is within the size limit.
- Exactly one valid JSON value is decoded.

Return false otherwise.

IMPORTANT:

Handlers must stop processing immediately when
decodeJSON returns false.

Do not allow a handler to continue creating or updating
a task after request decoding has failed.
*/
func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	defer r.Body.Close()
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, errors.New("content type must be application/json"))
		return false
	}

	// Limit the request body size to prevent abuse.
	const maxBodySize = 1 << 20 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Create a JSON decoder.
	decoder := json.NewDecoder(r.Body)

	// Configure the decoder to reject unknown JSON fields.
	decoder.DisallowUnknownFields()

	// Decode the request body into the supplied destination.
	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			http.Error(w, "empty request body", http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "invalid character") {
			http.Error(w, "malformed JSON", http.StatusBadRequest)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
		}
		return false
	}

	// Attempt to decode one additional value.
	if err := decoder.Decode(new(struct{})); !errors.Is(err, io.EOF) {
		http.Error(w, "multiple JSON values", http.StatusBadRequest)
		return false
	}

	return true
}

/*
SECTION 3: JSON RESPONSE WRITING

Function:
writeJSON

Suggested signature:

writeJSON(
    w http.ResponseWriter,
    status int,
    value any,
)

PURPOSE:

Convert a Go value into JSON and send it to the client.

Examples of values:

- A single Task.
- A slice of Tasks.
- A map containing response information.
- A structured error response.

REQUIREMENTS:

1. Set the Content-Type response header to application/json.
2. Set the appropriate HTTP status code.
3. Encode the supplied value as JSON.
4. Ensure the response uses valid JSON formatting.
5. Handle encoding failures appropriately.

IMPORTANT HTTP BEHAVIOUR:

HTTP response headers are generally committed when:

- WriteHeader is called.
- Response body writing begins.

After committing a response, you cannot reliably replace
it with a different HTTP status code.

Therefore:

- Do not attempt to write a second error response
  after JSON encoding has already failed.
- Handle unexpected encoding errors through logging
  or an appropriate server-side error mechanism.
- Avoid hiding operational failures without any visibility.

SPECIAL CASE:

HTTP 204 No Content must not use writeJSON.

A 204 response should have no response body.

The handler should write the status directly.
*/

func writeJSON(w http.ResponseWriter, status int, value any) {
	// Set the Content-Type header to application/json.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode the supplied value as JSON.
	if err := json.NewEncoder(w).Encode(value); err != nil {
		// Log the error and return a 500 Internal Server Error.
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

/*
SECTION 4: STANDARDISED ERROR RESPONSE

Function:
writeError

Suggested signature:

writeError(
    w http.ResponseWriter,
    status int,
    err error,
)

PURPOSE:

Return a consistent JSON error representation.

Example:

HTTP/1.1 404 Not Found
Content-Type: application/json

{
    "error": "task not found"
}

REQUIREMENTS:

1. Accept the desired HTTP status code.
2. Accept the application or request error.
3. Construct a consistent JSON error response.
4. Delegate JSON serialization to writeJSON.

DESIGN:

All endpoints should return errors in the same format.

Do not expose:
- Filesystem paths.
- Internal stack traces.
- Database connection strings.
- Sensitive implementation details.

The public error message should explain the problem
without revealing unnecessary internal information.

IMPORTANT:

Do not automatically assume that every internal error
message is safe to return to the client.

Consider separating public error messages from
internal diagnostic information in future improvements.
*/

func writeError(w http.ResponseWriter, status int, err error) {
	// Construct a consistent JSON error response.
	errorResponse := map[string]string{
		"error": err.Error(),
	}
	writeJSON(w, status, errorResponse)
}

/*
SECTION 5: METHOD NOT ALLOWED RESPONSE

Function:
methodNotAllowed

Suggested signature:

methodNotAllowed(

	w http.ResponseWriter,
	methods ...string,

)

PURPOSE:

Respond when a client uses an unsupported HTTP method
for an otherwise recognised endpoint.

Example:

POST /health

Expected response:

HTTP/1.1 405 Method Not Allowed
Allow: GET

	{
	    "error": "method not allowed"
	}

REQUIREMENTS:

1. Accept a list of supported HTTP methods.
2. Construct the Allow response header.
3. Include all methods supported by the endpoint.
4. Return HTTP 405.
5. Return the standardised JSON error response.

HEADER REQUIREMENT:

Use the HTTP Allow header.

Example values:

- GET
- GET, POST
- GET, PATCH, DELETE

Ensure that the listed methods accurately reflect
the endpoint's actual supported methods.

Do not include methods that the endpoint does not support.
*/
func methodNotAllowed(w http.ResponseWriter, methods ...string) {
	// Set the Allow header to the supported methods.
	w.Header().Set("Allow", strings.Join(methods, ", "))

	// Return a standardised error response.
	writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
}

/*
SECTION 6: APPLICATION ERROR TO HTTP STATUS MAPPING

Function:
statusFromError

Suggested signature:

statusFromError(err error) int

PURPOSE:

Translate service-layer application errors into
appropriate HTTP status codes.

The service layer must not know about HTTP.

For example:

The service returns:
- Task not found.
- Duplicate task ID.
- Invalid task input.
- Invalid state transition.
- Unexpected storage failure.

The HTTP layer determines how these errors
should be communicated to clients.

SUGGESTED MAPPING:

A. TASK NOT FOUND

	HTTP 404 Not Found.

B. DUPLICATE TASK ID

	HTTP 409 Conflict.

C. INVALID TASK INPUT

	HTTP 400 Bad Request.

D. INVALID STATE TRANSITION

	HTTP 409 Conflict.

E. UNEXPECTED STORAGE FAILURE

	HTTP 500 Internal Server Error.

REQUIREMENTS:

1. Use errors.Is() or errors.As().
2. Identify known service errors.
3. Return the corresponding HTTP status.
4. Return HTTP 500 for unexpected errors.

IMPORTANT:

Do not compare errors using their string messages.

Do not make the service layer return HTTP status codes.

The mapping belongs in the HTTP API layer.

DESIGN CONSIDERATION:

This function should only determine the status code.

The handler or a separate response helper should
determine the appropriate public error message.

This keeps error classification separate from
response formatting.
*/
var (
	ErrNotFound   = task.ErrNotFound
	ErrConflict   = task.ErrConflict
	ErrBadRequest = task.ErrBadRequest
)

func statusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, task.ErrInvalidTransition):
		return http.StatusConflict
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

/*
SECTION 7: ERROR HANDLING CONTRACT

All handlers should follow a consistent pattern.

1. Receive a request.
2. Validate request structure.
3. Call the service.
4. If the service returns an error:
   - Identify the application error.
   - Determine the HTTP status.
   - Return a standardised error response.
5. If successful:
   - Return the appropriate success status.
   - Return the response data when applicable.

EXAMPLES:

Invalid JSON:
400 Bad Request

Task not found:
404 Not Found

Duplicate ID:
409 Conflict

Unsupported method:
405 Method Not Allowed

Unsupported content type:
415 Unsupported Media Type

Oversized request:
413 Content Too Large

Unexpected storage failure:
500 Internal Server Error

Successful retrieval:
200 OK

Successful creation:
201 Created

Successful deletion:
204 No Content
*/

/*
SECTION 8: TESTING REQUIREMENTS

Create response_test.go or include the tests
in your existing server_test.go.

A. JSON DECODING TESTS

- Valid JSON is decoded successfully.
- Invalid JSON is rejected.
- Empty request body is rejected.
- Unknown fields are rejected.
- Multiple JSON values are rejected.
- Oversized request body is rejected.
- Unsupported Content-Type is rejected.
- Incorrect JSON field types are rejected.

B. JSON RESPONSE TESTS

- Correct Content-Type is returned.
- Correct status code is returned.
- Response body contains valid JSON.
- Structs are encoded correctly.
- Slices are encoded correctly.

C. ERROR RESPONSE TESTS

- Error response uses the standard JSON format.
- Correct status code is returned.
- Error message is present.
- Internal implementation details are not exposed.

D. METHOD NOT ALLOWED TESTS

- HTTP 405 is returned.
- Allow header is present.
- Allow header contains the supported methods.

E. ERROR MAPPING TESTS

- Not found maps to 404.
- Duplicate ID maps to 409.
- Invalid input maps to 400.
- Invalid state transition maps to 409.
- Unexpected errors map to 500.

Use httptest.NewRecorder() to inspect responses.
Use table-driven tests where appropriate.
*/

/*
DEFINITION OF DONE

[ ] JSON decoding works.
[ ] Request size is limited.
[ ] Unknown fields are rejected.
[ ] Multiple JSON values are rejected.
[ ] Content-Type is validated.
[ ] JSON responses are consistent.
[ ] Error responses are consistent.
[ ] Method-not-allowed responses include Allow.
[ ] Application errors map to HTTP status codes.
[ ] HTTP 204 responses contain no body.
[ ] All response utility tests pass.
*/
