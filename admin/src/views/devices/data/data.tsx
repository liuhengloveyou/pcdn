import {
  IconCircleOff,
  IconWifi,
} from '@tabler/icons-react'

export const labels = [
  {
    value: 'bug',
    label: 'Bug',
  },
  {
    value: 'feature',
    label: 'Feature',
  },
  {
    value: 'documentation',
    label: 'Documentation',
  },
]

export const statuses = [
  {
    value: "在线",
    label: "在线",
    icon: IconWifi,
  },
  {
    value: "离线",
    label: "离线",
    icon: IconCircleOff,
  }
]