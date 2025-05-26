<script setup lang="ts">
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { Icon } from "@iconify/vue";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import MinimumBrushRobotForm from "./robot-edit-forms/MinimumBrushRobotForm.vue";
import RandomBrushRobotForm from "./robot-edit-forms/RandomBrushRobotForm.vue";
import NaturalBrushRobotForm from "./robot-edit-forms/NaturalBrushRobotForm.vue";
import StableBrushRobotForm from "./robot-edit-forms/StableBrushRobotForm.vue";
import RangeBrushRobotForm from "./robot-edit-forms/RangeBrushRobotForm.vue";
import OrderbookFlashRobotForm from "./robot-edit-forms/OrderbookFlashRobotForm.vue";
import DepthRemoteRobotForm from "./robot-edit-forms/DepthRemoteRobotForm.vue";
import DepthNearRobotForm from "./robot-edit-forms/DepthNearRobotForm.vue";
import DepthSpreadRobotForm from "./robot-edit-forms/DepthSpreadRobotForm.vue";
import BalanceMonitorRobotForm from "./robot-edit-forms/BalanceMonitorRobotForm.vue";
import PriceMonitorRobotForm from "./robot-edit-forms/PriceMonitorRobotForm.vue";
import BalanceScheduledReportRobotForm from "./robot-edit-forms/BalanceScheduledReportRobotForm.vue";
import EmergencyBrakeRobotForm from "./robot-edit-forms/EmergencyBrakeRobotForm.vue";
import TGScheduledReportRobotForm from "./robot-edit-forms/TGScheduledReportRobotForm.vue";
import TGAlertRobotForm from "./robot-edit-forms/TGAlertRobotForm.vue";
import AMMModeRobotForm from "./robot-edit-forms/AMMModeRobotForm.vue";
import BulldozerRobotForm from "./robot-edit-forms/BulldozerRobotForm.vue";
import GradualLiftRobotForm from "./robot-edit-forms/GradualLiftRobotForm.vue";
import GradualDownRobotForm from "./robot-edit-forms/GradualDownRobotForm.vue";
import OscillationRobotForm from "./robot-edit-forms/RatioOscillationRobotForm.vue";
import HighThrowBargainRobotForm from "./robot-edit-forms/HighThrowBargainRobotForm.vue";
import EnhancementRangeOscillationRobotForm from "./robot-edit-forms/EnhancementRangeOscillationRobotForm.vue";
import RangeOscillationRobotForm from "./robot-edit-forms/RangeOscillationRobotForm.vue";
import CombinedFollowModeRobotForm from "./robot-edit-forms/CombinedFollowModeRobotForm.vue";
import StrictFollowModeRobotForm from "./robot-edit-forms/StrictFollowModeRobotForm.vue";
import ETFFollowModeRobotForm from "./robot-edit-forms/ETFFollowModeRobotForm.vue";
import RatioDownRobotForm from "./robot-edit-forms/RatioDownRobotForm.vue";
import RatioLiftRobotForm from "./robot-edit-forms/RatioLiftRobotForm.vue";
import RapidDownRobotForm from "./robot-edit-forms/RapidDownRobotForm.vue";
import RapidLiftRobotForm from "./robot-edit-forms/RapidLiftRobotForm.vue";
import { onMounted, ref } from "vue";

defineOptions({
  name: "RobotEditComponent",
});
const emits = defineEmits(["commited"]);
defineProps<{
  botType: number;
  botName: string;
}>();

const open = ref(false);

onMounted(() => {
  console.log("RobotEditComponent.onMounted");
});

function onCommited(data: any) {
  open.value = false;
  emits("commited", data);
}
</script>

<template>
  <Sheet v-model:open="open">
    <SheetTrigger as-child>
      <Button>
        <Icon
          icon="ic:baseline-tune"
          class="h-[1.2rem] w-[1.2rem] transition-all"
        />
        {{ $t("setting") }}
      </Button>
    </SheetTrigger>

    <SheetContent side="right" class="lg:w-[800px] lg:max-w-lg">
      <SheetHeader class="pb-4">
        <SheetTitle>{{ $t(botName) }}</SheetTitle>
      </SheetHeader>
      <ScrollArea class="h-[90vh] w-full">
      <div class="grid gap-4 py-4">
        <MinimumBrushRobotForm v-if="botType === 1" @commited="onCommited" />
        <RandomBrushRobotForm v-if="botType === 2" @commited="onCommited" />
        <NaturalBrushRobotForm v-if="botType === 3" @commited="onCommited" />
        <StableBrushRobotForm v-if="botType === 4" @commited="onCommited" />
        <RangeBrushRobotForm v-if="botType === 5" @commited="onCommited" />
        <OrderbookFlashRobotForm v-if="botType === 7" @commited="onCommited" />
        <DepthRemoteRobotForm v-if="botType === 8" @commited="onCommited" />
        <DepthNearRobotForm v-if="botType === 9" @commited="onCommited" />
        <DepthSpreadRobotForm v-if="botType === 10" @commited="onCommited" />

        <!-- StrategicBots -->
        <GradualLiftRobotForm v-if="botType === 100" @commited="onCommited" />
        <GradualDownRobotForm v-if="botType === 101" @commited="onCommited" />
        <RapidLiftRobotForm v-if="botType === 102" @commited="onCommited" />
        <RapidDownRobotForm v-if="botType === 103" @commited="onCommited" />
        <RatioLiftRobotForm v-if="botType === 104" @commited="onCommited" />
        <RatioDownRobotForm v-if="botType === 105" @commited="onCommited" />
        <ETFFollowModeRobotForm v-if="botType === 106" @commited="onCommited" />
        <StrictFollowModeRobotForm
          v-if="botType === 107"
          @commited="onCommited"
        />
        <CombinedFollowModeRobotForm
          v-if="botType === 108"
          @commited="onCommited"
        />
        <RangeOscillationRobotForm
          v-if="botType === 109"
          @commited="onCommited"
        />
        <EnhancementRangeOscillationRobotForm
          v-if="botType === 110"
          @commited="onCommited"
        />
        <HighThrowBargainRobotForm
          v-if="botType === 111"
          @commited="onCommited"
        />
        <OscillationRobotForm v-if="botType === 112" @commited="onCommited" />

        <!-- AI Bot -->
        <AMMModeRobotForm v-if="botType === 200" @commited="onCommited" />
        <BulldozerRobotForm v-if="botType === 201" @commited="onCommited" />

        <!-- RiskControlBots -->
        <BalanceMonitorRobotForm
          v-if="botType === 300"
          @commited="onCommited"
        />
        <PriceMonitorRobotForm v-if="botType === 301" @commited="onCommited" />
        <BalanceScheduledReportRobotForm
          v-if="botType === 302"
          @commited="onCommited"
        />
        <EmergencyBrakeRobotForm
          v-if="botType === 303"
          @commited="onCommited"
        />
        <TGScheduledReportRobotForm
          v-if="botType === 304"
          @commited="onCommited"
        />
        <TGAlertRobotForm v-if="botType === 305" @commited="onCommited" />
      </div>
    </ScrollArea>
    </SheetContent>

  </Sheet>
</template>

<style scoped>
.read-the-docs {
  color: #888;
}
</style>
