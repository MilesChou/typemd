---
title: Relation
description: Object 之間如何透過帶型別的連結互相連接。
sidebar:
  order: 4
---

Relation（關聯）是兩個 Object 之間的結構化連結。它用來表達一本書有作者、一個專案有成員、或一個想法連結到一場會議。

## 什麼是 Relation？

Relation 是在 Type schema 中定義的 `relation` 類型屬性。不同於普通的超連結或 wiki-link，Relation 具備以下特性：

- **具名** — 它有一個屬性名稱，例如 `author` 或 `books`，而不只是一個 URL
- **帶型別** — 它知道目標 Type（例如 `author` 指向 `person`）
- **可選雙向** — 將書連結到作者時，可以自動將作者連結回這本書

## 單向 vs. 雙向

**單向** Relation 是單方面的連結。在書上設定 `author` 會指向一個人，但那個 person Object 不知道這個連結。

**雙向** Relation 讓兩端保持同步。當你把書的 `author` 連結到一個人，那個人的 `books` 屬性也會自動更新——反之亦然。

```yaml
# book.yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books

# person.yaml
- name: books
  type: relation
  target: book
  multiple: true
  bidirectional: true
  inverse: author
```

## 單值 vs. 多值

- **單值** — 一本書有一個 `author`。重新連結會覆蓋先前的值。
- **多值** — 一個人可以有很多 `books`。連結會附加到清單中。

這由 schema 中的 `multiple` 欄位控制。

## Relation vs. Wiki-link

TypeMD 同時支援 Relation 和 wiki-link（`[[type/slug]]`）。它們的用途不同：

| | Relation | Wiki-link |
|---|----------|-----------|
| 定義於 | Type schema（frontmatter） | Markdown 內文 |
| 結構化 | 是——具名、帶型別、可查詢 | 否——自由格式的行內引用 |
| 雙向 | 依 schema 設定 | 自動追蹤 backlink |
| 使用場景 | 正式連結（作者、專案成員） | 非正式引用（另見、提及） |

用 Relation 處理屬於資料模型的連結。用 wiki-link 處理筆記中的隨意引用。
