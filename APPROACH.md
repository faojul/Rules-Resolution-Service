
```markdown
# Approach

## 1. Problem Understanding

The problem models a multi-dimensional configuration system where behavior varies based on:

- state
- client
- investor
- caseType

Each override applies conditionally using these dimensions, and resolution is determined using specificity rules.

---

## 2. Schema Design

### Selector Modeling

I modeled the selector as a `map[string]string` stored in JSON format.

#### Why this approach?

- Flexible for adding new dimensions
- Avoids schema explosion (columns for each dimension)
- Aligns with dynamic business requirements

#### Trade-offs

| Approach | Pros | Cons |
|---------|------|------|
| JSON (chosen) | flexible, extensible | slower queries |
| Columns | fast filtering | rigid, harder to evolve |

---

## 3. Resolution Algorithm

Steps:

1. Fetch overrides from repository
2. Filter:
   - Matches selector
   - Active status
   - Within effective date range
3. Sort:
   - Specificity DESC
   - EffectiveDate DESC
4. Select top candidate

Each trait is resolved independently.

### Complexity

- Filtering: O(N)
- Sorting: O(N log N)

---

## 4. Conflict Detection

Two overrides conflict if:

- Same stepKey and traitKey
- Same specificity
- Overlapping effective date ranges
- Selector overlap (matching dimensions)

Currently implemented using O(N²) pairwise comparison.

---

## 5. Edge Case Handling

| Case | Handling |
|------|--------|
| Equal specificity | latest effectiveDate wins |
| Same specificity + same date | conflict |
| Draft overrides | ignored |
| Expired overrides | ignored |
| Empty selector | treated as wildcard |

---

## 6. Design Decisions

- Followed clean architecture to maintain loose coupling
- Business logic kept in the service layer
- Repository abstracts DB access
- Resolver handles rule evaluation

---

## 7. Limitations & Future Improvements

If given more time, I would:

- Move filtering logic to the database layer for scalability
- Improve the explain API to include a full reasoning trace
- Add scenario-based testing using provided test_scenarios.json
- Implement caching for frequently accessed configurations
- Improve conflict detection efficiency

---

## 8. Nice-to-Have Features

Implemented:

- Bulk resolution (partial)
- Conflict detection endpoint

Planned:

- Simulation API
- Diff between contexts

---

## 9. Conclusion

The system provides a flexible and extensible way to manage complex configuration using specificity-based resolution. It balances simplicity and correctness while leaving room for scalability improvements.
