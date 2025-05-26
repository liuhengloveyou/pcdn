import api from "@/axios";
import type { HttpResponse } from ".";

// // 机器人下单模式
// type BotOrderModeType string

// const (
// 	BotOrderRandom    BotOrderModeType = "random"
// 	BotOrderBuyFirst  BotOrderModeType = "buyFirst"
// 	BotOrderSellFirst BotOrderModeType = "sellFirst"
// )

// // 机器人价格模式
// type BotPriceModeType string

// const (
// 	BotPriceNormalMode BotPriceModeType = "normalMode"
// 	BotPriceSafeMode   BotPriceModeType = "safeMode"
// )

// // 机器人启动模式
// type BotLaunchModeType string

// const (
// 	BotImmediateLaunch BotLaunchModeType = "immediateStart" // 立即启动
// 	BotScheduledLaunch BotLaunchModeType = "scheduledStart" // 定时启动
// )

// // 机器人数量分布模式
// type BotQuantModeType string

// const (
// 	BotContinuousDistribution BotQuantModeType = "continuousDist" // 连续分布
// 	BotRandomDistribution     BotQuantModeType = "randomDist"     // 随机分布
// )

// 最低维持刷量
export interface MinimumBrushRobotModel {
  id: number;

  // 价格模式：
  priceMode: "normalMode" | "safeMode"; // 一般模式 / 安全模式

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 定时启动时间
  launchTime: string;

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 随机刷量
export interface RandomBrushBotModel {
  id: number;

  // 单笔交易量：
  minNumPeerOrder: number;
  maxNumPeerOrder: number;

  // 时间间隔 s
  minTimeInterval: number;
  maxTimeInterval: number;

  // 下单模式：
  orderMode: "random" | "buyFirst" | "sellFirst"; // 随机 / 买优先 / 卖优先

  // 价格模式：
  priceMode: "normalMode" | "safeMode"; // 一般模式 / 安全模式

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 自然刷量
export interface NaturalBrushRobotModel {
  id: number;

  // 单笔交易量：
  minNumPeerOrder: number;
  maxNumPeerOrder: number;

  // 时间间隔 s
  minTimeInterval: number;
  maxTimeInterval: number;

  // 下单模式：
  orderMode: "random" | "buyFirst" | "sellFirst"; // 随机 / 买优先 / 卖优先

  // 价格模式：
  priceMode: "normalMode" | "safeMode"; // 一般模式 / 安全模式

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 正弦刷量
export interface StableBrushBotModel {
  id: number;

  // 单笔交易量：
  minNumPeerOrder: number;
  maxNumPeerOrder: number;

  // 时间间隔 s
  minTimeInterval: number;
  maxTimeInterval: number;

  // 下单模式：
  orderMode: "random" | "buyFirst" | "sellFirst"; // 随机 / 买优先 / 卖优先

  // 价格模式：
  priceMode: "normalMode" | "safeMode"; // 一般模式 / 安全模式

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 区间刷量
export interface RangeBrushBotModel {
  id: number;

  // 单笔交易量：
  minNumPeerOrder: number;
  maxNumPeerOrder: number;

  // 价格范围
  minPrice: string;
  maxPrice: string;

  // 时间间隔 s
  minTimeInterval: number;
  maxTimeInterval: number;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 盤口閃爍
export interface OrderbookFlashBotModel {
  id: number;

  // 时间间隔 s
  minTimeInterval: number;
  maxTimeInterval: number;

  minGear: number;
  maxGear: number;

  minOrderSize: number;
  maxOrderSize: number;

  minDivtions: number;
  maxDivtions: number;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

export interface DepthRemoteBotModel {
  id: number;

  minBuyOrderSize: number;
  maxBuyOrderSize: number;

  minBuySheet: number;
  maxBuySheet: number;

  minSellOrderSize: number;
  maxSellOrderSize: number;

  minSellSheet: number;
  maxSellSheet: number;

  minBuySpread: number;
  maxBuySpread: number;

  minSellSpread: number;
  maxSellSpread: number;

  minDuration: number;
  maxDuration: number;

  // 下单模式：
  orderMode: "random" | "buyFirst" | "sellFirst"; // 随机 / 买优先 / 卖优先

  // 价格模式：
  priceMode: "normalMode" | "safeMode"; // 一般模式 / 安全模式

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

export interface DepthNearBotModel {
  id: number;

  minBuyOrderSize: number;
  maxBuyOrderSize: number;

  minBuySheet: number;
  maxBuySheet: number;

  minSellOrderSize: number;
  maxSellOrderSize: number;

  minSellSheet: number;
  maxSellSheet: number;

  minBuySpread: number;
  maxBuySpread: number;

  minSellSpread: number;
  maxSellSpread: number;

  minDuration: number;
  maxDuration: number;

  // 下单模式：
  orderMode: "random" | "buyFirst" | "sellFirst"; // 随机 / 买优先 / 卖优先

  // 价格模式：
  priceMode: "normalMode" | "safeMode"; // 一般模式 / 安全模式

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

export interface DepthSpreadBotModel {
  id: number;

  minBuyOrderSize: number;
  maxBuyOrderSize: number;

  minBuySheet: number;
  maxBuySheet: number;

  minSellOrderSize: number;
  maxSellOrderSize: number;

  minSellSheet: number;
  maxSellSheet: number;

  minBuySpread: number;
  maxBuySpread: number;

