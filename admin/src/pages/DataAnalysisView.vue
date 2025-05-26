<script setup lang="ts">
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import DateRangePicker from "@/components/DateRangePicker.vue";
import RecentTransaction from "@/components/RecentTransaction.vue";
import { onMounted, onUnmounted } from "vue";
import AnalysisTransaction from "@/components/TransactionAnalysis.vue";

defineOptions({
  name: "DataAnalysisView",
});

onMounted(async () => {});

onUnmounted(() => {});
</script>

<template>
  <div class="flex-1 p-2">
    <Tabs default-value="OrderAnalysis">
      <TabsList>
        <TabsTrigger value="OrderAnalysis"> Order Analysis </TabsTrigger>
        <TabsTrigger value="BalanceAnalysis"> Balance Analysis </TabsTrigger>
      </TabsList>
      <TabsContent value="OrderAnalysis" class="space-y-4" >
        <AnalysisTransaction />
      </TabsContent>
      <TabsContent value="TimeQuery" class="space-y-4">
        <Card>
          <CardHeader>
            <div class="w-full flex items-center justify-between space-y-2">
              <div class="flex row">
                <DateRangePicker />
                <DateRangePicker />
              </div>

              <Button
                class="bg-primary shadow hover:bg-primary/90 h-9 px-4 py-2"
                >Query</Button
              >
            </div>
          </CardHeader>
          <CardContent>
            <RecentTransaction />
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
