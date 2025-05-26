<script setup lang="ts">
import RiskControlBots from "@/components/RiskControlBots.vue";
import { onMounted, onUnmounted } from "vue";

defineOptions({
  name: "RiskManagerView",
});

// const martParamStore = useMartParamStore();
// const activeMartParam = ref<MartParamModel | null>(null);
// const ziranShualiangDialog = ref(false);

// const ziranShualingModel = ref<ZiranShualingModel>({
//   id: 0,
//   minNumPeerOrder: 0,
//   maxNumPeerOrder: 0,
//   minTimeInterval: 0,
//   maxTimeInterval: 0,
//   orderModel: '',
//   priceModel: '',
//   startModel: '',
//   quantModel: '',
//   description: '',
//   isRunning: false,
//   startTime: '',
// });

onMounted(async () => {
  // await load();
});

onUnmounted(() => {});

// async function load() {
//   if (martParamStore.getActiveMartParam == null) {
//     setTimeout(load, 100);
//     return;
//   }

//   activeMartParam.value = martParamStore.getActiveMartParam;
//   const resp = await BotService.Load(
//     activeMartParam.value.domain,
//     activeMartParam.value.symbol,
//   );
//   console.log('>>>bots: ', resp, activeMartParam.value);
//   if (resp) {
//     if (resp.code == 0) {
//       resp.data.forEach((one) => {
//         if (one.botType == 1) {
//           one.ziranShualing.isRunning = one.isRunning == 1 ? true : false;
//           ziranShualingModel.value = one.ziranShualing;
//         }
//       });
//     }
//   }
// }

// function onCommited() {
//   load();
//   ziranShualiangDialog.value = false;
// }
</script>

<template>
  <h2 class="flex row items-center text-3xl font-bold tracking-tight mt-8 mb-4">
    <div class="bg-primary m-4 h-[20px] w-2"></div>
    {{ $t("RiskControlBot") }}
  </h2>
  <RiskControlBots />
</template>
