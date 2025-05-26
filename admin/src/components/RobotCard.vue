<script setup lang="ts">
import { toast } from "@/components/ui/toast";
import { Switch } from "@/components/ui/switch";
import {
  Card,
  CardContent,
  CardHeader,
  CardFooter,
  CardTitle,
} from "@/components/ui/card";
import { cn } from '@/lib/utils'
import RobotEditSheet from "./RobotEditSheet.vue";
import { computed } from "vue";

const props = defineProps({
  id: {
    type: Number,
    required: true,
  },

  botType: {
    type: Number,
    required: true,
  },
  botName: {
    type: String,
    required: true,
  },
  isRunning: {
    type: Number,
    required: true,
  },
  runTime:{
    type: String
  }
});
const emits = defineEmits(["RunOrStop", "commited"]);

const runSwitch = computed(() => (props.isRunning === 1 ? true : false));

function onCommited() {
  emits("commited", { id: props.id, botType: props.botType });
}

function RunOrStopBot(value: boolean) {
  console.log("RunOrStopBot: ", props.id)
  if (!props.id || props.id <= 0) {
    toast({
      variant: "destructive",
      title: "Please set first.",
    });
    return;
  }
  emits("RunOrStop", { id: props.id, botType: props.botType, stat: value });
}
</script>

<template>
  <Card class="h-[200px]">
    <CardHeader
      class="flex flex-row items-center justify-between space-y-0 pb-2 "
    >
      <CardTitle class="text-xl font-medium">
        {{ $t(props.botName) }}
      </CardTitle>
      <RobotEditSheet :bot-type="botType" :bot-name="props.botName" @commited="onCommited" />
    </CardHeader>
    <CardContent class="h-[70px]">
      <p class="text-xs text-muted-foreground">
        {{ $t("Descriptionisempty") }}
      </p>
    </CardContent>
    <CardFooter class="flex justify-between px-6 pb-0">
      <div>
        <div :class="cn('text-sm font-bold', (props.isRunning === 1) && 'text-primary')"> {{ $t((props.isRunning === 1) ? "BotRun" : "BotOff") }}</div>
        <p v-if="props.isRunning === 1" class="text-xs text-muted-foreground">{{ runTime }}</p>
      </div>
      <Switch
        id="airplane-mode"
        v-model:checked="runSwitch"
        @update:checked="RunOrStopBot"
      />
    </CardFooter>
  </Card>
</template>
