package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ConvertWordToPDF は、指定されたWord文書(.doc, .docx)をPDFに変換します。
func ConvertWordToPDF(inputPath, outputPath string) error {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("入力パスの解決に失敗: %w", err)
	}
	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("出力パスの解決に失敗: %w", err)
	}

	// 既存のPDFが存在する場合は先に削除しておく
	if _, err := os.Stat(absOutput); err == nil {
		if err := os.Remove(absOutput); err != nil {
			return fmt.Errorf("既存PDFの削除に失敗しました: %w", err)
		}
	}

	switch runtime.GOOS {
	case "windows":
		return convertOnWindows(absInput, absOutput)
	case "darwin": // macOS
		return convertOnMac(absInput, absOutput)
	default:
		return fmt.Errorf("未対応のOSです: %s", runtime.GOOS)
	}
}

// convertOnWindows はPowerShellとCOMオブジェクトを使用して変換します
func convertOnWindows(inputPath, outputPath string) error {
	psCmd := fmt.Sprintf(`
$word = New-Object -ComObject Word.Application
$word.Visible = $false
try {
    $doc = $word.Documents.Open('%s')
    $doc.SaveAs('%s', 17)
    $doc.Close()
} catch {
    Write-Error $_.Exception.Message
    exit 1
} finally {
    $word.Quit()
}
`, inputPath, outputPath)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Windows変換エラー: %v\n詳細: %s", err, string(out))
	}
	return nil
}

// convertOnMac はAppleScript(osascript)を使用して変換します
func convertOnMac(inputPath, outputPath string) error {
	appleScript := fmt.Sprintf(`
set inPath to POSIX file "%s"

tell application "Microsoft Word"
	activate
	open inPath
	
	-- Sandbox(アクセス権)エラー(-1708)を回避するため、Word自身のテンポラリフォルダに一旦保存する
	set tempFolder to path to temporary items as string
	set tempFileName to "export_" & (random number from 1000 to 9999) & ".pdf"
	set tempFilePath to tempFolder & tempFileName
	
	try
		-- 権限が保証されている一時フォルダへ書き出し
		save as active document file name tempFilePath file format format PDF
		set saveSuccess to true
	on error errMsg number errNum
		set saveSuccess to false
		set eMsg to errMsg
		set eNum to errNum
	end try
	
	-- エラーが起きても、絶対に保存せずにドキュメントを閉じる（ダイアログ防止）
	close active document saving no
	
	if not saveSuccess then
		error eMsg number eNum
	end if
	
	-- 成功したら一時ファイルのパスをGoへ返す
	return POSIX path of tempFilePath
end tell
`, inputPath)

	cmd := exec.Command("osascript")
	cmd.Stdin = strings.NewReader(appleScript)

	// コマンドを実行し、一時PDFのパス（標準出力）を受け取る
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("macOS変換エラー: %v\n詳細: %s", err, string(out))
	}

	tempPdfPath := strings.TrimSpace(string(out))
	if tempPdfPath == "" {
		return fmt.Errorf("一時PDFファイルのパスが取得できませんでした")
	}

	// Go側で一時ファイルを最終的な出力先(outputPath)へ移動する
	err = moveFile(tempPdfPath, outputPath)
	if err != nil {
		return fmt.Errorf("ファイルの移動に失敗しました: %v", err)
	}

	return nil
}

// moveFile はファイルを移動します（別ドライブ間でも動作するようにコピー＆削除を使用）
func moveFile(src, dst string) error {
	bytesRead, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, bytesRead, 0644)
	if err != nil {
		return err
	}
	return os.Remove(src)
}

// --- 動作テスト用 ---
func main() {
	fmt.Println("=== Word to PDF Converter Prototype ===")
	fmt.Printf("実行環境OS: %s\n", runtime.GOOS)

	inputFile := "./test_document.docx"
	outputFile := "./test_output.pdf"

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("【警告】テスト用入力ファイルが見つかりません: %s\n", inputFile)
		return
	}

	fmt.Printf("変換処理を開始します...\n入力: %s\n出力: %s\n", inputFile, outputFile)

	err := ConvertWordToPDF(inputFile, outputFile)
	if err != nil {
		fmt.Printf("❌ 変換失敗:\n%v\n", err)
		return
	}

	fmt.Println("✅ 変換に成功しました！")
}
