<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the GitHub issue reference in the following format:

* (<tag>) \#<issue-number> message

Tag must include `sql` if having any changes relate to schema

The issue numbers will later be link-ified during the release process,
so you do not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Schema Breaking" for breaking SQL Schema.
"API Breaking" for breaking API.

If any PR belong to multiple types of change, reference it into all types with only ticket id, no need description (convention)

Ref: https://keepachangelog.com/en/1.0.0/
-->

<!--
Templates for Unreleased:

## Unreleased

### Features

### Improvements

### Bug Fixes

### Schema Breaking

### API Breaking
-->

# Changelog

## Unreleased

### Improvements
- (deps) [#8](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/8) Bumps `block-explorer-rpc-cosmos` to v1.1.2

## v1.1.1 - 2024-04-14

### Improvements

- (tx) [#6](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/6) Extract involvers for Wasm txs
- (deps) [#7](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/7) Bumps `block-explorer-rpc-cosmos` to v1.1.1

### API Breaking

- (deps) [#5](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/5) Bumps `block-explorer-rpc-cosmos` to v1.1.0

## v1.0.3 - 2024-04-12

### Improvements

- (deps) [#3](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/3) Bumps `block-explorer-rpc-cosmos` to v1.0.3
- (contract) [#4](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/4) Intercept account info, return contract code id if is contract

## v1.0.2 - 2024-04-05

### Improvements

- (deps) [#1](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/1) Bumps `block-explorer-rpc-cosmos` to v1.0.2
- (contract) [#2](https://github.com/bcdevtools/wasm-block-explorer-rpc-cosmos/pull/2) Limit number of contracts per CW-20 balances query to 50
