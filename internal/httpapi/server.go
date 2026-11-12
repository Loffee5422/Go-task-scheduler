
package httpapi

// ============================================================
// 1. PACKAGE DECLARATION
// ============================================================
//
// `package httpapi` declares the package name.
//
// In Go, a package is a collection of related source files.
// Files in the same directory normally belong to the same package.
//
// Your project structure might look like:
//
// task_manager/
// ├── cmd/
// │   └── taskflow/
// │       └── main.go
// ├── internal/
// │   ├── service/
// │   │   └── manager.go
// │   └── httpapi/
// │       └── server.go
//
// The purpose of the httpapi package is to expose your
// existing application functionality through HTTP.
//
// IMPORTANT ARCHITECTURAL DISTINCTION:
//
// service package:
//     Contains business logic.
//     Example: CreateTask(), DeleteTask(), CompleteTask().
//
// httpapi package:
//     Handles HTTP requests and responses.
//     Example: receives POST /tasks and calls CreateTask().
//
// The HTTP layer should not independently implement
// business rules that already belong in the service layer.


// ============================================================
// 2. IMPORTS
// ============================================================

import (
    // errors provides functions for creating and manipulating errors.
    //
    // Here it is used to reject a nil TaskManager dependency.
    // A nil pointer means that a pointer does not reference
    // an actual object.
    "errors"

    // net/http is Go's standard HTTP package.
    //
    // It provides:
    // - HTTP server functionality.
    // - Request and response types.
    // - URL routing.
    // - HTTP methods.
    // - HTTP status codes.
    // - Request handlers.
    //
    // This is the primary package you need to understand
    // when building HTTP APIs in Go.
    "net/http"

    // strconv provides conversions between strings and numbers.
    //
    // HTTP URL parameters are strings.
    // Your TaskManager expects an integer ID.
    //
    // Example:
    // URL:       /tasks/123
    // Path value: "123" (string)
    // Required:   123 (int)
    //
    // Atoi means ASCII to integer.
    "strconv"

    // strings provides string manipulation functions.
    //
    // Here it is used to remove whitespace from query parameters.
    // Example: "  Pending  " becomes "Pending".
    "strings"

    // Imports your existing service layer.
    //
    // The HTTP API will call TaskManager methods rather than
    // implementing task management logic itself.
    "task_manager/internal/service"
)


// ============================================================
// 3. SERVER STRUCT
// ============================================================

// Server represents the HTTP API server.
//
// In Go, a struct is a user-defined data type that groups
// related fields together.
//
// This struct holds the dependencies required by the HTTP API.
//
// Notice that it does not contain a separate []Task collection.
// It delegates task management to the service layer.
//
// This avoids having two different copies of application state.
//
// Example:
//
//    manager := service.NewTaskManager(...)
//    server := Server{manager: manager}
//
// All HTTP handlers belonging to this Server will use
// the same TaskManager instance.
//
// `manager` is a pointer because we want to reference the
// existing TaskManager rather than create a copy.
//
// `*service.TaskManager` means:
// A pointer to a TaskManager type defined in the service package.
//
// The field is unexported because its name starts with
// a lowercase letter.
//
// Other packages cannot directly access server.manager.
//
// This is an example of encapsulation.
type Server struct {
    manager *service.TaskManager
}


// ============================================================
// 4. CONSTRUCTOR / DEPENDENCY INJECTION
// ============================================================

