import Checkbox from '@/components/ui/checkbox/Checkbox.vue'
import type { ColumnDef } from '@tanstack/vue-table'
import { h } from 'vue'

export const SelectColumn: ColumnDef<any> = {
  id: 'select',
  header: ({ table }: { table: any }) => h(Checkbox, {
    'modelValue': table.getIsAllPageRowsSelected() || (table.getIsSomePageRowsSelected() && 'indeterminate'),
    'onUpdate:modelValue': (value: any) => table.toggleAllPageRowsSelected(!!value),
    'ariaLabel': 'Select all',
  }),
  cell: ({ row }: { row: any }) => h(Checkbox, {
    'modelValue': row.getIsSelected(),
    'onUpdate:modelValue': (value: any) => row.toggleSelected(!!value),
    'ariaLabel': 'Select row',
  }),
  enableSorting: false,
  enableHiding: false,
}
