const prefixInput = document.getElementById('prefix');
const fileListContainer = document.getElementById('file-list');
const dropZone = document.getElementById('drop-zone');
const placeholderText = document.getElementById('placeholder-text');

// 追加したオプション設定用のUI要素の取得
const destDirInput = document.getElementById('dest-dir');
const selectDirBtn = document.getElementById('select-dir-btn');
const clearDirBtn = document.getElementById('clear-dir-btn');
const copyModeCheckbox = document.getElementById('copy-mode');

// フルパスを保存する配列
let filePaths = [];

// 1. 添字(prefix)が入力・変更されるたびにリストを更新
prefixInput.addEventListener('input', updateList);

// 2. Wailsのネイティブ機能を使って「フルパス」を取得する
window.runtime.EventsOn("wails:file-drop", (x, y, paths) => {
    placeholderText.style.display = 'none';
    
    for (const path of paths) {
        if (!filePaths.includes(path)) {
            filePaths.push(path); // フルパス（Macの /Users/... 等）を保存
        }
    }
    updateList();
});

// 見た目のためのドラッグイベント（背景色を変えるだけ）
dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropZone.classList.add('dragover');
});
dropZone.addEventListener('dragleave', (e) => {
    e.preventDefault();
    dropZone.classList.remove('dragover');
});
dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('dragover');
});

// 3. フォルダ選択ボタンの処理
selectDirBtn.addEventListener('click', async () => {
    // GoのSelectDirectory関数を呼び出してOSのダイアログを開く
    const dir = await window.go.main.App.SelectDirectory();
    if (dir) {
        destDirInput.value = dir;
    }
});

// 4. フォルダクリアボタンの処理
clearDirBtn.addEventListener('click', () => {
    destDirInput.value = "";
});

// 5. リストを描画する関数
function updateList() {
    fileListContainer.innerHTML = '';
    const prefix = prefixInput.value;

    filePaths.forEach(fullPath => {
        // フルパスからファイル名だけを切り出す (Macの / と Windowsの \ 両対応)
        const fileName = fullPath.split(/[/\\]/).pop();

        const row = document.createElement('div');
        row.className = 'list-row';

        const oldNameCol = document.createElement('div');
        oldNameCol.className = 'col col-old';
        oldNameCol.textContent = fileName;

        const newNameCol = document.createElement('div');
        newNameCol.className = 'col col-new';
        newNameCol.textContent = prefix + fileName;

        row.appendChild(oldNameCol);
        row.appendChild(newNameCol);
        fileListContainer.appendChild(row);
    });
}

// 6. アプリ終了ボタンの処理
const quitBtn = document.getElementById('quit-btn');
quitBtn.addEventListener('click', () => {
    window.runtime.Quit();
});

// 7. リネーム実行ボタンの処理
const executeBtn = document.getElementById('execute-btn');
executeBtn.addEventListener('click', async () => {
    const prefix = prefixInput.value;
    const destDir = destDirInput.value;
    const copyMode = copyModeCheckbox.checked;

    if (prefix === "") {
        alert("添字(prefix)を入力してください。");
        return;
    }
    if (filePaths.length === 0) {
        alert("ファイルがドロップされていません。");
        return;
    }

    try {
        // 引数に destDir と copyMode を追加してGoに渡す
        const resultMessage = await window.go.main.App.RenameFiles(prefix, filePaths, destDir, copyMode);
        alert(resultMessage);
        
        // 成功したらリストを空にして初期状態に戻す
        filePaths = [];
        updateList();
        placeholderText.style.display = 'block';
        prefixInput.value = ""; 
        // （保存先やチェックボックスの状態は連続使用を考慮してそのまま残します）
    } catch (error) {
        alert("エラーが発生しました: " + error);
    }
});