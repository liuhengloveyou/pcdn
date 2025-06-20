import { createFileRoute } from '@tanstack/react-router'
import Dashboard from '@/views/dashboard'

export const Route = createFileRoute('/_authenticated/')({
  beforeLoad: () => {
  },
  component: Dashboard,
})
