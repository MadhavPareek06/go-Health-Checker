### Web Service Health Monitor

This project is a small, but robust, web service health monitor built in Go. Its core purpose is to continuously check the status and response time of a list of web services (HTTP endpoints), aggregate the results, and report on their health in real time.

The application is a practical case study in Go's native concurrency features and idiomatic design patterns, showcasing how to build a resilient and efficient monitoring tool.

-----

### Project Architecture and File Structure

The project is structured into a logical, component-based file layout, where each file is responsible for a single, well-defined task.

```
health-monitor/
├── config.go       // Loads application settings
├── checker.go      // Contains the health check logic
├── reporter.go     // Defines the reporting interface
├── aggregator.go   // Gathers and summarizes results
├── main.go         // Wires everything together
├── config.json     // User-defined configuration file
└── go.mod          // Go module file for dependency management
```

-----

### Core Concepts

This project demonstrates several fundamental Go concepts in a real-world application:

  * **Goroutines and Channels:** The program uses a **worker pool** pattern to perform concurrent health checks. The `main` function sends URLs to a `jobs` channel, and a fixed number of **goroutines** pull from this channel, perform the checks, and send results back on a `results` channel. This approach is a hallmark of Go's concurrency model, which prioritizes communication over shared memory.
  * **Interfaces:** The `reporter.go` file uses an **interface** (`HealthReporter`) to define the contract for reporting. This allows you to easily swap the current console output for a different reporting mechanism, such as a Slack alert or a log file, without changing the core application logic.
  * **Thread Safety:** Since multiple goroutines are running at the same time, the program uses synchronization primitives to safely manage shared data. For simple counters, the `sync/atomic` package provides efficient, lock-free operations. For the main service status map, it uses `sync.Map`, a concurrent map designed for high-performance, thread-safe data access.
  * **Graceful Shutdown:** The application is designed to shut down gracefully when it receives a termination signal, such as pressing `Ctrl+C`. It uses a dedicated channel to listen for this signal, allowing it to stop accepting new work and wait for all in-flight health checks to complete before exiting cleanly.
  * **Context:** Each HTTP request uses a `context` with a timeout. This is a crucial practice for network-based applications, as it prevents the program from hanging indefinitely on a slow or unresponsive service.

-----

### How to Run the Project

Follow these steps to get the health monitor running on your machine.

**1. Go Environment:**
Ensure you have Go installed on your system. You can verify this by running `go version` in your terminal.

**2. Initialize the Project:**
Navigate to the root of your project directory in the terminal and initialize a Go module.

```bash
cd health-monitor
go mod init health-monitor
```

**3. Create the Configuration File:**
Create a file named `config.json` in your project's root directory. This file tells the program which URLs to monitor and how often to check them.

```json
{
  "ping_interval_ms": 5000,
  "concurrency": 10,
  "request_timeout_ms": 2000,
  "services": [
    "https://www.google.com",
    "https://www.example.com",
    "https://www.github.com"
  ]
}
```

**4. Run the Program:**
With the configuration file in place, you can now run the application.

```bash
go run.
```

The `go run.` command will find and execute the `main.go` file and its related packages within the current directory. You will see a report printed to the console every 5 seconds (based on the example configuration).

**5. Stop the Program:**
To stop the program, simply press `Ctrl+C` in the terminal. The graceful shutdown mechanism will take over, ensuring all ongoing tasks are completed before the application exits.

-----

### Future Enhancements

The current design provides a solid foundation that can be easily extended. Potential future enhancements could include:

  * **Adding New Reporters:** Use the `HealthReporter` interface to implement a reporter that sends alerts to a messaging service (e.g., Slack) or logs results to a file.
  * **Persistent Storage:** Integrate a database (e.g., SQLite, PostgreSQL) to store historical health data, allowing for long-term trend analysis.
  * **Dynamic Configuration:** Enhance the `config.go` component to support reloading the configuration at runtime without needing to restart the application.
