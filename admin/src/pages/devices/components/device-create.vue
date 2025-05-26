<script lang="ts" setup>
import { Plus } from "lucide-vue-next";
import TaskResourceDialog from "./task-resource-dialog.vue";
import { ref } from "vue";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

const isOpen = ref(false);
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogTrigger as-child>
      <Button>
        添加设备
        <Plus />
      </Button>
    </DialogTrigger>

    <DialogContent>
      <TaskResourceDialog :task="null" @close="isOpen = false" />
    </DialogContent>
  </Dialog>
</template>

<style scoped></style>
