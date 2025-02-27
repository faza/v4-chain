import * as _5 from "./assets/asset";
import * as _6 from "./assets/genesis";
import * as _7 from "./assets/query";
import * as _8 from "./assets/tx";
import * as _9 from "./blocktime/blocktime";
import * as _10 from "./blocktime/genesis";
import * as _11 from "./blocktime/params";
import * as _12 from "./blocktime/query";
import * as _13 from "./blocktime/tx";
import * as _14 from "./bridge/bridge_event_info";
import * as _15 from "./bridge/bridge_event";
import * as _16 from "./bridge/genesis";
import * as _17 from "./bridge/params";
import * as _18 from "./bridge/query";
import * as _19 from "./bridge/tx";
import * as _20 from "./clob/block_rate_limit_config";
import * as _21 from "./clob/clob_pair";
import * as _22 from "./clob/equity_tier_limit_config";
import * as _23 from "./clob/genesis";
import * as _24 from "./clob/liquidations_config";
import * as _25 from "./clob/liquidations";
import * as _26 from "./clob/matches";
import * as _27 from "./clob/mev";
import * as _28 from "./clob/operation";
import * as _29 from "./clob/order_removals";
import * as _30 from "./clob/order";
import * as _31 from "./clob/process_proposer_matches_events";
import * as _32 from "./clob/query";
import * as _33 from "./clob/tx";
import * as _34 from "./daemons/bridge/bridge";
import * as _35 from "./daemons/liquidation/liquidation";
import * as _36 from "./daemons/pricefeed/price_feed";
import * as _37 from "./delaymsg/block_message_ids";
import * as _38 from "./delaymsg/delayed_message";
import * as _39 from "./delaymsg/genesis";
import * as _40 from "./delaymsg/query";
import * as _41 from "./delaymsg/tx";
import * as _42 from "./epochs/epoch_info";
import * as _43 from "./epochs/genesis";
import * as _44 from "./epochs/query";
import * as _45 from "./feetiers/genesis";
import * as _46 from "./feetiers/params";
import * as _47 from "./feetiers/query";
import * as _48 from "./feetiers/tx";
import * as _49 from "./govplus/genesis";
import * as _50 from "./govplus/query";
import * as _51 from "./govplus/tx";
import * as _52 from "./indexer/events/events";
import * as _53 from "./indexer/indexer_manager/event";
import * as _54 from "./indexer/off_chain_updates/off_chain_updates";
import * as _55 from "./indexer/protocol/v1/clob";
import * as _56 from "./indexer/protocol/v1/subaccount";
import * as _57 from "./indexer/redis/redis_order";
import * as _58 from "./indexer/shared/removal_reason";
import * as _59 from "./indexer/socks/messages";
import * as _60 from "./perpetuals/genesis";
import * as _61 from "./perpetuals/params";
import * as _62 from "./perpetuals/perpetual";
import * as _63 from "./perpetuals/query";
import * as _64 from "./perpetuals/tx";
import * as _65 from "./prices/genesis";
import * as _66 from "./prices/market_param";
import * as _67 from "./prices/market_price";
import * as _68 from "./prices/query";
import * as _69 from "./prices/tx";
import * as _70 from "./ratelimit/capacity";
import * as _71 from "./ratelimit/genesis";
import * as _72 from "./ratelimit/limit_params";
import * as _73 from "./ratelimit/query";
import * as _74 from "./ratelimit/tx";
import * as _75 from "./rewards/genesis";
import * as _76 from "./rewards/params";
import * as _77 from "./rewards/query";
import * as _78 from "./rewards/reward_share";
import * as _79 from "./rewards/tx";
import * as _80 from "./sending/genesis";
import * as _81 from "./sending/query";
import * as _82 from "./sending/transfer";
import * as _83 from "./sending/tx";
import * as _84 from "./stats/genesis";
import * as _85 from "./stats/params";
import * as _86 from "./stats/query";
import * as _87 from "./stats/stats";
import * as _88 from "./stats/tx";
import * as _89 from "./subaccounts/asset_position";
import * as _90 from "./subaccounts/genesis";
import * as _91 from "./subaccounts/perpetual_position";
import * as _92 from "./subaccounts/query";
import * as _93 from "./subaccounts/subaccount";
import * as _94 from "./vest/genesis";
import * as _95 from "./vest/query";
import * as _96 from "./vest/tx";
import * as _97 from "./vest/vest_entry";
import * as _105 from "./assets/query.lcd";
import * as _106 from "./blocktime/query.lcd";
import * as _107 from "./bridge/query.lcd";
import * as _108 from "./clob/query.lcd";
import * as _109 from "./delaymsg/query.lcd";
import * as _110 from "./epochs/query.lcd";
import * as _111 from "./feetiers/query.lcd";
import * as _112 from "./perpetuals/query.lcd";
import * as _113 from "./prices/query.lcd";
import * as _114 from "./ratelimit/query.lcd";
import * as _115 from "./rewards/query.lcd";
import * as _116 from "./stats/query.lcd";
import * as _117 from "./subaccounts/query.lcd";
import * as _118 from "./vest/query.lcd";
import * as _119 from "./assets/query.rpc.Query";
import * as _120 from "./blocktime/query.rpc.Query";
import * as _121 from "./bridge/query.rpc.Query";
import * as _122 from "./clob/query.rpc.Query";
import * as _123 from "./delaymsg/query.rpc.Query";
import * as _124 from "./epochs/query.rpc.Query";
import * as _125 from "./feetiers/query.rpc.Query";
import * as _126 from "./govplus/query.rpc.Query";
import * as _127 from "./perpetuals/query.rpc.Query";
import * as _128 from "./prices/query.rpc.Query";
import * as _129 from "./ratelimit/query.rpc.Query";
import * as _130 from "./rewards/query.rpc.Query";
import * as _131 from "./sending/query.rpc.Query";
import * as _132 from "./stats/query.rpc.Query";
import * as _133 from "./subaccounts/query.rpc.Query";
import * as _134 from "./vest/query.rpc.Query";
import * as _135 from "./blocktime/tx.rpc.msg";
import * as _136 from "./bridge/tx.rpc.msg";
import * as _137 from "./clob/tx.rpc.msg";
import * as _138 from "./delaymsg/tx.rpc.msg";
import * as _139 from "./feetiers/tx.rpc.msg";
import * as _140 from "./govplus/tx.rpc.msg";
import * as _141 from "./perpetuals/tx.rpc.msg";
import * as _142 from "./prices/tx.rpc.msg";
import * as _143 from "./ratelimit/tx.rpc.msg";
import * as _144 from "./rewards/tx.rpc.msg";
import * as _145 from "./sending/tx.rpc.msg";
import * as _146 from "./stats/tx.rpc.msg";
import * as _147 from "./vest/tx.rpc.msg";
import * as _148 from "./lcd";
import * as _149 from "./rpc.query";
import * as _150 from "./rpc.tx";
export namespace dydxprotocol {
  export const assets = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._105,
    ..._119
  };
  export const blocktime = { ..._9,
    ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._106,
    ..._120,
    ..._135
  };
  export const bridge = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._107,
    ..._121,
    ..._136
  };
  export const clob = { ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._30,
    ..._31,
    ..._32,
    ..._33,
    ..._108,
    ..._122,
    ..._137
  };
  export namespace daemons {
    export const bridge = { ..._34
    };
    export const liquidation = { ..._35
    };
    export const pricefeed = { ..._36
    };
  }
  export const delaymsg = { ..._37,
    ..._38,
    ..._39,
    ..._40,
    ..._41,
    ..._109,
    ..._123,
    ..._138
  };
  export const epochs = { ..._42,
    ..._43,
    ..._44,
    ..._110,
    ..._124
  };
  export const feetiers = { ..._45,
    ..._46,
    ..._47,
    ..._48,
    ..._111,
    ..._125,
    ..._139
  };
  export const govplus = { ..._49,
    ..._50,
    ..._51,
    ..._126,
    ..._140
  };
  export namespace indexer {
    export const events = { ..._52
    };
    export const indexer_manager = { ..._53
    };
    export const off_chain_updates = { ..._54
    };
    export namespace protocol {
      export const v1 = { ..._55,
        ..._56
      };
    }
    export const redis = { ..._57
    };
    export const shared = { ..._58
    };
    export const socks = { ..._59
    };
  }
  export const perpetuals = { ..._60,
    ..._61,
    ..._62,
    ..._63,
    ..._64,
    ..._112,
    ..._127,
    ..._141
  };
  export const prices = { ..._65,
    ..._66,
    ..._67,
    ..._68,
    ..._69,
    ..._113,
    ..._128,
    ..._142
  };
  export const ratelimit = { ..._70,
    ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._114,
    ..._129,
    ..._143
  };
  export const rewards = { ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._115,
    ..._130,
    ..._144
  };
  export const sending = { ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._131,
    ..._145
  };
  export const stats = { ..._84,
    ..._85,
    ..._86,
    ..._87,
    ..._88,
    ..._116,
    ..._132,
    ..._146
  };
  export const subaccounts = { ..._89,
    ..._90,
    ..._91,
    ..._92,
    ..._93,
    ..._117,
    ..._133
  };
  export const vest = { ..._94,
    ..._95,
    ..._96,
    ..._97,
    ..._118,
    ..._134,
    ..._147
  };
  export const ClientFactory = { ..._148,
    ..._149,
    ..._150
  };
}