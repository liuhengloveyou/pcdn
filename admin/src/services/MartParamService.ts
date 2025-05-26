import api from "@/axios";
import type { HttpResponse } from ".";

export interface MartParamModel {
  id: number;
  uid: number;

  domain: string;
  symbol: string;

  memo: string;
  accessKey: string;
  secretKey: string;

  active: number;
}

export class MartParamService {
  static async Set(model: MartParamModel) {
    const resp = await api.post<HttpResponse<number>>("/api/mart/set", model);
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async Del(id: number) {
    const resp = await api.get<HttpResponse<number>>("/api/mart/del", {
      params: { id: id },
    });
    if (!resp || resp.status != 200) {
      return;
    }

    return resp.data;
  }

  static async LoadLite() {
    const resp = await api.get<HttpResponse<MartParamModel[]>>(
      "/api/mart/loadlite"
    );
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async Load() {
    const resp = await api.get<HttpResponse<MartParamModel[]>>(
      "/api/mart/load"
    );
    if (!resp || resp.status != 200) {
      return;
    }

    if (!resp || !resp.data) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async LoadOne(id: number) {
    const resp = await api.get<HttpResponse<MartParamModel>>("/api/mart/take", {
      params: {
        id: id,
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

  static async Active(id: number) {
    const resp = await api.get<HttpResponse<MartParamModel>>(
      "/api/mart/active",
      {
        params: {
          id: id,
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
}
