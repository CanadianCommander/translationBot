package actions

//==========================================================================
// Public
//==========================================================================

// ProjectSelectProxyRequestDto is sent to actions downstream of a project select action
type ProjectSelectProxyRequestDto struct {
	Project string
	Value   string
}
