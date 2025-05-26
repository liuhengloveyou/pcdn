import api from '@/axios';
import type { HttpResponse } from '.';

export interface OrderModel {
  id: number;
  uid: number;
  createTime: number;
  updateTime: number;
  createTimeStr: string;

  martDomain: string;
  symbol: string;
  side: string;
  type: string;

  client_order_id: string;
  quantity: string;
  price: string;
  notional: string;
  delayed: number;
}

export class OrderService {
  static async List(martDomain: string, symbol: string, page: number) {
    const resp = await api.get<HttpResponse<OrderModel[]>>('/api/order/list', {
      params: {
        domain: martDomain,
        symbol: symbol,
        page: page,
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
