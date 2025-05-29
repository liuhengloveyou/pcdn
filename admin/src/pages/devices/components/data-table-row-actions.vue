<script setup lang="ts">
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import type { Row } from "@tanstack/vue-table";
import { computed, ref, shallowRef, type Component } from "vue";
import { Ellipsis, FilePenLine, Trash2 } from "lucide-vue-next";
import { labels } from "../data/data";
import TaskDelete from "./task-delete.vue";
import TaskResourceDialog from "./task-resource-dialog.vue";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from "@/components/ui/dropdown-menu";
import { deviceSchema, type Device } from "./columns";
import { Toaster } from "@/components/ui/sonner";
import "vue-sonner/style.css"; // vue-sonner v2 requires this import
import { toast } from "vue-sonner";
import api from "@/axios";
import type { HttpResponse } from "@/services";

const props = defineProps<DataTableRowActionsProps>();

interface DataTableRowActionsProps {
  row: Row<Device>;
}
const task = computed(() => deviceSchema.parse(props.row.original));

const isOpen = ref(false);
const taskLabel = ref(task.value.sn);
const showComponent = shallowRef<Component | null>(null);

type TCommand = "edit" | "create" | "delete";
function handleSelect(command: TCommand) {
  switch (command) {
    case "edit":
      showComponent.value = TaskResourceDialog;
      break;
    case "create":
      showComponent.value = TaskResourceDialog;
      break;
    case "delete":
      showComponent.value = TaskDelete;
      break;
  }
}

async function handleResetPWD() {
  console.log("reset pwd", props.row.getValue("sn"));

  const resp = await api.get<HttpResponse>(`/api/device/resetpwd`, {
    params: { sn: props.row.getValue("sn") },
  });
  if (resp.data.code === 0) {
    toast.success("重置密码任务下发成功");
  } else {
    toast.error("重置密码任务下发失败");
  }
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="ghost"
        class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
      >
        <Ellipsis class="w-4 h-4" />
        <span class="sr-only">打开菜单</span>
      </Button>
    </DropdownMenuTrigger>

    <DropdownMenuContent class="w-56">
      <DropdownMenuGroup>
        <DropdownMenuItem @select.stop="handleResetPWD">
          <span>重置密码</span>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <span>Profile</span>
          <DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <span>Billing</span>
          <DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <span>Settings</span>
          <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <span>Keyboard shortcuts</span>
          <DropdownMenuShortcut>⌘K</DropdownMenuShortcut>
        </DropdownMenuItem>
      </DropdownMenuGroup>
      <DropdownMenuSeparator />
      <DropdownMenuGroup>
        <DropdownMenuItem>
          <span>Team</span>
        </DropdownMenuItem>
        <DropdownMenuSub>
          <DropdownMenuSubTrigger>
            <span>Invite users</span>
          </DropdownMenuSubTrigger>
          <DropdownMenuPortal>
            <DropdownMenuSubContent>
              <DropdownMenuItem>
                <span>Email</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <span>Message</span>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <span>More...</span>
              </DropdownMenuItem>
            </DropdownMenuSubContent>
          </DropdownMenuPortal>
        </DropdownMenuSub>
        <DropdownMenuItem>
          <span>New Team</span>
          <DropdownMenuShortcut>⌘+T</DropdownMenuShortcut>
        </DropdownMenuItem>
      </DropdownMenuGroup>
      <DropdownMenuSeparator />
      <DropdownMenuItem>
        <span>GitHub</span>
      </DropdownMenuItem>
      <DropdownMenuItem>
        <span>Support</span>
      </DropdownMenuItem>
      <DropdownMenuItem disabled>
        <span>API</span>
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem>
        <span>Log out</span>
        <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>

  <DropdownMenu>
    <!-- <DropdownMenuTrigger as-child>
      <Button
        variant="ghost"
        class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
      >
        <Ellipsis class="w-4 h-4" />
        <span class="sr-only">打开菜单</span>
      </Button>
    </DropdownMenuTrigger> -->
    <DropdownMenuContent align="end" class="w-[160px]">
      <DialogTrigger as-child>
        <DropdownMenuItem @select.stop="handleSelect('edit')">
          <span>编辑</span>
          <DropdownMenuShortcut>
            <FilePenLine class="w-4 h-4" />
          </DropdownMenuShortcut>
        </DropdownMenuItem>
      </DialogTrigger>

      <!-- <DropdownMenuItem disabled>
          Make a copy
        </DropdownMenuItem>
        <DropdownMenuItem disabled>
          Favorite
        </DropdownMenuItem> -->

      <DropdownMenuSeparator />

      <DropdownMenuSub>
        <DropdownMenuSubTrigger>标签</DropdownMenuSubTrigger>
        <DropdownMenuSubContent>
          <DropdownMenuRadioGroup v-model="taskLabel">
            <DropdownMenuRadioItem
              v-for="label in labels"
              :key="label.value"
              :value="label.value"
            >
              {{ label.label }}
            </DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuSubContent>
      </DropdownMenuSub>

      <DialogTrigger as-child> </DialogTrigger>

      <DropdownMenuSeparator />
      <DialogTrigger as-child>
        <DropdownMenuItem @select.stop="handleSelect('delete')">
          <span>删除</span>
          <DropdownMenuShortcut>
            <Trash2 class="w-4 h-4" />
          </DropdownMenuShortcut>
        </DropdownMenuItem>
      </DialogTrigger>
    </DropdownMenuContent>
  </DropdownMenu>

  <Dialog v-model:open="isOpen">
    <DialogContent>
      <component :is="showComponent" :task="task" @close="isOpen = false" />
    </DialogContent>
  </Dialog>

  <Toaster />
</template>
