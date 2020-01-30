package executionlayer

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"path"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc/transforms"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
	sdk "github.com/hdac-io/friday/types"
	"github.com/hdac-io/friday/x/executionlayer/types"
	abci "github.com/hdac-io/tendermint/abci/types"
	"github.com/stretchr/testify/assert"
)

//-------------------------------------------

func TestQueryKeyToBytes(t *testing.T) {
	input := setupTestInput()

	_, err := toBytes("address", "02014a87d1ec490005f85bb4296596ed741411f673a35317543439971c7c7731bb", input.elk.ReadableNameKeeper, input.ctx)
	assert.Nil(t, err)
	_, err = toBytes("address", "invalid address", input.elk.ReadableNameKeeper, input.ctx)
	assert.NotNil(t, err)

	expected := []byte("test-data")

	got, err := toBytes("uref", hex.EncodeToString(expected), input.elk.ReadableNameKeeper, input.ctx)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	_, err = toBytes("hash", hex.EncodeToString(expected), input.elk.ReadableNameKeeper, input.ctx)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	_, err = toBytes("local", hex.EncodeToString(expected), input.elk.ReadableNameKeeper, input.ctx)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)

	_, err = toBytes("invalid key type", "", input.elk.ReadableNameKeeper, input.ctx)
	assert.True(t, strings.Contains(err.Error(), "Unknown QueryKey type:"))
}

func TestUnitHashMapNormalInput(t *testing.T) {
	input := setupTestInput()

	blockHash := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	eeState := []byte{31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	result := input.elk.SetEEState(input.ctx, blockHash, eeState)
	assert.Equal(t, true, result)

	unitHash := input.elk.GetUnitHashMap(input.ctx, blockHash)
	assert.Equal(t, eeState, unitHash.EEState)
}

func TestUnitHashMapInCorrectInput(t *testing.T) {
	input := setupTestInput()

	blockHash := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	eeState := []byte{31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	result := input.elk.SetEEState(input.ctx, blockHash, eeState)
	assert.Equal(t, false, result)

	unitHash := input.elk.GetUnitHashMap(input.ctx, blockHash)
	assert.NotEqual(t, eeState, unitHash.EEState)
}

func TestMustGetProtocolVersion(t *testing.T) {
	expected, err := types.ToProtocolVersion(types.DefaultGenesisState().GenesisConf.Genesis.ProtocolVersion)
	assert.Nil(t, err)

	input := setupTestInput()
	got := input.elk.MustGetProtocolVersion(input.ctx)
	assert.Equal(t, *expected, got)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustGetProtocolVersion below should panic!")
		}
	}()
	input.elk.MustGetProtocolVersion(sdk.Context{})
}

func TestCreateBlock(t *testing.T) {
	input := setupTestInput()
	genesis(input)

	pubkey := sdk.MustGetSecp256k1FromBech32AccPubKey(GenesisPubKeyString)

	// Block #1
	blockHash1 := []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	beginBlockABCI1 := abci.RequestBeginBlock{
		Hash:   blockHash1,
		Header: abci.Header{LastBlockId: abci.BlockID{Hash: input.ctx.CandidateBlock().Hash}},
	}
	BeginBlocker(input.ctx, beginBlockABCI1, input.elk)

	counterDefineMSG := NewMsgExecute(
		GenesisPubKeyString,
		*pubkey,
		util.LoadWasmFile(path.Join(contractPath, counterDefineWasm)),
		[]byte{},
		util.LoadWasmFile(path.Join(contractPath, standardPaymentWasm)),
		util.MakeArgsStandardPayment(new(big.Int).SetUint64(200000000)),
		uint64(10),
		GenesisAccountAddress,
	)
	handlerMsgExecute(input.ctx, input.elk, counterDefineMSG)
	EndBloker(input.ctx, input.elk)

	// Block #2
	blockHash2 := []byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	nextBlockABCI2 := abci.RequestBeginBlock{
		Hash:   blockHash2,
		Header: abci.Header{LastBlockId: abci.BlockID{Hash: blockHash1}},
	}
	BeginBlocker(input.ctx, nextBlockABCI2, input.elk)

	counterCallMSG := NewMsgExecute(
		GenesisPubKeyString,
		*pubkey,
		util.LoadWasmFile(path.Join(contractPath, counterCallWasm)),
		[]byte{},
		util.LoadWasmFile(path.Join(contractPath, standardPaymentWasm)),
		util.MakeArgsStandardPayment(new(big.Int).SetUint64(200000000)),
		uint64(10),
		GenesisAccountAddress,
	)
	handlerMsgExecute(input.ctx, input.elk, counterCallMSG)
	EndBloker(input.ctx, input.elk)

	queryPath := "counter/count"
	arrPath := strings.Split(queryPath, "/")

	unitHash1 := input.elk.GetUnitHashMap(input.ctx, blockHash1)
	pv := input.elk.MustGetProtocolVersion(input.ctx)

	genesisAddr := sdk.MustGetEEAddressFromCryptoPubkey(GenesisPubKey)
	res1, _ := grpc.Query(input.elk.client, unitHash1.EEState, "address", genesisAddr.Bytes(), arrPath, &pv)
	assert.Equal(t, int32(0), res1.GetIntValue())

	unitHash2 := input.elk.GetUnitHashMap(input.ctx, blockHash2)
	res2, _ := grpc.Query(input.elk.client, unitHash2.EEState, "address", genesisAddr.Bytes(), arrPath, &pv)
	assert.Equal(t, int32(1), res2.GetIntValue())
}

