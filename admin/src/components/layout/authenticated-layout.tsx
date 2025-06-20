/* eslint-disable no-console */
import React from 'react'
import { useEffect, useState } from 'react'
import Cookies from 'js-cookie'
import { Outlet, useNavigate } from '@tanstack/react-router'
import { AccountApi } from '@/apis/account_api'
import { useAuthStore } from '@/stores/authStore'
import { cn } from '@/lib/utils'
import { SearchProvider } from '@/context/search-context'
import { SidebarProvider } from '@/components/ui/sidebar'
import { Spinner } from '@/components/ui/spinner'
import { AppSidebar } from '@/components/layout/app-sidebar'
import SkipToMain from '@/components/skip-to-main'

interface Props {
  children?: React.ReactNode
}

export function AuthenticatedLayout({ children }: Props) {
  const defaultOpen = Cookies.get('sidebar_state') !== 'false'
  const { user, setUser } = useAuthStore((state) => state.auth)
  const navigate = useNavigate()
  const [isChecking, setIsChecking] = useState(true)

  useEffect(() => {
    const initializeAuth = async () => {
      try {
        // 如果没有用户信息，从服务端获取
        if (!user) {
          const response = await AccountApi.userAuth()
          // console.log('AccountApi.userAuth response', response)
          if (response.code === 0 && response.data) {
            setUser(response.data)
          } else {
            // 获取用户信息失败，可能token已过期
            localStorage.removeItem('token')
            navigate({ to: '/sign-in-2' })
          }
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        localStorage.removeItem('token')
        navigate({ to: '/sign-in-2' })
      } finally {
        setIsChecking(false)
      }
    }

    initializeAuth()
  }, [user, setUser, navigate])

  // 如果正在检查认证状态，显示加载指示器
  if (isChecking) {
    return (
      <div className='flex h-screen w-full items-center justify-center'>
        <Spinner size='lg' />
      </div>
    )
  }

  return (
    <SearchProvider>
      <SidebarProvider defaultOpen={defaultOpen}>
        <SkipToMain />
        <AppSidebar />
        <div
          id='content'
          className={cn(
            'ml-auto w-full max-w-full',
            'peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon)-1rem)]',
            'peer-data-[state=expanded]:w-[calc(100%-var(--sidebar-width))]',
            'sm:transition-[width] sm:duration-200 sm:ease-linear',
            'flex h-svh flex-col',
            'group-data-[scroll-locked=1]/body:h-full',
            'has-[main.fixed-main]:group-data-[scroll-locked=1]/body:h-svh'
          )}
        >
          {children ? children : <Outlet />}
        </div>
      </SidebarProvider>
    </SearchProvider>
  )
}
