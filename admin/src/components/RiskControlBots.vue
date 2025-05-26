<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import RobotCard from "./RobotCard.vue";
import { BotModel } from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { MartParamModel } from "@/services/MartParamService";
import { useI18n } from "vue-i18n";
import { RiskManService } from "@/services/RiskManService";
import moment from "moment";
import { toast } from "@/components/ui/toast";
import { HttpResponse } from "@/services";

const riskControlBots = ref<BotModel[]>([
  {
    id: 0,
    botType: 300,
    botName: "BalanceMonitor",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 301,
    botName: "PriceMonitor",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 302,
    botName: "BalanceScheduledReport",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 303,
    botName: "EmergencyBrake",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 304,
    botName: "TGScheduledReport",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 305,
    botName: "TGAlert",
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
  console.log("RiskControlBots.martParam changed.");
  await load();
});

async function load() {
  if (martParamStore.getActiveMartParam == null) {
    setTimeout(load, 100);
    return;
  }

  activeMartParam.value = martParamStore.getActiveMartParam;
  const resp = await RiskManService.Load(
    activeMartParam.value.domain,
    activeMartParam.value.symbol
  );
  console.log("RiskManBots>>>bots: ", resp);

  if (resp && resp.code == 0) {
    resp.data.forEach((one) => {
      const idx = riskControlBots.value.findIndex(
        (card) => card.botType == one.botType
      );
      console.log("RiskManBots>>>bot: ", idx, one);
      if (idx >= 0) {
        riskControlBots.value[idx].id = one.id;
        riskControlBots.value[idx].isRunning = one.isRunning;
        if (riskControlBots.value[idx].isRunning === 1) {
          riskControlBots.value[idx].startTimeStr = moment(
            one.startTime
          ).format("YYYY-MM-DD HH:mm:ss");
        } else {
          riskControlBots.value[idx].startTimeStr = "";
        }
      }
    });
  }
}

async function RunOrStopBot(data: any) {
  const idx = riskControlBots.value.findIndex(
    (card) => card.botType == data.botType
  );

  if (martParamStore.getActiveMartParam == null) {
    return;
  }

  let resp: HttpResponse<string> | undefined;
  if (data.stat == true) {
    resp = await RiskManService.Run(
      data.id,
      data.botType,
      martParamStore.getActiveMartParam.domain,
      martParamStore.getActiveMartParam.symbol
    );
  } else if (data.stat == false) {
    resp = await RiskManService.Stop(
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

async function onCommited() {
 await load();
}

</script>

<template>
  <div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
    <RobotCard
      v-for="(bot, idx) in riskControlBots"
      :key="idx"
      :id="bot.id"
      :botType="bot.botType"
      :botName="bot.botName"
      :isRunning="bot.isRunning"
      :runTime="bot.startTimeStr"
      @RunOrStop="RunOrStopBot"
      @commited="onCommited"
    />
  </div>
</template>