func TestTransfer(t *testing.T) {
	input := setupTestInput()
	genesis(input)

	genpubkey := sdk.MustGetSecp256k1FromBech32AccPubKey(GenesisPubKeyString)
	receippubkey := sdk.MustGetSecp256k1FromBech32AccPubKey(RecipientPubKeyString)

	// Block #1
	blockHash1 := []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	nextBlockABCI1 := abci.RequestBeginBlock{
		Hash:   blockHash1,
		Header: abci.Header{LastBlockId: abci.BlockID{Hash: input.ctx.CandidateBlock().Hash}},
	}
	BeginBlocker(input.ctx, nextBlockABCI1, input.elk)

	transferMSG := NewMsgTransfer(
		GenesisPubKeyString,
		*genpubkey,
		*receippubkey,
		util.LoadWasmFile(path.Join(contractPath, transferWasm)),
		util.MakeArgsTransferToAccount(sdk.MustGetEEAddressFromBech32(RecipientPubKeyString).Bytes(), 100000000),
		util.LoadWasmFile(path.Join(contractPath, standardPaymentWasm)),
		util.MakeArgsStandardPayment(new(big.Int).SetUint64(200000000)),
		uint64(200000000),
		GenesisAccountAddress,
	)
	handlerMsgTransfer(input.ctx, input.elk, transferMSG)
	EndBloker(input.ctx, input.elk)

	// Block #2
	blockHash2 := []byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	nextBlockABCI2 := abci.RequestBeginBlock{
		Hash:   blockHash2,
		Header: abci.Header{LastBlockId: abci.BlockID{Hash: blockHash1}},
	}
	BeginBlocker(input.ctx, nextBlockABCI2, input.elk)
	input.ctx = input.ctx.WithBlockHeader(nextBlockABCI2.Header)
	EndBloker(input.ctx, input.elk)

	res, err := input.elk.GetQueryBalanceResultSimple(input.ctx, *receippubkey)
	queriedRes, _ := strconv.Atoi(res)

	assert.Equal(t, queriedRes, 100000000)
	assert.Equal(t, err, nil)

	res2, err := input.elk.GetQueryBalanceResultSimple(input.ctx, *genpubkey)
	queriedRes2, _ := strconv.Atoi(res2)
	fmt.Println(queriedRes)
	fmt.Println(queriedRes2)
}

func TestMarsahlAndUnMarshal(t *testing.T) {
	src := &transforms.TransformEntry{
		Transform: &transforms.Transform{TransformInstance: &transforms.Transform_Write{Write: &transforms.TransformWrite{Value: &state.Value{Value: &state.Value_IntValue{IntValue: 1}}}}}}
	bz, _ := proto.Marshal(src)

	obj := &transforms.TransformEntry{}
	proto.Unmarshal(bz, obj)

	assert.Equal(t, src.Transform.String(), obj.Transform.String())
}

