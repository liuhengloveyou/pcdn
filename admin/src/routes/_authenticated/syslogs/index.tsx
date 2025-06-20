import { createFileRoute } from '@tanstack/react-router'
import SysLogs from '@/views/syslogs'

export const Route = createFileRoute('/_authenticated/logs/')({  
  component: SysLogs,
})