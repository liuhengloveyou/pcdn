import type { SidebarData, Team, User } from '../types'
import { useSidebar } from '@/composables/use-sidebar'
import {
  AudioWaveform,
  GalleryVerticalEnd,
} from 'lucide-vue-next'

const user: User = {
  name: 'liuheng',
  email: 'liuheng@qq.com',
  avatar: '/avatars/shadcn.jpg',
}

const teams: Team[] = [
  {
    name: '智算未来 Inc',
    logo: GalleryVerticalEnd,
    plan: 'Enterprise',
  },
  {
    name: '智算未来 Corp.',
    logo: AudioWaveform,
    plan: 'Startup',
  },
  // {
  //   name: 'Evil Corp.',
  //   logo: Command,
  //   plan: 'Free',
  // },
]

const { navData } = useSidebar()

export const sidebarData: SidebarData = {
  user,
  teams,
  navMain: navData.value!,
}
