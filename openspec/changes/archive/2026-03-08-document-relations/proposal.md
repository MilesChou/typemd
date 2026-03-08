## Why

Relations are the structured linking mechanism in typemd, allowing objects to be connected through typed properties defined in type schemas. This feature supports single/multiple cardinality, bidirectional linking with inverse properties, type validation, and display integration. The feature is fully implemented but lacks a formal OpenSpec specification documenting its behavioral contract.

## What Changes

- Establish formal specification for the existing relations feature
- No code changes — documentation only

## Capabilities

### New Capabilities

- `object-relations`: Typed relation properties in schemas, linking/unlinking objects, bidirectional relations with inverse, type target validation, single/multiple cardinality, and reverse relation display

### Modified Capabilities

(none)

## Impact

- No code impact
- Creates `openspec/specs/object-relations/spec.md` as a behavioral contract
