import { createFileRoute } from '@tanstack/react-router'
import Chats from '@/views/chats'

export const Route = createFileRoute('/_authenticated/chats/')({
  component: Chats,
})
