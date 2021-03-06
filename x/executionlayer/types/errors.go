package types

import (
	"fmt"
	"strings"

	sdk "github.com/hdac-io/friday/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodePublicKeyDecode            sdk.CodeType = 101
	CodeProtocolVersionParse       sdk.CodeType = 102
	CodeTomlParse                  sdk.CodeType = 103
	CodeInvalidValidator           sdk.CodeType = 201
	CodeInvalidDelegation          sdk.CodeType = 202
	CodeInvalidInput               sdk.CodeType = 203
	CodeInvalidAddress             sdk.CodeType = sdk.CodeInvalidAddress
	CodeGRpcExecuteMissingParent   sdk.CodeType = 301
	CodeGRpcExecuteDeployGasError  sdk.CodeType = 302
	CodeGRpcExecuteDeployExecError sdk.CodeType = 303
)

// ErrPublicKeyDecode is an error
func ErrPublicKeyDecode(codespace sdk.CodespaceType, publicKey string) sdk.Error {
	return sdk.NewError(
		codespace, CodePublicKeyDecode, "Could not decode public key as Base64 : %v", publicKey)
}

// ErrProtocolVersionParse is an error
func ErrProtocolVersionParse(codespace sdk.CodespaceType, protocolVersion string) sdk.Error {
	return sdk.NewError(
		codespace, CodeProtocolVersionParse,
		"Could not parse Protocol Version : %v", protocolVersion)
}

// ErrTomlParse is an error
func ErrTomlParse(codespace sdk.CodespaceType, keyString string) sdk.Error {
	return sdk.NewError(
		codespace, CodeTomlParse,
		"Could not parse Toml with : %v", keyString)
}

func ErrNilValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "validator address is nil")
}

func ErrBadValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, "validator address is invalid")
}

func ErrDescriptionLength(codespace sdk.CodespaceType, descriptor string, got, max int) sdk.Error {
	msg := fmt.Sprintf("bad description length for %v, got length %v, max is %v", descriptor, got, max)
	return sdk.NewError(codespace, CodeInvalidValidator, msg)
}

func ErrNilDelegatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "delegator address is nil")
}

func ErrBadDelegationAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "unexpected address length for this (address, validator) pair")
}

func ErrBadDelegationAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "amount must be > 0")
}

func ErrGRpcExecuteMissingParent(codespace sdk.CodespaceType, hash string) sdk.Error {
	return sdk.NewError(codespace, CodeGRpcExecuteMissingParent, "execution engine - missing parent state %s", hash)
}

func ErrGRpcExecuteDeployGasError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeGRpcExecuteDeployGasError, "execution engine - deploy error - gas")
}

func ErrGRpcExecuteDeployExecError(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeGRpcExecuteDeployExecError, "execution engine - deploy error - execute : ", msg)
}

func ErrValidatorOwnerExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator already exist for this operator address, must use new validator operator address")
}

func ErrValidatorPubKeyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator already exist for this pubkey, must use new validator pubkey")
}

func ErrValidatorPubKeyTypeNotSupported(codespace sdk.CodespaceType, keyType string, supportedTypes []string) sdk.Error {
	msg := fmt.Sprintf("validator pubkey type %s is not supported, must use %s", keyType, strings.Join(supportedTypes, ","))
	return sdk.NewError(codespace, CodeInvalidValidator, msg)
}
