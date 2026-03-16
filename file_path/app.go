package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ★ここから下を一番最後に追加：実際にファイル名を変更する処理
func (a *App) RenameFiles(prefix string, filePaths []string) string {
	successCount := 0

	for _, oldPath := range filePaths {
		// フォルダのパス(dir)と、元のファイル名(base)に分ける
		dir := filepath.Dir(oldPath)
		base := filepath.Base(oldPath)

		// 新しいフルパスを作成
		newPath := filepath.Join(dir, prefix+base)

		// OSの機能を使って実際にファイル名を変更
		err := os.Rename(oldPath, newPath)
		if err == nil {
			successCount++
		}
	}

	// 完了メッセージを返す
	return fmt.Sprintf("%d個のファイルのファイル名を変更しました！", successCount)
}
