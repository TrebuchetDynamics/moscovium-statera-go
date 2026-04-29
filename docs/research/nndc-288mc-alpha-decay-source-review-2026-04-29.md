# Source Review: NNDC/NuDat 288Mc Alpha Decay Dataset

Date: 2026-04-29 12:17:39 CST -0600
Reviewer: Riju
Target URL: `https://www.nndc.bnl.gov/nudat3/getdecaydataset.jsp?dsid=288mc+a+decay+%280.17+s%29&nucleus=284NH`
Source class: evaluated-data
Access status: `source-content-readable-via-browser`, direct scripted fetch rate-limited

This document records browser-harness source evidence only. It does not change Statera engine data, isotope values, coefficients, defaults, or accepted records.

## Access evidence

Browser tool result:

- `browser_navigate` status: success
- Browser title: `ENSDF decay dataset`
- Browser snapshot: reported empty page, but DOM body text was readable through browser console.
- `document.body.innerText` returned ENSDF-style dataset text.

PDF target behavior:

- URL tested: `https://www.nndc.bnl.gov/ensnds/288/Mc/adopted.pdf`
- Browser result: navigation failed because a download started.
- Interpretation: not a rendered browser page in this harness pass; treat as downloadable source candidate requiring separate file-capture policy before acceptance.

Direct scripted fetch behavior:

- Method: Python `urllib.request.urlopen` against the NuDat dataset URL.
- Result: HTTP `429` rate limit.
- Interpretation: browser-accessible page should be used carefully; avoid rapid scripted polling.

## Browser-readable excerpt

The browser DOM text began with:

```text
284NH    288MC A DECAY (0.17 S)        2004OG03,2013OG01,2013RU1119NDS    201902
284NH  H TYP=FUL$AUT=BALRAJ SINGH$CIT=NDS 156, 148 (2019)$CUT=31-Jan-2019$
284NH c  See {+288}Mc Adopted Levels for details of production of the isotope.
284NH c  Possible decay scheme with |g transition(s) is proposed by 2015Ru11.
284NH c  2015Ru11 propose three scenarios for the decay scheme of {+288}Mc:
284NH2c  1. 100% |a decay to a level at 105 keV, which then decays by 105-keV
284NH3c  transition to the g.s. 2. 25% |a decay to each of the four excited
284NH4c  states, one at 105 keV (as in scenario 1), second at 290 keV,
284NH5c  and two levels between the 105- and 290-keV levels, with a cascade of
284NH6c  four |g transitions, but only the 105-keV transition is possibly seen
284NH7c  in 2013Ru11. 3. Same as scenario 2, except that 100% |a decay to the
284NH8c  proposed 290-keV level.
288MC  P 0                             0.17 S    2              10750     50
288MC cP T$From {+288}Mc Adopted Levels
288MC cP QP$From 2017Wa10. Note that experimental Q(|a)=10.65 MeV {I1} in
288MC2cP 2017Og01 and 10.63 MeV {I1} in 2015Og05.
284NH  N                       1.0     AP
284NH  L 0                               0.97 S  +12-10
284NH cL T$from Adopted Levels
284NH  L 105       1                                                           ?
284NH cL $Tentative level proposed by 2015Ru11
284NH cL E$in order to explain photon events, 2013Ru11 argue that {+288}Mc
284NH2cL decays by one |a branch to an excited state in {+284}Nh, which then
284NH3cL cascades down to the ground state of {+284}Nh via a few highly
284NH4cL converted |g-transitions with energies below the K-shell binding
284NH5cL energy of Nh atom.
284NH  A 10.3E3    1  100                                                      ?
```

## Candidate evidence table

These entries are candidate observations from readable browser text. They are not newly accepted values.

| Datum | Candidate value | Unit | Uncertainty / qualifier | Source location | Status |
| --- | ---: | --- | --- | --- | --- |
| Parent dataset | `288MC A DECAY (0.17 S)` | text | ENSDF dataset label | first line | candidate-only |
| Parent half-life text | `0.17 S    2` | seconds | ENSDF uncertainty text `2`; needs parser interpretation before acceptance | `288MC P` line | candidate-only |
| Q value text | `10750     50` | likely keV in ENSDF format | `50`; requires parser/unit confirmation | `288MC P` line | candidate-only |
| Experimental Q-alpha note | `10.65 MeV {I1}` and `10.63 MeV {I1}` | MeV | `{I1}` notation requires interpretation | `288MC cP QP` comment lines | candidate-only |
| Normalization qualifier | `AP` | ENSDF qualifier | approximate/qualifier semantics need preservation | `284NH N` line | candidate-only |
| Daughter ground-state half-life text | `0.97 S +12-10` | seconds | asymmetric uncertainty text | `284NH L 0` line | candidate-only |
| Tentative level | `105 1 ?` | likely keV | tentative marker `?` | `284NH L 105` line | candidate-only |

## Candidate use in Statera

Allowed use now:

- browser-harness proof that NuDat/ENSDF HTML can expose readable source text;
- ENSDF candidate schema design;
- parser test fixture planning;
- source review documentation.

Not allowed use now:

- changing `data/research.seed.json`;
- changing accepted `Q_alpha`, half-life, or uncertainty values;
- changing model coefficients;
- treating this browser excerpt as a parsed accepted dataset.

Required tests before acceptance:

1. Candidate parser preserves raw text.
2. Candidate parser preserves ENSDF qualifiers such as `AP` and `?`.
3. Candidate parser distinguishes central value, symmetric uncertainty, asymmetric uncertainty, and approximate/limit notation.
4. Candidate record requires source URL and retrieval/access evidence.
5. Accepted record conversion requires human/test review and cannot happen directly from browser text.

## Browser harness decision

Decision: `candidate-for-manual-review`

Reason: source text is readable through browser DOM, but direct scripted access returned HTTP `429`, and ENSDF notation requires a dedicated candidate parser before any value can be accepted.

## Next action

Implement candidate-only ENSDF source registration and parser fixtures using this page as a source-review reference, not as accepted engine data.
