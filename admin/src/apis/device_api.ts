/* eslint-disable no-console */
import api from '@/utils/axios';
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
  label: string;
  
  lastHeartbearStr: string;
  status?: string;
}

export async function addDevice(data: DeviceModel) {
  const resp = await api.post('/api/device/add', data);
  if (resp.status != 200) {
    return;
  }

  return resp.data;
}

export async function listDevice(page: number, pageSize: number) {
  const resp = await api.get<HttpResponse<DeviceModel[]>>('/api/device/list', {params: {page, pageSize}});
  // if (!resp || resp.status != 200) {
  // } else if (resp.data.code != 0) {
  // }

  return resp.data;
}

export function updateDevice(id: number, data: DeviceModel) {
  return api.put(`/api/device/update/${id}`, data);
}

export function deleteDevice(id: number) {
  return api.delete(`/api/device/delete/${id}`);
}

// 在现有的API函数中添加
export const resetDevicePassword = async (deviceId: string): Promise<void> => {
  try {
    const response = await api.post(`/api/devices/${deviceId}/reset-password`)
    return response.data
  } catch (error) {
    console.error('重置设备密码失败:', error)
    throw error
  }
}