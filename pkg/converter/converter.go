package converter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/adelgado0723/hw-adelgado0723/pkg/lib"
)

const maxCarbonChainLength = 10
const carbonCharacter = 'C'

type branch struct {
	smiles               string
	attachedCarbonNumber int
	numCarbons           int
}

type iupacName struct {
	smiles              string //SMILES representation of the molecule
	branchMolecules     []branch
	baseMolecule        string
	carbonMoleculeOrder map[int]int // map of carbon index in smiles to carbon index in IUPAC name
}

type GroupedBranches struct {
	length                int
	attachedCarbonIndeces []int
}

type FormattedBranch struct {
	name                 string
	carbonMoleculePrefix string
}

// validateSMILES function validates the SMILES string
func validateSMILES(smiles string) error {
	// validate that the SMILES string only contains "C", "(", and ")"
	if !lib.StrIsComposedOfBytes(smiles, map[byte]bool{carbonCharacter: true, '(': true, ')': true}) {
		return fmt.Errorf("SMILES string contains invalid characters")
	}
	// validate that the parentheses are balanced
	if !lib.ParenthesesAreBalanced(smiles) {
		return fmt.Errorf("Invalid input string. Parentheses are not balanced.")
	}
	return nil
}

// getCarbonChainLength function returns the length of the carbon chain, skipping over branches
func getCarbonChainLength(name string) int {
	var carbonChainLength int
	var inBranch = false

	for i := range name {
		if name[i] == '(' {
			inBranch = true
		} else if name[i] == ')' {
			inBranch = false
		}
		if inBranch {
			continue
		}
		if name[i] == carbonCharacter {
			carbonChainLength++
		}
	}
	return carbonChainLength
}

// groupBranchesByLength function takes a slice of branches and returns a map of grouped branches
func groupBranchesByLength(branches []branch) map[int]GroupedBranches {
	var grouped = make(map[int]GroupedBranches)
	for i := range branches {
		b, ok := grouped[branches[i].numCarbons]
		if !ok {
			b = GroupedBranches{
				length:                branches[i].numCarbons,
				attachedCarbonIndeces: []int{},
			}
			b.attachedCarbonIndeces = append(b.attachedCarbonIndeces, branches[i].attachedCarbonNumber)
			grouped[branches[i].numCarbons] = b
			continue
		}
		b.attachedCarbonIndeces = append(
			grouped[branches[i].numCarbons].attachedCarbonIndeces,
			branches[i].attachedCarbonNumber,
		)
		grouped[branches[i].numCarbons] = b
	}
	return grouped
}

// getFormattedBranchNames function takes a map of grouped branches and returns a slice of IUAPC formatted branch names
func getFormattedBranchNames(groupedBranches map[int]GroupedBranches) []FormattedBranch {
	formatted := make([]FormattedBranch, 0, len(groupedBranches))
	for i := range groupedBranches {

		indeces := lib.RemoveDuplicates(groupedBranches[i].attachedCarbonIndeces)
		var branchIndeces = make([]string, len(indeces))
		for j := range indeces {
			branchIndeces[j] = strconv.Itoa(indeces[j])
		}

		numLikeBranches := len(groupedBranches[i].attachedCarbonIndeces)
		formatted = append(
			formatted,
			FormattedBranch{
				name: fmt.Sprintf(
					"%s-%s%syl",
					strings.Join(branchIndeces, ","),
					lib.MultipleLikeMoleculePrefixes[numLikeBranches],
					lib.CarbonMoleculePrefixes[groupedBranches[i].length],
				),
				carbonMoleculePrefix: lib.CarbonMoleculePrefixes[groupedBranches[i].length],
			},
		)
	}
	return formatted
}

// print method returns the IUPAC name of the molecule
func (iupac *iupacName) print() string {
	var name string
	if len(iupac.branchMolecules) > 0 {
		var groupedBranches = groupBranchesByLength(iupac.branchMolecules)
		formattedGroupedBranchNames := getFormattedBranchNames(groupedBranches)

		// alphabetize the branches not accounting for numerical prefixes
		sort.Slice(formattedGroupedBranchNames, func(i, j int) bool {
			return formattedGroupedBranchNames[i].carbonMoleculePrefix < formattedGroupedBranchNames[j].carbonMoleculePrefix
		})

		branchNames := make([]string, len(formattedGroupedBranchNames))
		for i := range formattedGroupedBranchNames {
			branchNames[i] = formattedGroupedBranchNames[i].name
		}
		var branches = strings.Join(branchNames, "-")
		name = branches + iupac.baseMolecule
	} else {
		name = iupac.baseMolecule
	}
	return name
}