// NewServer creates and returns a new HTTP Server.
//
// This is a constructor-style function.
//
// Go does not have a dedicated constructor keyword.
// Developers commonly create functions named NewType()
// to initialize structs.
//
// PARAMETERS:
//
// manager *service.TaskManager
//
// The caller must provide a pointer to an existing TaskManager.
//
// RETURN VALUES:
//
// (*Server, error)
//
// Go functions can return multiple values.
//
// Here:
// - *Server: the successfully constructed server.
// - error: describes whether construction failed.
//
// Example usage:
//
//    server, err := httpapi.NewServer(manager)
//    if err != nil {
//        log.Fatal(err)
//    }
//
// WHY RETURN AN ERROR?
//
// A Server without a valid TaskManager cannot perform
// its primary responsibility.
//
// Rejecting invalid dependencies during construction
// prevents later nil-pointer panics.
//
// This is called fail-fast design.
func NewServer(manager *service.TaskManager) (*Server, error) {

    // Check whether the supplied pointer is nil.
    //
    // `nil` represents the absence of a value for a pointer,
    // map, slice, channel, function or interface.
    //
    // Without this check, a handler might eventually attempt:
    //
    //    s.manager.ListTasks()
    //
    // If manager is nil, this can cause a runtime panic.
    //
    // Returning an error allows the caller to handle the
    // problem before starting the HTTP server.
    if manager == nil {

        // errors.New creates a new error with the supplied message.
        //
        // An error is a built-in Go interface type.
        // Functions conventionally return errors as their
        // final return value.
        return nil, errors.New("task manager dependency cannot be nil")
    }

    // Create and return a Server instance.
    //
    // `&Server{...}` creates a struct value and returns
    // a pointer to that value.
    //
    // This is a struct literal.
    //
    // Field initialization:
    //
    // manager: manager
    //
    // Left side:
    //     The Server struct's manager field.
    //
    // Right side:
    //     The manager pointer received as the function argument.
    //
    // The same TaskManager instance is now referenced
    // by the HTTP Server.
    //
    // The second return value is nil because no error occurred.
    return &Server{
        manager: manager,
    }, nil
}


// ============================================================
// 5. ROUTE REGISTRATION
// ============================================================

// Handler registers all HTTP routes and returns an HTTP handler.
//
// This function configures the API's routing table.
//
// WHAT IS A ROUTE?
//
// A route maps an HTTP method and URL pattern to a function.
//
// For example:
//
// GET /tasks
//
// Means:
//
// When a client sends an HTTP GET request to /tasks,
// execute the function registered for that route.
//
// WHAT IS http.Handler?
//
// http.Handler is an interface defined by net/http.
//
// It represents something capable of handling an HTTP request.
//
// Its essential method is:
//
//    ServeHTTP(ResponseWriter, *Request)
//
// Any type implementing that method satisfies http.Handler.
//
// Returning an interface rather than a concrete type
// makes this function flexible and easy to test.
//
// WHY RETURN A HANDLER INSTEAD OF STARTING THE SERVER HERE?
//
// This separates route configuration from network startup.
//
// It allows you to test HTTP handlers using httptest
// without opening a real network port.
func (s *Server) Handler() http.Handler {

    // Create a new ServeMux.
    //
    // ServeMux is Go's standard HTTP request multiplexer.
    //
    // Think of it as a router.
    //
    // It examines incoming requests and determines which
    // handler function should receive each request.
    //
    // mux is a local variable.
    // := performs declaration and initialization.
    //
    // Each invocation creates a new routing configuration.
    mux := http.NewServeMux()

    // --------------------------------------------------------
    // HEALTH CHECK
    // --------------------------------------------------------

    // Register a handler for:
    //
    // GET /health
    //
    // A health endpoint is commonly used by:
    // - Deployment systems.
    // - Monitoring systems.
    // - Load balancers.
    // - Container orchestration platforms.
    //
    // It provides a simple indication that the API responds.
    //
    // HandleFunc accepts:
    // 1. A route pattern.
    // 2. A function that handles matching requests.
    //
    // `s.handleHealth` is a method value.
    //
    // It refers to the handleHealth method together with
    // the particular Server instance represented by s.
    //
    // The method will receive:
    // - ResponseWriter: where the response is written.
    // - Request: information about the incoming request.
    //
    // Go 1.22 introduced enhanced ServeMux routing patterns,
    // including HTTP method matching.
    mux.HandleFunc("GET /health", s.handleHealth)


    // --------------------------------------------------------
    // TASK COLLECTION ENDPOINTS
    // --------------------------------------------------------

    // GET /tasks
    //
    // Retrieves the collection of tasks.
    //
    // POST /tasks
    //
    // Creates a new task.
    //
    // Both routes point to the same handler function.
    // The handler determines what to do based on r.Method.
    //
    // This is possible because HTTP methods express
    // different operations on the same resource.
    //
    // Resource:
    //     /tasks
    //
    // Operations:
    //     GET  = retrieve
    //     POST = create
    //
    // This is a common REST API design pattern.
    mux.HandleFunc("GET /tasks", s.handleTasks)
    mux.HandleFunc("POST /tasks", s.handleTasks)


    // --------------------------------------------------------
    // INDIVIDUAL TASK ENDPOINTS
    // --------------------------------------------------------

    // These routes operate on one particular task.
    //
    // Examples:
    //
    // GET    /tasks/123
    // PATCH  /tasks/123
    // DELETE /tasks/123
    //
    // The {id} part is a wildcard.
    //
    // It captures one segment of the URL path.
    //
    // For example:
    //
    // Request URL: /tasks/123
    //
    // Captured value:
    //     id = "123"
    //
    // The handler can retrieve this value using:
    //
    //     r.PathValue("id")
    //
    // Notice that PathValue returns a string.
    // You must convert it into an integer before using it
    // with the TaskManager.
    //
    // GET:
    //     Retrieve an individual task.
    //
    // PATCH:
    //     Partially update an existing task.
    //
    // DELETE:
    //     Remove an existing task.
    mux.HandleFunc("GET /tasks/{id}", s.handleTask)
    mux.HandleFunc("PATCH /tasks/{id}", s.handleTask)
    mux.HandleFunc("DELETE /tasks/{id}", s.handleTask)


    // --------------------------------------------------------
    // TASK ACTION ENDPOINTS
    // --------------------------------------------------------

    // These endpoints represent specific actions on a task.
    //
    // POST /tasks/123/complete
    // POST /tasks/123/reopen
    //
    // The purpose of these routes is to expose business
    // operations that already exist in the service layer.
    //
    // Why POST?
    //
    // These operations change the state of a resource.
    // POST is commonly used for action-oriented endpoints.
    //
    // ServeMux selects the most specific matching pattern.
    //
    // Therefore:
    //
    // /tasks/{id}/complete
    //
    // is distinct from:
    //
    // /tasks/{id}
    //
    // These patterns do not conflict.
    mux.HandleFunc("POST /tasks/{id}/complete", s.handleComplete)
    mux.HandleFunc("POST /tasks/{id}/reopen", s.handleReopen)

    // Return the configured router.
    //
    // The caller can pass this handler to:
    //
    //    http.ListenAndServe(":8080", server.Handler())
    //
    // That starts a server listening on port 8080.
    //
    // This function itself does not start a network listener.
    return mux
}