  minSellSpread: number;
  maxSellSpread: number;

  minDuration: number;
  maxDuration: number;

  // 下单模式：
  orderMode: "random" | "buyFirst" | "sellFirst"; // 随机 / 买优先 / 卖优先

  // 价格模式：
  priceMode: "normalMode" | "safeMode"; // 一般模式 / 安全模式

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 振荡拉升
export interface GradualLiftBotModel {
  targetPrice: string;
  totalCost: string;
  totalTime: number;
  minOrderSize: number;
  maxOrderSize: number;

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 振荡下跌
export interface GradualDownBotModel {
  targetPrice: string;
  totalCost: string;
  totalTime: number;
  minOrderSize: number;
  maxOrderSize: number;

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 急速拉升
export interface RapidLiftBotModel {
  targetPrice: string;
  totalCost: string;
  totalTime: number;

  minOrderSize: number;
  maxOrderSize: number;

  // 时间间隔 s
  minInterval: number;
  maxInterval: number;

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 急速下跌
export interface RapidDownBotModel {
  targetPrice: string;
  totalCost: string;
  totalTime: number;

  minOrderSize: number;
  maxOrderSize: number;

  // 时间间隔 s
  minInterval: number;
  maxInterval: number;

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 比例拉升
export interface RatioLiftBotModel {
  targetPrice: string;
  totalCost: string;
  totalTime: number;

  minOrderSize: number;
  maxOrderSize: number;

  minOscillate: number;
  maxOscillate: number;
  oscillateRatio: number;

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 比例下跌
export interface RatioDownBotModel {
  targetPrice: string;
  totalCost: string;
  totalTime: number;

  minOrderSize: number;
  maxOrderSize: number;

  minOscillate: number;
  maxOscillate: number;
  oscillateRatio: number;

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}


// 区间振荡
export interface RangeOscillationBotModel {
  totalCost: string;

  minPrice: string;
  maxPrice: string;

  minOrderSize: number;
  maxOrderSize: number;


  // 时间间隔 s
  minInterval: number;
  maxInterval: number;

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;
}


// 比例振荡
export interface RatioOscillationBotModel {
  totalCost: string;
  totalTime: number;

  minOrderSize: number;
  maxOrderSize: number;

  minOscillate: number;
  maxOscillate: number;  
  // oscillateRatio : number;  

  // 数量分布类型：
  quantMode: "continuousDist" | "randomDist"; // 连续分布 / 随机分布

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  // 机器人描述：
  description: string;
}

// 高抛低吸
export interface HighThrowBargainRobotModel {

  maxSell: string;
  sellPrice: string;
  minSellSize: number;
  maxSellSize: number;
  minSellInterval  : number;
  maxSellInterval  : number;

  maxBuy: string;
  buyPrice: string;
  minBuySize: number;
  maxBuySize: number;
  minBuyInterval : number; 
  maxBuyInterval : number;

   // 启动模式：
   launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

}

export interface BotModel {
  id: number;
  botType: number;
  botName: string;
  isRunning: number;
  martDomain: string;
  symbol: string;
  startTime: number; // 啟動時間
  startTimeStr?: string;

  minimumBrushBot?: MinimumBrushRobotModel;
  randomBrushBot?: RandomBrushBotModel;
  naturalBrushBot?: NaturalBrushRobotModel;
  stableBrushBot?: StableBrushBotModel;
  rangeBrushBot?: RangeBrushBotModel;
  orderbookFlashBot?: OrderbookFlashBotModel;
  depthRemoteBot?: DepthRemoteBotModel;
  depthNearBot?: DepthNearBotModel;
  depthSpreadBot?: DepthSpreadBotModel;

  gradualLiftBot?: GradualLiftBotModel;
  gradualDownBot?: GradualDownBotModel;
  rapidLiftBot?: RapidLiftBotModel;
  rapidDownBot?: RapidDownBotModel;
  ratioLiftBot?: RatioLiftBotModel;
  ratioDownBot?: RatioDownBotModel;
  rangeOscillationBot?: RangeOscillationBotModel;
  ratioOscillationBot?: RatioOscillationBotModel;
  highThrowBargainBot?: HighThrowBargainRobotModel;
}

export class BotService {
  static async Set(model: BotModel) {
    const resp = await api.post<HttpResponse<number>>("/api/bot/set", model);
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async Load(martDomain: string, symbol: string) {
    const resp = await api.get<HttpResponse<BotModel[]>>("/api/bot/load", {
      params: {
        domain: martDomain,
        symbol: symbol,
      },
    });
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async LoadOne(botType: number, martDomain: string, symbol: string) {
    const resp = await api.get<HttpResponse<BotModel>>("/api/bot/take", {
      params: {
        t: botType,
        domain: martDomain,
        symbol: symbol,
      },
    });
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async Run(
    botId: number,
    botType: number,
    martDomain: string,
    symbol: string
  ) {
    const resp = await api.get<HttpResponse<string>>("/api/bot/run", {
      params: {
        id: botId,
        t: botType,
        domain: martDomain,
        symbol: symbol,
      },
    });
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async Stop(
    botId: number,
    botType: number,
    martDomain: string,
    symbol: string
  ) {
    const resp = await api.get<HttpResponse<string>>("/api/bot/stop", {
      params: {
        id: botId,
        t: botType,
        domain: martDomain,
        symbol: symbol,
      },
    });
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }
}
