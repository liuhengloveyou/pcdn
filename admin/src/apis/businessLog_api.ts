import api from '@/utils/axios'
import type { HttpResponse } from '.'

export interface BusinessLogModel {
  id: number;
  uid: number;
  createTime: number;
  updateTime: number;
  createTimeStr: string;

  businessType: string;
  payload: string;
}

export class BusinessLogApi {
  static async List(businessType: string) {
    const resp = await api.get<HttpResponse<BusinessLogModel[]>>(
      '/api/log/list',
      {
        params: {
          businessType: businessType,
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
