/*
 * 
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ampchain/go-amp/utxo"
)

// NativeInvokeCommand native invoke cmd
type NativeInvokeCommand struct {
	cli *Cli
	cmd *cobra.Command

	args       string
	account    string
	fee        string
	methodName string
}

// NewNativeInvokeCommand new native invoke cmd
func NewNativeInvokeCommand(cli *Cli) *cobra.Command {
	c := new(NativeInvokeCommand)
	c.cli = cli
	c.cmd = &cobra.Command{
		Use:     "invoke",
		Short:   "[Deprecated] Invoke a native method.",
		Example: c.example(),
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.TODO()
			return c.invoke(ctx, args[0])
		},
	}
	c.addFlags()
	return c.cmd
}

func (c *NativeInvokeCommand) addFlags() {
	c.cmd.Flags().StringVarP(&c.args, "args", "a", "{}", "invoke args")
	c.cmd.Flags().StringVarP(&c.account, "account", "", "", "account name")
	c.cmd.Flags().StringVar(&c.fee, "fee", "", "fee of one tx")
	c.cmd.Flags().StringVarP(&c.methodName, "method", "", "invoke", "method name")
}

func (c *NativeInvokeCommand) example() string {
	return `
AmpChain native invoke $codename  -a '{"to":"abc"} --method invoke or increase or others'
`
}

func (c *NativeInvokeCommand) invoke(ctx context.Context, codeName string) error {
	ct := &CommTrans{
		Fee:          c.fee,
		FrozenHeight: 0,
		Version:      utxo.TxVersion,
		From:         c.account,
		ModuleName:   "native",
		ContractName: codeName,
		MethodName:   c.methodName,
		Args:         make(map[string][]byte),
		IsQuick:      false,
		ChainName:    c.cli.RootOptions.Name,
		Keys:         c.cli.RootOptions.Keys,
		AChainClient: c.cli.AChainClient(),
		CryptoType:   c.cli.RootOptions.CryptoType,
	}

	// generate preExe params
	args := make(map[string]interface{})
	err := json.Unmarshal([]byte(c.args), &args)
	if err != nil {
		return err
	}
	ct.Args, err = convertToAmper3Args(args)
	if err != nil {
		return err
	}

	return ct.Transfer(ctx)
}

func convertToAmper3Args(args map[string]interface{}) (map[string][]byte, error) {
	argmap := make(map[string][]byte)
	for k, v := range args {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("bad key %s, expect string value, got %v", k, v)
		}
		argmap[k] = []byte(s)
	}
	return argmap, nil
}
