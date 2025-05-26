<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { BotModel, BotService } from "@/services/BotService";
import RobotCard from "./RobotCard.vue";
import { useMartParamStore } from "@/stores/martParam-store";
import { MartParamModel } from "@/services/MartParamService";
import { useI18n } from "vue-i18n";
import moment from "moment";
import { HttpResponse } from "@/services";
import { toast } from "@/components/ui/toast";

const strategicCards = ref<BotModel[]>([
  {
    id: 0,
    botType: 100,
    botName: "GradualLift",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 101,
    botName: "GradualDown",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 102,
    botName: "RapidLift",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 103,
    botName: "RapidDown",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 104,
    botName: "RatioLift",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,  
    botType: 105,       
    botName: "RatioDown",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },

  {
    id: 0,
    botType: 109,
    botName: "RangeOscillation",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  // {
  //   id: 0,
  //   botType: 110,
  //   botName: "EnhancementRangeOscillation",
  //   isRunning: 0,
  //   martDomain: "",
  //   symbol: "",
  //   startTime: 0,
  // },
  {
    id: 0,
    botType: 112,
    botName: "RatioOscillation",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 111,
    botName: "HighThrowBargain",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 106,
    botName: "ETFFollowMode",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 107,
    botName: "StrictFollowMode",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 108,
    botName: "CombinedFollowMode",
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
  console.log("Strategic>>>bots: ", resp);
  if (resp) {
    if (resp.code == 0) {
      resp.data.forEach((one) => {
        const idx = strategicCards.value.findIndex(
          (card) => card.botType == one.botType
        );

        if (idx >= 0) {
          strategicCards.value[idx].id = one.id;
          strategicCards.value[idx].isRunning = one.isRunning;
          if (strategicCards.value[idx].isRunning === 1) {
            strategicCards.value[idx].startTimeStr = moment(
              one.startTime
            ).format("YYYY-MM-DD HH:mm:ss");
          } else {
            strategicCards.value[idx].startTimeStr = '';
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
  const idx = strategicCards.value.findIndex(
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
      v-for="(bot, idx) in strategicCards"
      :key="idx" :id="bot.id" 
     :botType="bot.botType"
     :botName="bot.botName" 
     :isRunning="bot.isRunning" 
     :runTime="bot.startTimeStr"
     @RunOrStop="RunOrStopBot"
     @commited="onBotCommited" />
  </div>
</template>

<style scoped>
.read-the-docs {
  color: #888;
}
</style>
