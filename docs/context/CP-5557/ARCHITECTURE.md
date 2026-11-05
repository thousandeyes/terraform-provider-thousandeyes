# CP-5557: Architecture notes

## Repository conventions

- Provider data sources are implemented in `thousandeyes/data_source_*.go` and registered in `thousandeyes/provider.go` under `DataSourcesMap`.
- Schema uses Terraform Plugin SDK v2 `schema.Resource`; reads receive the configured `*client.APIClient` through `meta`.
- API services are created by converting the SDK's shared `client.Common` field to the service type, following `data_source_thousandeyes_agent.go`.
- Account-group context is applied to SDK requests with `SetAidFromContext`.
- Unit tests use `schema.TestResourceDataRaw`, a custom `http.RoundTripper`, and JSON responses to verify API requests and resulting state.
- Examples live in `examples/data-sources/<provider-name>/data-source.tf`; generated data source pages live in `docs/data-sources/`.

## Integration points

- Use `github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointagents` and its `EndpointAgentsAPIService`.
- Call `GetEndpointAgents().AgentName(name)` and add account-group context using `SetAidFromContext`.
- Endpoint Agent listing is cursor-paginated. Read the response `_links.next` cursor and continue requests until no next link remains; use the API's next cursor rather than inferring from page length.
- Treat API errors as read errors. Search all returned items for exact `AgentName` equality before choosing the first match, to guard against API filter behavior.
- Store the matched Endpoint Agent ID as Terraform's data source ID and populate `agent_id` and `agent_name`.

## Layering and scope

- Keep lookup and pagination inside the new data source read path; do not modify SDK models or introduce a new API abstraction for this single operation.
- Keep the exposed schema limited to the required lookup name and computed ID/name fields.
- Register the data source and add its example, generated docs, and focused HTTP-backed unit tests.
