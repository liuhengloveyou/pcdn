/* eslint-disable no-console */
import { useState } from 'react'
import { HttpResponse } from '@/apis'
import { toast } from 'sonner'
import api from '@/utils/axios'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useTasks } from '../context/tasks-context'
import { NetworkLimitDialog } from './network-limit-dialog'
import { TasksImportDialog } from './tasks-import-dialog'
import { TasksMutateDrawer } from './tasks-mutate-drawer'
import SystemMonitorDrawer from './system-monitor'
import RouterAdminDialog from './router-admin-dialog'

export function TasksDialogs() {
  const { open, setOpen, currentRow, setCurrentRow } = useTasks()
  const [loading, setLoading] = useState(false)

  const handleResetPassword = async () => {
    if (!currentRow) return

    setLoading(true)
    try {
      // 调用重置密码API
      const resp = await api.get<HttpResponse>(`/api/device/resetpwd`, {
        params: { sn: currentRow.sn },
      })
      if (resp.data.code === 0) {
        toast.success(`设备 ${currentRow.sn} 重置密码任务成功`)
      } else {
        toast.error(`设备 ${currentRow.sn} 重置密码任务失败`, {description: resp.data.msg})
      }

      setOpen(null)
      setCurrentRow(null)
    } catch (error) {
      console.error('重置密码失败:', error)
      toast.error('密码重置失败，请重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <TasksMutateDrawer
        key='task-create'
        open={open === 'create'}
        onOpenChange={() => setOpen('create')}
      />

      <TasksImportDialog
        key='tasks-import'
        open={open === 'import'}
        onOpenChange={() => setOpen('import')}
      />

      {currentRow && (
        <>
          <TasksMutateDrawer
            key={`task-update-${currentRow.id}`}
            open={open === 'update'}
            onOpenChange={() => {
              setOpen('update')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />

          <NetworkLimitDialog
            key={`network-limit-${currentRow.id}`}
            open={open === 'network-limit'}
            onOpenChange={(isOpen) => {
              if (!isOpen) {
                setOpen(null)
                setTimeout(() => {
                  setCurrentRow(null)
                }, 500)
              }
            }}
            device={currentRow}
          />

      <SystemMonitorDrawer
        open={open === 'system-monitor'}
        onOpenChange={(isOpen) => {
          if (!isOpen) {
            setOpen(null)
          }
        }}
        device={currentRow}
      />
          <ConfirmDialog
            key='task-delete'
            destructive
            open={open === 'delete'}
            onOpenChange={() => {
              setOpen('delete')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            handleConfirm={() => {
              setOpen(null)
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            title='确认删除设备？'
            desc='此操作无法撤销，设备数据将被永久删除。'
          />

          <ConfirmDialog
            key='reset-password'
            open={open === 'reset-password'}
            onOpenChange={() => {
              setOpen('reset-password')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            handleConfirm={handleResetPassword}
            isLoading={loading}
            cancelBtnText='取消'
            confirmText={loading ? '重置中...' : '确认重置'}
            title='确认重置密码？'
            desc={`即将重置设备 ${currentRow.sn} 的密码，重置后设备将使用默认密码。此操作无法撤销，请确认是否继续？`}
          />

          <RouterAdminDialog
            key={`router-admin-${currentRow.id}`}
            open={open === 'router-admin'}
            onOpenChange={(isOpen) => {
              if (!isOpen) {
                setOpen(null)
                setTimeout(() => {
                  setCurrentRow(null)
                }, 500)
              }
            }}
            device={currentRow}
          />
        </>
      )}
    </>
  )
}
