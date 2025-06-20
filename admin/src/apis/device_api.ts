/* eslint-disable no-console */
import api from '@/utils/axios'
import type { HttpResponse } from '.'

export interface DeviceModel {
  id: number
  uid: number
  sn: string
  createTime: number
  updateTime: number
  remoteAddr: string
  version: string
  timestamp: number
  lastHeartbear: number
  label: string

  lastHeartbearStr: string
  status?: string
}

export interface SystemMonitorProcess {
  pid: number
  name: string
  exe: string
  cpu: number
  memory: number
  status: string
}

export interface SystemMonitorCpu {
  usage: number
  cores: number
  temperature: number
}

export interface SystemMonitorMemory {
  used: number
  total: number
  available: number
}

export interface SystemMonitorDisk {
  used: number
  total: number
  free: number
}

export interface SystemMonitorNetwork {
  name: string;
  recv_rate: number; // bytes/s
  send_rate: number; // bytes/s
}

// 获取设备系统监控数据
export interface SystemMonitorData {
  cpu?: SystemMonitorCpu
  memory?: SystemMonitorMemory
  disk?: SystemMonitorDisk
  network?: Array<SystemMonitorNetwork>
  processes?: Array<SystemMonitorProcess>
}

export async function addDevice(data: DeviceModel) {
  const resp = await api.post('/api/device/add', data)
  if (resp.status != 200) {
    return
  }

  return resp.data
}

export async function listDevice(page: number, pageSize: number) {
  const resp = await api.get<HttpResponse<DeviceModel[]>>('/api/device/list', {
    params: { page, pageSize },
  })
  // if (!resp || resp.status != 200) {
  // } else if (resp.data.code != 0) {
  // }

  return resp.data
}

export function updateDevice(id: number, data: DeviceModel) {
  return api.put(`/api/device/update/${id}`, data)
}

export function deleteDevice(id: number) {
  return api.delete(`/api/device/delete/${id}`)
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

export async function getDeviceMonitor(sn: string) {
  try {
    const resp = await api.get<HttpResponse<SystemMonitorData>>(
      `/api/device/monitor`,
      {
        params: { sn: sn },
      }
    )
    return resp.data
  } catch (error) {
    console.error('获取设备监控数据失败:', error)
    throw error
  }
}

// 获取路由器管理界面URL
export async function getRouterAdminUrl(sn: string) {
  try {
    const resp = await api.get<HttpResponse<{ url: string }>>(
      `/api/device/router-admin`,
      {
        params: { sn: sn },
      }
    )
    return resp.data
  } catch (error) {
    console.error('获取路由器管理界面URL失败:', error)
    throw error
  }
}
