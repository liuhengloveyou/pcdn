<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { Row } from '@tanstack/vue-table'
import { computed, ref, shallowRef, type Component } from 'vue'
import { Ellipsis, FilePenLine, Trash2 } from 'lucide-vue-next'
import { labels } from '../data/data'
import TaskDelete from './task-delete.vue'
import TaskResourceDialog from './task-resource-dialog.vue'
import { Button } from '@/components/ui/button'
import DropdownMenuSub from '@/components/ui/dropdown-menu/DropdownMenuSub.vue'
import DropdownMenuSubTrigger from '@/components/ui/dropdown-menu/DropdownMenuSubTrigger.vue'
import DropdownMenuSubContent from '@/components/ui/dropdown-menu/DropdownMenuSubContent.vue'
import DropdownMenuRadioGroup from '@/components/ui/dropdown-menu/DropdownMenuRadioGroup.vue'
import DropdownMenuRadioItem from '@/components/ui/dropdown-menu/DropdownMenuRadioItem.vue'
import DropdownMenuItem from '@/components/ui/dropdown-menu/DropdownMenuItem.vue'
import DropdownMenuShortcut from '@/components/ui/dropdown-menu/DropdownMenuShortcut.vue'
import { deviceSchema, type Device } from './columns'


const props = defineProps<DataTableRowActionsProps>()

interface DataTableRowActionsProps {
  row: Row<Device>
}
const task = computed(() => deviceSchema.parse(props.row.original))

const taskLabel = ref(task.value.sn)

const showComponent = shallowRef<Component | null>(null)

type TCommand = 'edit' | 'create' | 'delete'
function handleSelect(command: TCommand) {
  switch (command) {
    case 'edit':
      showComponent.value = TaskResourceDialog
      break
    case 'create':
      showComponent.value = TaskResourceDialog
      break
    case 'delete':
      showComponent.value = TaskDelete
      break
  }
}

const isOpen = ref(false)
</script>


<template>
  <Dialog v-model:open="isOpen">
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button
          variant="ghost"
          class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
        >
          <Ellipsis class="w-4 h-4" />
          <span class="sr-only">Open menu</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" class="w-[160px]">
        <DialogTrigger as-child>
          <DropdownMenuItem @select.stop="handleSelect('edit')">
            <span>Edit</span>
            <DropdownMenuShortcut> <FilePenLine class="w-4 h-4" /> </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DialogTrigger>

        <DropdownMenuItem disabled>
          Make a copy
        </DropdownMenuItem>
        <DropdownMenuItem disabled>
          Favorite
        </DropdownMenuItem>

        <DropdownMenuSeparator />

        <DropdownMenuSub>
          <DropdownMenuSubTrigger>Labels</DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            <DropdownMenuRadioGroup v-model="taskLabel">
              <DropdownMenuRadioItem v-for="label in labels" :key="label.value" :value="label.value">
                {{ label.label }}
              </DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
          </DropdownMenuSubContent>
        </DropdownMenuSub>

        <DropdownMenuSeparator />

        <DialogTrigger as-child>
          <DropdownMenuItem @select.stop="handleSelect('delete')">
            <span>Delete</span>
            <DropdownMenuShortcut> <Trash2 class="w-4 h-4" /> </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DialogTrigger>
      </DropdownMenuContent>
    </DropdownMenu>

    <DialogContent>
      <component :is="showComponent" :task="task" @close="isOpen = false" />
    </DialogContent>
  </Dialog>
</template>
