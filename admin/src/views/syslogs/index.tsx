/* eslint-disable no-console */
import { useState, useEffect, useCallback } from 'react'
import { BusinessLogApi, type BusinessLogModel } from '@/apis/businessLog_api'
import { toast } from 'sonner'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'

export default function Syslogs() {
  const [logs, setLogs] = useState<BusinessLogModel[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [activeTab, setActiveTab] = useState<string>('CREATE_DEVICE')
  const [currentPage, setCurrentPage] = useState<number>(1)
  const [itemsPerPage] = useState<number>(30)

  // 定义业务类型
  const businessTypes = [
    { value: 'CREATE_DEVICE', label: '绑定设备' },
    { value: 'TC', label: '限速' },
  ]

  const loadLogs = useCallback(
    async (businessType: string, page: number = 1) => {
      setLoading(true)
      try {
        const resp = await BusinessLogApi.List(businessType)

        if (!resp) {
          toast.error('网络错误')
          setLoading(false)
          return
        }

        if (resp.code !== 0) {
          toast.error(resp.msg)
          setLoading(false)
          return
        }

        if (resp.code === 0 && resp.data) {
          // 处理数据，添加格式化时间
          const processedData = resp.data.map((log: BusinessLogModel) => ({
            ...log,
            createTimeStr: new Date(log.createTime).toLocaleString(),
          }))

          setLogs(processedData)
          setCurrentPage(page)
        }
      } catch (error) {
        toast.error('加载日志数据失败')
        console.error('Failed to load logs:', error)
      } finally {
        setLoading(false)
      }
    },
    []
  )

  // 初始加载
  useEffect(() => {
    loadLogs(activeTab, 1)
  }, [activeTab, loadLogs])

  const handlePageChange = (page: number) => {
    loadLogs(activeTab, page)
  }

  const handleRefresh = () => {
    loadLogs(activeTab, currentPage)
  }

  const handleExport = () => {
    if (logs.length === 0) {
      toast.error('没有数据可导出')
      return
    }

    const csvContent = [
      ['ID', '用户ID', '业务类型', '创建时间'].join(','),
      ...logs.map((log) =>
        [log.id, log.uid, log.businessType, log.createTimeStr].join(',')
      ),
    ].join('\n')

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `business_logs_${new Date().getTime()}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    toast.success('数据导出成功')
  }

  const handleTabChange = (value: string) => {
    setActiveTab(value)
  }

  return (
    <>
      <Header fixed>
        <div className='mb-6'>
          <h2 className='text-2xl font-bold tracking-tight'>业务日志</h2>
        </div>

        <div className='ml-auto flex items-center space-x-4'>
          <Search />
          <ThemeSwitch />
        </div>
      </Header>

      <Main>
        <div className='space-y-6'>
          <div className='flex flex-col space-y-4'>
            <div className='flex items-center justify-between'>
              <div className='w-full'>
                <div className='mb-4 flex items-center justify-between'>
                  <Tabs
                    value={activeTab}
                    onValueChange={handleTabChange}
                    className='w-auto'
                  >
                    <TabsList>
                      {businessTypes.map((type) => (
                        <TabsTrigger key={type.value} value={type.value}>
                          {type.label}
                        </TabsTrigger>
                      ))}
                    </TabsList>
                  </Tabs>

                  <div className='flex gap-2'>
                    <Button
                      variant='secondary'
                      size='sm'
                      onClick={handleRefresh}
                      className='border-blue-200 bg-blue-50 text-blue-700 hover:bg-blue-100'
                    >
                      刷新
                    </Button>
                    <Button
                      variant='secondary'
                      size='sm'
                      onClick={handleExport}
                      className='border-green-200 bg-green-50 text-green-700 hover:bg-green-100'
                    >
                      导出
                    </Button>
                  </div>
                </div>

                {loading && (
                  <div className='mb-4 flex items-center gap-2 p-2 text-sm'>
                    <div className='h-2 w-2 animate-pulse rounded-full bg-yellow-500'></div>
                    <span>加载中...</span>
                  </div>
                )}
              </div>
            </div>

            <div className='rounded-md border'>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className='w-[80px]'>ID</TableHead>
                    <TableHead>用户ID</TableHead>
                    <TableHead>业务类型</TableHead>
                    <TableHead>创建时间</TableHead>
                    <TableHead>操作详情</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {logs.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={4} className='h-24 text-center'>
                        {loading ? '正在加载数据...' : '暂无数据'}
                      </TableCell>
                    </TableRow>
                  ) : (
                    // 计算当前页显示的数据
                    logs
                      .slice(
                        (currentPage - 1) * itemsPerPage,
                        currentPage * itemsPerPage
                      )
                      .map((log) => (
                        <TableRow key={log.id}>
                          <TableCell className='font-medium'>
                            {log.id}
                          </TableCell>
                          <TableCell>{log.uid}</TableCell>
                          <TableCell>
                            <Badge variant='secondary'>
                              {log.businessType}
                            </Badge>
                          </TableCell>
                          <TableCell>{log.createTimeStr}</TableCell>
                          <TableCell>{log.payload}</TableCell>
                        </TableRow>
                      ))
                  )}
                </TableBody>
              </Table>
            </div>

            <div className='mt-4 flex flex-nowrap items-center justify-between gap-4'>
              {logs.length > 0 && (
                <Pagination>
                  <PaginationContent>
                    <PaginationItem>
                      <PaginationPrevious
                        onClick={() => handlePageChange(currentPage - 1)}
                        className={
                          currentPage === 1
                            ? 'pointer-events-none opacity-50'
                            : ''
                        }
                      ></PaginationPrevious>
                    </PaginationItem>

                    {Array.from(
                      { length: Math.ceil(logs.length / itemsPerPage) },
                      (_, i) => i + 1
                    ).map((page) => (
                      <PaginationItem key={page}>
                        <PaginationLink
                          onClick={() => handlePageChange(page)}
                          isActive={currentPage === page}
                        >
                          {page}
                        </PaginationLink>
                      </PaginationItem>
                    ))}

                    <PaginationItem>
                      <PaginationNext
                        onClick={() => handlePageChange(currentPage + 1)}
                        className={
                          currentPage === Math.ceil(logs.length / itemsPerPage)
                            ? 'pointer-events-none opacity-50'
                            : ''
                        }
                      ></PaginationNext>
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              )}

              <div
                className='text-muted-foreground text-sm'
                style={{ width: '200px' }}
              >
                {logs.length === 0
                  ? loading
                    ? '正在加载数据...'
                    : '暂无日志数据'
                  : `共 ${logs.length} 条记录，第 ${currentPage} 页`}
              </div>
            </div>
          </div>
        </div>
      </Main>
    </>
  )
}
