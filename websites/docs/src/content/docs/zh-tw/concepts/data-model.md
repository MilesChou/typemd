---
title: 資料模型
description: TypeMD 如何儲存和索引資料。
sidebar:
  order: 6
---

## 儲存

Object 以 Markdown 檔案搭配 YAML frontmatter 儲存在 `objects/<type>/` 底下。完整的 Object ID 格式為 `type/<slug>-<ulid>`，例如 `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`。透過 CLI 建立的 Object，檔名會自動附加 ULID（26 位小寫字元）以確保唯一性。

### 系統屬性

所有 Object 都有五個由 TypeMD 管理的系統屬性：

| 屬性 | 說明 | 可變性 |
|------|------|--------|
| `name` | 顯示名稱，建立時從 slug 自動填入 | 使用者撰寫 |
| `description` | 選填的單行摘要，用於列表顯示和搜尋結果 | 使用者撰寫 |
| `created_at` | 建立時間戳記，RFC 3339 格式（僅設定一次，不會修改） | 自動管理 |
| `updated_at` | 最後修改時間戳記，RFC 3339 格式（每次儲存時更新） | 自動管理 |
| `tags` | 標籤參照的陣列（關聯到內建 `tag` 型別，支援多值） | 使用者撰寫 |

系統屬性在 frontmatter 中永遠依上述順序排在最前面，接著才是 schema 定義的屬性。這些名稱是保留的，不能在 type schema 或 shared properties 中使用。

**使用者撰寫**的屬性（`name`、`description`、`tags`）可以被 object template 覆蓋。**自動管理**的屬性（`created_at`、`updated_at`）無法被覆蓋——它們永遠反映實際的建立和修改時間。

```
vault/
├── .typemd/
│   ├── types/              # 使用者定義的 type schema（YAML）
│   │   ├── book.yaml
│   │   └── person.yaml
│   ├── properties.yaml     # 共用屬性定義（選填）
│   ├── index.db            # SQLite 索引（自動更新）
│   └── tui-state.yaml      # TUI 會話狀態（自動儲存）
├── templates/              # 依 type 分類的 object template（選填）
│   └── book/
│       └── review.md       # 新 object 的預設 frontmatter 和內文
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

只有 `tag` type 是內建的（支援 `tags` 系統屬性）。內建的 `tag` type 包含複數形式（「tags」）以便在群組標題中正確顯示，並且預設啟用 `unique: true` 以強制 name 唯一性。其他所有 type 都必須透過 `.typemd/types/*.yaml` 檔案定義。

### Object Template

每個 type 可以在 `templates/<type>/<name>.md` 定義一個或多個 template。Template 是一般的 Markdown 檔案（frontmatter + 內文），用於提供建立 object 時的預設內容。

```bash
# 單一 template — 自動套用
tmd object create book clean-code

# 明確指定 template
tmd object create book clean-code -t review

# 多個 template，未指定 flag — 互動式選擇
tmd object create book clean-code
```

Template 的 frontmatter 屬性會覆蓋 schema 預設值。Template 的內文會成為 object 的初始內文。Template 中的自動管理系統屬性（`created_at`、`updated_at`）會被忽略。

### 唯一性約束

Type schema 可透過設定 `unique: true` 來啟用 name 唯一性檢查。啟用後，TypeMD 會阻止建立同一 type 下擁有相同 `name` 值的多個 object。內建的 `tag` type 預設啟用此功能。

```yaml
# .typemd/types/person.yaml
name: person
unique: true  # 同一個人名稱只能有一個
properties:
  - name: role
    type: string
```

唯一性在建立時強制執行，並由 `tmd type validate` 驗證。若 type 沒有設定 `unique: true`（預設值），則允許多個 object 共用相同的 name。

## 索引

TypeMD 使用 SQLite 搭配 FTS5 進行索引。索引儲存在 `.typemd/index.db`，包含：

- Object 中繼資料（Type、檔名、屬性）
- 從 Object 內文提取的 Wiki-link 記錄（用於反向連結追蹤）
- 涵蓋檔名、屬性和內文的全文搜尋索引

開啟 vault 時，若資料庫為空或缺失（例如 clone 後首次開啟），索引會自動同步。使用 TUI 或 CLI 指令時也會保持更新。在 TypeMD 外部手動編輯檔案後，使用 `tmd --reindex` 重建。

## 查詢

TypeMD 提供兩種尋找 Object 的方式：結構化查詢和全文搜尋。

### 結構化查詢

使用 `tmd query` 依屬性篩選 Object。條件使用 `key=value` 格式，以空格分隔（AND 邏輯）。

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

### 全文搜尋

使用 `tmd search` 搜尋檔名、屬性和內文。由 SQLite FTS5 驅動。

```bash
tmd search "concurrency"
tmd search "golang" --json
```

### TUI 搜尋

在 TUI 中，按 `/` 進入搜尋模式。結果會即時篩選。按 `Esc` 清除結果並回到完整列表。
