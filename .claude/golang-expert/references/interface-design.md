# Interface Design in Go

## Table of Contents
1. [Interface Segregation](#interface-segregation)
2. [Accept Interfaces Return Structs](#accept-interfaces-return-structs)
3. [Consumer-Defined Interfaces](#consumer-defined-interfaces)
4. [Composition Patterns](#composition-patterns)
5. [Common Mistakes](#common-mistakes)

## Interface Segregation

### Small, Focused Interfaces

```go
// GOOD - single responsibility interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Compose when needed
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

### Real-World Example

```go
// BAD - fat interface
type Repository interface {
    Get(ctx context.Context, id string) (*Entity, error)
    List(ctx context.Context, filter Filter) ([]*Entity, error)
    Create(ctx context.Context, e *Entity) error
    Update(ctx context.Context, e *Entity) error
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context, filter Filter) (int, error)
    Exists(ctx context.Context, id string) (bool, error)
    Transaction(ctx context.Context, fn func(Repository) error) error
}

// GOOD - segregated interfaces
type EntityReader interface {
    Get(ctx context.Context, id string) (*Entity, error)
}

type EntityLister interface {
    List(ctx context.Context, filter Filter) ([]*Entity, error)
}

type EntityWriter interface {
    Create(ctx context.Context, e *Entity) error
    Update(ctx context.Context, e *Entity) error
}

type EntityDeleter interface {
    Delete(ctx context.Context, id string) error
}

// Functions accept only what they need
func ProcessEntity(ctx context.Context, r EntityReader, id string) error {
    entity, err := r.Get(ctx, id)
    // ...
}
```

## Accept Interfaces Return Structs

### The Pattern

```go
// Function accepts interface
func ProcessData(r io.Reader) (*Result, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, err
    }
    return &Result{Data: data}, nil // Returns concrete struct
}

// Constructor returns concrete type
func NewService(logger Logger, repo Repository) *Service {
    return &Service{
        logger: logger,
        repo:   repo,
    }
}
```

### Why This Pattern?

```go
// Accepting interfaces = flexibility
// Any type implementing io.Reader can be used
func Parse(r io.Reader) (*Document, error) {
    // Works with files, buffers, network connections, etc.
}

// Can call with:
Parse(os.Stdin)
Parse(strings.NewReader("data"))
Parse(bytes.NewBuffer(data))
Parse(resp.Body)

// Returning structs = no unnecessary abstraction
// Caller knows exactly what they get
func NewParser() *Parser { // Not ParserInterface
    return &Parser{}
}
```

## Consumer-Defined Interfaces

### Define Interfaces Where Used

```go
// package database - defines concrete type
package database

type UserStore struct {
    db *sql.DB
}

func (s *UserStore) GetUser(ctx context.Context, id int) (*User, error) { ... }
func (s *UserStore) SaveUser(ctx context.Context, u *User) error { ... }
func (s *UserStore) DeleteUser(ctx context.Context, id int) error { ... }
func (s *UserStore) ListUsers(ctx context.Context) ([]*User, error) { ... }

// package service - defines interface it needs
package service

// Only defines what it actually uses
type UserGetter interface {
    GetUser(ctx context.Context, id int) (*User, error)
}

type Service struct {
    users UserGetter // Depends on minimal interface
}

func NewService(users UserGetter) *Service {
    return &Service{users: users}
}

// In tests - easy to mock with minimal interface
type mockUserGetter struct {
    user *User
    err  error
}

func (m *mockUserGetter) GetUser(ctx context.Context, id int) (*User, error) {
    return m.user, m.err
}
```

## Composition Patterns

### Embedding for Composition

```go
// Base types
type Logger interface {
    Log(msg string)
}

type Metrics interface {
    Record(name string, value float64)
}

// Composed interface
type Observable interface {
    Logger
    Metrics
}

// Struct embedding
type BaseService struct {
    logger  Logger
    metrics Metrics
}

type UserService struct {
    BaseService // Embedded - gains logger and metrics
    repo UserRepository
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    s.logger.Log("getting user")
    defer s.metrics.Record("user.get", 1)
    return s.repo.Get(ctx, id)
}
```

### Decorator Pattern

```go
type Handler interface {
    Handle(ctx context.Context, req Request) (Response, error)
}

// Logging decorator
type LoggingHandler struct {
    next   Handler
    logger Logger
}

func (h *LoggingHandler) Handle(ctx context.Context, req Request) (Response, error) {
    h.logger.Log("handling request")
    resp, err := h.next.Handle(ctx, req)
    h.logger.Log("request complete")
    return resp, err
}

// Metrics decorator
type MetricsHandler struct {
    next    Handler
    metrics Metrics
}

func (h *MetricsHandler) Handle(ctx context.Context, req Request) (Response, error) {
    start := time.Now()
    resp, err := h.next.Handle(ctx, req)
    h.metrics.Record("request.duration", time.Since(start).Seconds())
    return resp, err
}

// Compose decorators
handler := &LoggingHandler{
    logger: logger,
    next: &MetricsHandler{
        metrics: metrics,
        next:    &CoreHandler{},
    },
}
```

## Common Mistakes

### Mistake 1: Empty Interface Abuse

```go
// BAD - loses type safety
func Process(data interface{}) interface{} {
    // Type assertions everywhere
    switch v := data.(type) {
    case string:
        return processString(v)
    case int:
        return processInt(v)
    }
    return nil
}

// GOOD - use generics or specific types
func Process[T Processable](data T) Result {
    return data.Process()
}

// Or specific functions
func ProcessString(s string) StringResult { ... }
func ProcessInt(i int) IntResult { ... }
```

### Mistake 2: Interface Pollution

```go
// BAD - interface defined "just in case"
type UserService interface {
    GetUser(id int) (*User, error)
    // Only one implementation exists
}

type userServiceImpl struct { ... }

// GOOD - concrete type until abstraction needed
type UserService struct { ... }

func (s *UserService) GetUser(id int) (*User, error) { ... }

// Add interface ONLY when:
// 1. Multiple implementations exist
// 2. Testing requires mocking
// 3. Package boundary requires abstraction
```

### Mistake 3: Large Interfaces in Parameters

```go
// BAD - requires full implementation for testing
func SendNotification(repo UserRepository, id int) error {
    user, err := repo.GetUser(id) // Only uses GetUser
    if err != nil {
        return err
    }
    return notify(user.Email)
}

// GOOD - minimal interface
type UserEmailGetter interface {
    GetUserEmail(id int) (string, error)
}

func SendNotification(getter UserEmailGetter, id int) error {
    email, err := getter.GetUserEmail(id)
    if err != nil {
        return err
    }
    return notify(email)
}
```

### Mistake 4: Returning Interfaces

```go
// BAD - returns interface hiding implementation
func NewStore() Store { // Returns interface
    return &concreteStore{}
}

// Caller doesn't know what they get
// Can't access concrete methods
// Harder to debug

// GOOD - returns concrete type
func NewStore() *ConcreteStore {
    return &ConcreteStore{}
}

// Caller can use as interface if needed:
var store Store = NewStore()
```