func TestGenesisState(t *testing.T) {
	testMock := setupTestInput()

	expected := types.DefaultGenesisState()
	var got types.GenesisState

	// GenesisConf test
	testMock.elk.SetGenesisConf(testMock.ctx, expected.GenesisConf)
	got.GenesisConf = testMock.elk.GetGenesisConf(testMock.ctx)

	if !reflect.DeepEqual(expected.GenesisConf.WasmCosts, got.GenesisConf.WasmCosts) {
		t.Errorf("expected: %v, but got: %v", expected.GenesisConf.WasmCosts, got.GenesisConf.WasmCosts)
	}
	assert.Equal(t, expected.GenesisConf.Genesis.Timestamp, got.GenesisConf.Genesis.Timestamp)
	assert.Equal(t, expected.GenesisConf.Genesis.ProtocolVersion, got.GenesisConf.Genesis.ProtocolVersion)
	assert.Equal(t, expected.GenesisConf.Genesis.MintWasm, got.GenesisConf.Genesis.MintWasm)
	assert.Equal(t, expected.GenesisConf.Genesis.PosWasm, got.GenesisConf.Genesis.PosWasm)

	// GenesisAccounts test
	expected.Accounts = make([]types.Account, 1)
	pubkey := sdk.MustGetSecp256k1FromBech32AccPubKey(GenesisPubKeyString)
	expected.Accounts[0].PublicKey = *pubkey
	expected.Accounts[0].InitialBalance = "2"
	expected.Accounts[0].InitialBondedAmount = "1"

	testMock.elk.SetGenesisAccounts(testMock.ctx, expected.Accounts)
	gottonAccounts := testMock.elk.GetGenesisAccounts(testMock.ctx)
	if !reflect.DeepEqual(expected.Accounts, gottonAccounts) {
		t.Errorf("expected: %v, but got: %v", expected.Accounts, gottonAccounts)
	}

	// ChainName test
	expected.ChainName = "keeper-test-chain-name"
	testMock.elk.SetChainName(testMock.ctx, expected.ChainName)
	gottonChainName := testMock.elk.GetChainName(testMock.ctx)
	assert.Equal(t, expected.ChainName, gottonChainName)
}

func TestValidator(t *testing.T) {
	input := setupTestInput()

	valPubKeyStr := "fridayvaloperpub1addwnpepqfaxrvy4f95duln3t6vvtd0qd0sdpwfsn3fh9snpnq06w25qualj6vczad0"
	valPubKey, _ := sdk.GetValPubKeyBech32(valPubKeyStr)
	valAddr, _ := sdk.GetEEAddressFromCryptoPubkey(valPubKey)

	consPubKey, _ := sdk.GetConsPubKeyBech32("fridayvalconspub16jrl8jvqq98x7jjxfcm8252pwd4nv6fetpzk6nzx2ddyc3fn0p2rz4mwf44nqjtfga5k5at4xad82sjhx9r9zdfcwuc5uvt90934jjr4d4xk242909rxks28v9erv3jvwfcx2wp4fe8h54fsddu9zar5v3tyknrs8pykk2mw2p29j4n6w455c7j2d3x4ykft9akx6s24gsu8ys2nvayrykqst965z")
	val := types.NewValidator(valAddr, consPubKey, types.Description{
		Website: "https://validator.friday",
		Details: "Test validator",
	}, "0")

	input.elk.SetValidator(input.ctx, valAddr, val)

	resVal, _ := input.elk.GetValidator(input.ctx, valAddr)

	assert.Equal(t, valAddr, resVal.OperatorAddress)
	assert.Equal(t, consPubKey, resVal.ConsPubKey)
	assert.Equal(t, val.Description.Website, resVal.Description.Website)
	assert.Equal(t, val.Description.Details, resVal.Description.Details)

	val.Description.Moniker = "friday"
	input.elk.SetValidatorDescription(input.ctx, valAddr, val.Description)
	assert.Equal(t, "friday", input.elk.GetValidatorDescription(input.ctx, valAddr).Moniker)

	validators := input.elk.GetAllValidators(input.ctx)
	assert.Equal(t, 1, len(validators))
}
