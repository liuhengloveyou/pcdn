import { createFileRoute } from '@tanstack/react-router'
import SignUp from '@/views/auth/sign-up'

export const Route = createFileRoute('/(auth)/sign-up')({
  component: SignUp,
})
