import { createFileRoute } from '@tanstack/react-router'
import Apps from '@/views/apps'

export const Route = createFileRoute('/_authenticated/apps/')({
  component: Apps,
})
