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
  remote_addr: z.string(),
  version: z.string(),
  timestamp: z.number(),
  last_heartbear: z.number(),
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
    accessorKey: 'status',
    header: ({ column }) => h(DataTableColumnHeader<Device>, { column, title: '状态' }),

    cell: ({ row }) => {
      const status = statuses.find(
        status => status.value === row.getValue('status'),
      )

      if (!status)
        return null

      return h('div', { class: 'flex w-[100px] items-center' }, [
        status.icon && h(status.icon, { class: 'mr-2 h-4 w-4 text-muted-foreground' }),
        h('span', status.label),
      ])
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
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
    cell: ({ row }) => h(DataTableRowActions, { row }),
  },
]
