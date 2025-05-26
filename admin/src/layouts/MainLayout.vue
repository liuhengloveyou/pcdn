<template>
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset>
      <header
        class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
      >
        <div class="flex items-center gap-2 px-4 w-full justify-between">
          <div class="flex items-center gap-2">
            <SidebarTrigger class="-ml-1" />
            <CommandMenuPanel />
            <Separator orientation="vertical" class="mr-2 h-4" />
          </div>

          <div class="flex items-center gap-2">
            <ToggleTheme />
            <ThemePopover />
          </div>
        </div>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4 pt-0">
        <RouterView />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>

<script setup lang="ts">
import CommandMenuPanel from "@/components/command-menu-panel/index.vue";
import AppSidebar from "@/components/app-sidebar/index.vue";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { onMounted } from "vue";
import ToggleTheme from "@/components/ToggleTheme.vue";
import ThemePopover from "@/components/custom-theming/theme-popover.vue";

defineOptions({
  name: "MainLayout",
});

onMounted(() => {
  console.log("MainLayout.onMounted");
});
</script>
