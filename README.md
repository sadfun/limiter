# Rate Limiter

+ Custom storage (Map, Redis, KeyDB, etc)
+ Very precise algorithm (Sliding window log)
+ Generics-driven (store limits by any `comparable` type)
+ No external dependencies

## Performance considerations
Please keep in mind that sliding window log is very accurate, but can use up to `O(nm)` memory, where `m` is max limit and `n` is count of connected clients. This is a trade-off between precision and performance.
Thus, this limiter is recommended for counting very costly operations, such as credit card transactions, account creations, AI model calls etc.

## Installation

```bash
go get -u github.com/sadfun/limiter
```

## Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/sadfun/limiter"
    "github.com/sadfun/limiter/storage"
)

func main() {
    limiter := limiter.NewLimiter[string](&limiter.Config[string]{
        Limit: 10,
        Duration: time.Second,
        Storage: storage.NewMapStorage[string](),
    })

    if !limiter.Use("key") {
        panic("Limit exceeded")
    }
}
```

## Testing

```bash
go test -v .
```