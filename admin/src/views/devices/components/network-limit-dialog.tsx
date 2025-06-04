/* eslint-disable no-console */
import { useState } from 'react'
import { IconWifi, IconRefresh } from '@tabler/icons-react'
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
  const [uploadLimit, setUploadLimit] = useState('')
  // const [downloadLimit, setDownloadLimit] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  // 模拟当前限速规则数据
  const currentRules = {
    upload: '100 Mbps',
    download: '500 Mbps',
    status: '已启用',
    lastUpdated: '2024-01-15 14:30:25',
  }

  const handleSubmit = async () => {
    if (!device || !device.sn) return

    if (!uploadLimit) {
      toast.error('请填写上行限速值')
      return
    }

    setIsLoading(true)
    try {
      // 这里调用API更新限速规则
      console.log('更新设备限速:', {
        deviceId: device?.id,
        sn: device?.sn,
        uploadLimit: uploadLimit + ' Mbps',
        // downloadLimit: downloadLimit + ' Mbps',
      })

      const respone = await api.post<HttpResponse>(`/api/device/tc`, {
        sn: device?.sn,
      })
      if (respone.data.code === 0) {
        toast.success(`设备 ${device?.sn} 限速任务下发成功`)
      } else {
        toast.error(`设备 ${device?.sn} 限速任务下发成功`, {
          description: respone.data.msg,
        })
      }

      onOpenChange(false)
      setUploadLimit('')
      // setDownloadLimit('')
    } catch (error) {
      console.error('更新失败:', error)
      toast.error('更新失败，请重试')
    } finally {
      setIsLoading(false)
    }
  }

  const handleRefreshRules = async () => {
    setIsLoading(true)
    try {
      // 这里调用API刷新当前规则
      console.log('刷新限速规则:', device?.id)
      await new Promise((resolve) => setTimeout(resolve, 500))
      alert('规则已刷新')
    } catch (error) {
      console.error('更新失败:', error)
      alert('刷新失败')
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

        <div className='space-y-4'>
          {/* 当前限速规则 */}
          <Card>
            <CardHeader className='pb-3'>
              <div className='flex items-center justify-between'>
                <CardTitle className='text-sm'>当前限速规则</CardTitle>
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
            </CardHeader>
            <CardContent className='space-y-2'>
              <div className='flex items-center justify-between'>
                <span className='text-muted-foreground text-sm'>上行限速:</span>
                <Badge variant='secondary'>{currentRules.upload}</Badge>
              </div>
              {/* <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">下载限速:</span>
                <Badge variant="secondary">{currentRules.download}</Badge>
              </div> */}
              <div className='flex items-center justify-between'>
                <span className='text-muted-foreground text-sm'>状态:</span>
                <Badge variant='default'>{currentRules.status}</Badge>
              </div>
              <div className='text-muted-foreground text-xs'>
                最后更新: {currentRules.lastUpdated}
              </div>
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
                onChange={(e) => setUploadLimit(e.target.value)}
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
            {isLoading ? '设置中...' : '设置限速'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
