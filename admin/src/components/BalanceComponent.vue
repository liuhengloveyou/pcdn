<script setup lang="ts">
import { ScrollArea } from "@/components/ui/scroll-area";
import type { MartParamModel } from "@/services/MartParamService";
import { MartService, type WalletBalance } from "@/services/MartService";
import { useMartParamStore } from "../stores/martParam-store";
import { onMounted, ref, watch } from "vue";

defineOptions({
  name: "BalanceComponent",
});

const martParamStore = useMartParamStore();
const activeMartParam = ref<MartParamModel | null>(null);
const balances = ref<WalletBalance[]>([]);

onMounted(async () => {
  await load();
});

watch(martParamStore, async () => {
  await load();
});

async function load() {
  if (martParamStore.getActiveMartParam == null) {
    setTimeout(load, 100);
    return;
  }

  await loadSpotWallet();
}

async function loadSpotWallet() {
  activeMartParam.value = martParamStore.getActiveMartParam;
  if (activeMartParam.value === null) {
    return;
  }

  MartService.LoadSpotWallet(activeMartParam.value.domain).then((resp) => {
    console.log("loadSpotWallet>>>", resp);
    if (resp?.code === 0) {
      balances.value = resp.data;
    }
  });
}
</script>

<template>
  <ScrollArea class="w-full h-[380px] rounded-md">
    <div class="space-y-4">
      <div v-for="(b, idx) in balances" :key="idx" class="flex items-center">
        <div class="space-y-1">
          <p class="text-sm font-medium leading-none">
            {{ b.currency }}
          </p>
          <p class="text-sm text-muted-foreground">
            {{ $t("available") }}:
            {{ b.available }}
          </p>
          <p class="text-sm text-muted-foreground">
            {{ $t("frozen") }}:
            {{ b.frozen }}
          </p>
          <p v-if="b.updateAt" class="text-sm text-muted-foreground">
            {{ $t("syncAt") }}:
            {{ b.updateAt }}
          </p>
        </div>
        <div class="ml-auto font-medium">
         
        </div>
      </div>
    </div>
  </ScrollArea>
</template>
