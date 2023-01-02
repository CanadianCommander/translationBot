package git

import (
	"strconv"
	"time"
)

//==========================================================================
// Public
//==========================================================================

// GenerateNewBranchName creates a new branch name in a standard format
func GenerateNewBranchName() string {
	return GeneratedBranchPrefix + time.Now().Format("2006-01-02_") + strconv.FormatInt(time.Now().Unix(), 10)
}
