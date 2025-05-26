<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import RobotCard from "./RobotCard.vue";
import { BotModel, BotService } from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { MartParamModel } from "@/services/MartParamService";
import moment from "moment";
import { HttpResponse } from "@/services";
import { useI18n } from "vue-i18n";
import Toaster from "@/components/ui/toast/Toaster.vue";
import { toast } from "@/components/ui/toast";

const commonBotCards = ref<BotModel[]>([
  {
    id: 0,
    botType: 7,
    botName: "OrderbookFlash",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 8,
    botName: "DepthRemote",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 9,
    botName: "DepthNear",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 10,
    botName: "DepthSpread",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
]);

const martParamStore = useMartParamStore();
const activeMartParam = ref<MartParamModel | null>(null);
  const { t } = useI18n();
  
onMounted(async () => {
  await load();
});

watch(martParamStore, async () => {
  console.log("ChangguiBotComponent.martParam changed.");
  await load();
});

async function load() {
  if (martParamStore.getActiveMartParam == null) {
    setTimeout(load, 100);
    return;
  }

  activeMartParam.value = martParamStore.getActiveMartParam;
  const resp = await BotService.Load(
    activeMartParam.value.domain,
    activeMartParam.value.symbol
  );
  console.log("CommonBots>>>bots: ", resp);
  if (resp) {
    if (resp.code == 0) {
      resp.data.forEach((one) => {
        const idx = commonBotCards.value.findIndex(
          (card) => card.botType == one.botType
        );

        if (idx >= 0) {
          commonBotCards.value[idx].id = one.id;
          commonBotCards.value[idx].isRunning = one.isRunning;
          if (commonBotCards.value[idx].isRunning === 1) {
            commonBotCards.value[idx].startTimeStr = moment(
              one.startTime
            ).format("YYYY-MM-DD HH:mm:ss");
          } else {
            commonBotCards.value[idx].startTimeStr = '';
          }
        }
      });
    }
  }
}

async function onBotCommited(data: any) {
  console.log(">>>onBotCommited", data)
  await load();
}

async function RunOrStopBot(data: any) {
  const idx = commonBotCards.value.findIndex(
    (card) => card.botType == data.botType
  );
 

  if (martParamStore.getActiveMartParam == null) {
    return;
  }

  let resp: HttpResponse<string> | undefined;
  if (data.stat == true) {
    resp = await BotService.Run(
      data.id,
      data.botType,
      martParamStore.getActiveMartParam.domain,
      martParamStore.getActiveMartParam.symbol
    );
  } else if (data.stat == false) {
    resp = await BotService.Stop(
      data.id,
      data.botType,
      martParamStore.getActiveMartParam.domain,
      martParamStore.getActiveMartParam.symbol
    );
  }
  console.log("RunOrStopBot.resp>>>>>>>>", idx, data, resp);
  if (resp) {
    if (resp.code !== 0) {
      toast({
      variant: "destructive",
      title: t("error." + resp.code),
    });
    }
  }

  await load();
}
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
    <RobotCard
      v-for="(bot, idx) in commonBotCards"
      :key="idx"
      :id="bot.id"
      :botType="bot.botType"
      :botName="bot.botName"
      :isRunning="bot.isRunning"
      :runTime="bot.startTimeStr"
      @RunOrStop="RunOrStopBot"
      @commited="onBotCommited"
    />
  </div>

  <Toaster />
</template>
