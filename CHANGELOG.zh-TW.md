# 變更日誌

本檔案記錄本專案所有值得注意的變更。

格式依循 [Keep a Changelog](https://keepachangelog.com/)。

## [v0.2.0] - 2026-03-11

### 新增

- 屬性型別系統 — 在型別 schema 中定義 9 種屬性型別（`string`、`text`、`number`、`bool`、`date`、`datetime`、`url`、`enum`、`relation`）(#8)
- 共用屬性 — 在 `.typemd/properties.yaml` 定義可重用的屬性，並透過 `use` 在型別 schema 中參照 (#188)
- 型別 Emoji — 在型別 schema 加入 `emoji` 欄位，於 TUI 中視覺化識別型別 (#145)
- 屬性 Emoji — 在屬性 schema 加入 `emoji` 欄位，用於緊湊顯示 (#144)
- TUI 標題面板 — 瀏覽物件時顯示型別 emoji 與物件名稱的專用標題列 (#169)
- TUI 置頂屬性 — 在 schema 中標記 `pinned: true`，使屬性在 TUI 詳細檢視中突出顯示 (#168)
- TUI Session 持久化 — 游標位置、選取物件與面板狀態在 TUI 重新啟動後恢復 (#82)
- `--readonly` 旗標 — 以唯讀模式啟動 TUI，停用所有編輯功能 (#107)
- `--reindex` 旗標 — 全域旗標，啟動時強制重建 SQLite 索引，取代原本的 `tmd reindex` 子指令 (#159)
- 前綴比對 — 可用 ULID 後綴的短前綴解析物件，不需輸入完整 ID (#72)
- Homebrew 安裝 — 透過 `brew install typemd/tap/tmd` 安裝 (#140)

### 變更

- `name` 屬性 — 現為必要系統屬性，自動從物件 slug 填入；型別 schema 不可自行定義名為 `name` 的屬性 (#187)
- TUI 物件列表 — 群組標頭中顯示型別 emoji (#163)
- 未定義屬性 — 型別 schema 未宣告的屬性在同步時會被靜默過濾 (#174)

### 修正

- Relation 顯示 — 移除 relation 屬性顯示值中的 ULID 後綴

[v0.2.0]: https://github.com/typemd/typemd/releases/tag/v0.2.0

## [v0.1.0] - 2026-03-08

### 新增

- 物件與型別 — 在 YAML 中定義型別 schema，透過 `tmd object create` 建立 Markdown 物件檔案 (#18)
- ULID 檔名 — 唯一的 ULID 後綴，避免物件命名衝突 (#48)
- Relation — 透過 `tmd relation link` / `tmd relation unlink` 建立雙向連結，支援單值覆寫與多值附加
- Wiki-links 與反向連結 — 在 Markdown 內文中使用 `[[target]]` 語法，自動追蹤反向連結 (#10)
- 查詢 — `tmd query` 依型別與屬性篩選，`tmd search` 全文搜尋，皆支援 `--json` 輸出
- 驗證 — `tmd type validate` 檢查 schema 完整性、屬性型別、孤立 relation 與壞掉的 wiki-links (#20)
- 遷移 — `tmd migrate` 在 schema 演進時更新既有物件 (#22)
- 自動重建索引 — SQLite 索引為空或遺失時自動重建 (#41)
- 孤立清理 — 重新索引時偵測並移除過期的 relation (#21)
- CLI 重組 — 指令依資源類型分組：`tmd object`、`tmd type`、`tmd relation` (#141)
- TUI — 三面板介面 (#47)、原地內文編輯 (#85)、編輯模式視覺指示 (#84)、退出時自動儲存 (#86)、快捷鍵說明 (#104)
- TUI 顯示 — 移除顯示名稱中的 ULID (#75)、縮減縮排 (#57)、群組化物件列表 (#43)
- MCP Server — `tmd mcp` 將 vault 開放給 AI 助手使用
- `.gitignore` 初始化 — `tmd init` 建立 `.typemd/.gitignore` 排除 `index.db` (#1)
- `tmd` 執行檔 — `go install` 產生 `tmd` binary (#61)
- 支援英文與繁體中文的文件網站 (#50, #54)
- 使用 Godog 與 Gherkin feature 檔案的 BDD 測試框架 (#111, #112)
- GitHub Actions 跨平台編譯發布流程 (#39)
- 程式碼重構 — 統一命名慣例、抽取 helper、改善錯誤處理 (#56)
- Vault 結構重構 — 移除 `objects/` 目錄層 (#117)

[v0.1.0]: https://github.com/typemd/typemd/releases/tag/v0.1.0
