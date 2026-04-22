# GoTestKit - Go Testing Toolkit

A comprehensive testing utilities library for Go that makes writing tests easier, more expressive, and more maintainable.

## Features

- **assert** - Fluent assertion library with rich error messages
- **mock** - Simple and powerful mocking framework
- **testdata** - Test data generation and management
- **fake** - Fake data generators (Faker)
- **testcontainer** - Testcontainers integration helpers
- **httptest** - HTTP testing utilities

## Installation

```bash
go get github.com/atop0914/gotestkit
```

## Quick Start

```go
import (
    "testing"
    "github.com/atop0914/gotestkit/assert"
    "github.com/atop0914/gotestkit/fake"
)

func TestExample(t *testing.T) {
    // Assertions
    assert.Equal(t, 3, 3, "values should match")
    assert.Contains(t, "hello world", "world")
    
    // Fake data generation
    name := fake.Person().Name()
    email := fake.Internet().Email()
}
```

## License

MIT
