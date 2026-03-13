package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 1. アプリとウィンドウの作成
	a := app.New()
	w := a.NewWindow("成績資料PDFコンバーター")

	// 2. 画面に表示するテキスト
	statusLabel := widget.NewLabel("ここにWordファイルをドラッグ＆ドロップしてください")
	statusLabel.Alignment = fyne.TextAlignCenter

	// 3. ウィンドウにファイルがドロップされた時の処理
	w.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		// 複数のファイルがドロップされても1つずつ処理
		for _, uri := range uris {
			inputPath := uri.Path()
			
			// 拡張子がWordかどうか簡易チェック
			ext := filepath.Ext(inputPath)
			if ext != ".doc" && ext != ".docx" {
				statusLabel.SetText("エラー: Wordファイル以外がドロップされました\n" + inputPath)
				continue
			}

			statusLabel.SetText(fmt.Sprintf("処理中...\n%s", inputPath))

			// --------------------------------------------------------
			// ★ここに先ほど完成した ConvertWordToPDF 関数を呼び出す処理を書きます★
			// 例: err := ConvertWordToPDF(inputPath, outputPath)
			// --------------------------------------------------------

			statusLabel.SetText(fmt.Sprintf("✅ 変換完了！\n%s", inputPath))
		}
	})

	// 4. 画面のレイアウトとサイズ調整
	w.SetContent(container.NewCenter(statusLabel))
	w.Resize(fyne.NewSize(500, 300))

	// 5. アプリの実行
	w.ShowAndRun()
}