// ============================================================
// 6. HEALTH CHECK HANDLER
// ============================================================

// handleHealth verifies that the API responds to requests.
//
// This is a method belonging to Server.
//
// RECEIVER:
//
// (s *Server)
//
// A receiver connects a method to a particular type.
//
// This means the method can access Server fields:
//
//     s.manager
//
// The receiver is a pointer, allowing the method to access
// the original Server instance.
//
// PARAMETERS:
//
// w http.ResponseWriter
//
// Represents the outgoing HTTP response.
//
// You use it to:
// - Set response headers.
// - Set the HTTP status code.
// - Write response data.
//
// r *http.Request
//
// Represents the incoming HTTP request.
//
// It contains:
// - HTTP method.
// - URL.
// - Headers.
// - Request body.
// - Query parameters.
// - Path values.
//
// IMPORTANT:
//
// The handler does not return an HTTP response object.
// It writes the response through w.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {

    // Check the HTTP method.
    //
    // r.Method contains the incoming method as a string.
    //
    // Example:
    //     "GET"
    //     "POST"
    //     "PATCH"
    //
    // http.MethodGet is a predefined constant equal to "GET".
    //
    // The route already specifies GET, so this check is
    // technically redundant when using this ServeMux setup.
    //
    // It is retained as an additional defensive check.
    if r.Method != http.MethodGet {

        // Respond with HTTP 405 Method Not Allowed.
        //
        // methodNotAllowed is a helper function defined
        // elsewhere in your httpapi package.
        //
        // It should set the appropriate response status
        // and potentially include an Allow header.
        methodNotAllowed(w, http.MethodGet)

        // Stop executing this handler.
        //
        // Without return, execution would continue and
        // potentially write another response.
        return
    }

    // Return a successful JSON response.
    //
    // HTTP status:
    //     200 OK
    //
    // JSON body:
    //     {"status":"ok"}
    //
    // map[string]string means:
    //
    //     Key type:   string
    //     Value type: string
    //
    // A Go map is a collection of key-value pairs.
    //
    // writeJSON is another helper function in your package.
    // It should:
    // 1. Set Content-Type to application/json.
    // 2. Write the HTTP status code.
    // 3. Encode the supplied value as JSON.
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


