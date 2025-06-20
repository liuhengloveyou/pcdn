/* eslint-disable no-console */
import { useState, useEffect, useCallback } from 'react'
import { type DeviceModel, listDevice } from '@/apis/device_api'
import { toast } from 'sonner'
import { Switch } from '@/components/ui/switch'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { columns } from './components/columns'
import { DataTable } from './components/data-table'
import { TasksDialogs } from './components/tasks-dialogs'
import { TasksPrimaryButtons } from './components/tasks-primary-buttons'
import TasksProvider from './context/tasks-context'

export default function Device() {
  const [devices, setDevices] = useState<DeviceModel[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [currentPage, setCurrentPage] = useState<number>(1)
  const [pageSize, setPageSize] = useState<number>(20)
  const [totalItems, setTotalItems] = useState<number>(0)

  // 自动刷新
  const [isAutoRefresh, setIsAutoRefresh] = useState(false)
  const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(
    null
  )

  const onLoad = useCallback(async () => {
    setLoading(true)
    try {
      const resp = await listDevice(currentPage, pageSize)
      console.log('resp', resp)
      setLoading(false)

      if (!resp) {
        toast.error('网络错误')
        return
      }

      if (resp.code !== 0) {
        toast.error(resp.msg)
        return
      }

      if (resp.code === 0 && resp.data) {
        const processedData = resp.data.map((device: DeviceModel) => ({
          ...device,
          status: Date.now() - device.lastHeartbear <= 60000 ? '在线' : '离线',
          lastHeartbearStr:
            device.lastHeartbear === 0
              ? '-'
              : new Date(device.lastHeartbear).toLocaleString(),
        }))
        setTotalItems(resp.total as number)
        setDevices(processedData)
      }
    } catch (error) {
      setLoading(false)
      toast.error('加载设备数据失败')
      console.error('Failed to load devices:', error)
    }
  }, [currentPage, pageSize])

  useEffect(() => {
    onLoad()
  }, [currentPage, onLoad, pageSize]) // 依赖 currentPage 和 pageSize，当它们变化时重新加载数据

  useEffect(() => {
    if (isAutoRefresh) {
      const interval = setInterval(() => {
        onLoad()
      }, 5000) // Refresh every 5 seconds
      setRefreshInterval(interval)

      return () => {
        clearInterval(interval)
      }
    } else {
      if (refreshInterval) {
        clearInterval(refreshInterval)
        setRefreshInterval(null)
      }
    }
  }, [isAutoRefresh, onLoad]) // 移除 refreshInterval 依赖，避免无限循环

  // 组件卸载时清理定时器
  useEffect(() => {
    return () => {
      if (refreshInterval) {
        clearInterval(refreshInterval)
      }
    }
  }, [refreshInterval])

  return (
    <TasksProvider refreshData={onLoad}>
      <Header fixed>
        <div className='ml-auto flex items-center space-x-4'>
          <Search />
          <ThemeSwitch />
          {/* <ProfileDropdown /> */}
        </div>
      </Header>

      <Main>
        <div className='mb-2 flex flex-wrap items-center justify-between space-y-2 gap-x-4'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>我的设备</h2>
            {/* <p className='text-muted-foreground'>
              这是您本月任务的列表！
            </p> */}
          </div>
          <div className='flex items-center gap-4'>
            <div className='bg-muted/50 flex items-center gap-3 rounded-lg border px-2 py-2'>
              <Switch
                id='auto-refresh'
                checked={isAutoRefresh}
                onCheckedChange={setIsAutoRefresh}
                className='data-[state=checked]:bg-green-600'
              />
              <span className='text-muted-foreground text-xs'>
                {isAutoRefresh ? (
                  <span className='flex items-center gap-1'>
                    <span className='h-1 w-1 animate-ping rounded-full bg-green-500'></span>
                    每5秒刷新
                  </span>
                ) : (
                  '自动刷新关闭'
                )}
              </span>
            </div>
            <TasksPrimaryButtons />
          </div>
        </div>
        <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
          <DataTable
            data={devices}
            columns={columns}
            loading={loading}
            totalItems={totalItems}
            currentPage={currentPage}
            pageSize={pageSize}
            onPageChange={setCurrentPage}
            onPageSizeChange={setPageSize}
          />
        </div>
      </Main>

      <TasksDialogs />
    </TasksProvider>
  )
}
