## Key Design Decisions

- Used column-based selector instead of JSON for efficient SQL filtering
- Specificity computed dynamically from non-null fields
- Resolution algorithm implemented as deterministic function

## Algorithm

1. Fetch overrides by step+trait
2. Filter by:
   - matching selector
   - active status
   - effective date
3. Sort by:
   - specificity DESC
   - effectiveDate DESC
4. Pick winner

## Edge Cases

- equal specificity conflict
- expired overrides ignored
- draft overrides ignored
- null dimension = wildcard

## Improvements

- Redis caching
- batch resolution
- explain API
- conflict detection optimization