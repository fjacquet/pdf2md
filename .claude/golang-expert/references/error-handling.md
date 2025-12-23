# Error Handling in Go

## Table of Contents
1. [Error Wrapping](#error-wrapping)
2. [Sentinel Errors](#sentinel-errors)
3. [Custom Error Types](#custom-error-types)
4. [Error Checking](#error-checking)
5. [Patterns and Anti-Patterns](#patterns-and-anti-patterns)

## Error Wrapping

### Basic Wrapping with Context

```go
func ProcessFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open file %s: %w", path, err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("read file %s: %w", path, err)
    }

    if err := process(data); err != nil {
        return fmt.Errorf("process file %s: %w", path, err)
    }

    return nil
}

// Result: "process file config.yaml: parse line 42: unexpected token"
```

### When to Wrap vs Return As-Is

```go
// WRAP - when adding context is useful
func GetUser(id int) (*User, error) {
    user, err := db.Query(id)
    if err != nil {
        return nil, fmt.Errorf("get user %d: %w", id, err) // Adds context
    }
    return user, nil
}

// RETURN AS-IS - when context is already clear
func (s *Service) GetUser(id int) (*User, error) {
    return s.repo.GetUser(id) // repo.GetUser already has good context
}

// DON'T WRAP - at package boundary when returning sentinel
func FindUser(name string) (*User, error) {
    user, err := db.FindByName(name)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrUserNotFound // Don't wrap - this IS the error
    }
    if err != nil {
        return nil, fmt.Errorf("find user %q: %w", name, err)
    }
    return user, nil
}
```

## Sentinel Errors

### Defining Sentinel Errors

```go
package mypackage

import "errors"

// Package-level sentinel errors
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
    ErrConflict     = errors.New("conflict")
)
```

### Using Sentinel Errors

```go
func GetUser(id int) (*User, error) {
    user, err := repo.Find(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("find user: %w", err)
    }
    return user, nil
}

// Caller checks with errors.Is
user, err := GetUser(42)
if errors.Is(err, ErrNotFound) {
    // Handle not found case
}
```

## Custom Error Types

### Basic Custom Error

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Usage
func Validate(user User) error {
    if user.Email == "" {
        return &ValidationError{Field: "email", Message: "required"}
    }
    return nil
}
```

### Error with Multiple Fields

```go
type ValidationErrors struct {
    Errors []ValidationError
}

func (e *ValidationErrors) Error() string {
    var msgs []string
    for _, err := range e.Errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (e *ValidationErrors) Add(field, msg string) {
    e.Errors = append(e.Errors, ValidationError{Field: field, Message: msg})
}

func (e *ValidationErrors) HasErrors() bool {
    return len(e.Errors) > 0
}
```

### Wrappable Custom Error

```go
type OpError struct {
    Op  string
    Err error
}

func (e *OpError) Error() string {
    return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *OpError) Unwrap() error {
    return e.Err
}

// errors.Is and errors.As work through the chain
err := &OpError{Op: "read", Err: io.EOF}
errors.Is(err, io.EOF) // true
```

## Error Checking

### errors.Is - Check Error Identity

```go
if errors.Is(err, ErrNotFound) {
    // err is or wraps ErrNotFound
}

if errors.Is(err, context.Canceled) {
    // Request was canceled
}

if errors.Is(err, os.ErrNotExist) {
    // File doesn't exist
}
```

### errors.As - Check Error Type

```go
var validErr *ValidationError
if errors.As(err, &validErr) {
    fmt.Printf("field %s: %s\n", validErr.Field, validErr.Message)
}

var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Printf("path error on %s: %v\n", pathErr.Path, pathErr.Err)
}
```

### Multiple Error Checks

```go
func handleError(err error) {
    switch {
    case errors.Is(err, ErrNotFound):
        // 404
    case errors.Is(err, ErrUnauthorized):
        // 401
    case errors.Is(err, ErrInvalidInput):
        // 400
    default:
        // 500
    }
}
```

## Patterns and Anti-Patterns

### DO: Handle Errors Once

```go
// GOOD - handle at appropriate level
func main() {
    if err := run(); err != nil {
        log.Fatal(err) // Handle once at top level
    }
}

func run() error {
    return doWork() // Just propagate
}

// BAD - handle multiple times
func run() error {
    err := doWork()
    if err != nil {
        log.Printf("error: %v", err) // Logged here
        return err                    // AND returned (will be logged again)
    }
    return nil
}
```

### DO: Add Context When Wrapping

```go
// GOOD - adds useful context
return fmt.Errorf("process order %s: %w", orderID, err)

// BAD - no useful context
return fmt.Errorf("error: %w", err)
return fmt.Errorf("failed: %w", err)
```

### DON'T: Ignore Errors

```go
// BAD
result, _ := doSomething()

// BAD - error ignored
doSomething()

// GOOD - explicit ignore with comment
_ = writer.Close() // Best effort, already handled main error

// GOOD - handle it
result, err := doSomething()
if err != nil {
    return err
}
```

### DON'T: Check Error Strings

```go
// BAD - fragile, breaks with wrapping
if err.Error() == "not found" {
    // ...
}

// BAD - also fragile
if strings.Contains(err.Error(), "not found") {
    // ...
}

// GOOD - use sentinel or type
if errors.Is(err, ErrNotFound) {
    // ...
}
```

### DON'T: Panic for Expected Errors

```go
// BAD - panic for recoverable error
func GetConfig() Config {
    data, err := os.ReadFile("config.yaml")
    if err != nil {
        panic(err) // Don't do this!
    }
    // ...
}

// GOOD - return error
func GetConfig() (Config, error) {
    data, err := os.ReadFile("config.yaml")
    if err != nil {
        return Config{}, fmt.Errorf("read config: %w", err)
    }
    // ...
}

// OK - panic for programmer errors only
func MustCompile(pattern string) *regexp.Regexp {
    re, err := regexp.Compile(pattern)
    if err != nil {
        panic(err) // OK - invalid regex is a programmer error
    }
    return re
}
```

### Error Handling in Deferred Functions

```go
func WriteFile(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("create: %w", err)
    }

    defer func() {
        closeErr := f.Close()
        if err == nil {
            err = closeErr // Only set if no prior error
        }
    }()

    if _, err := f.Write(data); err != nil {
        return fmt.Errorf("write: %w", err)
    }

    return nil
}
```
