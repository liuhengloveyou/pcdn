/* eslint-disable no-console */
import { useState, useEffect } from 'react'
import { IconWifi, IconRefresh, IconChevronDown, IconChevronUp } from '@tabler/icons-react'
import { HttpResponse } from '@/apis'
import { toast } from 'sonner'
import api from '@/utils/axios'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Device } from '../data/schema'

interface NetworkLimitDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  device: Device | null
}

export function NetworkLimitDialog({
  open,
  onOpenChange,
  device,
}: NetworkLimitDialogProps) {
  const [uploadLimit, setUploadLimit] = useState(100)
  // const [downloadLimit, setDownloadLimit] = useState('')
  const [ifaceName, setIfaceName] = useState('eth0')
  const [isLoading, setIsLoading] = useState(false)

  // 当前限速规则数据
  const [currentRules, setCurrentRules] = useState({
    upload: '-- Mbps',
    detail: '',
    lastUpdated: '--',
  })
  
  // 控制详情展开/折叠的状态
  const [detailExpanded, setDetailExpanded] = useState(false)

  // 当对话框打开时，自动获取当前限速规则
  useEffect(() => {
    if (open && device?.sn) {
      handleRefreshRules()
    }
  }, [open, device?.sn])

  const handleSubmit = async () => {
    if (!device || !device.sn) return

    if (!uploadLimit) {
      toast.error('请填写上行限速值')
      return
    }

    // 这里调用API更新限速规则
    console.log('更新设备限速:', {
      deviceId: device?.id,
      sn: device?.sn,
      uploadLimit: uploadLimit // + ' Mbps',
      // downloadLimit: downloadLimit + ' Mbps',
    })

    setIsLoading(true)
    try {
      const respone = await api.post<HttpResponse>(`/api/device/tc`, {
        sn: device?.sn,
        uploadLimit: uploadLimit,
      })

      if (respone.data.code === 0) {
        toast.success(`${JSON.stringify(respone.data.data)}`)
      } else {
        toast.error(`限速失败`, {
          description: respone.data.msg,
        })
      }

      handleRefreshRules();
    } catch (error) {
      console.error('更新失败:', error)
      toast.error('更新失败，请重试')
    } finally {
      setIsLoading(false)
    }
  }

  const handleRefreshRules = async () => {
    if (!device || !device.sn) return

    setIsLoading(true)
    try {
      // 调用API获取当前限速规则
      const response = await api.post<HttpResponse>(`/api/device/tc/stat`, {
        sn: device.sn,
        ifaceName: ifaceName,
      })

      if (response.data.code === 0 && response.data.data) {
        // 更新当前规则显示
        const tcData = response.data.data
        // 根据返回数据更新currentRules
        setCurrentRules({
            upload: tcData['val'] ? `${tcData['val']}` : '-- Mbps',
            detail: tcData['detail'] ? (tcData['detail'] as string).replace(/\\n/g, '\n') : '',
            lastUpdated: new Date().toLocaleString('zh-CN'),
          })
        console.log('获取到限速规则:', JSON.stringify(tcData, null, 2))
      } else {
        toast.error('获取限速规则失败', {
          description: response.data.msg || '未知错误',
        })
      }
    } catch (error) {
      console.error('刷新失败:', error)
      toast.error('刷新限速规则失败')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-w-md'>
        <DialogHeader>
          <DialogTitle className='flex items-center gap-2'>
            <IconWifi size={20} />
            网卡限速设置
          </DialogTitle>
          <DialogDescription>
            为设备
            <span className='font-bold'>{device?.sn}</span>
            设置网络限速规则
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-6'>
          <Card>
            <CardHeader>
              <CardTitle className='text-sm'>当前限速规则</CardTitle>
            </CardHeader>
            <CardContent className='space-y-3'>
              <Label
                htmlFor='iface-name'
                className='text-muted-foreground text-sm'
              >
                网口名称
              </Label>
              <div className='flex items-center gap-2'>
                <Input
                  id='iface-name'
                  type='text'
                  placeholder='请输入网口名称'
                  value={ifaceName}
                  onChange={(e) => setIfaceName(e.target.value)}
                  className='flex-1'
                />
                <Button
                  variant='outline'
                  size='sm'
                  onClick={handleRefreshRules}
                  disabled={isLoading}
                >
                  <IconRefresh
                    size={14}
                    className={isLoading ? 'animate-spin' : ''}
                  />
                  刷新
                </Button>
              </div>

              <div className='flex items-center justify-between'>
                <span className='text-muted-foreground text-sm'>上行限速:</span>
                <Badge variant='secondary'>{currentRules.upload}</Badge>
              </div>
              {/* <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">下载限速:</span>
                <Badge variant="secondary">{currentRules.download}</Badge>
              </div> */}

              <div className='text-muted-foreground text-xs'>
                最后更新: {currentRules.lastUpdated}
              </div>
              
              {currentRules.detail && (
                <div className='mt-2'>
                  <button
                    type='button'
                    onClick={() => setDetailExpanded(!detailExpanded)}
                    className='flex items-center text-xs text-muted-foreground hover:text-primary transition-colors'
                  >
                    {detailExpanded ? (
                      <>
                        <IconChevronUp size={14} className='mr-1' />
                        收起详情
                      </>
                    ) : (
                      <>
                        <IconChevronDown size={14} className='mr-1' />
                        查看详情
                      </>
                    )}
                  </button>
                  
                  {detailExpanded && (
                    <div className='mt-2 p-2 bg-muted rounded-md text-xs font-mono whitespace-pre-wrap overflow-auto max-h-40'>
                      {currentRules.detail}
                    </div>
                  )}
                </div>
              )}
            </CardContent>
          </Card>

          {/* 新限速设置 */}
          <div className='space-y-3'>
            <div className='space-y-2'>
              <Label htmlFor='upload-limit'>上行限速 (Mbps)</Label>
              <Input
                id='upload-limit'
                type='number'
                placeholder='请输入上行限速'
                value={uploadLimit}
                onChange={(e) => setUploadLimit(Number(e.target.value))}
              />
            </div>
            {/* <div className="space-y-2">
              <Label htmlFor="download-limit">下载限速 (Mbps)</Label>
              <Input
                id="download-limit"
                type="number"
                placeholder="请输入下载限速"
                value={downloadLimit}
                onChange={(e) => setDownloadLimit(e.target.value)}
              />
            </div> */}
          </div>
        </div>

        <DialogFooter>
          <Button variant='outline' onClick={() => onOpenChange(false)}>
            取消
          </Button>
          <Button onClick={handleSubmit} disabled={isLoading}>
            {isLoading ? '加载中...' : '设置限速'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
