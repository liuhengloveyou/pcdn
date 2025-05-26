import api from "@/axios";
import type { HttpResponse } from ".";

// 余额监控
export interface BalanceMonitorRobotModel {
  id: number;

  mail1: string;
  mail2: string;
  mail3: string;

  remindInterval: number;
  balanceAlertThreshold: string;
  tokenAlertThreshold: string;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 价格监控
export interface PriceMonitorRobotModel {
  id: number;

  mail1: string;

  remindInterval: number;
  minPrice: string;
  maxPrice: string;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 余额监控
export interface BalanceScheduleReportRobotModel {
  id: number;

  mail1: string;
  mail2: string;
  mail3: string;

  remindInterval: number;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// 紧急制动
export interface EmergencyBrakeRobotModel {
  id: number;

  mail1: string;
  mail2: string;
  mail3: string;

  balanceAlertThreshold: string;
  tokenAlertThreshold: string;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// TG Scheduled Report
export interface TGScheduledReportRobotModel {
  id: number;

  chatId: string;
  userName: string;

  remindInterval: number;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

// TG Alert
export interface TGAlertRobotModel {
  id: number;

  chatId: string;
  userName: string;
  availableMoney: string;
  availableToken: string;
  totalMoney: string;
  totalToken: string;
  priceRangeMin: string;
  priceRangeMax: string;

  remindInterval: number;

  // 启动模式：
  launchMode: "immediateStart" | "scheduledStart"; // 立即启动 / 定时启动

  //
  isRunning: boolean;
  startTime: number | string; // 啟動時間
}

export interface RiskManBotModel {
  id: number;
  botType: number;
  botName: string;
  isRunning: number;
  martDomain: string;
  symbol: string;
  startTime: number; // 啟動時間
  startTimeStr?: string;

  balanceMonitorBot?: BalanceMonitorRobotModel;
  priceMonitorBot?: PriceMonitorRobotModel;
  balanceScheduleReportBot?: BalanceScheduleReportRobotModel;
  emergencyBrakeBot?: EmergencyBrakeRobotModel;
  tgScheduledReportBot?: TGScheduledReportRobotModel;
  tgAlertBot?: TGAlertRobotModel;
}

export class RiskManService {
  static async Set(model: RiskManBotModel) {
    const resp = await api.post<HttpResponse<number>>(
      "/api/riskman/set",
      model
    );
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async Load(martDomain: string, symbol: string) {
    const resp = await api.get<HttpResponse<RiskManBotModel[]>>(
      "/api/riskman/load",
      {
        params: {
          domain: martDomain,
          symbol: symbol,
        },
      }
    );
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async LoadOne(botType: number, martDomain: string, symbol: string) {
    const resp = await api.get<HttpResponse<RiskManBotModel>>(
      "/api/riskman/take",
      {
        params: {
          t: botType,
          domain: martDomain,
          symbol: symbol,
        },
      }
    );
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
    const resp = await api.get<HttpResponse<string>>("/api/riskman/run", {
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
    const resp = await api.get<HttpResponse<string>>("/api/riskman/stop", {
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
