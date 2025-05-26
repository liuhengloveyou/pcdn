import api from '@/axios';
import type { HttpResponse } from '.';

export interface BusinessLogModel {
  id: number;
  uid: number;
  createTime: number;
  updateTime: number;
  createTimeStr: string;

  domain: string;
  symbol: string;
  side: string;
  botType: string;
  businessType: string;
}

export class BusinessLogService {
  static async List(martDomain: string, symbol: string) {
    const resp = await api.get<HttpResponse<BusinessLogModel[]>>(
      '/api/log/list',
      {
        params: {
          domain: martDomain,
          symbol: symbol,
        },
      },
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
