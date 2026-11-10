#!/usr/bin/env bash
set -euo pipefail

: "${TE_ENDPOINT_AGENT_ID:?Set TE_ENDPOINT_AGENT_ID to an existing Endpoint Agent ID}"

terraform import thousandeyes_endpoint_agent.manual "$TE_ENDPOINT_AGENT_ID"