// ============================================================
// 7. TASK COLLECTION HANDLER
// ============================================================

// handleTasks supports two operations:
//
// GET  /tasks  -> retrieve tasks
// POST /tasks  -> create a task
//
// Both operations use the same URL but have different
// HTTP methods.
//
// This is an example of method-based routing.
//
// The handler is responsible for translating HTTP input
// into calls to the service layer.
//
// It should not independently implement task creation,
// validation, or persistence business rules.
func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {

    // switch evaluates r.Method and executes the matching case.
    //
    // It is similar to a switch statement in other languages.
    //
    // Each case handles a different HTTP operation.
    switch r.Method {

    // --------------------------------------------------------
    // GET /tasks
    // --------------------------------------------------------

    case http.MethodGet:

        // Retrieve all tasks from the service layer.
        //
        // The service is responsible for managing task data.
        //
        // The returned value is assigned to tasks.
        //
        // The handler does not need to know how the tasks
        // are stored internally.
        tasks := s.manager.ListTasks()

        // ----------------------------------------------------
        // QUERY PARAMETERS
        // ----------------------------------------------------

        // Query parameters are optional values supplied
        // after a question mark in the URL.
        //
        // Example:
        //
        // GET /tasks?status=Pending
        //
        // r.URL.Query() returns the query parameters
        // as a url.Values collection.
        //
        // Get("status") retrieves the value associated
        // with the status key.
        //
        // The result is a string.
        //
        // Example:
        //     "Pending"
        //
        // TrimSpace removes leading and trailing whitespace.
        //
        // Example:
        //     "  Pending  " -> "Pending"
        //
        // The := operator declares and initializes status.
        if status := strings.TrimSpace(r.URL.Query().Get("status")); status != "" {

            // If the user provided a non-empty status,
            // retrieve only tasks matching that status.
            //
            // This delegates filtering to the service layer.
            //
            // The original tasks variable is reassigned.
            //
            // Notice the difference:
            //
            // := declares a new variable.
            // = assigns a new value to an existing variable.
            //
            // Here, tasks already exists, so we use =.
            tasks = s.manager.FilterTasksByStatus(status)
        }

        // Convert the tasks value into a JSON response.
        //
        // HTTP status: 200 OK.
        //
        // Example response:
        //
        // [
        //   {
        //     "id": 1,
        //     "title": "Learn Go"
        //   }
        // ]
        //
        // The exact JSON fields depend on your Task struct
        // and its JSON struct tags.
        //
        // JSON encoding is handled by writeJSON.
        writeJSON(w, http.StatusOK, tasks)


    // --------------------------------------------------------
    // POST /tasks
    // --------------------------------------------------------

    case http.MethodPost:

        // Define the expected JSON request structure.
        //
        // This is an anonymous struct.
        //
        // Anonymous means it is declared directly here
        // rather than given a separate type name.
        //
        // The API expects two fields:
        //
        // {
        //     "id": 1,
        //     "title": "Learn Go"
        // }
        //
        // JSON TAGS:
        //
        // `json:"id"` tells encoding/json to map the JSON
        // property "id" to the Go field ID.
        //
        // Go exported field:
        //     ID
        //
        // JSON property:
        //     id
        //
        // JSON tags define how Go values map to JSON names.
        //
        // The fields are exported because their names start
        // with uppercase letters.
        //
        // encoding/json generally needs exported struct
        // fields to decode JSON into them.
        var request struct {
            ID    int    `json:"id"`
            Title string `json:"title"`
        }

        // Decode the incoming JSON request body.
        //
        // Example incoming request:
        //
        // POST /tasks
        // Content-Type: application/json
        //
        // {
        //     "id": 1,
        //     "title": "Learn Go"
        // }
        //
        // decodeJSON is a helper function.
        //
        // The second argument is the HTTP request.
        //
        // The third argument is a pointer to request.
        //
        // Why a pointer?
        //
        // Decoding needs to modify the request struct.
        // Passing &request allows the function to populate
        // the actual variable rather than a copy.
        //
        // If decoding fails, the helper should write an
        // appropriate HTTP error response.
        //
        // false means decoding failed.
        //
        // Return immediately to avoid continuing with
        // invalid or incomplete request data.
        if !decodeJSON(w, r, &request) {
            return
        }

        // Call the service layer to create the task.
        //
        // request.ID:
        //     Integer ID received from JSON.
        //
        // request.Title:
        //     Title received from JSON.
        //
        // CreateTask returns an error.
        //
        // This follows Go's standard explicit error handling:
        //
        //     if err != nil {
        //         handle the error
        //     }
        //
        // The short declaration creates err within this
        // if statement's scope.
        if err := s.manager.CreateTask(request.ID, request.Title); err != nil {

            // Convert the application error into an HTTP
            // status code.
            //
            // statusFromError is a helper function.
            //
            // For example:
            //     Task not found       -> 404
            //     Invalid input        -> 400
            //     Duplicate task ID    -> possibly 409
            //
            // The actual mapping depends on your helper.
            //
            // writeError serializes the error into JSON
            // and writes the corresponding HTTP status.
            writeError(w, statusFromError(err), err)

            // Stop because task creation failed.
            return
        }

        // Persist the modified task data.
        //
        // CreateTask may only modify the in-memory state.
        //
        // Save() writes the current state to persistent storage.
        //
        // This is an important distinction:
        //
        // In-memory state:
        //     Current application data held in RAM.
        //
        // Persistent state:
        //     Data saved to disk or another durable storage system.
        //
        // If Save fails, the API returns an internal server error.
        if err := s.manager.Save(); err != nil {

            // 500 Internal Server Error.
            //
            // This means the server could not complete
            // the requested operation due to an internal failure.
            writeError(w, http.StatusInternalServerError, err)
            return
        }

        // Retrieve the newly created task.
        //
        // Why retrieve it again?
        //
        // CreateTask returns only an error.
        // It does not return the created Task object.
        //
        // FindTask obtains the actual task so the API
        // can return it to the client.
        created, err := s.manager.FindTask(request.ID)

        // If the task cannot be retrieved after creation,
        // return an internal server error.
        //
        // This is not expected under normal operation,
        // assuming the service maintains its invariants.
        if err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }

        // Return the created task.
        //
        // 201 Created means that the requested resource
        // was successfully created.
        //
        // This is more semantically accurate than returning
        // 200 OK for a successful resource creation.
        writeJSON(w, http.StatusCreated, created)


    // --------------------------------------------------------
    // UNSUPPORTED HTTP METHOD
    // --------------------------------------------------------

    default:

        // Return 405 Method Not Allowed.
        //
        // For this route, only GET and POST are supported.
        //
        // Example:
        // DELETE /tasks
        //
        // This is a valid URL but an unsupported operation.
        //
        // Note that ServeMux's method-aware patterns normally
        // reject unsupported methods before reaching this
        // switch. The default remains useful as defensive logic.
        methodNotAllowed(w, http.MethodGet, http.MethodPost)
    }
}


