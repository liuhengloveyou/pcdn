<script setup lang="ts">
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Toaster from "@/components/ui/toast/Toaster.vue";
import RecentTransaction from "@/components/RecentTransaction.vue";
import StrategicBots from "@/components/StrategicBots.vue";
import AIBots from "@/components/AIBotsComponent.vue";
import RiskControlBots from "@/components/RiskControlBots.vue";
import BalanceComponent from "@/components/BalanceComponent.vue";
import VolumeBots from "@/components/VolumeBots.vue";
import DepthBots from "@/components/DepthBots.vue";

defineOptions({
  name: "HomeView",
});

// const martParamStore = useMartParamStore();
// const activeMartParam = ref<MartParamModel | null>(null);
// const orders = ref<OrderModel[]>([]);
// const balances = ref<WalletBalance[]>([]);

// onMounted(async () => {
//   await load();
// });

// onUnmounted(() => {});

// watch(martParamStore, async () => {
//   console.log("martParam changed.");
//   await load();
// });

</script>

<template>
  <div class="flex-1 p-6 pt-6">
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <Kline />
      <div class="grid gap-4 grid-cols-2 col-span-3">
        <Card class="col-span-1 h-[500px]">
          <CardHeader>
            <CardTitle>{{ $t("orderRecordTitle") }}</CardTitle>
            <!-- <CardDescription> You made 265 orders. </CardDescription> -->
          </CardHeader>
          <CardContent>
            <RecentTransaction />
          </CardContent>
        </Card>

        <Card class="col-span-1 h-[500px]">
          <CardHeader>
            <CardTitle>{{ $t("balance") }}</CardTitle>
            <!-- <CardDescription> You made 265 sales this month. </CardDescription> -->
          </CardHeader>
          <CardContent >
            <BalanceComponent />
          </CardContent>
        </Card>
      </div>
    </div>

    <h2
      class="flex row items-center text-3xl font-bold tracking-tight mt-8 mb-4"
    >
      <div class="bg-primary m-4 h-[20px] w-2"></div>
      {{ $t("volumeBot") }}
    </h2>
    <VolumeBots />

    <h2
      class="flex row items-center text-3xl font-bold tracking-tight mt-8 mb-4"
    >
      <div class="bg-primary m-4 h-[20px] w-2"></div>
      {{ $t("depthBot") }}
    </h2>
    <DepthBots />

    <h2
      class="flex row items-center text-3xl font-bold tracking-tight mt-8 mb-4"
    >
      <div class="bg-primary m-4 h-[20px] w-2"></div>
      {{ $t("StrategicBot") }}
    </h2>
    <StrategicBots />

    <h2
      class="flex row items-center text-3xl font-bold tracking-tight mt-8 mb-4"
    >
      <div class="bg-primary m-4 h-[20px] w-2"></div>
      {{ $t("AIBot") }}
    </h2>
    <AIBots />

    <h2
      class="flex row items-center text-3xl font-bold tracking-tight mt-8 mb-4"
    >
      <div class="bg-primary m-4 h-[20px] w-2"></div>
      {{ $t("RiskControlBot") }}
    </h2>
    <RiskControlBots />
  </div>

  <Toaster />
</template>
