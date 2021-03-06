package ovn

import (
	"yunion.io/x/onecloud/pkg/vpcagent/ovnutil"
)

// cmp scans the database for irows.  For those present, mark them with ocver.
// If all rows are found, return true to indicate this.  Otherwise return as
// 2nd value the args to destroy these found records
func cmp(db *ovnutil.OVNNorthbound, ocver string, irows ...ovnutil.IRow) (bool, []string) {
	irowsFound := make([]ovnutil.IRow, 0, len(irows))
	irowsDiff := make([]ovnutil.IRow, 0)

	for _, irow := range irows {
		irowFound := db.FindOneMatchNonZeros(irow)
		if irowFound != nil {
			irowsFound = append(irowsFound, irowFound)
		} else {
			// TODO db.FindByIndex()
			switch row := irow.(type) {
			case *ovnutil.LogicalSwitchPort:
				rowQ := &ovnutil.LogicalSwitchPort{
					Name: row.Name,
				}
				if irow := db.FindOneMatchNonZeros(rowQ); irow != nil {
					irowsDiff = append(irowsDiff, irow)
				}
			}
		}
	}
	// mark them anyway even if not all found, to avoid the destroy
	// call at sweep stage
	for _, irowFound := range irowsFound {
		irowFound.OvnSetExternalIds(externalKeyOcVersion, ocver)
	}
	if len(irowsFound) == len(irows) {
		return true, nil
	}
	args := ovnutil.OvnNbctlArgsDestroy(irowsFound)
	args = append(args, ovnutil.OvnNbctlArgsDestroy(irowsDiff)...)
	return false, args
}
