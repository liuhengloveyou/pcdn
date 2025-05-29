import type { ColumnDef } from '@tanstack/vue-table'
import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import { SelectColumn } from '@/components/data-table/table-columns'
import { Badge } from '@/components/ui/badge'
import { h } from 'vue'
import { labels,  statuses } from '../data/data'
import DataTableRowActions from './data-table-row-actions.vue'
import { z } from 'zod'

// We're keeping a simple non-relational schema here.
// IRL, you will have a schema for your data models.
export const deviceSchema = z.object({
  id: z.number(),
  uid: z.number(),
  sn: z.string(),
  createTime: z.number(),
  updateTime: z.number(),
  remoteAddr: z.string(),
  version: z.string(),
  timestamp: z.number(),
  lastHeartbear: z.number(),
  lastHeartbearStr: z.string(),
})

export type Device = z.infer<typeof deviceSchema>

export const columns: ColumnDef<Device>[] = [
  SelectColumn as ColumnDef<Device>,
  {
    accessorKey: 'id',
    header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: 'ID' }),
    cell: ({ row }) => h('div', { class: 'w-20' }, row.getValue('id')),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'sn',
    header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: '设备编号(SN)' }),

    cell: ({ row }) => {
      const label = labels.find(label => label.value === row.original.sn)

      return h('div', { class: 'flex space-x-2' }, [
        label ? h(Badge, { variant: 'outline' }, () => label.label) : null,
        h('span', { class: 'max-w-[500px] truncate font-medium' }, row.getValue('sn')),
      ])
    },
  },
  {
    accessorKey: 'version',
    header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: '版本' }),
    cell: ({ row }) => h('div', { class: 'w-20' }, row.getValue('version')),
    enableSorting: false,
    enableHiding: true,
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: '状态' }),

    cell: ({ row }) => {
      const status = statuses.find(
        status => status.value === row.getValue('status'),
      )

      if (!status)
        return null

      return h('div', { class: 'flex w-[100px] items-center' }, [
        status.icon && h(status.icon, { 
          class: `mr-2 h-4 w-4 ${status.value === '在线' ? 'text-green-500' : 'text-red-500'}` 
        }),
        h('span', status.label),
      ])
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: 'lastHeartbearStr',
    header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: '最后活跃时间' }),
    cell: ({ row }) => h('div', { class: 'w-20' }, row.getValue('lastHeartbearStr')),
    enableSorting: true,
    enableHiding: true,
  },
  // {
  //   accessorKey: 'priority',
  //   header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: 'Priority' }),
  //   cell: ({ row }) => {
  //     const priority = priorities.find(
  //       priority => priority.value === row.getValue('priority'),
  //     )

  //     if (!priority)
  //       return null

  //     return h('div', { class: 'flex items-center' }, [
  //       priority.icon && h(priority.icon, { class: 'mr-2 h-4 w-4 text-muted-foreground' }),
  //       h('span', {}, priority.label),
  //     ])
  //   },
  //   filterFn: (row, id, value) => {
  //     return value.includes(row.getValue(id))
  //   },
  // },
  {
    id: 'actions',
    header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: '操作' }),
    cell: ({ row }) => h(DataTableRowActions, { row }),
  },
]
