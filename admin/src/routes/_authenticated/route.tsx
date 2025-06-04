import { createFileRoute } from '@tanstack/react-router'
import { AuthenticatedLayout } from '@/components/layout/authenticated-layout'

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: () => {
    // eslint-disable-next-line no-console
    console.log('_authenticated.route >>> load')
  },
  component: AuthenticatedLayout,
})
