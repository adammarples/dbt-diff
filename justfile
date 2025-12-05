# Build the binary
build:
    go build -o bin/dbt-diff

# Install to PATH
install:
    go build -o ~/.local/bin/dbt-diff
    @echo "Installed to ~/.local/bin/dbt-diff"

# Run tests
test:
    go test ./...

# Clean build artifacts and generated manifests
clean:
    rm -rf bin/
    rm -rf target/main/
    rm -rf target/local/

# Run go mod tidy
sync:
    go mod tidy

# Show project info
inspect:
    @echo "Go version:"
    @go version
    @echo ""
    @echo "Dependencies:"
    @go list -m all
    @echo ""
    @echo "Build status:"
    @if [ -f bin/dbt-diff ]; then echo "  Binary: bin/dbt-diff (built)"; else echo "  Binary: not built"; fi
    @if [ -f ~/.local/bin/dbt-diff ]; then echo "  Installed: ~/.local/bin/dbt-diff"; else echo "  Installed: no"; fi
