<script lang="ts" setup>
import { Download } from 'lucide-vue-next'
import { h, ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { toast } from 'vue-sonner'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'

const isOpen = ref(false)
const file = ref()
const error = ref()

watch(file, () => {
  error.value = null
})
watch(isOpen, () => {
  file.value = null
})

function onSubmit() {
  error.value = null

  if (!file.value) {
    error.value = 'File is reqred'
    return
  }

  toast({
    title: 'You submitted the following values:',
    description: h('pre', { class: 'mt-2 w-[340px] rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white' }, JSON.stringify(file.value, null, 2))),
  })
  isOpen.value = false
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button variant="outline">
        Import
        <Download />
      </Button>
    </DialogTrigger>

    <DialogContent>
      <DialogHeader>
        <DialogTitle>Import Tasks</DialogTitle>
        <DialogDescription>
          Import tasks qckly from a CSV file.
        </DialogDescription>
      </DialogHeader>

      <div class="grid w-full max-w-sm items-center gap-1.5">
        <Label>File</Label>
        <Input id="file" v-model="file" type="file" />
        <span v-if="error" class="text-destructive">{{ error }}</span>
      </div>
      <DialogFooter>
        <Button variant="secondary" @click="isOpen = false">
          Cancel
        </Button>
        <Button @click="onSubmit">
          Import
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<style scoped>

</style>
