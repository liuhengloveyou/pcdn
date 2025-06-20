import {
  IconBrowserCheck,
  IconHelp,
  IconLayoutDashboard,
  IconNotification,
  IconPalette,
  IconSettings,
  IconTool,
  IconUserCog,
  IconUsers,
} from '@tabler/icons-react'
import {  Command,  ListTodo } from 'lucide-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  teams: [
    {
      name: '智算 PCDN 管理系统',
      logo: Command,
      plan: '为您提供一站式的PCDN服务',
    }
  ],
  navGroups: [
    {
      title: '菜单',
      items: [
        {
          title: '仪表盘',
          url: '/',
          icon: IconLayoutDashboard,
        },
        {
          title: '设备管理',
          url: '/devices',
          icon: ListTodo,
        },
        {
          title: '账号管理',
          url: '/users',
          icon: IconUsers,
        },
        {
          title: '系统日志',
          url: '/logs',
          icon: IconNotification,
        },
        // {
        //   title: '应用',
        //   url: '/apps',
        //   icon: IconPackages,
        // },
        // {
        //   title: '聊天',
        //   url: '/chats',
        //   badge: '3',
        //   icon: IconMessages,
        // },
      ],
    },
    // {
    //   title: '页面',
    //   items: [
    //     {
    //       title: '认证',
    //       icon: IconLockAccess,
    //       items: [
    //         {
    //           title: '登录',
    //           url: '/sign-in',
    //         },
    //         {
    //           title: '登录 (双栏)',
    //           url: '/sign-in-2',
    //         },
    //         {
    //           title: '注册',
    //           url: '/sign-up',
    //         },
    //         {
    //           title: '忘记密码',
    //           url: '/forgot-password',
    //         },
    //         {
    //           title: '验证码',
    //           url: '/otp',
    //         },
    //       ],
    //     },
    //     {
    //       title: '错误',
    //       icon: IconBug,
    //       items: [
    //         {
    //           title: '未授权',
    //           url: '/401',
    //           icon: IconLock,
    //         },
    //         {
    //           title: '禁止访问',
    //           url: '/403',
    //           icon: IconUserOff,
    //         },
    //         {
    //           title: '未找到',
    //           url: '/404',
    //           icon: IconError404,
    //         },
    //         {
    //           title: '服务器内部错误',
    //           url: '/500',
    //           icon: IconServerOff,
    //         },
    //         {
    //           title: '维护错误',
    //           url: '/503',
    //           icon: IconBarrierBlock,
    //         },
    //       ],
    //     },
    //   ],
    // },
    {
      title: '系统',
      items: [
        {
          title: '设置',
          icon: IconSettings,
          items: [
            {
              title: '个人资料',
              url: '/settings',
              icon: IconUserCog,
            },
            {
              title: '账户',
              url: '/settings/account',
              icon: IconTool,
            },
            {
              title: '外观',
              url: '/settings/appearance',
              icon: IconPalette,
            },
            {
              title: '通知',
              url: '/settings/notifications',
              icon: IconNotification,
            },
            {
              title: '显示',
              url: '/settings/display',
              icon: IconBrowserCheck,
            },
          ],
        },
        {
          title: '帮助中心',
          url: '/help-center',
          icon: IconHelp,
        },
      ],
    },
  ],
}