// ============================================================
// 8. INDIVIDUAL TASK HANDLER
// ============================================================

// handleTask operates on one task identified by its ID.
//
// Supported endpoints:
//
// GET    /tasks/{id}
// PATCH  /tasks/{id}
// DELETE /tasks/{id}
//
// This handler introduces two important HTTP concepts:
//
// 1. Path parameters.
// 2. Different operations on the same resource.
//
// The resource is an individual task.
// The HTTP method determines what happens to it.
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {

    // --------------------------------------------------------
    // EXTRACT AND VALIDATE PATH PARAMETER
    // --------------------------------------------------------

    // Extract the id captured by the ServeMux router.
    //
    // Example:
    //
    // Request:
    //     GET /tasks/123
    //
    // r.PathValue("id"):
    //     "123"
    //
    // IMPORTANT:
    //
    // PathValue returns a string.
    // It does not automatically convert the value into an int.
    idStr := r.PathValue("id")

    // Convert the string into an integer.
    //
    // strconv.Atoi:
    //
    // Input:  "123"
    // Output: 123
    //
    // Atoi returns two values:
    //
    //     (int, error)
    //
    // If conversion fails, the error is non-nil.
    //
    // Examples:
    //
    // "123" -> 123, nil
    // "abc" -> 0, error
    // ""    -> 0, error
    id, err := strconv.Atoi(idStr)

    // Validate the ID.
    //
    // The condition uses logical OR: ||
    //
    // If conversion failed OR the ID is negative,
    // the request is rejected.
    //
    // Why check err first?
    //
    // A failed conversion produces an invalid ID.
    // We must not use that ID in business operations.
    if err != nil || id < 0 {

        // Return HTTP 400 Bad Request.
        //
        // This indicates that the client supplied
        // an invalid request parameter.
        writeError(w, http.StatusBadRequest, errors.New("invalid task ID"))
        return
    }


    // --------------------------------------------------------
    // SELECT OPERATION
    // --------------------------------------------------------

    // The ID is now validated.
    //
    // Determine which operation the client requested.
    switch r.Method {


    // --------------------------------------------------------
    // GET /tasks/{id}
    // --------------------------------------------------------

    case http.MethodGet:

        // Retrieve the task with the supplied ID.
        //
        // FindTask belongs to the service layer.
        //
        // It is responsible for locating the task and
        // returning an error if the task does not exist.
        task, err := s.manager.FindTask(id)

        // If the task does not exist or retrieval fails,
        // convert the error into an appropriate HTTP response.
        if err != nil {
            writeError(w, statusFromError(err), err)
            return
        }

        // Return the task as JSON.
        //
        // HTTP 200 OK.
        //
        // Example:
        //
        // GET /tasks/123
        //
        // Response:
        // {
        //     "id": 123,
        //     "title": "Learn Go"
        // }
        writeJSON(w, http.StatusOK, task)


    // --------------------------------------------------------
    // PATCH /tasks/{id}
    // --------------------------------------------------------

    case http.MethodPatch:

        // Define the fields accepted for this update.
        //
        // PATCH is generally used for partial modification.
        //
        // Example:
        //
        // PATCH /tasks/123
        // Content-Type: application/json
        //
        // {
        //     "title": "Learn HTTP APIs"
        // }
        //
        // Unlike PUT, PATCH does not necessarily replace
        // the entire resource.
        //
        // Here, the API only allows the title to be updated.
        //
        // This is another anonymous struct.
        var update struct {
            Title string `json:"title"`
        }

        // Decode the JSON request body into update.
        //
        // The function writes an error response and returns
        // false if decoding fails.
        if !decodeJSON(w, r, &update) {
            return
        }

        // Ask the service layer to update the task.
        //
        // The HTTP handler does not directly modify
        // the task object.
        //
        // The service layer should enforce business rules
        // such as checking whether the task exists.
        if err := s.manager.UpdateTask(id, update.Title); err != nil {

            // Map the service error to an HTTP status code.
            writeError(w, statusFromError(err), err)
            return
        }

        // Save the updated task state to persistent storage.
        if err := s.manager.Save(); err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }

        // Retrieve the updated task so it can be returned
        // to the client.
        //
        // This ensures the response contains the actual
        // resulting task state.
        updated, err := s.manager.FindTask(id)

        if err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }

        // Return the updated resource.
        //
        // HTTP 200 OK.
        writeJSON(w, http.StatusOK, updated)


    // --------------------------------------------------------
    // DELETE /tasks/{id}
    // --------------------------------------------------------

    case http.MethodDelete:

        // Ask the service layer to delete the task.
        //
        // The service is responsible for determining whether
        // the task exists and whether deletion is permitted.
        if err := s.manager.DeleteTask(id); err != nil {

            // For example, a missing task should normally
            // map to 404 Not Found.
            writeError(w, statusFromError(err), err)
            return
        }

        // Persist the deletion.
        //
        // Without this operation, the task might disappear
        // from memory but reappear after restarting the program.
        if err := s.manager.Save(); err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }

        // Return HTTP 204 No Content.
        //
        // 204 means the operation succeeded but the response
        // intentionally contains no response body.
        //
        // WriteHeader sets the HTTP status code.
        //
        // No JSON encoding is necessary here.
        w.WriteHeader(http.StatusNoContent)


    // --------------------------------------------------------
    // UNSUPPORTED METHOD
    // --------------------------------------------------------

    default:

        // Only GET, PATCH and DELETE are supported here.
        //
        // Return HTTP 405 for other methods.
        methodNotAllowed(w, http.MethodGet, http.MethodPatch, http.MethodDelete)
    }
}


