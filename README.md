# yamlctl

`yamlctl` is a small command-line utility for reading and updating YAML scalar values while preserving surrounding comments and layout as much as possible.

It is useful when you want to automate small YAML edits without rewriting the whole file.

## Features

- Read a scalar value from a YAML file with `get`
- Update a scalar value in place with `set`
- Preserve nearby comments, spacing, and quoted scalar style when possible
- Keep the CLI simple with dot-separated paths such as `root.child`

## Installation

### Build from source

```bash
go build -o yamlctl .
```

### Run without building

```bash
go run . --help
```

## Usage

### Get a value

Given a file:

```yaml
app:
  name: demo
  port: "8080"
```

Run:

```bash
yamlctl get config.yaml app.name
```

Output:

```text
demo
```

### Set a value

Given a file:

```yaml
# application settings
app:
  name: demo
  port: "8080" # public port
```

Run:

```bash
yamlctl set config.yaml app.port 9090
```

Result:

```yaml
# application settings
app:
  name: demo
  port: "9090" # public port
```

## Path Syntax

- Paths use dot-separated keys, for example `root.child.value`
- The current implementation supports mapping traversal only
- Array indexing is not supported

## Behavior Notes

- `get` returns only scalar values
- `set` updates only existing scalar values
- `set` does not create missing keys
- Quoted values keep their original quote style when updated
- Plain scalar values may be automatically quoted when required to keep valid YAML
- Block scalars such as `|` and `>` are not supported for in-place updates

## Limitations

- Only scalar values can be updated without rewriting file formatting
- Invalid paths return an error
- Missing keys return an error
- Non-scalar targets return an error

## Development

Run tests with:

```bash
go test ./...
```

Check the available commands with:

```bash
go run . --help
```
