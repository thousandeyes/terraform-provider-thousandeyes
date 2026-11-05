# CP-5557: Endpoint Agent data source

## Goal

Allow Terraform configurations to look up an existing ThousandEyes Endpoint Agent by its name and reference its identifier in other resources.

## User journey

1. A user configures `data "thousandeyes_endpoint_agent"` with the Endpoint Agent's name.
2. The provider queries Endpoint Agents using the API's exact `agentName` filter and follows all pages needed.
3. When a matching agent exists, Terraform state exposes its ID as `agent_id` and its name as `agent_name`.
4. When no matching agent exists, Terraform receives a clear diagnostic naming the requested agent. API failures are returned as diagnostics.

## Acceptance criteria

- Register a `thousandeyes_endpoint_agent` data source.
- Require `agent_name`; return computed `agent_id` and `agent_name`.
- Search Endpoint Agents with the API `agentName` filter, traverse pagination as needed, and verify exact-name equality before selecting the first match.
- Set the Terraform data source ID from the matched Endpoint Agent ID.
- Refresh computed state when a later read resolves the configured name to an Endpoint Agent with a different ID.
- Return a clear no-match error and propagate API errors.
- Include automated coverage for success, drift refresh, pagination, no-match behavior, and API errors, plus a Terraform example and generated provider documentation.

## Out of scope

- Endpoint Agent lifecycle management (CP-5558).
- Endpoint Test data sources.
- Additional Endpoint Agent attributes beyond the ID and name.

## References

- Jira: [CP-5557](https://thousandeyes.atlassian.net/browse/CP-5557)
- Existing lookup: `thousandeyes/data_source_thousandeyes_agent.go`
- SDK operations: `endpointagents/api_endpoint_agents.go`
