# CP-5557 implementation plan

## Tasks

1. **Add focused data source tests** — complete
   - File: `thousandeyes/data_source_thousandeyes_endpoint_agent_test.go`
   - Cover successful exact-name lookup and state values, name-filter and account-group query parameters, pagination to a later page, no exact match, API error propagation, and drift refresh when the same configured name resolves to a new agent ID.
   - Verify tests fail before implementing the data source.

2. **Implement and register the data source** — complete
   - Add `thousandeyes/data_source_thousandeyes_endpoint_agent.go`.
   - Define required `agent_name` and computed `agent_id` schema fields.
   - Query Endpoint Agents with `agentName`, preserve account-group context, follow `_links.next` cursors, verify exact equality, set state, and report errors.
   - Register `thousandeyes_endpoint_agent` in `thousandeyes/provider.go`.

3. **Add usage example** — complete
   - Add `examples/data-sources/thousandeyes_endpoint_agent/data-source.tf` demonstrating lookup and use of `agent_id`.

4. **Generate provider documentation** — complete
   - Add `docs/data-sources/endpoint_agent.md` using the repository's Terraform Plugin Docs workflow.
   - Include success usage and schema for the required and computed attributes.

5. **Review and validate** — complete
   - Run the focused unit test and the repository's relevant validation/lint commands.
   - Review the diff against this spec and architecture notes; resolve any gaps before completion.

## Dependency graph

```text
1 (tests) -> 2 (implementation and registration) -> 3 (example) -> 4 (docs) -> 5 (review/validation)
```

## Definition of done

- All acceptance criteria in `SPEC.md` are met.
- Tests cover success, pagination, drift refresh, no exact match, and API errors.
- The example and generated data source documentation are present.
- Focused validation passes and the final diff is reviewed.
