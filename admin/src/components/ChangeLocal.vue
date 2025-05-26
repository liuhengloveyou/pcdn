<script setup lang="ts">
import { Button } from "@/components/ui/button";
import { Icon } from "@iconify/vue";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { useI18n } from "vue-i18n";
import { onMounted } from "vue";

defineOptions({
  name: "ChangeLocal",
});
let { locale } = useI18n({ useScope: "global" });

onMounted(() => {
  const val = localStorage.getItem("localValue");
  if (val) {
    locale.value = val;
  }
});

function changeLocal(v: string) {
  console.log(locale.value);
  if (v === "en-US") {
    locale.value = "en-US";
  } else if (v === "zh-Hant") {
    locale.value = "zh-Hant";
  }

  localStorage.setItem("localValue", locale.value);
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost">
        <Icon
          icon="ic:baseline-g-translate"
          class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"
        />
        <Icon
          icon="ic:baseline-g-translate"
          class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"
        />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuItem @click="changeLocal('en-US')">
        ENGLISH
      </DropdownMenuItem>
      <DropdownMenuItem @click="changeLocal('zh-Hant')">
        繁體中文
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
