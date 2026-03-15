---
title: 概念總覽
description: TypeMD 的設計理念。
sidebar:
  order: 1
---

TypeMD 是一個 local-first 的知識管理工具，圍繞一個核心理念打造：**用物件思考，而非檔案**。

## 物件，不是檔案

大多數筆記工具用資料夾中的檔案來組織內容。TypeMD 採取不同的方式——你的知識由**物件**組成。一本書、一個人、一個想法、一場會議——每個都是具有定義結構和型別連結的物件。

檔案系統只是儲存層。你操作的是物件，不是檔名。

## 用類型定義結構

每個物件都屬於一個**類型**，類型定義了它的結構。`book` 有 `author` 和 `status`；`person` 有 `role` 和 `company`。類型以簡單的 YAML schema 檔案定義——你可以設計適合自己思路的結構。

## 用關聯和連結建立聯繫

物件之間有兩種連結方式：

- **關聯（Relation）**——定義在 schema 中的結構化連結。書的 `author` 關聯指向一個 `person` 物件。關聯可以是雙向的：更新一端會自動更新另一端。
- **連結（Link）**——在 Markdown 內文中使用 `[[type/slug-ulid]]` 語法的行內引用。比關聯輕量，適合隨意的連接。TypeMD 會自動追蹤反向連結。

## Local-first，原生 Markdown

所有內容以純 Markdown 檔案搭配 YAML frontmatter 儲存。不依賴雲端，沒有鎖定。你可以用任何工具讀取、編輯和版本控制你的知識庫。

## 深入了解

- [物件](/zh-tw/concepts/objects)——物件是什麼、如何運作
- [類型](/zh-tw/concepts/types)——用 schema 定義結構
- [關聯](/zh-tw/concepts/relations)——物件之間的結構化連結
- [連結](/zh-tw/concepts/links)——行內引用和反向連結
- [術語表](/zh-tw/concepts/glossary)——所有關鍵術語的定義
