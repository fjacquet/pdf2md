# Functional Patterns in Go

## Table of Contents
1. [Dependency Injection](#dependency-injection)
2. [Immutability](#immutability)
3. [Pure Functions](#pure-functions)
4. [No Global State](#no-global-state)
5. [Option Pattern](#option-pattern)

## Dependency Injection

### Constructor Injection (Preferred)

```go
// Define interface for dependency
type Logger interface {
    Info(msg string, args ...any)
    Error(msg string, args ...any)
}

type Repository interface {
    Get(ctx context.Context, id string) (*Entity, error)
    Save(ctx context.Context, e *Entity) error
}

// Service with injected dependencies
type Service struct {
    logger Logger
    repo   Repository
}

// NewService - configure everything at construction
func NewService(logger Logger, repo Repository) *Service {
    return &Service{
        logger: logger,
        repo:   repo,
    }
}

// Methods use injected dependencies
func (s *Service) Process(ctx context.Context, id string) error {
    s.logger.Info("processing", "id", id)
    entity, err := s.repo.Get(ctx, id)
    if err != nil {
        return fmt.Errorf("get entity: %w", err)
    }
    // process...
    return s.repo.Save(ctx, entity)
}
```

### Container Pattern

```go
// Container holds all dependencies
type Container struct {
    logger     Logger
    config     Config
    repo       Repository
    service    *Service
}

func NewContainer(cfg Config) (*Container, error) {
    c := &Container{config: cfg}

    // Initialize in dependency order
    c.logger = NewLogger(cfg.LogLevel)

    repo, err := NewRepository(cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("create repository: %w", err)
    }
    c.repo = repo

    c.service = NewService(c.logger, c.repo)

    return c, nil
}

// Getters provide access (immutable after construction)
func (c *Container) GetService() *Service { return c.service }
func (c *Container) GetLogger() Logger    { return c.logger }
```

## Immutability

### Return New Values Instead of Mutating

```go
// BAD - mutates input
func AddTag(user *User, tag string) {
    user.Tags = append(user.Tags, tag)
}

// GOOD - returns new value
func WithTag(user User, tag string) User {
    result := user
    result.Tags = make([]string, len(user.Tags)+1)
    copy(result.Tags, user.Tags)
    result.Tags[len(user.Tags)] = tag
    return result
}

// GOOD - builder pattern for complex mutations
type UserBuilder struct {
    user User
}

func NewUserBuilder(base User) *UserBuilder {
    return &UserBuilder{user: base}
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
    b.user.Name = name
    return b
}

func (b *UserBuilder) WithTag(tag string) *UserBuilder {
    b.user.Tags = append(b.user.Tags, tag)
    return b
}

func (b *UserBuilder) Build() User {
    return b.user
}
```

### Immutable Configuration

```go
// Config is immutable after creation
type Config struct {
    host     string
    port     int
    timeout  time.Duration
}

func NewConfig(host string, port int, timeout time.Duration) Config {
    return Config{host: host, port: port, timeout: timeout}
}

// Only getters, no setters
func (c Config) Host() string          { return c.host }
func (c Config) Port() int             { return c.port }
func (c Config) Timeout() time.Duration { return c.timeout }
func (c Config) Address() string       { return fmt.Sprintf("%s:%d", c.host, c.port) }
```

## Pure Functions

### Characteristics of Pure Functions

1. Same input always produces same output
2. No side effects (no I/O, no mutation of external state)
3. Doesn't depend on external mutable state

```go
// PURE - deterministic, no side effects
func CalculateTotal(items []Item) float64 {
    var total float64
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }
    return total
}

// PURE - transformation function
func FilterActive(users []User) []User {
    result := make([]User, 0, len(users))
    for _, u := range users {
        if u.Active {
            result = append(result, u)
        }
    }
    return result
}

// IMPURE - has side effects (logging)
func ProcessWithLogging(data []byte, logger Logger) error {
    logger.Info("processing data") // Side effect!
    // ...
}

// IMPURE - depends on external state
var globalConfig Config // BAD
func GetTimeout() time.Duration {
    return globalConfig.Timeout // Depends on global!
}
```

### Separating Pure and Impure Code

```go
// Pure business logic
func ValidateOrder(order Order) []ValidationError {
    var errs []ValidationError
    if order.Total <= 0 {
        errs = append(errs, ValidationError{Field: "total", Msg: "must be positive"})
    }
    if len(order.Items) == 0 {
        errs = append(errs, ValidationError{Field: "items", Msg: "cannot be empty"})
    }
    return errs
}

// Impure orchestration layer
func (s *Service) CreateOrder(ctx context.Context, order Order) error {
    // Pure validation
    if errs := ValidateOrder(order); len(errs) > 0 {
        return &ValidationErrors{Errors: errs}
    }

    // Impure I/O
    return s.repo.Save(ctx, &order)
}
```

## No Global State

### Eliminating Globals

```go
// BAD - package-level mutable state
var (
    db     *sql.DB
    logger *log.Logger
)

func init() {
    db, _ = sql.Open("postgres", os.Getenv("DB_URL"))
    logger = log.New(os.Stdout, "", log.LstdFlags)
}

func GetUser(id int) (*User, error) {
    return queryUser(db, id) // Uses global!
}

// GOOD - explicit dependencies
type UserRepository struct {
    db     *sql.DB
    logger *log.Logger
}

func NewUserRepository(db *sql.DB, logger *log.Logger) *UserRepository {
    return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*User, error) {
    r.logger.Printf("fetching user %d", id)
    return r.queryUser(ctx, id)
}
```

### Constants Over Variables

```go
// BAD - mutable package variables
var (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// GOOD - constants
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// For complex "constants", use functions
func DefaultConfig() Config {
    return Config{
        Timeout:    DefaultTimeout,
        MaxRetries: MaxRetries,
    }
}
```

## Option Pattern

### Functional Options for Flexible Configuration

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
    logger  Logger
}

type ServerOption func(*Server)

func WithHost(host string) ServerOption {
    return func(s *Server) { s.host = host }
}

func WithPort(port int) ServerOption {
    return func(s *Server) { s.port = port }
}

func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) { s.timeout = d }
}

func WithLogger(l Logger) ServerOption {
    return func(s *Server) { s.logger = l }
}

func NewServer(opts ...ServerOption) *Server {
    s := &Server{
        host:    "localhost",
        port:    8080,
        timeout: 30 * time.Second,
        logger:  DefaultLogger(),
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(
    WithHost("0.0.0.0"),
    WithPort(9000),
    WithLogger(customLogger),
)
```
