package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/snapshots"
	snapshottypes "cosmossdk.io/store/snapshots/types"
	storetypes "cosmossdk.io/store/types"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	rosettaCmd "github.com/cosmos/rosetta/cmd"
	dydxapp "github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/constants"
	protocolflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Unnamed import of statik for swagger UI support.
	// Used in cosmos-sdk when registering the route for swagger docs.
	_ "github.com/dydxprotocol/v4-chain/protocol/client/docs/statik"
)

const (
	EnvPrefix = "DYDX"

	flagIAVLCacheSize = "iavl-cache-size"
)

// NewRootCmd creates a new root command for `dydxprotocold`. It is called once in the main function.
// TODO(DEC-1097): improve `cmd/` by adding tests, custom app configs, custom init cmd, and etc.
func NewRootCmd(
	option *RootCmdOption,
	homeDir string,
) *cobra.Command {
	return NewRootCmdWithInterceptors(
		option,
		homeDir,
		func(serverCtxPtr *server.Context) {

		},
		func(s string, appConfig *DydxAppConfig) (string, *DydxAppConfig) {
			return s, appConfig
		},
		func(app *dydxapp.App) *dydxapp.App {
			return app
		},
	)
}

func NewRootCmdWithInterceptors(
	option *RootCmdOption,
	homeDir string,
	serverCtxInterceptor func(serverCtxPtr *server.Context),
	appConfigInterceptor func(string, *DydxAppConfig) (string, *DydxAppConfig),
	appInterceptor func(app *dydxapp.App) *dydxapp.App,
) *cobra.Command {
	tempApp := dydxapp.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		viper.New(),
	)
	defer func() {
		if err := tempApp.Close(); err != nil {
			panic(err)
		}
	}()

	initClientCtx := client.Context{}.
		WithCodec(tempApp.AppCodec()).
		WithInterfaceRegistry(tempApp.InterfaceRegistry()).
		WithTxConfig(tempApp.TxConfig()).
		WithLegacyAmino(tempApp.LegacyAmino()).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(homeDir).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   constants.AppDaemonName,
		Short: "Start dydxprotocol app",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context()).WithViper("")
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := appConfigInterceptor(initAppConfig())
			customTMConfig := initTendermintConfig()

			if err := server.InterceptConfigsPreRunHandler(
				cmd,
				customAppTemplate,
				customAppConfig,
				customTMConfig,
			); err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)

			// Format logs for error tracking if it is enabled via flags.
			if ddErrorTrackingFormatterEnabled :=
				serverCtx.Viper.Get(protocolflags.DdErrorTrackingFormat); ddErrorTrackingFormatterEnabled != nil {
				if enabled, err := cast.ToBoolE(ddErrorTrackingFormatterEnabled); err == nil && enabled {
					dydxapp.SetZerologDatadogErrorTrackingFormat()
				}
			}
			serverCtxInterceptor(serverCtx)

			return nil
		},
		SilenceUsage: true,
	}

	initRootCmd(tempApp, rootCmd, option, appInterceptor)
	initClientCtx, err := config.ReadFromClientConfig(initClientCtx)
	if err != nil {
		panic(err)
	}
	// [ Home Dir Temp Fix ] (also see protocol/cmd/dydxprotocold/main.go)
	// Unset the temp home dir. This must be done after `ReadFromClientConfig`, otherwise it will
	// create a temp dir in cwd.
	initClientCtx.HomeDir = ""
	if err := autoCliOpts(tempApp, initClientCtx).EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}

// initRootCmd initializes the app's root command with useful commands.
func initRootCmd(
	tempApp *dydxapp.App,
	rootCmd *cobra.Command,
	option *RootCmdOption,
	appInterceptor func(app *dydxapp.App) *dydxapp.App,
) {
	valOperAddressCodec := address.NewBech32Codec(sdktypes.GetConfig().GetBech32ValidatorAddrPrefix())
	rootCmd.AddCommand(
		genutilcli.InitCmd(tempApp.ModuleBasics, dydxapp.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(
			banktypes.GenesisBalancesIterator{},
			dydxapp.DefaultNodeHome,
			genutiltypes.DefaultMessageValidator,
			valOperAddressCodec,
		),
		genutilcli.MigrateGenesisCmd(genutilcli.MigrationMap),
		genutilcli.GenTxCmd(
			tempApp.ModuleBasics,
			tempApp.TxConfig(),
			banktypes.GenesisBalancesIterator{},
			dydxapp.DefaultNodeHome,
			valOperAddressCodec,
		),
		genutilcli.ValidateGenesisCmd(tempApp.ModuleBasics),
		AddGenesisAccountCmd(dydxapp.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
	)

	server.AddCommands(
		rootCmd,
		dydxapp.DefaultNodeHome,
		func(logger log.Logger, db dbm.DB, writer io.Writer, options servertypes.AppOptions) servertypes.Application {
			return appInterceptor(newApp(logger, db, writer, options))
		},
		appExport,
		func(cmd *cobra.Command) {
			addModuleInitFlags(cmd)

			if option.startCmdCustomizer != nil {
				option.startCmdCustomizer(cmd)
			}
		},
	)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(),
	)

	rootCmd.AddCommand(rosettaCmd.RosettaCommand(tempApp.InterfaceRegistry(), tempApp.AppCodec()))
}

// autoCliOpts returns options based upon the modules in the dYdX v4 app.
//
// Creates an instance of the application that is discarded to enumerate the modules.
func autoCliOpts(tempApp *dydxapp.App, initClientCtx client.Context) autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range tempApp.ModuleManager.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	cliKR, err := keyring.NewAutoCLIKeyring(initClientCtx.Keyring)
	if err != nil {
		panic(err)
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(tempApp.ModuleManager.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdktypes.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdktypes.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdktypes.GetConfig().GetBech32ConsensusAddrPrefix()),
		Keyring:               cliKR,
		ClientCtx:             initClientCtx,
	}
}

// addModuleInitFlags adds module specific init flags.
func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// queryCommand adds transaction and account querying commands.
func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	// Module specific query sub-commands are added by AutoCLI

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// txCommand adds transaction signing, encoding / decoding, and broadcasting commands.
func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	// Module specific tx sub-commands are added by AutoCLI

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// newApp initializes and returns a new app.
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) *dydxapp.App {
	var cache storetypes.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		appGenesis, err := genutiltypes.AppGenesisFromFile(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}

		chainID = appGenesis.ChainID
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := dbm.NewDB("metadata", server.GetAppDBBackend(appOpts), snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent)),
	)

	return dydxapp.New(
		logger,
		db,
		traceStore,
		true,
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetIAVLCacheSize(int(cast.ToUint64(appOpts.Get(flagIAVLCacheSize)))),
		baseapp.SetChainID(chainID),
	)
}

// appExport creates and exports a new app, returns the state of the app for a genesis file.
//
// Deprecated: this feature relies on the use of known unstable, legacy export functionality
// from cosmos.
func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	dydxApp := dydxapp.New(
		logger,
		db,
		traceStore,
		height == -1, // -1: no height provided
		appOpts,
	)

	if height != -1 {
		if err := dydxApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return dydxApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
