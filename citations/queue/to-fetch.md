# Sources To Fetch

Purpose: source records and source content still needed. Fetching must respect access controls and licenses. Metadata-only access is useful for identity but does not validate data or coefficients.

## Existing DOI/source leads

- Oganessian et al. 2022, "First experiment at the Super Heavy Element Factory: High cross section of 288Mc in the 243Am+48Ca reaction and identification of the new isotope 264Lr", DOI: `10.1103/PhysRevC.106.L031301`.
- Oganessian et al. 2022, "New isotope 286Mc produced in the 243Am+48Ca reaction", DOI: `10.1103/PhysRevC.106.064306`.
- Giuliani et al. 2019, "Colloquium: Superheavy elements: Oganesson and beyond", DOI: `10.1103/RevModPhys.91.011001`.
- Oganessian and Utyonkov 2015, "Super-heavy element research", DOI: `10.1088/0034-4885/78/3/036301`.
- Huang et al. 2021, "The AME 2020 atomic mass evaluation (II). Tables, graphs and references", DOI: `10.1088/1674-1137/abddaf`.
- Kondev et al. 2021, "The NUBASE2020 evaluation of nuclear physics properties", DOI: `10.1088/1674-1137/abddae`.
- Wang et al. 2015, "Systematic study of alpha-decay energies and half-lives of superheavy nuclei", DOI: `10.1103/PhysRevC.92.064301`.
- Royer and Zhang 2008, "Recent alpha decay half-lives and analytic expression predictions including superheavy nuclei", DOI: `10.1103/PhysRevC.77.037602`.
- Hosseini and Hassanabadi 2017, "Theoretical approaches to alpha decay half-lives of super-heavy nuclei", DOI: `10.1088/1674-1137/41/6/064101`.
- Ishizuka et al. 2023, "Nuclear fission properties of super heavy nuclei described within the four-dimensional Langevin model", DOI: `10.3389/fphy.2023.1111868`.
- Yakushev et al. 2024, "Manifestation of relativistic effects in the chemical properties of nihonium and moscovium revealed by gas chromatography studies", DOI: `10.3389/fchem.2024.1474820`.
- IUPAC 2016, "Names and symbols of the elements with atomic numbers 113, 115, 117 and 118", DOI: `10.1515/pac-2016-0501`.

## New arXiv source leads sampled on 2026-04-29

- arXiv:`0803.4151v2`, "Superheavy Elements in the Magic Islands", categories `nucl-th`, `nucl-ex`.
- arXiv:`1111.0505v1`, "New approach for alpha decay half-lives of superheavy nuclei and applicability of WKB approximation", categories `nucl-th`, `nucl-ex`.
- arXiv:`0811.1619v1`, "Alpha decay chains from superheavy nuclei", category `nucl-th`.
- arXiv:`1810.04421v1`, "Study on alpha decay chains of Z = 122 superheavy nuclei with deformation effects and Langer modification", category `nucl-th`.
- arXiv:`1704.06334v1`, "Assessing theoretical uncertainties in fission barriers of superheavy nuclei", category `nucl-th`.
- arXiv:`1902.10108v1`, hyperheavy/nuclear landscape with covariant density functional theory.
- arXiv:`2211.05671v1`, relativistic/nonrelativistic mean-field study for superheavy `Z=120`.

## Evaluated-data targets

- NNDC/NuDat/ENSDF HTML/PDF for `288Mc` adopted and alpha-decay datasets.
- NNDC/NuDat/ENSDF HTML/PDF for `290Mc` adopted and alpha-decay datasets.
- NNDC/NuDat/ENSDF entries for daughter-chain nuclides used in workbook/decay context.

## Open-source project references to fetch or refresh

- Geant4 reference papers and docs.
- OpenMC methods/docs for reproducible Monte Carlo output.
- PyNE nuclear-data handling docs.
- pace_ensdf and nuclei parser documentation/source licenses.
- nuclear chart plotting tools using AME/NUBASE data.
- DeepXDE and SciML ModelingToolkit docs for differential-equation/model metadata lessons.

## Fetch blockers already observed

- APS source pages/PDFs for some alpha-model papers can return HTTP `403`; if blocked, record blocker and do not infer coefficients.
- NNDC direct scripted fetch can return HTTP `429`; browser-readable text may still be used for candidate review.
- NNDC adopted PDF URL may trigger download instead of browser rendering.

## Context-only / not-for-simulation

- House Oversight Committee letter to FBI Director, April 20, 2026, URL: `https://oversight.house.gov/wp-content/uploads/2026/04/FBI-Missing-Scientists-Letter_4.20.26.pdf`.
