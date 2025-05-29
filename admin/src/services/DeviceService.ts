import api from '@/axios';
import type { HttpResponse } from '.';

export interface DeviceModel {
  id: number;
  uid: number;
  sn: string;
  createTime: number;
  updateTime: number;
  remoteAddr: string;
  version: string;
  timestamp: number;
  lastHeartbear: number;

  lastHeartbearStr?: string;
  status?: string;
}

export async function addDevice(data: DeviceModel) {
  const resp = await api.post('/api/device/add', data);
  if (resp.status != 200) {
    return;
  }

  if (resp.data.code != 0) {
  }

  return resp.data;
}

export async function listDevice(params?: any) {
  const resp = await api.get<HttpResponse<DeviceModel[]>>('/api/device/list', { params });
  if (!resp || resp.status != 200) {
  } else if (resp.data.code != 0) {
  }

  return resp.data;
}

export function updateDevice(id: number, data: DeviceModel) {
  return api.put(`/api/device/update/${id}`, data);
}

export function deleteDevice(id: number) {
  return api.delete(`/api/device/delete/${id}`);
}

