---
title: Type
description: Type 如何定義 Object 的結構。
sidebar:
  order: 3
---

Type（類型）是定義 Object 種類和屬性的藍圖。

## 什麼是 Type？

把 Type 想成一個帶有結構的分類。`book` Type 知道書有標題、閱讀狀態和評分。`person` Type 知道人有名字和角色。

Type 以 YAML schema 檔案定義，放在 `.typemd/types/`：

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: status
    type: enum
    values: [to-read, reading, done]
    default: to-read
  - name: rating
    type: number
```

## 為什麼需要 Type

Type 為你的知識庫帶來**一致性**和**可查詢性**：

- **一致性** — 每個 book Object 都有相同的屬性集合，不會出現一本書用 `status`、另一本用 `reading_state` 的情況。
- **可查詢性** — 因為屬性有型別，你可以精確查詢：「顯示所有 status 是 reading 且 rating > 4 的書」。
- **驗證** — TypeMD 會根據 schema 驗證屬性值（例如 enum 值必須在允許清單中），及早發現錯誤。

## 內建 Type

TypeMD 內建三種 Type 讓你快速開始：

| Type | 屬性 |
|------|------|
| `book` | title (string)、status (enum: to-read/reading/done)、rating (number) |
| `person` | name (string)、role (string) |
| `note` | title (string)、tags (string) |

你可以修改這些內建 Type，或建立自己的 Type 來符合你的知識領域。

## Type vs. 標籤

許多工具使用自由格式的標籤（tag）來分類筆記。標籤很彈性但也很脆弱——打錯字就會產生幽靈分類，而且沒有結構可以查詢。

Type 是明確的：一個 `book` 永遠是 `book`，帶有可以篩選、排序和連結的定義屬性。你用一點前期定義的工夫，換來更強大的組織能力。

## 屬性型別

Type schema 中每個屬性都有一個資料型別：

| 型別 | 說明 | 範例 |
|------|------|------|
| `string` | 文字 | `"Go in Action"` |
| `number` | 整數或浮點數 | `42`、`3.14` |
| `enum` | 固定集合中的值 | `"reading"` |
| `relation` | 連結到另一個 Object | `"person/alan-donovan"` |
