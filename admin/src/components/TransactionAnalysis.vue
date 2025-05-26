<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { MartParamModel } from "@/services/MartParamService";
import { OrderModel, OrderService } from "@/services/OrderService";
import { useMartParamStore } from "@/stores/martParam-store";
import { onMounted, ref, watch } from "vue";
import moment from "moment";

defineOptions({
  name: "AnalysisTransaction",
});

const martParamStore = useMartParamStore();
const activeMartParam = ref<MartParamModel | null>(null);
const orders = ref<OrderModel[]>([]);
let page = 1;

onMounted(async () => {
  await load();
});

watch(martParamStore, async () => {
  await load();
});

async function nextPage(n: number) {
  page = page + n
  if (page < 1) {
    page  =1
  }

  await load()
}
async function load() {
  if (martParamStore.getActiveMartParam == null) {
    setTimeout(load, 100);
    return;
  }

  await queryMyOrder();
}

async function queryMyOrder() {
  if (martParamStore.getActiveMartParam == null) {
    return;
  }

  activeMartParam.value = martParamStore.getActiveMartParam;

  const resp = await OrderService.List(
    "",
    "",
    page
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
  console.log("orders:", orders.value);
  if (page > 1 && orders.value.length == 0) {
    nextPage(-1);
  }
}

</script>

<template>
  <Card>
    <CardHeader  class="flex row gap-2 items-end justify-end py-4">
      <!-- <div class="w-full flex items-center justify-between space-y-2">
        <DateRangePicker @updated="onDateUpdated" />
        <Button
          class="bg-primary shadow hover:bg-primary/90 h-9 px-4 py-2"
          @click="onQuery"
          >Query
        </Button>
      </div> -->
      <div class="space-x-2">
        <Button size="sm"  variant="outline" @click="nextPage(-1)"> Previous </Button>
        <Button size="sm"  variant="outline" @click="nextPage(1)"> Next </Button>
      </div>
    </CardHeader>
    <CardContent>
      <div class="space-y-4">
        <Table>
        <!-- <TableCaption>No data.</TableCaption> -->
        <TableHeader>
          <TableRow>
            <TableHead> Symbol </TableHead>
            <TableHead>Mart</TableHead>
            <TableHead>Side</TableHead>
            <TableHead> Quantity </TableHead>
            <TableHead> Price </TableHead>
            <TableHead> Time </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="order in orders">
            <TableCell class="font-medium"> {{ order.symbol }} </TableCell>
            <TableCell>{{ order.martDomain }}</TableCell>
            <TableCell>{{ order.side }}</TableCell>
            <TableCell> {{ order.quantity }}</TableCell>
            <TableCell> {{ order.price }}</TableCell>
            <TableCell> {{ order.createTimeStr }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      </div>
    </CardContent>
  </Card>
</template>
