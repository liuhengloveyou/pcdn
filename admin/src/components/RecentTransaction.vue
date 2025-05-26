<script setup lang="ts">
import { ScrollArea } from "@/components/ui/scroll-area";
import { MartParamModel } from "@/services/MartParamService";
import { OrderModel, OrderService } from "@/services/OrderService";
import { useMartParamStore } from "@/stores/martParam-store";
import { onMounted, ref, watch } from "vue";
import moment from "moment";

defineOptions({
  name: "RecentTransaction",
});

const martParamStore = useMartParamStore();
const activeMartParam = ref<MartParamModel | null>(null);
const orders = ref<OrderModel[]>([]);

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

  await loadMyOrder();
}

async function loadMyOrder() {
  if (martParamStore.getActiveMartParam == null) {
    setTimeout(loadMyOrder, 100);
    return;
  }

  activeMartParam.value = martParamStore.getActiveMartParam;

  const resp = await OrderService.List(
    activeMartParam.value.domain,
    activeMartParam.value.symbol,
    1
  );
  if (resp) {
    if (resp.code === 0) {
      resp.data.forEach((one) => {
        one.createTimeStr = moment(one.createTime).format(
          "YYYY-MM-DD HH:mm:ss"
        );
      });

      orders.value = resp.data;
    }
  }
  // console.log("orders:", orders.value);
}
</script>

<template>
  <ScrollArea class="w-full h-[380px] rounded-md p-0">
    <div class="space-y-4">
      <div v-for="(order, idx) in orders" :key="idx" class="flex items-center">
        <div class="space-y-1">
          <p class="text-sm font-medium leading-none">
            {{ order.symbol }} @
            {{ order.martDomain }}
          </p>
          <p class="text-sm text-muted-foreground">
            {{ order.price }} / {{ order.quantity }}
          </p>
          <p class="text-sm text-muted-foreground">
            {{ order.client_order_id }}
          </p>
          <p class="text-sm text-muted-foreground">
            {{ order.createTimeStr }}
          </p>
   
        </div>
        <div class="ml-auto font-medium">
          <span :style="{ color: order.side === 'BUY' ? 'green' : 'red' }">{{
            order.side
          }}</span>
        </div>
      </div>
    </div>
  </ScrollArea>
</template>
