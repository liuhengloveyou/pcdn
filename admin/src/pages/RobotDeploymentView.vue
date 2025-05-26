<script setup lang="ts">
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { onMounted, onUnmounted } from "vue";
import AIBots from "@/components/AIBotsComponent.vue";
import StrategicBots from "@/components/StrategicBots.vue";
import VolumeBots from "@/components/VolumeBots.vue";
import DepthBots from "@/components/DepthBots.vue";

defineOptions({
  name: "RobotDeploymentView",
});

// const martParamStore = useMartParamStore();
// const activeMartParam = ref<MartParamModel | null>(null);

// const ziranShualingModel = ref<ZiranShualingModel>({
//   id: 0,
//   minNumPeerOrder: 0,
//   maxNumPeerOrder: 0,
//   minTimeInterval: 0,
//   maxTimeInterval: 0,
//   orderModel: "",
//   priceModel: "",
//   startModel: "",
//   quantModel: "",
//   description: "",
//   isRunning: false,
//   startTime: "",
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
//     activeMartParam.value.symbol
//   );
//   console.log(">>>bots: ", resp, activeMartParam.value);
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
  <div class="flex-1 space-y-4 p-8 pt-6">
    <Tabs default-value="volumeBot" class="space-y-4">
      <TabsList>
        <TabsTrigger value="volumeBot">
          {{ $t("volumeBot") }}
        </TabsTrigger>
        <TabsTrigger value="depthBot">
          {{ $t("depthBot") }}
        </TabsTrigger>
        <TabsTrigger value="StrategicBot">
          {{ $t("StrategicBot") }}
        </TabsTrigger>
        <TabsTrigger value="AIBot"> {{ $t("AIBot") }} </TabsTrigger>
      </TabsList>
      <TabsContent value="volumeBot" class="space-y-4">
        <VolumeBots />
      </TabsContent>

      <TabsContent value="depthBot" class="space-y-4">
        <DepthBots />
      </TabsContent>

      <TabsContent value="StrategicBot" class="space-y-4">
        <StrategicBots />
      </TabsContent>

      <TabsContent value="AIBot" class="space-y-4">
        <AIBots />
      </TabsContent>
    </Tabs>
  </div>
</template>
