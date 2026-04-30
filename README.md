# GoTestKit - Go Testing Toolkit

A comprehensive testing utilities library for Go that makes writing tests easier, more expressive, and more maintainable.

## Features

- **assert** - Fluent assertion library with rich error messages and composable matchers
- **mock** - Simple and powerful mocking framework with expectation tracking
- **fake** - Fake data generators (Faker) for realistic test data
- **httptest** - HTTP testing utilities for REST API testing
- **testdata** - Test data loading and golden file testing support
- **testcontainer** - Docker container integration for database integration tests
- **benchmark** - Benchmark utilities for performance testing
- **suite** - Test suite framework with Setup/Teardown lifecycle hooks

## Installation

```bash
go get github.com/atop0914/gotestkit
```

## Quick Start

### Assertions

```go
import "github.com/atop0914/gotestkit/assert"

func TestExample(t *testing.T) {
    assert.Equal(t, 3, 3, "values should match")
    assert.Contains(t, "hello world", "world")
    assert.Nil(t, nil)
    assert.NotNil(t, &struct{}{})
    assert.True(t, true)
    assert.Error(t, func() error { return errors.New("failed") })
    assert.Panics(t, func() { panic("oops") })
}
```

### Fake Data Generation

```go
import "github.com/atop0914/gotestkit/fake"

func TestExample(t *testing.T) {
    name := fake.Person().Name()
    email := fake.Internet().Email()
    phone := fake.Phone().Number()
    company := fake.Company().Name()
}
```

### HTTP Testing

```go
import "github.com/atop0914/gotestkit/httptest"

func TestAPI(t *testing.T) {
    server := httptest.NewServer(t, handler)
    defer server.Close()

    resp, body := server.Get(t, "/api/users")
    assert.Equal(t, 200, resp.StatusCode)
    assert.JSONContains(t, body, "users")
}
```

### Mocking

```go
import "github.com/atop0914/gotestkit/mock"

func TestMock(t *testing.T) {
    m := mock.NewMock(t)
    m.ExpectCall("Save").With("data").Returns(nil)
    
    err := SaveData(m, "data")
    assert.Nil(t, err)
    assert.True(t, m.AllExpectationsMet())
}
```

### Test Suite

```go
import "github.com/atop0914/gotestkit/suite"

type MySuite struct {
    suite.TestSuite
    db *sql.DB
}

func (s *MySuite) SetupSuite() {
    s.db = connectDB()
}

func (s *MySuite) TearDownSuite() {
    s.db.Close()
}

func (s *MySuite) TestSomething() {
    // test code
}

func TestMySuite(t *testing.T) {
    suite.RunSuite(t, new(MySuite))
}
```

### Benchmark Utilities

```go
import "github.com/atop0914/gotestkit/benchmark"

func BenchmarkExample(b *testing.B) {
    result := benchmark.Run(b, func(b *testing.B) {
        // benchmark code
    })
    
    stats := result.Stats()
    fmt.Printf("%.2f ns/op\n", stats.Mean())
}
```

### Container Testing

```go
import (
    "github.com/atop0914/gotestkit/testcontainer"
    "context"
)

func TestWithDB(t *testing.T) {
    ctx := context.Background()
    
    db, err := testcontainer.StartPostgres(ctx,
        testcontainer.WithPostgresDatabase("testdb"),
        testcontainer.WithPostgresUsername("user"),
        testcontainer.WithPostgresPassword("pass"),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer db.Terminate(ctx)
    
    // Use db.ConnectionString(ctx) to connect
}
```

## Modules

| Module | Description |
|--------|-------------|
| `assert` | 30+ assertion functions with clear error messages |
| `mock` | Mock framework with Expect/Called pattern |
| `fake` | 100+ fake data generators (names, emails, addresses, etc.) |
| `httptest` | HTTP server/client testing helpers |
| `testdata` | Test data loading and golden file comparison |
| `testcontainer` | PostgreSQL, MySQL, Redis, MongoDB containers |
| `benchmark` | Benchmark result tracking and comparison |
| `suite` | Test suite with Setup/Teardown lifecycle |

## License

MIT
