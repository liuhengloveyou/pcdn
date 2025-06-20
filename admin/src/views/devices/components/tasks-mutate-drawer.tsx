/* eslint-disable no-console */
import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { addDevice } from '@/apis/device_api.ts'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { useTasks } from '../context/tasks-context'
import { Device } from '../data/schema.ts'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentRow?: Device
}

const formSchema = z.object({
  sn: z.string().min(1, '设备SN是必需的。'),
  label: z.string().min(1, '请选择标签。'),
})
type DeviceForm = z.infer<typeof formSchema>

export function TasksMutateDrawer({ open, onOpenChange, currentRow }: Props) {
  const isUpdate = !!currentRow
  const [isSubmitting, setIsSubmitting] = useState(false)
  const { refreshData } = useTasks()

  const form = useForm<DeviceForm>({
    resolver: zodResolver(formSchema),
    defaultValues: currentRow ?? {
      sn: '',
      label: '',
    },
  })

  const onSubmit = async (data: DeviceForm) => {
    setIsSubmitting(true)
    try {
      const response = await addDevice({
        sn: data.sn,
        label: data.label,
        id: 0,
        uid: 0,
        createTime: 0,
        updateTime: 0,
        remoteAddr: '',
        version: '',
        timestamp: 0,
        lastHeartbear: 0,
        lastHeartbearStr: '',
      })
      console.log('添加设备成功:', response)
      if (response.code!== 0) {
        toast.error('添加失败', {description: response.msg})
        return
      }

      toast.success('添加成功')
      onOpenChange(false)
      form.reset()
      // 刷新表格数据
      refreshData()
    } catch (error) {
      console.error('添加失败:', error)
      // 这里可以添加错误处理，比如显示toast
      toast.error('添加失败')
    } finally {
      setIsSubmitting(false)
    }
  }

  // 处理Sheet关闭事件，在提交过程中阻止关闭
  const handleOpenChange = (newOpen: boolean) => {
    // 如果正在提交，阻止关闭
    if (isSubmitting && !newOpen) {
      return
    }
    onOpenChange(newOpen)
    if (!newOpen) {
      form.reset()
    }
  }

  return (
    <Sheet open={open} onOpenChange={handleOpenChange}>
      <SheetContent className='flex flex-col'>
        <SheetHeader className='text-left'>
          <SheetTitle>{isUpdate ? '更新' : '创建'}设备</SheetTitle>
          <SheetDescription>
            {isUpdate
              ? '通过提供必要信息来更新设备。'
              : '通过提供必要信息来添加新设备。'}
            完成后点击保存。
          </SheetDescription>
        </SheetHeader>
        <Form {...form}>
          <form
            id='tasks-form'
            onSubmit={form.handleSubmit(onSubmit)}
            className='flex-1 space-y-5 px-4'
          >
            <FormField
              control={form.control}
              name='sn'
              render={({ field }) => (
                <FormItem className='space-y-1'>
                  <FormLabel>设备SN</FormLabel>
                  <FormControl>
                    <Input {...field} placeholder='请输入设备SN' />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='label'
              render={({ field }) => (
                <FormItem className='relative space-y-3'>
                  <FormLabel>标签</FormLabel>
                  <FormControl>
                    <RadioGroup
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      className='flex flex-col space-y-1'
                    >
                      <FormItem className='flex items-center space-y-0 space-x-3'>
                        <FormControl>
                          <RadioGroupItem value='documentation' />
                        </FormControl>
                        <FormLabel className='font-normal'>机顶盒</FormLabel>
                      </FormItem>
                      <FormItem className='flex items-center space-y-0 space-x-3'>
                        <FormControl>
                          <RadioGroupItem value='feature' />
                        </FormControl>
                        <FormLabel className='font-normal'>电源</FormLabel>
                      </FormItem>
                      <FormItem className='flex items-center space-y-0 space-x-3'>
                        <FormControl>
                          <RadioGroupItem value='bug' />
                        </FormControl>
                        <FormLabel className='font-normal'>小盒子</FormLabel>
                      </FormItem>
                    </RadioGroup>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </form>
        </Form>
        <SheetFooter className='gap-2'>
          <SheetClose asChild>
            <Button variant='outline' disabled={isSubmitting}>
              关闭
            </Button>
          </SheetClose>
          <Button form='tasks-form' type='submit' disabled={isSubmitting}>
            {isSubmitting ? (isUpdate ? '更新中...' : '创建中...') : '保存'}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  )
}
