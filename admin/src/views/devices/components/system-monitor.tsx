/* eslint-disable no-console */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState, useEffect } from 'react'
import {
  IconCpu,
  IconNetwork,
  IconActivity,
  IconServer,
  IconChartLine,
  IconDeviceFloppy,
} from '@tabler/icons-react'
import { getDeviceMonitor, SystemMonitorData } from '@/apis/device_api'
import {
  XAxis,
  YAxis,
  CartesianGrid,
  ResponsiveContainer,
  Area,
  AreaChart,
  Tooltip,
} from 'recharts'
import { toast } from 'sonner'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Device } from '../data/schema'

interface HistoryData {
  time: string
  cpu: number
  memory: number
  diskUsage: number
  // 动态网络接口数据
  [key: string]: string | number
}

interface NetworkHistoryData {
  time: string
  [interfaceName: string]: string | number // 支持动态网络接口名称
}

interface SystemMonitorDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  device: Device | null
}

export default function SystemMonitorDrawer({
  open,
  onOpenChange,
  device,
}: SystemMonitorDrawerProps) {
  const [stats, setStats] = useState<SystemMonitorData | null>(null)
  const [loading, setLoading] = useState(true)
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date())
  const [historyData, setHistoryData] = useState<HistoryData[]>([])
  const [networkHistoryData, setNetworkHistoryData] = useState<
    NetworkHistoryData[]
  >([])

  // 格式化网络速度单位
  const formatNetworkSpeed = (bytesPerSecond: number): string => {
    if (bytesPerSecond === 0) return '0 B/s'

    const units = ['B/s', 'KB/s', 'MB/s', 'GB/s']
    const k = 1024
    const i = Math.floor(Math.log(bytesPerSecond) / Math.log(k))

    if (i >= units.length) {
      return (
        (bytesPerSecond / Math.pow(k, units.length - 1)).toFixed(2) +
        ' ' +
        units[units.length - 1]
      )
    }

    return (
      (bytesPerSecond / Math.pow(k, i)).toFixed(i === 0 ? 0 : 2) +
      ' ' +
      units[i]
    )
  }

  // 自动刷新
  const [isAutoRefresh, setIsAutoRefresh] = useState(true)
  const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(
    null
  )

  // 从API获取系统监控数据
  // 在fetchSystemStats函数中添加调试日志
  const fetchSystemStats = async () => {
    if (!device) return
    setLoading(true)

    try {
      const response = await getDeviceMonitor(device.sn)
      if (response.code === 0 && response.data) {
        setStats(response.data)
        // 更新历史数据
        const now = new Date()
        const timeStr = now.toLocaleTimeString('zh-CN', { hour12: false })

        const newDataPoint: HistoryData = {
          time: timeStr,
          cpu: response.data.cpu?.usage || 0,
          memory: response.data.memory
            ? (response.data.memory.used / response.data.memory.total) * 100
            : 0,
          diskUsage: response.data.disk
            ? (response.data.disk.used / response.data.disk.total) * 100
            : 0,
        }

        setHistoryData((prev) => {
          const updated = [...prev, newDataPoint]
          // 保持最近60个数据点（5分钟的数据）
          return updated.slice(-60)
        })

        // 更新网络历史数据
        if (response.data.network && response.data.network.length > 0) {
          const networkDataPoint: NetworkHistoryData = {
            time: timeStr,
          }

          // 为每个网络接口添加上传和下载速度数据
          response.data.network.forEach((net) => {
            networkDataPoint[`${net.name}_recv`] =
              (net.recv_rate || 0) / (1024 * 1024) // 转换为 MB/s
            networkDataPoint[`${net.name}_send`] =
              (net.send_rate || 0) / (1024 * 1024) // 转换为 MB/s
          })

          setNetworkHistoryData((prev) => {
            const updated = [...prev, networkDataPoint]
            // 保持最近60个数据点（5分钟的数据）
            return updated.slice(-60)
          })
        }

        setLastUpdate(now)
      } else {
        // 处理API返回成功但数据无效的情况
        console.warn('获取系统监控数据无效:', response)
        toast.warning('获取系统监控数据无效')
      }
    } catch (error) {
      console.error('获取系统监控数据失败:', error)
      toast.error(`获取系统监控数据失败: ${(error as Error).message}`)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (open && device) {
      fetchSystemStats()
    }
  }, [open, device])

  // 自动刷新效果
  useEffect(() => {
    if (isAutoRefresh && open && device) {
      const interval = setInterval(() => {
        fetchSystemStats()
      }, 5000) // 每5秒刷新一次
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
  }, [isAutoRefresh, open, device])

  // 组件卸载时清理定时器
  useEffect(() => {
    return () => {
      if (refreshInterval) {
        clearInterval(refreshInterval)
      }
    }
  }, [refreshInterval])

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running':
        return 'bg-green-500'
      case 'sleeping':
        return 'bg-yellow-500'
      case 'stopped':
        return 'bg-red-500'
      default:
        return 'bg-gray-500'
    }
  }

  // 自定义工具提示
  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className='bg-background rounded border p-2 text-xs shadow-md'>
          <p className='font-medium'>{label}</p>
          {payload.map((entry: any, index: number) => {
            // 确定显示单位
            let unit = '%'
            if (
              entry.name &&
              (entry.name.includes('上传') || entry.name.includes('下载'))
            ) {
              unit = ' MB/s'
            }

            return (
              <p key={`item-${index}`} style={{ color: entry.stroke }}>
                {entry.name}: {entry.value?.toFixed(2) || '0.00'}
                {unit}
              </p>
            )
          })}
        </div>
      )
    }
    return null
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side='bottom'
        className='h-[80vh] max-h-[90vh] overflow-hidden'
      >
        <SheetHeader className='pb-2'>
          <SheetTitle className='flex items-center justify-between'>
            <div className='flex items-center space-x-2'>
              <IconChartLine className='h-5 w-5' />
              <span>系统监控 - {device?.sn || '未知设备'}</span>
            </div>
          </SheetTitle>
        </SheetHeader>

        <div className='flex-1 overflow-auto'>
          {loading && !stats ? (
            <div className='flex h-96 items-center justify-center'>
              <div className='flex items-center space-x-2'>
                <IconActivity className='animate-spin' size={24} />
                <span>加载系统监控数据...</span>
              </div>
            </div>
          ) : (
            <div className='h-full'>
              {/* 监控标签页 - 修复了Tabs结构 */}
              <Tabs
                defaultValue='processes'
                className='flex h-full w-full flex-col'
              >
                <div className='flex items-center justify-between'>
                  <TabsList className='grid w-full max-w-md grid-cols-5'>
                    <TabsTrigger
                      value='processes'
                      className='flex items-center space-x-1 text-xs'
                    >
                      <IconServer className='h-3 w-3' />
                      <span className='hidden sm:inline'>进程</span>
                      <span className='sm:hidden'>PROC</span>
                    </TabsTrigger>
                    <TabsTrigger
                      value='network'
                      className='flex items-center space-x-1 text-xs'
                    >
                      <IconNetwork className='h-3 w-3' />
                      <span className='hidden sm:inline'>网络</span>
                      <span className='sm:hidden'>NET</span>
                    </TabsTrigger>

                    <TabsTrigger
                      value='cpu'
                      className='flex items-center space-x-1 text-xs'
                    >
                      <IconCpu className='h-3 w-3' />
                      <span className='hidden sm:inline'>CPU</span>
                    </TabsTrigger>
                    <TabsTrigger
                      value='memory'
                      className='flex items-center space-x-1 text-xs'
                    >
                      <IconDeviceFloppy className='h-3 w-3' />
                      <span className='hidden sm:inline'>内存</span>
                      <span className='sm:hidden'>MEM</span>
                    </TabsTrigger>
                    <TabsTrigger
                      value='disk'
                      className='flex items-center space-x-1 text-xs'
                    >
                      <IconDeviceFloppy className='h-3 w-3' />
                      <span className='hidden sm:inline'>磁盘</span>
                      <span className='sm:hidden'>DISK</span>
                    </TabsTrigger>
                  </TabsList>

                  <div className='flex items-center space-x-2'>
                    <div className='flex items-center justify-between'>
                      <div className='bg-muted/50 ml-2 flex items-center gap-2 rounded-lg border px-2 py-1'>
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
                            '自动刷新'
                          )}
                        </span>
                      </div>
                      <Badge variant='outline' className='text-xs'>
                        最后更新: {lastUpdate.toLocaleTimeString()}
                      </Badge>
                    </div>
                  </div>
                </div>

                {/* CPU 监控 - 现在在Tabs内部 */}
                <TabsContent value='cpu' className='flex-1 overflow-auto'>
                  <div className='grid h-full gap-4'>
                    <Card className='flex-1'>
                      <CardHeader className='pb-3'>
                        <div className='flex items-center justify-between'>
                          <CardTitle className='flex items-center space-x-2 text-lg'>
                            <IconCpu className='h-5 w-5 text-blue-500' />
                            <span>CPU 使用率</span>
                          </CardTitle>
                          <div className='text-right'>
                            <div className='text-2xl font-bold text-blue-500'>
                              {stats?.cpu?.usage?.toFixed(1) || '0.0'}%
                            </div>
                            <div className='text-muted-foreground text-xs'>
                              当前使用率
                            </div>
                          </div>
                        </div>

                        <div className='mt-4 grid grid-cols-2 gap-4'>
                          <div className='bg-muted/50 rounded-lg p-3 text-center'>
                            <div className='text-lg font-semibold'>
                              {stats?.cpu?.cores || 0}
                            </div>
                            <div className='text-muted-foreground text-xs'>
                              核心数
                            </div>
                          </div>
                          <div className='bg-muted/50 rounded-lg p-3 text-center'>
                            <div className='text-lg font-semibold'>
                              {stats?.cpu?.temperature || 0}°C
                            </div>
                            <div className='text-muted-foreground text-xs'>
                              温度
                            </div>
                          </div>
                        </div>
                      </CardHeader>

                      <CardContent className='pt-0'>
                        <div className='h-[250px]'>
                          <ResponsiveContainer width='100%' height='100%'>
                            <AreaChart data={historyData}>
                              <defs>
                                <linearGradient
                                  id='cpuGradient'
                                  x1='0'
                                  y1='0'
                                  x2='0'
                                  y2='1'
                                >
                                  <stop
                                    offset='5%'
                                    stopColor='#3b82f6'
                                    stopOpacity={0.8}
                                  />
                                  <stop
                                    offset='95%'
                                    stopColor='#3b82f6'
                                    stopOpacity={0.1}
                                  />
                                </linearGradient>
                              </defs>
                              <CartesianGrid
                                strokeDasharray='3 3'
                                stroke='#374151'
                                opacity={0.3}
                              />
                              <XAxis
                                dataKey='time'
                                axisLine={false}
                                tickLine={false}
                                tick={{ fontSize: 12, fill: '#6b7280' }}
                                interval='preserveStartEnd'
                              />
                              <YAxis
                                domain={[0, 100]}
                                axisLine={false}
                                tickLine={false}
                                tick={{ fontSize: 12, fill: '#6b7280' }}
                                width={40}
                              />
                              <Tooltip content={<CustomTooltip />} />
                              <Area
                                type='monotone'
                                dataKey='cpu'
                                name='CPU使用率'
                                stroke='#3b82f6'
                                strokeWidth={2}
                                fill='url(#cpuGradient)'
                              />
                            </AreaChart>
                          </ResponsiveContainer>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </TabsContent>

                {/* 内存监控 */}
                <TabsContent value='memory' className='space-y-4'>
                  <Card>
                    <CardHeader>
                      <CardTitle className='flex items-center space-x-2'>
                        <IconDeviceFloppy className='h-5 w-5 text-green-500' />
                        <span>内存使用</span>
                      </CardTitle>
                      <div className='grid grid-cols-3 gap-4 text-sm'>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.memory?.used
                              ? `${stats.memory.used}GB`
                              : '0GB'}
                          </div>
                          <div className='text-muted-foreground'>已使用</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.memory?.total
                              ? `${stats.memory.total}GB`
                              : '0GB'}
                          </div>
                          <div className='text-muted-foreground'>总容量</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.memory?.available
                              ? `${stats.memory.available}GB`
                              : '0GB'}
                          </div>
                          <div className='text-muted-foreground'>可用</div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className='h-[250px]'>
                        <ResponsiveContainer width='100%' height='100%'>
                          <AreaChart data={historyData}>
                            <defs>
                              <linearGradient
                                id='memoryGradient'
                                x1='0'
                                y1='0'
                                x2='0'
                                y2='1'
                              >
                                <stop
                                  offset='5%'
                                  stopColor='#10b981'
                                  stopOpacity={0.8}
                                />
                                <stop
                                  offset='95%'
                                  stopColor='#10b981'
                                  stopOpacity={0.1}
                                />
                              </linearGradient>
                            </defs>
                            <CartesianGrid
                              strokeDasharray='3 3'
                              stroke='#374151'
                              opacity={0.3}
                            />
                            <XAxis
                              dataKey='time'
                              axisLine={false}
                              tickLine={false}
                              tick={{ fontSize: 12, fill: '#6b7280' }}
                              interval='preserveStartEnd'
                            />
                            <YAxis
                              domain={[0, 100]}
                              axisLine={false}
                              tickLine={false}
                              tick={{ fontSize: 12, fill: '#6b7280' }}
                              width={40}
                            />
                            <Tooltip content={<CustomTooltip />} />
                            <Area
                              type='monotone'
                              dataKey='memory'
                              name='内存使用率'
                              stroke='#10b981'
                              strokeWidth={2}
                              fill='url(#memoryGradient)'
                            />
                          </AreaChart>
                        </ResponsiveContainer>
                      </div>
                    </CardContent>
                  </Card>
                </TabsContent>

                {/* 磁盘监控 */}
                <TabsContent value='disk' className='space-y-4'>
                  <Card>
                    <CardHeader>
                      <CardTitle className='flex items-center space-x-2'>
                        <IconDeviceFloppy className='h-5 w-5 text-amber-500' />
                        <span>磁盘使用</span>
                      </CardTitle>
                      <div className='grid grid-cols-3 gap-4 text-sm'>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.disk?.used ? `${stats.disk.used}GB` : '0GB'}
                          </div>
                          <div className='text-muted-foreground'>已使用</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.disk?.total
                              ? `${stats.disk.total}GB`
                              : '0GB'}
                          </div>
                          <div className='text-muted-foreground'>总容量</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.disk?.free ? `${stats.disk.free}GB` : '0GB'}
                          </div>
                          <div className='text-muted-foreground'>剩余</div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className='space-y-4'>
                        <div className='h-4 w-full rounded-full bg-gray-200 dark:bg-gray-700'>
                          <div
                            className='flex h-4 items-center justify-center rounded-full bg-amber-500 text-xs font-medium text-white transition-all duration-300'
                            style={{
                              width: `${((stats?.disk?.used || 0) / (stats?.disk?.total || 1)) * 100}%`,
                            }}
                          >
                            {Math.round(
                              ((stats?.disk?.used || 0) /
                                (stats?.disk?.total || 1)) *
                                100
                            )}
                            %
                          </div>
                        </div>

                        {/* 添加磁盘使用率曲线图 */}
                        <div className='mt-4 h-[250px]'>
                          <ResponsiveContainer width='100%' height='100%'>
                            <AreaChart data={historyData}>
                              <defs>
                                <linearGradient
                                  id='diskGradient'
                                  x1='0'
                                  y1='0'
                                  x2='0'
                                  y2='1'
                                >
                                  <stop
                                    offset='5%'
                                    stopColor='#f59e0b'
                                    stopOpacity={0.8}
                                  />
                                  <stop
                                    offset='95%'
                                    stopColor='#f59e0b'
                                    stopOpacity={0.1}
                                  />
                                </linearGradient>
                              </defs>
                              <CartesianGrid
                                strokeDasharray='3 3'
                                stroke='#374151'
                                opacity={0.3}
                              />
                              <XAxis
                                dataKey='time'
                                axisLine={false}
                                tickLine={false}
                                tick={{ fontSize: 12, fill: '#6b7280' }}
                                interval='preserveStartEnd'
                              />
                              <YAxis
                                domain={[0, 100]}
                                axisLine={false}
                                tickLine={false}
                                tick={{ fontSize: 12, fill: '#6b7280' }}
                                width={40}
                              />
                              <Tooltip content={<CustomTooltip />} />
                              <Area
                                type='monotone'
                                dataKey='diskUsage'
                                name='磁盘使用率'
                                stroke='#f59e0b'
                                strokeWidth={2}
                                fill='url(#diskGradient)'
                              />
                            </AreaChart>
                          </ResponsiveContainer>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </TabsContent>

                {/* 网络监控 */}
                <TabsContent value='network' className='space-y-4'>
                  <div className='grid gap-4'>
                    {stats?.network && stats.network.length > 0 && (
                      <Card>
                        <CardHeader>
                          <div className='flex items-center justify-between'>
                            <CardTitle className='flex items-center space-x-2'>
                              <IconNetwork className='h-5 w-5 text-blue-500' />
                              <span>网络流量监控</span>
                            </CardTitle>
                            
                            {/* 图例 */}
                            <div className='flex flex-wrap gap-4'>
                              {stats.network.map((net, index) => {
                                const colors = [
                                  { recv: '#3b82f6', send: '#f59e0b' }, // 蓝色/橙色
                                  { recv: '#10b981', send: '#ef4444' }, // 绿色/红色
                                  { recv: '#8b5cf6', send: '#f97316' }, // 紫色/橙色
                                  { recv: '#06b6d4', send: '#ec4899' }, // 青色/粉色
                                ]
                                const colorSet = colors[index % colors.length]

                                return (
                                  <div
                                    key={net.name}
                                    className='flex items-center space-x-3'
                                  >
                                    <div className='flex items-center space-x-1'>
                                      <div
                                        className='h-0.5 w-3 rounded-full'
                                        style={{
                                          backgroundColor: colorSet.recv,
                                        }}
                                      ></div>
                                      <span className='text-muted-foreground text-xs'>
                                        {net.name} 下载
                                      </span>
                                    </div>
                                    <div className='flex items-center space-x-1'>
                                      <div
                                        className='h-0.5 w-3 rounded-full'
                                        style={{
                                          backgroundColor: colorSet.send,
                                        }}
                                      ></div>
                                      <span className='text-muted-foreground text-xs'>
                                        {net.name} 上传
                                      </span>
                                    </div>
                                  </div>
                                )
                              })}
                            </div>
                          </div>
                        </CardHeader>
                        <CardContent>
                          <div className='h-[300px]'>
                            <ResponsiveContainer width='100%' height='100%'>
                              <AreaChart data={networkHistoryData}>
                                <CartesianGrid
                                  strokeDasharray='3 3'
                                  stroke='#374151'
                                  opacity={0.3}
                                />
                                <XAxis
                                  dataKey='time'
                                  axisLine={false}
                                  tickLine={false}
                                  tick={{ fontSize: 12, fill: '#6b7280' }}
                                  interval='preserveStartEnd'
                                />
                                <YAxis
                                  axisLine={false}
                                  tickLine={false}
                                  tick={{ fontSize: 12, fill: '#6b7280' }}
                                  width={60}
                                  tickFormatter={(value) =>
                                    `${value.toFixed(1)}MB`
                                  }
                                />
                                <Tooltip content={<CustomTooltip />} />

                                {/* 为每个网络接口渲染上传和下载曲线 */}
                                {stats.network.map((net, index) => {
                                  const colors = [
                                    { recv: '#3b82f6', send: '#f59e0b' }, // 蓝色/橙色
                                    { recv: '#10b981', send: '#ef4444' }, // 绿色/红色
                                    { recv: '#8b5cf6', send: '#f97316' }, // 紫色/橙色
                                    { recv: '#06b6d4', send: '#ec4899' }, // 青色/粉色
                                  ]
                                  const colorSet = colors[index % colors.length]

                                  return [
                                    <Area
                                      key={`${net.name}_recv`}
                                      type='monotone'
                                      dataKey={`${net.name}_recv`}
                                      name={`${net.name} 下载`}
                                      stroke={colorSet.recv}
                                      strokeWidth={3}
                                      fill='transparent'
                                      dot={{ r: 0 }}
                                      activeDot={{
                                        r: 4,
                                        stroke: colorSet.recv,
                                        strokeWidth: 2,
                                        fill: '#fff',
                                      }}
                                    />,
                                    <Area
                                      key={`${net.name}_send`}
                                      type='monotone'
                                      dataKey={`${net.name}_send`}
                                      name={`${net.name} 上传`}
                                      stroke={colorSet.send}
                                      strokeWidth={3}
                                      fill='transparent'
                                      dot={{ r: 0 }}
                                      activeDot={{
                                        r: 4,
                                        stroke: colorSet.send,
                                        strokeWidth: 2,
                                        fill: '#fff',
                                      }}
                                    />,
                                  ]
                                })}
                              </AreaChart>
                            </ResponsiveContainer>
                          </div>
                        </CardContent>
                      </Card>
                    )}

                    {stats?.network && stats.network.length > 0 && (
                      <Card className='p-0'>
                        <CardContent className='p-0'>
                          <div className='overflow-hidden rounded-lg border'>
                            <table className='w-full'>
                              <thead className='bg-muted/50'>
                                <tr>
                                  <th className='text-muted-foreground px-4 py-3 text-left text-sm font-medium'>
                                    接口名称
                                  </th>
                                  <th className='text-muted-foreground px-4 py-3 text-right text-sm font-medium'>
                                    下载速度
                                  </th>
                                  <th className='text-muted-foreground px-4 py-3 text-right text-sm font-medium'>
                                    上传速度
                                  </th>
                                  <th className='text-muted-foreground px-4 py-3 text-right text-sm font-medium'>
                                    总流量
                                  </th>
                                </tr>
                              </thead>
                              <tbody className='divide-y'>
                                {stats.network.map((net, index) => {
                                  const colors = [
                                    'text-blue-600 dark:text-blue-400',
                                    'text-green-600 dark:text-green-400',
                                    'text-purple-600 dark:text-purple-400',
                                    'text-orange-600 dark:text-orange-400',
                                  ]
                                  const colorClass =
                                    colors[index % colors.length]

                                  return (
                                    <tr
                                      key={net.name}
                                      className='hover:bg-muted/30 transition-colors'
                                    >
                                      <td className='px-4 py-3'>
                                        <div className='flex items-center space-x-2'>
                                          <div
                                            className={`h-2 w-2 rounded-full ${colorClass.replace('text-', 'bg-')}`}
                                          ></div>
                                          <span
                                            className={`font-medium ${colorClass}`}
                                          >
                                            {net.name}
                                          </span>
                                        </div>
                                      </td>
                                      <td className='px-4 py-3 text-right font-mono text-sm'>
                                        {formatNetworkSpeed(net.recv_rate)}
                                      </td>
                                      <td className='px-4 py-3 text-right font-mono text-sm'>
                                        {formatNetworkSpeed(net.send_rate)}
                                      </td>
                                      <td className='px-4 py-3 text-right font-mono text-sm font-semibold'>
                                        {formatNetworkSpeed(
                                          net.recv_rate + net.send_rate
                                        )}
                                      </td>
                                    </tr>
                                  )
                                })}
                              </tbody>
                            </table>
                          </div>
                        </CardContent>
                      </Card>
                    )}
                  </div>
                </TabsContent>

                {/* 进程监控 */}
                <TabsContent value='processes' className='h-full'>
                  <Card className='flex h-full flex-col'>
                    <CardHeader className='p-0'>
                      <CardTitle className='flex items-center'>
                        <Badge variant='secondary' className='ml-auto'>
                          {stats?.processes?.length || 0} 个进程
                        </Badge>
                      </CardTitle>
                    </CardHeader>

                    <CardContent className='flex-1 pt-0'>
                      <div className='space-y-2'>
                        {/* 表头 */}
                        <div className='text-muted-foreground grid grid-cols-5 gap-4 border-b pb-2 text-xs font-medium'>
                          <div>PID</div>
                          <div>进程名</div>
                          <div>EXE</div>
                          <div className='text-right'>CPU %</div>
                          <div className='text-right'>内存 (MB)</div>
                          <div className='text-center'>状态</div>
                        </div>

                        {/* 进程列表 */}
                        <ScrollArea className='h-[calc(80vh-245px)]'>
                          <div className='space-y-1'>
                            {stats?.processes?.map((process) => (
                              <div
                                key={process.pid}
                                className='hover:bg-muted/50 grid grid-cols-5 gap-4 rounded px-2 py-2 text-sm transition-colors'
                              >
                                <div className='font-mono text-xs'>
                                  {process.pid}
                                </div>
                                <div className='truncate font-medium'>
                                  {process.name}
                                </div>
                                <div className='truncate font-medium'>
                                  {process.exe}
                                </div>
                                <div className='text-right font-mono'>
                                  {process.cpu?.toFixed(1) || '0.0'}%
                                </div>
                                <div className='text-right font-mono'>
                                  {process.memory}
                                </div>
                                <div className='text-center'>
                                  <Badge
                                    variant='outline'
                                    className={`text-xs ${getStatusColor(process.status)} border-none text-white`}
                                  >
                                    {process.status}
                                  </Badge>
                                </div>
                              </div>
                            )) || (
                              <div className='text-muted-foreground py-4 text-center'>
                                暂无进程数据
                              </div>
                            )}
                          </div>
                        </ScrollArea>
                      </div>
                    </CardContent>
                  </Card>
                </TabsContent>
              </Tabs>
            </div>
          )}
        </div>
      </SheetContent>
    </Sheet>
  )
}
