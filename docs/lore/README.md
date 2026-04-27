# Lore And Context Archive

This directory preserves Track B context: public discourse, media claims, government oversight documents, and user-submitted narratives that explain why people ask about Element 115.

Track B is not simulation data. Records here must never seed `data/research.seed.json`, physics constants, or model defaults.

## Labels

- `verified-event`: the event is supported by a reliable source.
- `media-claim`: a source reports or repeats a claim without proving it.
- `unsupported-linkage`: a claim connects facts without sufficient evidence.
- `physics-claim`: a claim can be expressed as a physics assertion for validation.
- `not-for-simulation`: the record is excluded from simulation defaults.

## Layout

| Path | Purpose |
|---|---|
| `intake-notes/` | Quarantined user-submitted material awaiting source review |
| `records/` | Normalized Track B records after source review |
