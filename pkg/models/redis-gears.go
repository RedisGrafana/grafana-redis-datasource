package models

/**
 * RedisGears Commands
 */
const (
	GearsPyStats           = "rg.pystats"
	GearsDumpRegistrations = "rg.dumpregistrations"
	GearsPyExecute         = "rg.pyexecute"
	GearsPyDumpReqs        = "rg.pydumpreqs"
)

/**
 * RG.PYSTATS Radix marshaling
 */
type PyStats struct {
	TotalAllocated int64 `redis:"TotalAllocated"`
	PeakAllocated  int64 `redis:"PeakAllocated"`
	CurrAllocated  int64 `redis:"CurrAllocated"`
}

/**
 * RG.DUMPREGISTRATIONS Radix marshaling
 */
type DumpRegistrations struct {
	ID               string           `redis:"id"`
	Reader           string           `redis:"reader"`
	Desc             string           `redis:"desc"`
	RegistrationData RegistrationData `redis:"RegistrationData"`
	PD               string           `redis:"PD"`
}

/**
 * Registration data for RG.DUMPREGISTRATIONS Radix marshaling
 */
type RegistrationData struct {
	Mode         string                 `redis:"mode"`
	NumTriggered int64                  `redis:"numTriggered"`
	NumSuccess   int64                  `redis:"numSuccess"`
	NumFailures  int64                  `redis:"numFailures"`
	NumAborted   int64                  `redis:"numAborted"`
	LastError    string                 `redis:"lastError"`
	Args         map[string]interface{} `redis:"args"`
	Status       string                 `redis:"status"`
}

/**
 * RG.PYDUMPREQS Radix marshaling
 */
type PyDumpReq struct {
	GearReqVersion int64       `redis:"GearReqVersion"`
	Name           string      `redis:"Name"`
	IsDownloaded   string      `redis:"IsDownloaded"`
	IsInstalled    string      `redis:"IsInstalled"`
	CompiledOs     string      `redis:"CompiledOs"`
	Wheels         interface{} `redis:"Wheels"`
}
