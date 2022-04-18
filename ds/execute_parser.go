package ds

import (
	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/ethereum/go-ethereum/common"
)

type ExecuteParams struct {
	SessionId     string
	Protocol      common.Address
	CreditAccount common.Address
	Borrower      common.Address
	Index         uint
	BlockNumber   int64
}

type FuncWithMultiCall struct {
	Name          string `json:"name"`
	MultiCallsLen int    `json:"len"`
}

type ExecuteParserI interface {
	GetExecuteCalls(txHash, creditManagerAddr string, paramsList []ExecuteParams) []*KnownCall
	GetMainEventLogs(txHash, creditFacade string) []*FuncWithMultiCall
	GetTransfers(txHash string, borrower, account, underlyingToken string, owner []string) core.Transfers
}

type KnownCall struct {
	// Input string
	Depth     uint8          `json:"depth"`
	Name      string         `json:"name"`
	Args      *core.Json     `json:"args"`
	Transfers core.Transfers `json:"transfers"`
}
