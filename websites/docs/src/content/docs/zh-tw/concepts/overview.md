---
title: 概念總覽
description: TypeMD 的關鍵名詞與概念。
sidebar:
  order: 1
---

TypeMD 使用一組核心概念來組織你的知識。本頁提供每個名詞的快速定義，更深入的說明請參閱各概念的獨立頁面。

## Vault

Vault 是由 TypeMD 管理的目錄。它包含 `.typemd/` 設定資料夾和 `objects/` 內容資料夾。每個 Vault 都是獨立的——你可以像一般資料夾一樣搬移、複製或用版本控制管理。

## Object

Object（物件）是 TypeMD 中知識的基本單位。你不需要用「檔案」或「筆記」來思考，而是用 Object——一本書、一個人、一個想法、一場會議。

每個 Object 以 Markdown 檔案儲存，包含 YAML frontmatter（屬性）和自由格式的內文。

[深入了解 Object →](/zh-tw/concepts/objects)

## Type

Type（類型）定義你要建立什麼樣的 Object。每個 Object 只屬於一個 Type（例如 `book`、`person`、`note`）。Type 由 schema 檔案定義，指定 Object 可以擁有哪些屬性。

[深入了解 Type →](/zh-tw/concepts/types)

## Relation

Relation（關聯）是兩個 Object 之間具名且帶型別的連結。不同於簡單的超連結，Relation 定義在 Type schema 中，可以設定為雙向——更新一端會自動更新另一端。

[深入了解 Relation →](/zh-tw/concepts/relations)

## Property

Property（屬性）是 Object 上的具名欄位，由 Type schema 定義。屬性有型別，例如 `string`、`number`、`enum` 或 `relation`，儲存在 Markdown 檔案的 YAML frontmatter 中。

## Slug

Slug 是 Object 檔名中人類可讀的部分，例如 `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md` 中的 `golang-in-action`。Slug 來自你建立 Object 時給的名稱。

## ULID

ULID（Universally Unique Lexicographically Sortable Identifier）是透過 CLI 建立 Object 時自動附加在 slug 後方的 26 字元識別碼，確保即使兩個 Object 的 slug 相同也不會衝突。手動建立的 Object 不需要 ULID。

## Wiki-link

Wiki-link 是在 Markdown 內文中使用 `[[type/slug]]` 語法來引用其他 Object 的方式。TypeMD 會追蹤這些連結及其反向連結（backlink），讓你發現 Object 之間的關聯。

[深入了解 Wiki-link →](/zh-tw/concepts/wiki-links)

## Index

Index（索引）是一個 SQLite 資料庫（`.typemd/index.db`），快取 Object 的中繼資料以供快速查詢和搜尋。它會自動建立和更新，你不需要手動編輯。