// ============================================================
// 9. COMPLETE TASK HANDLER
// ============================================================

// handleComplete marks an existing task as complete.
//
// Endpoint:
//
// POST /tasks/{id}/complete
//
// This is an action-oriented endpoint.
//
// Instead of requiring the client to understand how
// completion is represented internally, the client simply
// requests the completion operation.
//
// The service layer owns the actual state transition.
//
// This keeps HTTP-specific logic separate from business logic.
func (s *Server) handleComplete(w http.ResponseWriter, r *http.Request) {

    // Validate the HTTP method.
    //
    // Only POST is accepted.
    if r.Method != http.MethodPost {
        methodNotAllowed(w, http.MethodPost)
        return
    }

    // Extract the task ID from the URL.
    //
    // Example:
    // POST /tasks/123/complete
    //
    // PathValue("id") returns "123".
    idStr := r.PathValue("id")

    // Convert the string into an integer.
    //
    // The error is checked immediately after conversion.
    id, err := strconv.Atoi(idStr)

    // Reject invalid or negative IDs.
    if err != nil || id < 0 {
        writeError(w, http.StatusBadRequest, errors.New("invalid task ID"))
        return
    }

    // Execute the business operation.
    //
    // CompleteTask should contain the actual logic
    // for changing the task's status to completed.
    //
    // The HTTP layer does not need to know the internal
    // representation of task completion.
    if err := s.manager.CompleteTask(id); err != nil {

        // Translate application errors into HTTP errors.
        writeError(w, statusFromError(err), err)
        return
    }

    // Persist the updated state.
    if err := s.manager.Save(); err != nil {
        writeError(w, http.StatusInternalServerError, err)
        return
    }

    // Retrieve the resulting task.
    //
    // This allows the client to see the final state
    // after the completion operation.
    task, err := s.manager.FindTask(id)

    if err != nil {
        writeError(w, http.StatusInternalServerError, err)
        return
    }

    // Return the completed task.
    //
    // HTTP 200 OK.
    writeJSON(w, http.StatusOK, task)
}


