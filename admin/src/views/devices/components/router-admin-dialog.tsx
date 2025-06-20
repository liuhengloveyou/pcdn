/* eslint-disable no-console */
import { useState, useEffect } from 'react'
import { toast } from 'sonner'
import { Device } from '../data/schema'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
import { getRouterAdminUrl } from '@/apis/device_api'

interface RouterAdminDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  device: Device
}

export default function RouterAdminDialog({
  open,
  onOpenChange,
  device,
}: RouterAdminDialogProps) {
  const [loading, setLoading] = useState(false)
  const [routerUrl, setRouterUrl] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    if (open && device) {
      setLoading(true)
      setError('')
      // 获取路由器管理界面URL
      getRouterAdminUrl(device.sn)
        .then((resp) => {
          if (resp.code === 0 && resp.data) {
            setRouterUrl(resp.data.url)
          } else {
            setError(resp.msg || '获取路由器管理界面失败')
            toast.error('获取路由器管理界面失败', {
              description: resp.msg,
            })
          }
        })
        .catch((err) => {
          console.error('获取路由器管理界面失败:', err)
          setError('网络错误，请重试')
          toast.error('网络错误，请重试')
        })
        .finally(() => {
          setLoading(false)
        })
    }
  }, [open, device])

  const handleOpenRouter = () => {
    if (routerUrl) {
      window.open(routerUrl, '_blank')
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>路由器管理界面</DialogTitle>
        </DialogHeader>

        <div className="py-4">
          {loading ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <span className="ml-2">正在获取路由器管理界面...</span>
            </div>
          ) : error ? (
            <div className="text-center py-4">
              <p className="text-destructive">{error}</p>
              <Button
                variant="outline"
                className="mt-4"
                onClick={() => onOpenChange(false)}
              >
                关闭
              </Button>
            </div>
          ) : (
            <div className="text-center py-4">
              <p className="mb-4">
                点击下方按钮将在新窗口打开路由器管理界面
              </p>
              <Button onClick={handleOpenRouter}>打开路由器管理界面</Button>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}