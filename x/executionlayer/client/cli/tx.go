package cli

import (
	"math/big"
	"os"
	"regexp"
	"strconv"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/util"
	"github.com/hdac-io/friday/codec"
	"github.com/hdac-io/friday/x/executionlayer/types"

	"github.com/hdac-io/friday/client/context"
	sdk "github.com/hdac-io/friday/types"
	"github.com/hdac-io/friday/x/auth"
	"github.com/hdac-io/friday/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdTransfer is the CLI command for transfer
func GetCmdTransfer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [token_contract_address] [from_address] [to_address] [amount] [fee] [gas_price]",
		Short: "Transfer token",
		Args:  cobra.ExactArgs(6), // # of arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[1]).WithCodec(cdc)

			tokenOwnerAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			fromAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			toAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			toPublicKey := types.ToPublicKey(toAddress)
			amount, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			fee, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}
			gasPrice, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			transferCode := util.LoadWasmFile(os.ExpandEnv("$HOME/.nodef/contracts/transfer_to_account.wasm"))
			transferAbi := util.MakeArgsTransferToAccount(toPublicKey, amount)
			paymentCode := util.LoadWasmFile(os.ExpandEnv("$HOME/.nodef/contracts/standard_payment.wasm"))
			paymentAbi := util.MakeArgsStandardPayment(new(big.Int).SetUint64(fee))

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgTransfer(tokenOwnerAddress, fromAddress, toAddress, transferCode, transferAbi, paymentCode, paymentAbi, gasPrice)
			txBldr = txBldr.WithGas(gasPrice)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBonding is the CLI command for bonding
func GetCmdBonding(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bond [address] [amount] [fee] [gas_price]",
		Short: "Create and sign a bonding tx",
		Args:  cobra.ExactArgs(4), // # of arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			coins, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			fee, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			gasPrice, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			bondingCode := util.LoadWasmFile(os.ExpandEnv("$HOME/.nodef/contracts/bonding.wasm"))
			bondingAbi := util.MakeArgsBonding(coins)
			paymentCode := util.LoadWasmFile(os.ExpandEnv("$HOME/.nodef/contracts/standard_payment.wasm"))
			paymentAbi := util.MakeArgsStandardPayment(new(big.Int).SetUint64(fee))

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgExecute([]byte{0}, cliCtx.FromAddress, cliCtx.FromAddress, bondingCode, bondingAbi, paymentCode, paymentAbi, gasPrice)
			txBldr = txBldr.WithGas(gasPrice)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUnbonding is the CLI command for unbonding
func GetCmdUnbonding(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unbond [address] [amount] [fee] [gas_price]",
		Short: "Create and sign a unbonding tx",
		Args:  cobra.ExactArgs(4), // # of arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			coins, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			fee, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			gasPrice, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			unbondingCode := util.LoadWasmFile(os.ExpandEnv("$HOME/.nodef/contracts/unbonding.wasm"))
			unbondingAbi := util.MakeArgsUnBonding(coins)
			paymentCode := util.LoadWasmFile(os.ExpandEnv("$HOME/.nodef/contracts/standard_payment.wasm"))
			paymentAbi := util.MakeArgsStandardPayment(new(big.Int).SetUint64(fee))

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgExecute([]byte{0}, cliCtx.FromAddress, cliCtx.FromAddress, unbondingCode, unbondingAbi, paymentCode, paymentAbi, gasPrice)
			txBldr = txBldr.WithGas(gasPrice)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func BuildCreateValidatorMsg(cliCtx context.CLIContext) (sdk.Msg, error) {
	amounstStr := viper.GetString(FlagAmount)
	re := regexp.MustCompile("[0-9]+")
	coins, err := strconv.ParseUint(re.FindAllString(amounstStr, -1)[0], 10, 64)
	if err != nil {
		return types.MsgCreateValidator{}, err
	}

	valAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)

	pk, err := sdk.GetConsPubKeyBech32(pkStr)
	if err != nil {
		return types.MsgCreateValidator{}, err
	}

	description := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagDetails),
	)

	msg := types.NewMsgCreateValidator(sdk.ValAddress(valAddr), pk, coins, description)

	return msg, nil
}
