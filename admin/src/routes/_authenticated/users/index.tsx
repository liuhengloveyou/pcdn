import { createFileRoute } from '@tanstack/react-router'
import Users from '@/views/users'

export const Route = createFileRoute('/_authenticated/users/')({
  component: Users,
})
