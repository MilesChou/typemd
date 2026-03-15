---
title: 查詢
description: 依屬性值對 Object 進行結構化篩選。
sidebar:
  order: 5
---

查詢（Query）讓你使用結構化條件依屬性值篩選 Object。不同於[搜尋](/zh-tw/basics/search)會在所有內容中進行全文比對，查詢針對特定屬性以精確值進行篩選。

## 運作方式

查詢使用 `key=value` 格式來比對 Object。提供多個條件時，會以 AND 邏輯組合——只有符合所有條件的 Object 才會被回傳。

`type` 是一個特殊的篩選鍵，比對 Object 的 type 名稱。其他所有鍵都比對 frontmatter 屬性值。

## CLI 用法

```bash
# 依 type 篩選
tmd query "type=book"

# 依 type 和屬性篩選
tmd query "type=book status=reading"

# 以 JSON 格式輸出
tmd query "type=book" --json
```

| 選項 | 說明 |
|------|------|
| `--json` | 以 JSON 格式輸出結果 |

## 範例

| 查詢 | 比對結果 |
|------|----------|
| `type=book` | 所有 book Object |
| `type=book status=reading` | status 為「reading」的書 |
| `type=person role=author` | role 為「author」的人 |

## 相關頁面

- [搜尋](/zh-tw/basics/search) — 跨所有內容的全文搜尋
- [tmd query](/zh-tw/tui/query) — CLI 參考
- [資料模型](/zh-tw/developers/data-model#查詢架構) — 查詢的詳細說明
