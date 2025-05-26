<script lang="ts" setup>
import { computed } from 'vue';
import type { Task } from '../data/schema'
import TaskForm from './task-form.vue'
import { DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';

const props = defineProps<{
  task: Task | null
}>()
defineEmits(['close'])

const task = computed(() => props.task)
const title = computed(() => task.value?.id ? `Edit Task` : 'New Task')
const description = computed(() => task.value?.id ? `Edit task ${task.value.id}` : 'Create new task')
</script>

<template>
  <div>
    <DialogHeader>
      <DialogTitle>
        {{ title }}
      </DialogTitle>
      <DialogDescription>
        {{ description }}
      </DialogDescription>
    </DialogHeader>
    <TaskForm class="mt-2" :task="task" @close="$emit('close')" />
  </div>
</template>

<style scoped>

</style>
