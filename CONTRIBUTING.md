# Contributing to pdf2md

Thank you for your interest in contributing to `pdf2md`! We welcome contributions from the community to make this tool better.

## Getting Started

1.  **Fork the repository** on GitHub.
2.  **Clone your fork** locally:
    ```bash
    git clone https://github.com/your-username/pdf2md.git
    cd pdf2md
    ```
3.  **Install dependencies**:
    ```bash
    go mod download
    ```

## Development Workflow

1.  **Create a branch** for your feature or fix:
    ```bash
    git checkout -b feature/my-new-feature
    ```
2.  **Make your changes**. Please follow the existing code style.
3.  **Run tests** to ensure no regressions:
    ```bash
    go test ./...
    ```
4.  **Run the benchmark** (optional but recommended for layout changes):
    ```bash
    go run cmd/benchmark/main.go
    ```

## Code Style

- We use `gofumpt` for formatting. Please run `gofumpt -w .` before committing.
- We use `golangci-lint` for linting. Please run `golangci-lint run` to check for issues.

## Pull Requests

1.  Push your branch to your fork.
2.  Open a Pull Request against the `main` branch.
3.  Describe your changes clearly and link to any relevant issues.

## License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
