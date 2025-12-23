# KISS & DRY Principles in Go

## Table of Contents

1. [KISS - Keep It Simple](#kiss---keep-it-simple)
2. [DRY - Don't Repeat Yourself](#dry---dont-repeat-yourself)
3. [Rule of Three](#rule-of-three)
4. [Anti-Patterns to Avoid](#anti-patterns-to-avoid)

## KISS - Keep It Simple

### Prefer Direct Solutions

```go
// BAD - over-engineered factory pattern for simple case
type UserFactory interface {
    CreateUser(name string) User
}
type DefaultUserFactory struct{}
func (f DefaultUserFactory) CreateUser(name string) User {
    return User{Name: name}
}

// GOOD - direct function
func NewUser(name string) User {
    return User{Name: name}
}
```

### Avoid Unnecessary Abstractions

```go
// BAD - abstraction for single implementation
type StringFormatter interface {
    Format(s string) string
}
type UpperCaseFormatter struct{}
func (f UpperCaseFormatter) Format(s string) string {
    return strings.ToUpper(s)
}

// GOOD - just use the function
func FormatUpper(s string) string {
    return strings.ToUpper(s)
}
```

### Flat Over Nested

```go
// BAD - deeply nested
func ProcessOrder(order Order) error {
    if order.Valid {
        if order.HasItems() {
            if order.Customer != nil {
                if order.Customer.Active {
                    // finally do something
                    return nil
                }
                return ErrInactiveCustomer
            }
            return ErrNoCustomer
        }
        return ErrNoItems
    }
    return ErrInvalidOrder
}

// GOOD - early returns (guard clauses)
func ProcessOrder(order Order) error {
    if !order.Valid {
        return ErrInvalidOrder
    }
    if !order.HasItems() {
        return ErrNoItems
    }
    if order.Customer == nil {
        return ErrNoCustomer
    }
    if !order.Customer.Active {
        return ErrInactiveCustomer
    }

    // Do the actual work
    return nil
}
```

### Simple Error Handling

```go
// BAD - over-complicated error handling
func GetUser(id int) (*User, error) {
    user, err := repo.Find(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, &NotFoundError{Resource: "user", ID: id}
        }
        return nil, &DatabaseError{Operation: "find", Err: err}
    }
    return user, nil
}

// GOOD - simple and clear
func GetUser(id int) (*User, error) {
    user, err := repo.Find(id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("find user %d: %w", id, err)
    }
    return user, nil
}
```

## DRY - Don't Repeat Yourself

### Extract Common Logic

```go
// BAD - duplicated validation
func CreateUser(u User) error {
    if u.Email == "" {
        return errors.New("email required")
    }
    if !strings.Contains(u.Email, "@") {
        return errors.New("invalid email")
    }
    // create user...
}

func UpdateUser(u User) error {
    if u.Email == "" {
        return errors.New("email required")
    }
    if !strings.Contains(u.Email, "@") {
        return errors.New("invalid email")
    }
    // update user...
}

// GOOD - extracted validation
func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("email required")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email")
    }
    return nil
}

func CreateUser(u User) error {
    if err := ValidateEmail(u.Email); err != nil {
        return err
    }
    // create user...
}

func UpdateUser(u User) error {
    if err := ValidateEmail(u.Email); err != nil {
        return err
    }
    // update user...
}
```

### Use Constants

```go
// BAD - magic strings/numbers repeated
func ParseDate(s string) (time.Time, error) {
    return time.Parse("2006-01-02", s)
}
func FormatDate(t time.Time) string {
    return t.Format("2006-01-02")
}

// GOOD - single source of truth
const DateFormat = "2006-01-02"

func ParseDate(s string) (time.Time, error) {
    return time.Parse(DateFormat, s)
}
func FormatDate(t time.Time) string {
    return t.Format(DateFormat)
}
```

### Common Utilities Package

```go
// internal/common/strings.go
package common

// TruncateString truncates s to maxLen, adding "..." if truncated
func TruncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    if maxLen <= 3 {
        return s[:maxLen]
    }
    return s[:maxLen-3] + "..."
}

// internal/common/slice.go
func Contains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0, len(slice))
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}
```

## Rule of Three

Don't extract until you see three repetitions:

```go
// First occurrence - just write it
func HandleUserRequest(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "id required", http.StatusBadRequest)
        return
    }
    // ...
}

// Second occurrence - note the duplication but don't extract yet
func HandleOrderRequest(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "id required", http.StatusBadRequest)
        return
    }
    // ...
}

// Third occurrence - NOW extract
func RequireQueryParam(r *http.Request, key string) (string, error) {
    value := r.URL.Query().Get(key)
    if value == "" {
        return "", fmt.Errorf("%s required", key)
    }
    return value, nil
}
```

### Acceptable Duplication

Sometimes duplication is acceptable:

1. **Test code** - Clarity over DRY in tests
2. **Simple one-liners** - `if err != nil { return err }`
3. **Different domains** - Similar code that may evolve differently

```go
// OK - test clarity is more important than DRY
func TestCreateUser(t *testing.T) {
    user := User{Name: "Alice", Email: "alice@example.com"}
    err := CreateUser(user)
    if err != nil {
        t.Fatalf("CreateUser() error = %v", err)
    }
}

func TestUpdateUser(t *testing.T) {
    user := User{Name: "Bob", Email: "bob@example.com"}
    err := UpdateUser(user)
    if err != nil {
        t.Fatalf("UpdateUser() error = %v", err)
    }
}
```

## Anti-Patterns to Avoid

### Premature Abstraction

```go
// BAD - interface for single implementation
type UserStore interface {
    Get(id int) (*User, error)
    Save(user *User) error
}

// Only one implementation ever exists
type PostgresUserStore struct { /* ... */ }

// GOOD - start concrete, abstract when needed
type UserStore struct {
    db *sql.DB
}

// Add interface ONLY when you need multiple implementations
// (e.g., for testing with a mock, or supporting multiple DBs)
```

### Over-Configuration

```go
// BAD - everything configurable
type ProcessorConfig struct {
    BufferSize      int
    MaxWorkers      int
    RetryCount      int
    RetryDelay      time.Duration
    EnableLogging   bool
    LogLevel        string
    MetricsEnabled  bool
    MetricsInterval time.Duration
    // 20 more fields...
}

// GOOD - sensible defaults, minimal config
type ProcessorConfig struct {
    MaxWorkers int           // default: runtime.NumCPU()
    Timeout    time.Duration // default: 30s
}

func NewProcessor(cfg ProcessorConfig) *Processor {
    if cfg.MaxWorkers == 0 {
        cfg.MaxWorkers = runtime.NumCPU()
    }
    if cfg.Timeout == 0 {
        cfg.Timeout = 30 * time.Second
    }
    return &Processor{cfg: cfg}
}
```

### Feature Flags for Hypothetical Features

```go
// BAD - code for features that don't exist
func Process(data []byte, opts Options) error {
    if opts.EnableExperimentalMode {
        // Nobody uses this
    }
    if opts.UseNewAlgorithm {
        // This is never true
    }
    if opts.LegacyCompatMode {
        // Removed 2 years ago
    }
    // actual code...
}

// GOOD - only code that's needed
func Process(data []byte) error {
    // actual code...
}
```
