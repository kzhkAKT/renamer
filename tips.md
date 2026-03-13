How to compile on macOS for macOS
```zsh
go build -o mac_convert main.go
```

How to compile on macOS for Windows
```zsh
GOOS=darwin GOARCH=arm64 go build -o mac_convert_arm64 main.go
```

