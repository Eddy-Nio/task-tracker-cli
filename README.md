# Task Tracker CLI

A command-line interface tool for managing tasks, built with Go.

## Features

- Create, read, update, and delete tasks
- Filter tasks by status
- Persistent storage using JSON
- Simple and intuitive command interface
- Status aliases for quick updates

## Installation

1. Clone the repository:

```bash
git clone https://github.com/Eddy-Nio/task-tracker-cli.git
```

2. Navigate to the project directory:

```bash
cd task-tracker-cli
```

3. Build the project:

```bash
go build -o task-tracker cmd/main.go
```

4. Run the project:

```bash
./task-tracker
```

## Usage

To see the list of available commands, run:

```bash
./task-tracker
```
### Adding a Task

```bash
./task-tracker add -t "Task Title" -d "Task Description" -s "todo/t"
```

### Listing Tasks

```bash
./task-tracker list # List all tasks
./task-tracker list -s "todo/t" # List all tasks with status "todo"
./task-tracker list -s "in_progress/ip/p" # List all tasks with status "in_progress"
./task-tracker list -s "done/d" # List all tasks with status "done"
```

### Updating a Task

```bash
./task-tracker update -i "task_id" -t "New Title" -d "New Description" -s "todo/t"
```

### Deleting a Task

```bash
./task-tracker delete -i "task_id"
```

### Clearing the Task List

```bash
./task-tracker clear
```

### Showing Help for a Command

```bash
./task-tracker help add
```
## Task Status Options

- `TODO` (aliases: `todo`, `t`)
- `IN_PROGRESS` (aliases: `in_progress`, `ip`, `p`)
- `DONE` (aliases: `done`, `d`)

## Development

### Prerequisites
- Go 1.20 or higher
- [Cobra CLI](https://github.com/spf13/cobra)

### Building from Source

```bash
go build -o task-tracker cmd/main.go
```

### Running Tests

```bash
make test
```

### Generating Coverage Report

```bash
make coverage
```

### Cleaning Generated Files

```bash
make clean
```

### Running Tests and Coverage

```bash
make test-coverage
```
The report will be generated in:
- `coverage.html`: A detailed HTML report
- `coverage.out`: A text-based coverage report

## Configuration

The application can be configured using a `config.yaml` file. If no configuration file is present, default values will be used.

### Configuration Options

```yaml
storage:
  filePath: "tasks.json"    # Path to store tasks
  backupDir: "backups"      # Directory for automatic backups

task:
  maxTitleLength: 100       # Maximum length for task titles
  maxDescriptionLength: 500 # Maximum length for task descriptions
  dateFormat: "2006-01-02"  # Date format for display
  autoBackup: true         # Enable/disable automatic backups
  backupInterval: 24h      # Interval between backups
```

### Custom Configuration

To use custom configuration:

1. Create a `config.yaml` file in the application directory
2. Override any default values as needed
3. The application will automatically load your custom configuration

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
