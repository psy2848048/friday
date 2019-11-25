package executionlayer

import (
	"github.com/hdac-io/friday/x/executionlayer/types"
)

const (
	ModuleName      = types.ModuleName
	RouterKey       = types.RouterKey
	HashMapStoreKey = types.HashMapStoreKey
	DeployStoreKey  = types.DeployStoreKey
)

var (
	NewMsgExecute  = types.NewMsgExecute
	ModuleCdc      = types.ModuleCdc
	RegisterCodec  = types.RegisterCodec
	NewUnitHashMap = types.NewUnitHashMap

	ErrMalforemdAccountsCsv = types.ErrMalforemdAccountsCsv
	ErrProtocolVersionParse = types.ErrProtocolVersionParse
	ErrInvalidWasmPath      = types.ErrInvalidWasmPath
)

type (
	MsgExecute                = types.MsgExecute
	QueryExecutionLayer       = types.QueryExecutionLayer
	UnitHashMap               = types.UnitHashMap
	QueryExecutionLayerResp   = types.QueryExecutionLayerResp
	QueryExecutionLayerDetail = types.QueryExecutionLayerDetail
	QueryGetBalance           = types.QueryGetBalance
	QueryGetBalanceDetail     = types.QueryGetBalanceDetail
)