// ============================================================
// 10. REOPEN TASK HANDLER
// ============================================================

// handleReopen changes a completed task back into an
// incomplete or pending state.
//
// Endpoint:
//
// POST /tasks/{id}/reopen
//
// This is conceptually the reverse of handleComplete.
//
// The request lifecycle is almost identical:
//
// 1. Validate HTTP method.
// 2. Extract ID.
// 3. Convert ID.
// 4. Validate ID.
// 5. Execute service operation.
// 6. Persist state.
// 7. Retrieve resulting task.
// 8. Return JSON response.
//
// The only meaningful business operation difference
// is ReopenTask instead of CompleteTask.
func (s *Server) handleReopen(w http.ResponseWriter, r *http.Request) {

    // Reject unsupported HTTP methods.
    if r.Method != http.MethodPost {
        methodNotAllowed(w, http.MethodPost)
        return
    }

    // Extract the ID and convert it to an integer.
    //
    // Here we perform conversion directly instead of
    // assigning the string to a separate idStr variable.
    //
    // Both approaches are valid.
    //
    // The returned error is still checked below.
    id, err := strconv.Atoi(r.PathValue("id"))

    // Reject malformed or negative IDs.
    if err != nil || id < 0 {
        writeError(w, http.StatusBadRequest, errors.New("invalid task ID"))
        return
    }

    // Execute the business operation.
    //
    // ReopenTask is responsible for changing the task's
    // state according to the service's business rules.
    if err := s.manager.ReopenTask(id); err != nil {

        // Return the corresponding HTTP error.
        writeError(w, statusFromError(err), err)
        return
    }

    // Save the modified task state.
    if err := s.manager.Save(); err != nil {
        writeError(w, http.StatusInternalServerError, err)
        return
    }

    // Retrieve the task after reopening it.
    task, err := s.manager.FindTask(id)

    // Handle unexpected retrieval errors.
    if err != nil {
        writeError(w, http.StatusInternalServerError, err)
        return
    }

    // Return the updated task.
    writeJSON(w, http.StatusOK, task)
}
