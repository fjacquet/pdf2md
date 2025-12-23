# Go Code Review Checklist

## Quick Checklist

Use this checklist when reviewing Go code:

- [ ] **KISS**: Is this the simplest solution?
- [ ] **DRY**: Is there duplicated code?
- [ ] **Errors**: Are all errors handled?
- [ ] **Interfaces**: Are they small and focused?
- [ ] **Tests**: Is there adequate coverage?
- [ ] **Concurrency**: Any race conditions?
- [ ] **Performance**: Any obvious bottlenecks?

---

## Detailed Review Categories

### 1. Simplicity (KISS)

**Check for:**
- Unnecessary abstractions (interfaces for single implementations)
- Over-engineered patterns (factories, builders when not needed)
- Deeply nested code (prefer guard clauses)
- Complex one-liners (break into readable steps)

**Ask:**
- Can a junior developer understand this in 5 minutes?
- Is there a simpler way to achieve this?
- Would deleting code make it better?

```go
// RED FLAG - unnecessary complexity
type UserFactoryInterface interface {
    CreateUser(name string) UserInterface
}

// GREEN - simple and direct
func NewUser(name string) *User {
    return &User{Name: name}
}
```

### 2. DRY (Don't Repeat Yourself)

**Check for:**
- Duplicated validation logic
- Copy-pasted error handling
- Repeated string/number literals
- Similar functions that could be unified

**Ask:**
- Is this the third time I'm seeing this pattern?
- Should this be extracted to a shared function?
- Is there a constant for this magic value?

```go
// RED FLAG - duplicated validation
func CreateUser(u User) error {
    if u.Email == "" { return errors.New("email required") }
    // ...
}
func UpdateUser(u User) error {
    if u.Email == "" { return errors.New("email required") }
    // ...
}

// GREEN - shared validation
func ValidateUser(u User) error {
    if u.Email == "" { return errors.New("email required") }
    return nil
}
```

### 3. Error Handling

**Check for:**
- Ignored errors (`_` without justification)
- Missing error context in wrapping
- Error strings starting with capital letters
- Panic for recoverable errors

**Ask:**
- Is every error either handled or returned?
- Does the error message provide context?
- Can the caller distinguish error types if needed?

```go
// RED FLAGS
result, _ := DoSomething()           // Ignored error
return fmt.Errorf("error: %w", err)  // No context
return fmt.Errorf("Error: %w", err)  // Capitalized
panic(err)                           // For recoverable error

// GREEN
result, err := DoSomething()
if err != nil {
    return fmt.Errorf("process item %s: %w", id, err)
}
```

### 4. Interface Design

**Check for:**
- Large interfaces (more than 3-5 methods)
- Interfaces defined by implementer, not consumer
- Empty interface (`interface{}`) abuse
- Returning interfaces instead of concrete types

**Ask:**
- Does each consumer need all these methods?
- Is this interface defined where it's used?
- Could this be a smaller, focused interface?

```go
// RED FLAG - fat interface
type Repository interface {
    Get, List, Create, Update, Delete, Count, Exists, Tx...
}

// GREEN - segregated interfaces
type EntityReader interface {
    Get(ctx context.Context, id string) (*Entity, error)
}
```

### 5. Concurrency

**Check for:**
- Loop variable capture bugs
- Missing synchronization (maps, slices)
- Goroutine leaks (no way to exit)
- Data races (shared state without locks)

**Ask:**
- Is context used for cancellation?
- Can this goroutine exit cleanly?
- Is shared state properly protected?

```go
// RED FLAG - loop variable capture
for _, item := range items {
    go func() { process(item) }() // Bug!
}

// GREEN
for _, item := range items {
    item := item
    go func() { process(item) }()
}
```

### 6. Testing

**Check for:**
- Missing test cases (happy path only)
- Tests without error case coverage
- Non-deterministic tests (time, random)
- Tests that test implementation, not behavior

**Ask:**
- Are edge cases covered?
- Can tests run in parallel?
- Do tests use table-driven pattern?

```go
// RED FLAG - incomplete test
func TestParse(t *testing.T) {
    result := Parse("valid")
    if result != expected {
        t.Fail()
    }
}

// GREEN - table-driven, covers cases
func TestParse(t *testing.T) {
    tests := []struct{
        name string
        input string
        want Result
        wantErr bool
    }{
        {"valid", "abc", Result{...}, false},
        {"empty", "", Result{}, true},
        {"invalid", "!!!", Result{}, true},
    }
    // ...
}
```

### 7. Performance (if relevant)

**Check for:**
- Allocations in hot paths
- String concatenation in loops
- Missing pre-allocation for known sizes
- N+1 query patterns

**Ask:**
- Is this code in a hot path?
- Can we pre-allocate?
- Are we doing repeated work?

```go
// RED FLAG - repeated allocations
var s string
for _, item := range items {
    s += item.Name // O(n²)
}

// GREEN - pre-allocated builder
var b strings.Builder
b.Grow(estimatedSize)
for _, item := range items {
    b.WriteString(item.Name)
}
```

### 8. Functional Principles

**Check for:**
- Global mutable state
- Functions that mutate inputs
- Side effects in "getter" functions
- Dependencies not injected

**Ask:**
- Can this function be pure?
- Is state explicitly passed, not global?
- Are dependencies injected?

```go
// RED FLAG - global state
var db *sql.DB
func GetUser(id int) (*User, error) {
    return queryUser(db, id)
}

// GREEN - injected dependency
type Service struct { db *sql.DB }
func (s *Service) GetUser(id int) (*User, error) {
    return queryUser(s.db, id)
}
```

---

## Review Comment Examples

### Good Review Comments

```
// Suggestion: This could be simplified using a guard clause
// instead of nested if statements.

// Question: Is there a reason we're not using the existing
// ValidateEmail function from the common package?

// Issue: This goroutine has no way to exit if the context
// is canceled. Consider adding a ctx.Done() case.

// Nit: Error messages conventionally start with lowercase
// in Go: "process user" instead of "Process user"
```

### Avoid These Comments

```
// "Why did you do this?" (assumes negative intent)
// "This is wrong" (without explanation)
// "I would have done it differently" (not actionable)
// "LGTM" (when issues exist)
```

---

## Severity Levels

1. **Blocker**: Security issue, data corruption risk, guaranteed bug
2. **Major**: Logic error, missing error handling, race condition
3. **Minor**: Code style, naming, minor optimization
4. **Nit**: Formatting, typos, suggestions for future
