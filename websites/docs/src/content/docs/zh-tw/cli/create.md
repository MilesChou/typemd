---
title: tmd object create
description: 從 Type schema 建立新 Object。
sidebar:
  order: 2.5
---

根據 Type schema 建立新的 Object 檔案（Markdown + YAML frontmatter）。

```bash
tmd object create book clean-code
tmd object create person robert-martin
```

指令會在 `objects/<type>/` 下產生 `.md` 檔案，所有 schema 定義的屬性設為預設值（若無預設值則為 `null`）。Object 同時會被加入 SQLite 索引。

每個 Object 會被分配唯一的 ULID 後綴，因此同一 Type 下可以有相同名稱的物件（例如兩個 `book/clean-code-<ulid>` 但 ULID 不同）。如果 Type 不存在，會回傳錯誤。

## 單一參數用法

如果 type schema 定義了 name template，可以省略第二個參數（slug）。TypeMD 會從 name template 自動產生名稱：

```bash
# 使用 name template 自動產生名稱
tmd object create meeting
```

若 type 沒有定義 name template，則必須提供 slug 參數。

## Template

使用 `-t` / `--template` flag 指定要套用的 object template：

```bash
# 指定 template
tmd object create book clean-code -t review
tmd object create book clean-code --template review
```

Template 檔案位於 `templates/<type>/<name>.md`。如果該 type 只有一個 template，建立 object 時會自動套用，不需要 `-t` flag。如果有多個 template 且未指定 `-t`，會出現互動式選擇提示。

Template 的 frontmatter 屬性會覆蓋 schema 預設值，template 的內文會成為 object 的初始內文。自動管理的系統屬性（`created_at`、`updated_at`）在 template 中會被忽略。
