<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useUserSessionStore } from "@/stores/userSessionStore";
import { computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const UserSessionStore = useUserSessionStore();
const sessUser = computed(() => UserSessionStore.getSessionUser);

function logout() {
  document.cookie =
    "trade-sess=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;max-age=0";
  UserSessionStore.logout();
  router.replace("/login");
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" class="relative h-8 w-8 rounded-full">
        <Avatar class="h-8 w-8">
          <AvatarImage src="/avatars/02.png" alt="@moli.bot" />
          <AvatarFallback>MOLI</AvatarFallback>
        </Avatar>
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent class="w-56" align="end">
      <DropdownMenuLabel class="font-normal flex">
        <div class="flex flex-col space-y-1">
          <p class="text-sm font-medium leading-none">
            {{ sessUser?.nickname }}
          </p>
          <p class="text-xs leading-none text-muted-foreground">
            {{ sessUser?.email }}
          </p>
        </div>
      </DropdownMenuLabel>
      <DropdownMenuSeparator />
      <DropdownMenuGroup>
        <router-link
          to="/user"
          class="text-sm font-medium transition-colors hover:text-primary"
        >
          <DropdownMenuItem>
            Settings
            <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
          </DropdownMenuItem>
        </router-link>

        <router-link
          to="/syslog"
          class="text-sm font-medium transition-colors hover:text-primary"
        >
          <DropdownMenuItem> System log</DropdownMenuItem>
        </router-link>
      </DropdownMenuGroup>
      <DropdownMenuSeparator />
      <DropdownMenuItem @click="logout">
        Log out
        <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
