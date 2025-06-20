import React, { useState } from 'react'
import useDialogState from '@/hooks/use-dialog-state'
import { Device } from '../data/schema.ts'

type TasksDialogType = 'create' | 'update' | 'delete' | 'import' | 'network-limit' | 'reset-password' | 'system-monitor' | 'router-admin'

interface TasksContextType {
  open: TasksDialogType | null
  setOpen: (str: TasksDialogType | null) => void
  currentRow: Device | null
  setCurrentRow: React.Dispatch<React.SetStateAction<Device | null>>
  refreshData: () => void
}

const TasksContext = React.createContext<TasksContextType | null>(null)

interface Props {
  children: React.ReactNode
  refreshData: () => void
}

export default function TasksProvider({ children, refreshData }: Props) {
  const [open, setOpen] = useDialogState<TasksDialogType>(null)
  const [currentRow, setCurrentRow] = useState<Device | null>(null)
  return (
    <TasksContext.Provider value={{ open, setOpen, currentRow, setCurrentRow, refreshData }}>
      {children}
    </TasksContext.Provider>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export const useTasks = () => {
  const tasksContext = React.useContext(TasksContext)

  if (!tasksContext) {
    throw new Error('useTasks has to be used within <TasksContext>')
  }

  return tasksContext
}
