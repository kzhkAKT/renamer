const prefixInput = document.getElementById('prefix');
const fileListContainer = document.getElementById('file-list');
const dropZone = document.getElementById('drop-zone');
const placeholderText = document.getElementById('placeholder-text');

// ドロップされた元のファイル名を保持する配列
let originalFiles = [];

// 1. 添字(prefix)が入力・変更されるたびにリストを更新
prefixInput.addEventListener('input', updateList);

// 2. ドラッグ＆ドロップのイベント処理
dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropZone.classList.add('dragover'); // 背景色を変える
});

dropZone.addEventListener('dragleave', (e) => {
    e.preventDefault();
    dropZone.classList.remove('dragover'); // 背景色を戻す
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('dragover');

    if (e.dataTransfer.files.length > 0) {
        // 中央の「ここにドロップ」の文字を消す
        placeholderText.style.display = 'none';
        
        // ドロップされたファイルの名前を配列に追加
        for (const file of e.dataTransfer.files) {
            originalFiles.push(file.name);
        }
        updateList(); // 画面を更新
    }
});

// 3. リストを描画する関数
function updateList() {
    fileListContainer.innerHTML = ''; // 一度リストを空にする
    const prefix = prefixInput.value;

    originalFiles.forEach(fileName => {
        const row = document.createElement('div');
        row.className = 'list-row';

        // ★変更前（左側）のカラム作成： 'col-old' という目印を追加
        const oldNameCol = document.createElement('div');
        oldNameCol.className = 'col col-old';
        oldNameCol.textContent = fileName;

        // ★変更後（右側）のカラム作成： 'col-new' という目印を追加
        const newNameCol = document.createElement('div');
        newNameCol.className = 'col col-new';
        newNameCol.textContent = prefix + fileName;

        row.appendChild(oldNameCol);
        row.appendChild(newNameCol);
        fileListContainer.appendChild(row);
    });
}
