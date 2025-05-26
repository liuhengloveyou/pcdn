<script setup lang="ts">
import { ref } from "vue";
import RobotCard from "./RobotCard.vue";
import { BotModel } from "@/services/BotService";

const aiBots = ref<BotModel[]>([
  {
    id: 0,
    botType: 200,
    botName: "AMMMode",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
  {
    id: 0,
    botType: 201,
    botName: "Bulldozer",
    isRunning: 0,
    martDomain: "",
    symbol: "",
    startTime: 0,
  },
]);
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
    <RobotCard
      v-for="(bot, idx) in aiBots"
      :key="idx"
      :id="bot.id"
      :isRunning="bot.isRunning"
      :botType="bot.botType"
      :botName="bot.botName"
    />
  </div>
</template>

<style scoped>
.read-the-docs {
  color: #888;
}
</style>
