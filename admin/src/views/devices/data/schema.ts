import { z } from 'zod';
import type { DeviceModel } from '@/apis/device_api';

// 使用 DeviceModel 类型作为基础
export const deviceSchema = z.object({
  id: z.number(),
  uid: z.number(),
  sn: z.string(),
  label: z.string().optional(), // 改为可选字段
  status: z.string().optional(),
  createTime: z.number(),
  updateTime: z.number(),
  remoteAddr: z.string(),
  version: z.string(),
  timestamp: z.number(),
  lastHeartbear: z.number(),
  lastHeartbearStr: z.string()
});

// 重新导出 DeviceModel 类型，确保类型一致
export type Device = DeviceModel;
export type { DeviceModel };
