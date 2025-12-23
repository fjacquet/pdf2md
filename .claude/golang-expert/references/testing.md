# Testing in Go

## Table of Contents
1. [Table-Driven Tests](#table-driven-tests)
2. [Test Helpers](#test-helpers)
3. [Mocks and Fakes](#mocks-and-fakes)
4. [Integration Tests](#integration-tests)
5. [Benchmarks](#benchmarks)
6. [Test Patterns](#test-patterns)

## Table-Driven Tests

### Basic Pattern

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 5, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### With Error Cases

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Result
        wantErr bool
        errMsg  string // Optional: specific error message
    }{
        {
            name:  "valid input",
            input: "valid",
            want:  Result{Value: "valid"},
        },
        {
            name:    "empty input",
            input:   "",
            wantErr: true,
            errMsg:  "input cannot be empty",
        },
        {
            name:    "invalid format",
            input:   "!!!",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)

            if tt.wantErr {
                if err == nil {
                    t.Fatal("expected error, got nil")
                }
                if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
                    t.Errorf("error = %q, want containing %q", err.Error(), tt.errMsg)
                }
                return
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Parallel Tests

```go
func TestConcurrent(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"test1", "a", "A"},
        {"test2", "b", "B"},
        {"test3", "c", "C"},
    }

    for _, tt := range tests {
        tt := tt // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Run in parallel
            got := Process(tt.input)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

## Test Helpers

### Setup and Teardown

```go
func TestWithTempDir(t *testing.T) {
    // t.TempDir() auto-cleans up
    dir := t.TempDir()

    // Create test file
    path := filepath.Join(dir, "test.txt")
    if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
        t.Fatal(err)
    }

    // Test your code
    result, err := ProcessFile(path)
    if err != nil {
        t.Fatal(err)
    }
    // assertions...
}
```

### Helper Functions

```go
// Helper marks function as test helper
func createTestUser(t *testing.T, name string) *User {
    t.Helper() // Error line numbers point to caller

    user := &User{Name: name, Email: name + "@test.com"}
    if err := user.Validate(); err != nil {
        t.Fatalf("createTestUser: %v", err)
    }
    return user
}

func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

### Test Fixtures

```go
// testdata/ directory is ignored by go build
func TestParseConfig(t *testing.T) {
    data, err := os.ReadFile("testdata/config.yaml")
    if err != nil {
        t.Fatal(err)
    }

    config, err := ParseConfig(data)
    assertNoError(t, err)
    assertEqual(t, config.Name, "test")
}
```

## Mocks and Fakes

### Interface-Based Mocking

```go
// Production interface
type UserRepository interface {
    Get(ctx context.Context, id int) (*User, error)
    Save(ctx context.Context, u *User) error
}

// Mock implementation
type mockUserRepo struct {
    users map[int]*User
    err   error
}

func newMockUserRepo() *mockUserRepo {
    return &mockUserRepo{users: make(map[int]*User)}
}

func (m *mockUserRepo) Get(ctx context.Context, id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    user, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

func (m *mockUserRepo) Save(ctx context.Context, u *User) error {
    if m.err != nil {
        return m.err
    }
    m.users[u.ID] = u
    return nil
}

// Usage in tests
func TestService_GetUser(t *testing.T) {
    repo := newMockUserRepo()
    repo.users[1] = &User{ID: 1, Name: "Alice"}

    svc := NewService(repo)

    user, err := svc.GetUser(context.Background(), 1)
    assertNoError(t, err)
    assertEqual(t, user.Name, "Alice")
}

func TestService_GetUser_NotFound(t *testing.T) {
    repo := newMockUserRepo()
    svc := NewService(repo)

    _, err := svc.GetUser(context.Background(), 999)
    if !errors.Is(err, ErrNotFound) {
        t.Errorf("expected ErrNotFound, got %v", err)
    }
}
```

### Fake Services

```go
// Fake HTTP server
func TestHTTPClient(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/users/1" {
            json.NewEncoder(w).Encode(User{ID: 1, Name: "Test"})
            return
        }
        http.NotFound(w, r)
    }))
    defer server.Close()

    client := NewClient(server.URL)
    user, err := client.GetUser(1)
    assertNoError(t, err)
    assertEqual(t, user.Name, "Test")
}
```

## Integration Tests

### Build Tags

```go
//go:build integration

package mypackage

import "testing"

func TestDatabaseIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    db, err := sql.Open("postgres", os.Getenv("TEST_DB_URL"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // Run integration tests...
}
```

Run with: `go test -tags=integration ./...`

### TestMain for Setup

```go
var testDB *sql.DB

func TestMain(m *testing.M) {
    // Setup
    var err error
    testDB, err = sql.Open("postgres", os.Getenv("TEST_DB_URL"))
    if err != nil {
        log.Fatal(err)
    }

    // Run tests
    code := m.Run()

    // Teardown
    testDB.Close()

    os.Exit(code)
}

func TestWithDatabase(t *testing.T) {
    // Use testDB...
}
```

## Benchmarks

### Basic Benchmark

```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData(1000)

    b.ResetTimer() // Don't count setup time
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}
```

### With Sub-Benchmarks

```go
func BenchmarkParse(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            data := generateData(size)
            b.ResetTimer()

            for i := 0; i < b.N; i++ {
                Parse(data)
            }
        })
    }
}
```

### Memory Benchmarks

```go
func BenchmarkAllocation(b *testing.B) {
    b.ReportAllocs() // Report memory allocations

    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}
```

Run with: `go test -bench=. -benchmem ./...`

## Test Patterns

### Golden Files

```go
func TestRender(t *testing.T) {
    got := Render(testData)

    golden := filepath.Join("testdata", t.Name()+".golden")

    if *update { // -update flag
        os.WriteFile(golden, []byte(got), 0644)
    }

    want, err := os.ReadFile(golden)
    if err != nil {
        t.Fatal(err)
    }

    if got != string(want) {
        t.Errorf("mismatch:\ngot:\n%s\nwant:\n%s", got, want)
    }
}
```

### Test Environment Variables

```go
func TestWithEnv(t *testing.T) {
    // Save and restore
    orig := os.Getenv("API_KEY")
    defer os.Setenv("API_KEY", orig)

    os.Setenv("API_KEY", "test-key")

    // Or use t.Setenv (auto-restores)
    t.Setenv("API_KEY", "test-key")

    // Test code...
}
```

### Testing Time-Dependent Code

```go
// Inject time source
type Clock interface {
    Now() time.Time
}

type realClock struct{}
func (realClock) Now() time.Time { return time.Now() }

type mockClock struct {
    now time.Time
}
func (m mockClock) Now() time.Time { return m.now }

type Service struct {
    clock Clock
}

func TestExpiry(t *testing.T) {
    clock := mockClock{now: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
    svc := &Service{clock: clock}

    // Test with controlled time
}
```
