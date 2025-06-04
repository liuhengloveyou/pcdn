import { createFileRoute } from '@tanstack/react-router'
import Device from '@/views/devices'

export const Route = createFileRoute('/_authenticated/devices/')({
  component: Device,
})
