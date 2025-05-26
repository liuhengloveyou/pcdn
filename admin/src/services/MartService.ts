import api from "@/axios";
import type { HttpResponse } from ".";


export interface kline {
  // 时间戳，毫秒级别，必要字段
  timestamp: number;
  // 开盘价，必要字段
  open: string;
  // 收盘价，必要字段
  close: string;
  // 最高价，必要字段
  high: string;
  // 最低价，必要字段
  low: string;
  // 成交量，非必须字段
  volume: string;
  // 成交额，非必须字段，如果需要展示技术指标'EMV'和'AVP'，则需要为该字段填充数据。
  turnover: string;
}

export interface WalletBalance {
  available: string;
  currency: string;
  frozen: string;
  name: string;
  updateAt: string;
}

export interface DepthModel {
  ts: number; // Create time(Timestamp in milliseconds)
  asks: string[][]; // Order book on sell side
  bids: string[][]; // Order book on buy side
  amount: string; // Total number of current price depth
  price: string; // The price at current depth
  symbol: string; // Trading pair
}

export class MartService {
  static async LoadKline(martDomain: string, symbol: string) {
    const resp = await api.get<HttpResponse<kline[]>>("/api/spot/kline", {
      params: {
        mart: martDomain,
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

  static async LoadDepth(martDomain: string, symbol: string) {
    const resp = await api.get<HttpResponse<DepthModel>>("/api/spot/depth", {
      params: {
        mart: martDomain,
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

  static async LoadSpotWallet(martDomain: string) {
    const resp = await api.get<HttpResponse<WalletBalance[]>>(
      "/api/spot/wallet",
      {
        params: {
          domain: martDomain,
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

  static async NewSpotOrder(
    martDomain: string,
    symbol: string,
    side: string,
    quantity: string,
    price: string
  ) {
    const resp = await api.post<HttpResponse<WalletBalance[]>>(
      "/api/spot/neworder",
      {
        mart: martDomain,
        symbol: symbol,
        side: side,
        quantity: quantity,
        price: price,
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
}
