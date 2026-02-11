# jourlog

`jourlog` is a Go package for writing structured logs to `systemd-journald` and reading logs back using journal filters.

## Installation

```bash
go get github.com/yskomur/jourlog
```

## Quick start

```go
package main

import (
	"github.com/coreos/go-systemd/v22/journal"
	"github.com/yskomur/jourlog"
)

func main() {
	logger := jourlog.NewJourlog()
	logger.SetLogLevel(journal.PriInfo)
	logger.SetEcho(true)

	logger.Info("application started")
	logger.Error("database unavailable")
}
```

## Read from journal

```go
package main

import (
	"fmt"
	"log"

	"github.com/yskomur/jourlog"
)

func main() {
	reader, err := jourlog.NewJournalReader()
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	if err := reader.LastHour(); err != nil {
		log.Fatal(err)
	}

	for {
		message, err := reader.Retrieve()
		if err != nil {
			break
		}
		fmt.Println(message)
	}
}
```

## API overview

### Logging (`JourLog`)

- `NewJourlog()` creates a logger.
- `SetLogLevel(...)` sets the minimum accepted priority.
- `SetEcho(true)` mirrors messages to stdout.
- Priority helpers: `Emerge`, `Alert`, `Critical`, `Error`, `Warning`, `Notice`, `Info`, `Debug`.

Each log record includes caller metadata fields (`CODE_FILE`, `CODE_LINE`, `CODE_FUNC`) in journald.

### Journal reading (`JournalReader`)

- `NewJournalReader()` opens a journal reader.
- Time filters: `SetSince`, `SetUntil`, `LastHour`, `LastDay`, `LastWeek`, `Today`.
- Match filters: `SetUnit`, `SetService`, `SetPriority`, `SetHostname`, `SetExecutable`, `SetMessageFilter`, `AddFilter`.
- Cursor helpers: `SeekHead`, `SeekTail`, `SeekCursor`, `GetCursor`.
- Limits and state: `SetLimit`, `ResetCounter`, `ClearFilters`.

## Documentation

- Package docs: `go doc github.com/yskomur/jourlog`
- Rendered docs: `pkg.go.dev/github.com/yskomur/jourlog`

## License

See [LICENSE](LICENSE).
