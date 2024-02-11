package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	vestmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) TotalSupply(goCtx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	denom := k.GetParams(ctx).Denom
	totalMinted := k.bankKeeper.GetSupply(ctx, denom)
	rewardsVesterBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.VesterAccountName), denom)
	distributionBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(distrtypes.ModuleName), denom)
	communityVesterBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(vestmoduletypes.CommunityVesterAccountName), denom)
	bridgeBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(bridgetypes.ModuleName), denom)

	// Assuming the following balances in the module accounts need to be subtracted
	// from the total minted supply
	supply := totalMinted.
		Sub(rewardsVesterBalance).
		Sub(distributionBalance).
		Sub(communityVesterBalance).
		Sub(bridgeBalance)

	return &types.QueryTotalSupplyResponse{Supply: sdk.Coins{supply}}, nil
}
