/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState, useEffect } from 'react'
import {
  IconCpu,
  IconNetwork,
  IconActivity,
  IconServer,
  IconRefresh,
  IconChartLine,
  IconDeviceFloppy,
} from '@tabler/icons-react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  ResponsiveContainer,
  Area,
  AreaChart,
  Tooltip,
} from 'recharts'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Device } from '../data/schema'

interface SystemStats {
  cpu: {
    usage: number
    cores: number
    temperature: number
  }
  memory: {
    used: number
    total: number
    available: number
  }
  disk: {
    used: number
    total: number
    free: number
  }
  network: {
    upload: number
    download: number
    connections: number
  }
  processes: Array<{
    pid: number
    name: string
    cpu: number
    memory: number
    status: string
  }>
}

interface HistoryData {
  time: string
  cpu: number
  memory: number
  networkUp: number
  networkDown: number
  diskUsage: number
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
  const [stats, setStats] = useState<SystemStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date())
  const [historyData, setHistoryData] = useState<HistoryData[]>([])

  // 模拟数据获取
  const fetchSystemStats = async () => {
    if (!device) return

    setLoading(true)
    // 这里应该调用实际的API，传入device.id
    setTimeout(() => {
      const newStats = {
        cpu: {
          usage: Math.floor(Math.random() * 100),
          cores: 8,
          temperature: Math.floor(Math.random() * 30) + 40,
        },
        memory: {
          used: Math.floor(Math.random() * 8) + 4,
          total: 16,
          available: Math.floor(Math.random() * 4) + 2,
        },
        disk: {
          used: Math.floor(Math.random() * 200) + 100,
          total: 500,
          free: Math.floor(Math.random() * 100) + 50,
        },
        network: {
          upload: Math.floor(Math.random() * 1000),
          download: Math.floor(Math.random() * 5000),
          connections: Math.floor(Math.random() * 100) + 20,
        },
        processes: [
          {
            pid: 1234,
            name: 'nginx',
            cpu: 15.2,
            memory: 128,
            status: 'running',
          },
          {
            pid: 5678,
            name: 'mysql',
            cpu: 8.5,
            memory: 512,
            status: 'running',
          },
          { pid: 9012, name: 'redis', cpu: 3.1, memory: 64, status: 'running' },
          {
            pid: 3456,
            name: 'node',
            cpu: 12.8,
            memory: 256,
            status: 'running',
          },
          {
            pid: 7890,
            name: 'apache',
            cpu: 6.4,
            memory: 192,
            status: 'sleeping',
          },
        ],
      }

      setStats(newStats)

      // 更新历史数据
      const now = new Date()
      const timeStr = now.toLocaleTimeString('zh-CN', { hour12: false })
      const newDataPoint: HistoryData = {
        time: timeStr,
        cpu: newStats.cpu.usage,
        memory: (newStats.memory.used / newStats.memory.total) * 100,
        networkUp: newStats.network.upload / 1024, // 转换为MB/s
        networkDown: newStats.network.download / 1024,
        diskUsage: (newStats.disk.used / newStats.disk.total) * 100,
      }

      setHistoryData((prev) => {
        const updated = [...prev, newDataPoint]
        // 保持最近60个数据点（5分钟的数据）
        return updated.slice(-60)
      })

      setLastUpdate(now)
      setLoading(false)
    }, 1000)
  }

  useEffect(() => {
    if (open && device) {
      fetchSystemStats()
      const interval = setInterval(fetchSystemStats, 5000) // 每5秒刷新
      return () => clearInterval(interval)
    }
  }, [open, device])

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

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
          {payload.map((entry: any, index: number) => (
            <p key={`item-${index}`} style={{ color: entry.stroke }}>
              {entry.name}: {entry.value.toFixed(2)}
              {entry.dataKey.includes('network') ? ' MB/s' : '%'}
            </p>
          ))}
        </div>
      )
    }
    return null
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side='bottom'
        className='h-[660px] max-h-[90vh] overflow-hidden'
      >
        <SheetHeader className='pb-4'>
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
            <div className='space-y-4'>
              {/* 监控标签页 - 修复了Tabs结构 */}
              <Tabs
                defaultValue='processes'
                className='flex h-full w-full flex-col'
              >
                <div className='mb-4 flex items-center justify-between'>
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
                      <div className='flex items-center space-x-2'>
                        <Badge variant='outline' className='text-xs'>
                          最后更新: {lastUpdate.toLocaleTimeString()}
                        </Badge>
                        <Badge
                          variant={
                            device?.status === '在线' ? 'default' : 'secondary'
                          }
                        >
                          {device?.status || '未知'}
                        </Badge>
                      </div>
                      <Button
                        variant='outline'
                        size='sm'
                        onClick={fetchSystemStats}
                        disabled={loading}
                      >
                        <IconRefresh
                          className={`mr-2 h-4 w-4 ${loading ? 'animate-spin' : ''}`}
                        />
                        刷新
                      </Button>
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
                              {stats?.cpu.usage}%
                            </div>
                            <div className='text-muted-foreground text-xs'>
                              当前使用率
                            </div>
                          </div>
                        </div>

                        <div className='mt-4 grid grid-cols-2 gap-4'>
                          <div className='bg-muted/50 rounded-lg p-3 text-center'>
                            <div className='text-lg font-semibold'>
                              {stats?.cpu.cores}
                            </div>
                            <div className='text-muted-foreground text-xs'>
                              核心数
                            </div>
                          </div>
                          <div className='bg-muted/50 rounded-lg p-3 text-center'>
                            <div className='text-lg font-semibold'>
                              {stats?.cpu.temperature}°C
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
                            {stats?.memory.used}GB
                          </div>
                          <div className='text-muted-foreground'>已使用</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.memory.total}GB
                          </div>
                          <div className='text-muted-foreground'>总容量</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.memory.available}GB
                          </div>
                          <div className='text-muted-foreground'>可用</div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className='h-[300px]'>
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
                            {stats?.disk.used}GB
                          </div>
                          <div className='text-muted-foreground'>已使用</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.disk.total}GB
                          </div>
                          <div className='text-muted-foreground'>总容量</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.disk.free}GB
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
                              width: `${((stats?.disk.used || 0) / (stats?.disk.total || 1)) * 100}%`,
                            }}
                          >
                            {Math.round(
                              ((stats?.disk.used || 0) /
                                (stats?.disk.total || 1)) *
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
                  <Card>
                    <CardHeader>
                      <CardTitle className='flex items-center space-x-2'>
                        <IconNetwork className='h-5 w-5 text-purple-500' />
                        <span>网络状态</span>
                      </CardTitle>
                      <div className='grid grid-cols-3 gap-4 text-sm'>
                        <div>
                          <div className='text-2xl font-bold'>
                            {formatBytes((stats?.network.upload || 0) * 1024)}/s
                          </div>
                          <div className='text-muted-foreground'>上传速度</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {formatBytes((stats?.network.download || 0) * 1024)}
                            /s
                          </div>
                          <div className='text-muted-foreground'>下载速度</div>
                        </div>
                        <div>
                          <div className='text-2xl font-bold'>
                            {stats?.network.connections}
                          </div>
                          <div className='text-muted-foreground'>连接数</div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className='h-[300px]'>
                        <ResponsiveContainer width='100%' height='100%'>
                          <LineChart data={historyData}>
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
                              width={50}
                            />
                            <Tooltip content={<CustomTooltip />} />
                            <Line
                              type='monotone'
                              dataKey='networkUp'
                              stroke='#3b82f6'
                              strokeWidth={2}
                              dot={false}
                              name='上传'
                            />
                            <Line
                              type='monotone'
                              dataKey='networkDown'
                              stroke='#10b981'
                              strokeWidth={2}
                              dot={false}
                              name='下载'
                            />
                          </LineChart>
                        </ResponsiveContainer>
                      </div>
                      <div className='mt-4 flex justify-center space-x-6'>
                        <div className='flex items-center space-x-2'>
                          <div className='h-3 w-3 rounded-full bg-blue-500'></div>
                          <span className='text-sm'>上传</span>
                        </div>
                        <div className='flex items-center space-x-2'>
                          <div className='h-3 w-3 rounded-full bg-green-500'></div>
                          <span className='text-sm'>下载</span>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </TabsContent>

                {/* 进程监控 */}
                <TabsContent value='processes' className='space-y-4'>
                  <Card>
                    <CardHeader>
                      <CardTitle className='flex items-center space-x-2'>
                        <IconServer className='h-5 w-5 text-purple-500' />
                        <span>运行进程</span>
                        <Badge variant='secondary' className='ml-auto'>
                          {stats?.processes.length || 0} 个进程
                        </Badge>
                      </CardTitle>
                    </CardHeader>

                    <CardContent className='h-full pt-0'>
                      <div className='space-y-2'>
                        {/* 表头 */}
                        <div className='text-muted-foreground grid grid-cols-5 gap-4 border-b pb-2 text-xs font-medium'>
                          <div>PID</div>
                          <div>进程名</div>
                          <div className='text-right'>CPU %</div>
                          <div className='text-right'>内存 (MB)</div>
                          <div className='text-center'>状态</div>
                        </div>

                        {/* 进程列表 */}
                        <ScrollArea className='h-[300px]'>
                          <div className='space-y-1'>
                            {stats?.processes.map((process) => (
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
                                <div className='text-right font-mono'>
                                  {process.cpu.toFixed(1)}%
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
                            ))}
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
