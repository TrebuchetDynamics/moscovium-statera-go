# Research Seed Data

`research.seed.json` is a small bootstrap dataset, not a complete evaluated nuclear data library.

Each record must preserve:

- isotope identifier
- evaluated half-life and decay values
- daughter isotope
- evaluated-data URL
- DOI trail for peer-reviewed source material

Top-level `notes` are required planning/provenance context for the seed file. They must be non-empty text and must not be treated as scientific data defaults, evaluated measurements, model coefficients, or citation substitutes.

Do not add data from informal summaries, secondary websites, or remembered values unless the record is clearly marked as an intake candidate and excluded from simulation defaults.

