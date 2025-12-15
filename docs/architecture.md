# VIO – Core Design Invariants & API Stability

This document defines the core invariants, lifecycle rules, and API stability guarantees of VIO.
These rules are considered architectural constraints and must be preserved across refactors and feature additions.

## API Stability

- The API is considered v1 (unstable but contract-tested)
- Breaking changes require:
  - Updating invariant tests
  - Updating OpenAPI
- DTO fields are additive-only within v1
- Error response shape is stable

## Data Invariants
### Media Files
* A media file is either present or missing, never both.
* A missing media file:
  * Remains queryable via the API.
  * Does not participate in cleanup counts.
  * Can transition back to present if the file is rediscovered at its established path.
* VIO never deletes files from the filesystem. All destructive actions apply only to database records.

### Series / Seasons / Episodes
* A season exists if and only if it has one or more episodes.
* A series exists if and only if it has one or more seasons.
* An episode exists as long as it is linked to a season.

### Movies
* A movie exists if and only if it has at least one media file attached.
* A movie may have multiple media files attached (e.g. multiple versions, qualities, or encodes).
(Future capability — not assumed in the current implementation.)

### Scanning Behavior
* Incremental scans:
  * Never create duplicate entities.
  * Never delete database records.
  * Only process newly discovered files.
* Full rescans may update metadata and presence state but still never delete filesystem content.

## Entity Lifecycles
### Media Files
```
discovered → indexed → missing → (restored | purged)
```
* discovered: file detected on disk
* indexed: metadata stored and linked
* missing: file no longer present on disk
* restored: file reappears at the same path
* purged: database record removed (manual or automated policy)

### Movies
```
created → identified (metadata)
          ↘──────────── active ───────────↗
                    empty → deleted
```
* A movie becomes empty only when it has zero media files attached.
* Empty movies are eligible for deletion.

### Series
```
created → identified (metadata)
          ↘──────────── active ───────────↗
                    empty → deleted
```
* A series becomes empty only when it has zero seasons.

### Seasons
```
created → identified (metadata) → populated → empty → deleted
```
* A season becomes empty when it has zero episodes.

### Episodes
```
created → linked → identified (metadata)
                   ↘──────── active ───────↗
                         orphaned → deleted
```
* An episode becomes orphaned when it has no present media files linked to it.

## Observable vs Internal Data

* All fields exposed via API DTOs are considered part of the public contract unless explicitly documented otherwise.
* Internal fields and implementation details may change freely as long as public contracts remain stable.

## API Stability Guarantees
* Field names do not change casually.
* Endpoint request/response shapes do not change casually.
* HTTP response codes do not change casually.
* Endpoints do not change purpose.
* Breaking changes must be:
  * Explicitly documented.
  * Accompanied by a version bump.
* The API provides core media server functionality:
  * Library creation and management
  * Media discovery and querying
  * Incremental and full scanning
  * Metadata retrieval
  * Media access
