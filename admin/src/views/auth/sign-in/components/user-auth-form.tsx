/* eslint-disable no-console */
import { HTMLAttributes, useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Link, useNavigate } from '@tanstack/react-router'
import { AccountApi } from '@/apis/account_api'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/authStore'
import { cn } from '@/lib/utils'
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
import { PasswordInput } from '@/components/password-input'

type UserAuthFormProps = HTMLAttributes<HTMLFormElement>

const formSchema = z.object({
  cellphone: z
    .string()
    .min(1, { message: '请输入您的手机号' })
    .regex(/^1[3-9]\d{9}$/, { message: '手机号格式不正确' }),
  password: z
    .string()
    .min(1, {
      message: '请输入您的密码',
    })
    .min(6, {
      message: '密码长度至少为6个字符',
    }),
})

export function UserAuthForm({ className, ...props }: UserAuthFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const { setUser } = useAuthStore((state) => state.auth)
  const navigate = useNavigate()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      cellphone: '15360651247',
      password: '123456',
    },
  })

  async function onSubmit(data: z.infer<typeof formSchema>) {
    setIsLoading(true)
    try {
      // 调用登录API
      const response = await AccountApi.login(
        data.cellphone.trim(),
        data.password.trim(),
        ''
      )
      console.log(response)
 
      if (response.code !== 0) {
        // 登录失败
        toast("错误", {
          description: response.msg || "登录失败，请检查账号和密码",
        });
        return;
      } 

      // 保存用户信息和令牌
      setUser(response.data)
      // 显示成功消息
      toast.success('登录成功')

      // 重定向到主页
      navigate({ to: '/' })
    } catch (error) {
      console.error('登录失败:', error)
      toast.error('登录失败，请检查您的手机号和密码')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn('grid gap-3', className)}
        {...props}
      >
        <FormField
          control={form.control}
          name='cellphone'
          render={({ field }) => (
            <FormItem>
              <FormLabel>手机号</FormLabel>
              <FormControl>
                <Input placeholder='请输入手机号' {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem className='relative'>
              <FormLabel>密码</FormLabel>
              <FormControl>
                <PasswordInput placeholder='********' {...field} />
              </FormControl>
              <FormMessage />
              <Link
                to='/forgot-password'
                className='text-muted-foreground absolute -top-0.5 right-0 text-sm font-medium hover:opacity-75'
              >
                忘记密码?
              </Link>
            </FormItem>
          )}
        />
        <Button className='mt-2' disabled={isLoading}>
          {isLoading ? '登录中...' : '登录'}
        </Button>
      </form>
    </Form>
  )
}
