---
title: Shared Properties
description: 可在多個 type schema 中重複使用的共用屬性定義。
sidebar:
  order: 3.7
---

Shared Properties（共用屬性）讓你定義一次屬性，就能在多個 type schema 中重複使用。不必在每個 type 裡複製相同的屬性定義，只要在一個集中檔案中定義，然後用 `use` 關鍵字引用即可。

## 為什麼需要 Shared Properties？

想像你有一個 `due_date` 屬性，需要出現在 `project`、`task` 和 `milestone` 三個 type 中。沒有 Shared Properties 的話，你得複製三次相同的定義——如果要修改（例如加上 emoji），就得在三個地方都更新。

有了 Shared Properties，你只需要定義一次 `due_date`，然後到處引用：

```yaml
# .typemd/properties.yaml
properties:
  - name: due_date
    type: date
    emoji: "\U0001F4C5"
  - name: priority
    type: select
    options:
      - value: high
      - value: medium
      - value: low
```

## 檔案格式

Shared Properties 定義在 vault 的 `.typemd/` 目錄下的 `.typemd/properties.yaml` 檔案中。檔案使用頂層的 `properties` 陣列：

```yaml
properties:
  - name: due_date
    type: date
    emoji: "\U0001F4C5"
  - name: priority
    type: select
    options:
      - value: high
      - value: medium
      - value: low
  - name: assignee
    type: relation
    target: person
    multiple: true
```

每個項目支援與 type schema 屬性相同的所有欄位：`name`、`type`、`emoji`、`pin`、`options`、`target`、`default`、`multiple`、`bidirectional` 和 `inverse`。

此檔案為選填。如果 `.typemd/properties.yaml` 不存在或沒有 `properties` 陣列，TypeMD 會將其視為空的共用屬性集合。

## 在 type schema 中引用

在 type schema 中，使用 `use` 關鍵字依名稱引用共用屬性：

```yaml
# .typemd/types/project.yaml
name: project
emoji: "\U0001F4CB"
properties:
  - name: title
    type: string
  - use: due_date
  - use: priority
  - name: budget
    type: number
```

載入 type 時，每個 `use` 項目會被解析為 `.typemd/properties.yaml` 中完整的屬性定義。解析後的 type schema 行為與屬性直接寫在裡面完全相同。

你可以混合使用 `use` 項目和一般的屬性定義。`properties` 陣列中的順序會被保留。

## 覆寫欄位

透過 `use` 引用共用屬性時，你可以覆寫兩個欄位：

| 欄位 | 說明 |
|------|------|
| `pin` | 覆寫此 type 中的置頂位置 |
| `emoji` | 覆寫此 type 中的顯示 emoji |

其他所有欄位（`type`、`options`、`target`、`default`、`multiple`、`bidirectional`、`inverse`、`template`）不能覆寫——它們來自共用定義。

### 覆寫範例

```yaml
# .typemd/types/task.yaml
name: task
properties:
  - use: due_date
    pin: 1
    emoji: "\U0001F5D3\uFE0F"
  - use: priority
```

在此範例中，`task` type 使用共用的 `due_date` 屬性，但將它置頂到位置 1 並更換 emoji。`priority` 屬性則直接使用共用定義的原樣。

## 驗證規則

共用屬性的驗證規則與 type schema 屬性相同，另加以下限制：

- **名稱唯一** — 共用屬性之間不能有重複的名稱
- **不可使用保留名稱** — 不能使用 `name` 這個名稱（它是保留的[系統屬性](/zh-tw/concepts/data-model#系統屬性)）
- **型別必須有效** — 屬性型別必須是有效的[預定義屬性型別](/zh-tw/concepts/predefined-property)
- **選項必填** — `select` 和 `multi_select` 屬性必須包含 `options` 陣列
- **目標必填** — `relation` 屬性必須包含 `target` 欄位

當 type schema 使用 `use` 時：

- **`use` 和 `name` 互斥** — 屬性項目不能同時有 `use` 和 `name`
- **引用必須存在** — 引用的共用屬性必須存在於 `.typemd/properties.yaml` 中
- **名稱不能衝突** — 本地屬性不能與引用的共用屬性同名
- **不能重複引用** — 同一個共用屬性不能在同一個 type schema 中被引用兩次
- **僅允許 `pin` 和 `emoji` 覆寫** — 在 `use` 項目上設定其他欄位會產生錯誤

執行 `tmd type validate` 可以檢查所有共用屬性和 type schema 引用。

## Shared Properties vs. 內嵌屬性

| | 共用屬性 | 內嵌屬性 |
|---|---------|---------|
| 定義位置 | `.typemd/properties.yaml` | `.typemd/types/<type>.yaml` |
| 可重複使用 | 是——透過 `use` 引用 | 否——僅限單一 type |
| 可依 type 自訂 | 僅限 `pin` 和 `emoji` | 完全可自訂 |
| 適用情境 | 多個 type 共用的屬性 | 單一 type 獨有的屬性 |

**經驗法則**：如果某個屬性以相同的定義出現在兩個以上的 type 中，就把它做成共用屬性。
