### Coding convention with build tags
Project structured with much usage of build tags so the codebase can be compiled with different configurations like:
- `be_json_rpc_evm` for EVM compatible chains

Coding convention:
- Constructor of mandatory struct initialization must be solid across builds.
- `default_*` for structs and functions that are not-specific chain, used to get rid of warning in consumer projects.
- `compatible_*` for structs and functions that are not available on the other types of chains, used to get rid of warning in consumer projects.
- `evm_*` for functions that are EVM specific.