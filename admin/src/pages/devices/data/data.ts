import {
  ArrowDown,
  ArrowRight,
  ArrowUp,
  CheckCircle,
  Circle,
  CircleCheck,
  CircleHelp,
  CircleOff,
  CirclePlus,
  TimerOff,
  Wifi,
} from 'lucide-vue-next'
import { h } from 'vue'

export const labels = [
  {
    value: '优质设备',
    label: '优质设备',
  },
  {
    value: '问题设备',
    label: '问题设备',
  },
  {
    value: '移动设备',
    label: '移动设备',
  },
]

export const statuses = [
  {
    value: "在线",
    label: "在线",
    icon: h(Wifi),
  },
  {
    value: "离线",
    label: "离线",
    icon: h(CircleOff),
  },
  {
    value: 'in progress',
    label: 'In Progress',
    icon: h(TimerOff),
  },
  {
    value: 'done',
    label: 'Done',
    icon: h(CircleCheck),
  },
  {
    value: 'canceled',
    label: 'Canceled',
    icon: h(CirclePlus),
  },
]

export const priorities = [
  {
    value: 'low',
    label: 'Low',
    icon: h(ArrowDown),
  },
  {
    value: 'medium',
    label: 'Medium',
    icon: h(ArrowRight),
  },
  {
    value: 'high',
    label: 'High',
    icon: h(ArrowUp),
  },
]