// getAttachedCarbonIndex function returns the index of the carbon atom attached to the branch
// smiles[index] is expected to be either a carbonCharacter or a ')'
func getAttachedCarbonIndex(smiles string, index int) (int, error) {
	curr := index
	for curr > 0 {
		firstCharBeforeParen := curr - 1
		if smiles[firstCharBeforeParen] == carbonCharacter {
			return firstCharBeforeParen, nil
		} else if smiles[firstCharBeforeParen] == ')' {
			// go backwards until we find the previous matching opening parenthesis
			var openParenIndex, err = lib.GetMatchingOpeningParenIndex(smiles, firstCharBeforeParen)
			if err != nil {
				return 0, err
			}
			curr = openParenIndex
		}
	}
	return -1, fmt.Errorf("No attached carbon found")
}

// setBranchData method discovers and sets branch data for the IUPAC name
func (iupac *iupacName) setBranchData() error {
	if iupac.smiles == "" {
		return fmt.Errorf("Cannot set branch data. SMILES string is empty.")
	}
	if len(iupac.carbonMoleculeOrder) < 1 {
		return fmt.Errorf("Cannot set branch data. Carbon molecule order is not set.")
	}

	for i := range iupac.smiles {
		if iupac.smiles[i] == '(' {
			var branch branch
			var firstCharAfterParen = i + 1
			var closeParen, matchingParenErr = lib.GetMatchingClosingParenIndex(iupac.smiles, i)
			if matchingParenErr != nil {
				return matchingParenErr
			}

			branch.smiles = iupac.smiles[firstCharAfterParen:closeParen]
			branch.numCarbons = getCarbonChainLength(branch.smiles)
			if branch.numCarbons > maxCarbonChainLength {
				return fmt.Errorf(
					"Cannot determine carbon molecule for chain longer than %d",
					maxCarbonChainLength,
				)
			}

			var attachedCarbonIndex, getCarbonIndexErr = getAttachedCarbonIndex(iupac.smiles, i)
			if getCarbonIndexErr != nil {
				return getCarbonIndexErr
			}

			carbonMolecule, ok := iupac.carbonMoleculeOrder[attachedCarbonIndex]
			if !ok {
				return fmt.Errorf("Cannot determine carbon molecule order for branch %d", i)
			}
			branch.attachedCarbonNumber = carbonMolecule
			iupac.branchMolecules = append(iupac.branchMolecules, branch)
		}
	}
	return nil
}

// saveCarbonChainOrder method saves the order of the carbon molecule chain
func (iupac *iupacName) saveCarbonChainOrder() error {
	if iupac.smiles == "" {
		return fmt.Errorf("Cannot set carbon molecule order from empty SMILES string")
	}

	var carbonChainIndex int = 1
	var inBranch = false

	for i := range iupac.smiles {
		if iupac.smiles[i] == '(' {
			inBranch = true
		} else if iupac.smiles[i] == ')' {
			inBranch = false
			continue
		}

		if inBranch {
			continue
		}

		if iupac.smiles[i] == carbonCharacter {
			iupac.carbonMoleculeOrder[i] = carbonChainIndex
			carbonChainIndex++
		}
	}
	return nil
}

// setData method gets the IUPAC name by finding the length and location of all branches
func (iupac *iupacName) setData(smiles string) error {
	iupac.smiles = smiles

	var hydrocarbonLength = getCarbonChainLength(smiles)
	var prefix, ok = lib.CarbonMoleculePrefixes[hydrocarbonLength]
	if !ok {
		return fmt.Errorf("Cannot determine name for hydrocarbon of length %d", hydrocarbonLength)
	}
	iupac.baseMolecule = fmt.Sprintf("%sane", prefix)
	err := iupac.saveCarbonChainOrder()
	if err != nil {
		return err
	}

	err = iupac.setBranchData()
	if err != nil {
		return err
	}

	return nil
}

func NewIUPACName() iupacName {
	return iupacName{
		branchMolecules:     make([]branch, 0),
		carbonMoleculeOrder: make(map[int]int),
		baseMolecule:        "",
	}
}

// Convert function converts a SMILES string to an IUPAC name
func Convert(smiles string) (string, error) {
	iupac := NewIUPACName()

	if smiles == "" {
		return "", fmt.Errorf("Cannot set data from empty SMILES string")
	}
	smiles = strings.TrimSpace(smiles)
	smiles = strings.ToUpper(smiles)
	err := validateSMILES(smiles)
	if err != nil {
		return "", err
	}

	err = iupac.setData(smiles)
	if err != nil {
		return "", err
	}

	return iupac.print(), nil
}
