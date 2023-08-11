// Code generated by "make cli"; DO NOT EDIT.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package accountscmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/accounts"
	"github.com/hashicorp/boundary/internal/cmd/base"
	"github.com/hashicorp/boundary/internal/cmd/common"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

func initLdapFlags() {
	flagsOnce.Do(func() {
		extraFlags := extraLdapActionsFlagsMapFunc()
		for k, v := range extraFlags {
			flagsLdapMap[k] = append(flagsLdapMap[k], v...)
		}
	})
}

var (
	_ cli.Command             = (*LdapCommand)(nil)
	_ cli.CommandAutocomplete = (*LdapCommand)(nil)
)

type LdapCommand struct {
	*base.Command

	Func string

	plural string

	extraLdapCmdVars
}

func (c *LdapCommand) AutocompleteArgs() complete.Predictor {
	initLdapFlags()
	return complete.PredictAnything
}

func (c *LdapCommand) AutocompleteFlags() complete.Flags {
	initLdapFlags()
	return c.Flags().Completions()
}

func (c *LdapCommand) Synopsis() string {
	if extra := extraLdapSynopsisFunc(c); extra != "" {
		return extra
	}

	synopsisStr := "account"

	synopsisStr = fmt.Sprintf("%s %s", "ldap-type", synopsisStr)

	return common.SynopsisFunc(c.Func, synopsisStr)
}

func (c *LdapCommand) Help() string {
	initLdapFlags()

	var helpStr string
	helpMap := common.HelpMap("account")

	switch c.Func {

	default:

		helpStr = c.extraLdapHelpFunc(helpMap)

	}

	// Keep linter from complaining if we don't actually generate code using it
	_ = helpMap
	return helpStr
}

var flagsLdapMap = map[string][]string{

	"create": {"auth-method-id", "name", "description"},

	"update": {"id", "name", "description", "version"},
}

func (c *LdapCommand) Flags() *base.FlagSets {
	if len(flagsLdapMap[c.Func]) == 0 {
		return c.FlagSet(base.FlagSetNone)
	}

	set := c.FlagSet(base.FlagSetHTTP | base.FlagSetClient | base.FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")
	common.PopulateCommonFlags(c.Command, f, "ldap-type account", flagsLdapMap, c.Func)

	extraLdapFlagsFunc(c, set, f)

	return set
}

func (c *LdapCommand) Run(args []string) int {
	initLdapFlags()

	switch c.Func {
	case "":
		return cli.RunResultHelp

	}

	c.plural = "ldap-type account"
	switch c.Func {
	case "list":
		c.plural = "ldap-type accounts"
	}

	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}

	if strutil.StrListContains(flagsLdapMap[c.Func], "id") && c.FlagId == "" {
		c.PrintCliError(errors.New("ID is required but not passed in via -id"))
		return base.CommandUserError
	}

	var opts []accounts.Option

	if strutil.StrListContains(flagsLdapMap[c.Func], "auth-method-id") {
		switch c.Func {

		case "create":
			if c.FlagAuthMethodId == "" {
				c.PrintCliError(errors.New("AuthMethod ID must be passed in via -auth-method-id or BOUNDARY_AUTH_METHOD_ID"))
				return base.CommandUserError
			}

		}
	}

	client, err := c.Client()
	if c.WrapperCleanupFunc != nil {
		defer func() {
			if err := c.WrapperCleanupFunc(); err != nil {
				c.PrintCliError(fmt.Errorf("Error cleaning kms wrapper: %w", err))
			}
		}()
	}
	if err != nil {
		c.PrintCliError(fmt.Errorf("Error creating API client: %w", err))
		return base.CommandCliError
	}
	accountsClient := accounts.NewClient(client)

	switch c.FlagName {
	case "":
	case "null":
		opts = append(opts, accounts.DefaultName())
	default:
		opts = append(opts, accounts.WithName(c.FlagName))
	}

	switch c.FlagDescription {
	case "":
	case "null":
		opts = append(opts, accounts.DefaultDescription())
	default:
		opts = append(opts, accounts.WithDescription(c.FlagDescription))
	}

	if c.FlagFilter != "" {
		opts = append(opts, accounts.WithFilter(c.FlagFilter))
	}

	var version uint32

	switch c.Func {

	case "update":
		switch c.FlagVersion {
		case 0:
			opts = append(opts, accounts.WithAutomaticVersioning(true))
		default:
			version = uint32(c.FlagVersion)
		}

	}

	if ok := extraLdapFlagsHandlingFunc(c, f, &opts); !ok {
		return base.CommandUserError
	}

	var resp *api.Response
	var item *accounts.Account

	var createResult *accounts.AccountCreateResult

	var updateResult *accounts.AccountUpdateResult

	switch c.Func {

	case "create":
		createResult, err = accountsClient.Create(c.Context, c.FlagAuthMethodId, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = createResult.GetResponse()
		item = createResult.GetItem()

	case "update":
		updateResult, err = accountsClient.Update(c.Context, c.FlagId, version, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = updateResult.GetResponse()
		item = updateResult.GetItem()

	}

	resp, item, err = executeExtraLdapActions(c, resp, item, err, accountsClient, version, opts)
	if exitCode := c.checkFuncError(err); exitCode > 0 {
		return exitCode
	}

	output, err := printCustomLdapActionOutput(c)
	if err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}
	if output {
		return base.CommandSuccess
	}

	switch c.Func {

	}

	switch base.Format(c.UI) {
	case "table":
		c.UI.Output(printItemTable(item, resp))

	case "json":
		if ok := c.PrintJsonItem(resp); !ok {
			return base.CommandCliError
		}
	}

	return base.CommandSuccess
}

func (c *LdapCommand) checkFuncError(err error) int {
	if err == nil {
		return 0
	}
	if apiErr := api.AsServerError(err); apiErr != nil {
		c.PrintApiError(apiErr, fmt.Sprintf("Error from controller when performing %s on %s", c.Func, c.plural))
		return base.CommandApiError
	}
	c.PrintCliError(fmt.Errorf("Error trying to %s %s: %s", c.Func, c.plural, err.Error()))
	return base.CommandCliError
}

var (
	extraLdapActionsFlagsMapFunc = func() map[string][]string { return nil }
	extraLdapSynopsisFunc        = func(*LdapCommand) string { return "" }
	extraLdapFlagsFunc           = func(*LdapCommand, *base.FlagSets, *base.FlagSet) {}
	extraLdapFlagsHandlingFunc   = func(*LdapCommand, *base.FlagSets, *[]accounts.Option) bool { return true }
	executeExtraLdapActions      = func(_ *LdapCommand, inResp *api.Response, inItem *accounts.Account, inErr error, _ *accounts.Client, _ uint32, _ []accounts.Option) (*api.Response, *accounts.Account, error) {
		return inResp, inItem, inErr
	}
	printCustomLdapActionOutput = func(*LdapCommand) (bool, error) { return false, nil }
)
