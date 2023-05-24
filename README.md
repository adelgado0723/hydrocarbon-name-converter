# Hydrocarbon Name Converter

A command line tool written in Golang that converts names of simple hydrocarbons from:

- [x] SMILES format to IUPAC nomenclature

### Description

- Accepts the name of a compound in [Simplified Molecular-Input Line-Entry System](https://en.wikipedia.org/wiki/Simplified_molecular-input_line-entry_system) (SMILES) format
- Outputs must produce the name of input compounds in [IUPAC nomenclature](https://en.wikipedia.org/wiki/IUPAC_nomenclature_of_organic_chemistry)
- Developed using TDD 

### Constraints

- Program only handles alkanes (simple hydrocarbons up to 10 in a chain) to greatly reduce the scope and complexity
- Program only handles straight-chain and branched alkanes (no cyclic)
