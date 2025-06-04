import { createFileRoute } from '@tanstack/react-router'
import Dashboard from '@/views/dashboard'

export const Route = createFileRoute('/_authenticated/')({
  beforeLoad: () => {
   // eslint-disable-next-line no-console
   console.log('_authenticated.index >>> load');
  },
  component: Dashboard,
})
