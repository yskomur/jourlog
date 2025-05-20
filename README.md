# jourlog

A Go package for logging to systemd's journal with advanced features.

## Features

- Log to systemd journal with proper metadata
- Support for all standard log levels (Emergency, Alert, Critical, Error, Warning, Notice, Info, Debug)
- Context-aware logging to include request IDs, user IDs, and other metadata
- Capture file, line, and function information for each log entry
- Option to echo logs to console
- Journal reading capabilities with filtering options

## Installation

```bash
go get github.com/yskomur/jourlog
```

## Usage

### Basic Logging

```go
package main

import (
    "github.com/yskomur/jourlog"
)

func main() {
    // Use the global logger
    jourlog.JLog.Info("Application started")

    // Log with formatting
    jourlog.JLog.Info("User %s logged in", "john.doe")

    // Log an error
    jourlog.JLog.Error("Failed to connect to database: %v", err)

    // Set log level
    jourlog.JLog.SetLogLevel(journal.PriDebug)

    // Enable console output
    jourlog.JLog.SetEcho(true)
}
```

### Context-Aware Logging

```go
package main

import (
    "context"
    "github.com/yskomur/jourlog"
)

func handleRequest(ctx context.Context) {
    // Create a context with request ID
    ctx = context.WithValue(ctx, "request_id", "req-123456")
    ctx = context.WithValue(ctx, "user_id", "user-789")

    // Log with context
    jourlog.JLog.InfoWithContext(ctx, "Processing request")

    // Log an error with context
    if err := processData(); err != nil {
        jourlog.JLog.ErrorWithContext(ctx, "Failed to process data: %v", err)
    }
}
```

### Reading from Journal

```go
package main

import (
    "fmt"
    "github.com/yskomur/jourlog"
    "time"
)

func main() {
    // Create a journal reader
    reader, err := jourlog.NewJournalReader()
    if err != nil {
        fmt.Printf("Failed to create journal reader: %v\n", err)
        return
    }
    // Don't forget to close the reader when done
    defer reader.Close()

    // Configure to read SSH service logs from the last hour
    if err := reader.LastHour(); err != nil {
        fmt.Printf("Failed to set time filter: %v\n", err)
        return
    }

    // Read entries
    for {
        message, err := reader.Retrieve()
        if err != nil {
            if err.Error() == "no more entries" {
                break
            }
            fmt.Printf("Error reading journal: %v\n", err)
            break
        }

        fmt.Println(message)
    }
}
```

### Advanced Journal Reading

```go
package main

import (
    "fmt"
    "github.com/yskomur/jourlog"
    "time"
)

func main() {
    // Create a journal reader
    reader, err := jourlog.NewJournalReader()
    if err != nil {
        fmt.Printf("Failed to create journal reader: %v\n", err)
        return
    }
    // Always close the reader when done to release resources
    defer reader.Close()

    // Set filters for specific service and unit
    if err := reader.SetService("ssh"); err != nil {
        fmt.Printf("Failed to set service filter: %v\n", err)
        return
    }

    // Set time range (last 24 hours)
    if err := reader.LastDay(); err != nil {
        fmt.Printf("Failed to set time filter: %v\n", err)
        return
    }

    // Alternative: Set custom time range
    // startTime := time.Now().Add(-6 * time.Hour)
    // if err := reader.SetSince(startTime); err != nil {
    //     fmt.Printf("Failed to set start time: %v\n", err)
    //     return
    // }

    // Set priority filter (only warning and above)
    if err := reader.SetPriority(4); err != nil {
        fmt.Printf("Failed to set priority filter: %v\n", err)
        return
    }

    // Set hostname filter
    if err := reader.SetHostname("myserver"); err != nil {
        fmt.Printf("Failed to set hostname filter: %v\n", err)
        return
    }

    // Set limit to retrieve only 100 entries
    reader.SetLimit(100)

    // Move to the beginning of the journal
    if err := reader.SeekHead(); err != nil {
        fmt.Printf("Failed to seek to head: %v\n", err)
        return
    }

    // Read entries
    entryCount := 0
    for {
        message, err := reader.Retrieve()
        if err != nil {
            if err.Error() == "no more entries" || err.Error() == "reached limit of 100 entries" {
                break
            }
            fmt.Printf("Error reading journal: %v\n", err)
            break
        }

        fmt.Printf("%d: %s\n", entryCount+1, message)
        entryCount++
    }

    fmt.Printf("Retrieved %d entries\n", entryCount)

    // Clear all filters for a new query
    reader.ClearFilters()

    // Reset the counter
    reader.ResetCounter()

    // Set a new filter for today's logs only
    if err := reader.Today(); err != nil {
        fmt.Printf("Failed to set today filter: %v\n", err)
        return
    }

    // Get and save the current cursor position
    cursor, err := reader.GetCursor()
    if err != nil {
        fmt.Printf("Failed to get cursor: %v\n", err)
    } else {
        fmt.Printf("Current cursor position: %s\n", cursor)

        // Later, you can seek back to this position
        // if err := reader.SeekCursor(cursor); err != nil {
        //     fmt.Printf("Failed to seek to cursor: %v\n", err)
        // }
    }
}
```

## Log Levels

The package supports the following log levels, in order of decreasing severity:

1. **Emergency** (`Emerge`): System is unusable
2. **Alert** (`Alert`): Action must be taken immediately
3. **Critical** (`Critical`): Critical conditions
4. **Error** (`Error`): Error conditions
5. **Warning** (`Warning`): Warning conditions
6. **Notice** (`Notice`): Normal but significant condition
7. **Info** (`Info`): Informational messages
8. **Debug** (`Debug`): Debug-level messages

## License

See the [LICENSE](LICENSE) file for details.
