package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// フォルダ選択ダイアログを表示する関数
func (a *App) SelectDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "保存先フォルダを選択してください",
	})
	if err != nil {
		return ""
	}
	return dir
}

// 実際にファイルを処理する関数 (引数が増えました)
func (a *App) RenameFiles(prefix string, filePaths []string, destDir string, copyMode bool) string {
	successCount := 0

	for _, oldPath := range filePaths {
		dir := filepath.Dir(oldPath)
		base := filepath.Base(oldPath)

		// 保存先が指定されていればそこを使い、空なら元の場所を使う
		targetDir := dir
		if destDir != "" {
			targetDir = destDir
		}

		newPath := filepath.Join(targetDir, prefix+base)

		if copyMode {
			// コピーモード：元のファイルを残して新しいファイルを作成
			err := copyFileContents(oldPath, newPath)
			if err == nil {
				successCount++
			}
		} else {
			// 移動モード：ファイル名（または場所）を変更
			err := os.Rename(oldPath, newPath)
			if err == nil {
				successCount++
			}
		}
	}

	return fmt.Sprintf("%d個のファイルを処理しました！", successCount)
}

// ファイルの中身をコピーするための専用関数
func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}