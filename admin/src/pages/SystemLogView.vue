<script setup lang="ts">
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { onMounted, ref } from "vue";
import {
  BusinessLogModel,
  BusinessLogService,
} from "@/services/BusinessLogService";
import { useMartParamStore } from "@/stores/martParam-store";
import moment from "moment";

defineOptions({
  name: "BusinessLogPage",
});

const data = ref<BusinessLogModel[]>([]);
const martParamStore = useMartParamStore();

onMounted(async () => {
  await load();
});

async function load() {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    setTimeout(load, 100);
    return;
  }

  const resp = await BusinessLogService.List(
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    if (resp.code === 0) {
      resp.data.forEach((one) => {
        one.createTimeStr = moment(one.createTime).format(
          "YYYY-MM-DD HH:mm:ss"
        );

        one.businessType =
          one.businessType === "BusinessStopBot" ? "Stop Bot" : "Start Bot";
      });

      data.value = resp.data;
    }
  }
  console.log("system logs:", data);
}
</script>

<template>
  <div class="w-full p-6">
    <div class="flex row gap-2 items-center justify-end py-4">
      <!-- <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="outline" class="ml-auto">
            Columns <ChevronDown class="ml-2 h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuCheckboxItem
            v-for="column in table
              .getAllColumns()
              .filter((column) => column.getCanHide())"
            :key="column.id"
            class="capitalize"
            :checked="column.getIsVisible()"
            @update:checked="
              (value) => {
                column.toggleVisibility(!!value);
              }
            "
          >
            {{ column.id }}
          </DropdownMenuCheckboxItem>
        </DropdownMenuContent>
      </DropdownMenu> -->

      <div class="space-x-2">
        <!-- :disabled="!table.getCanPreviousPage()" -->
        <Button size="sm"> Previous </Button>
        <!-- :disabled="!table.getCanNextPage()" -->
        <Button size="sm"> Next </Button>
      </div>
    </div>
    <div class="rounded-md border">
      <Table>
        <!-- <TableCaption>No data.</TableCaption> -->
        <TableHeader>
          <TableRow>
            <TableHead> Business Type </TableHead>
            <TableHead>Create Time</TableHead>
            <TableHead>Account Name</TableHead>
            <TableHead>Robot Type </TableHead>
            <TableHead> Symbol </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="row in data">
            <TableCell class="font-medium"> {{ row.businessType }} </TableCell>
            <TableCell>{{ row.createTimeStr }}</TableCell>
            <TableCell>{{ row.uid }}</TableCell>
            <TableCell> {{ row.botType }}</TableCell>
            <TableCell> {{ row.symbol }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <div class="flex items-center justify-end space-x-2 py-4">
      <!-- <div class="flex-1 text-sm text-muted-foreground">
        {{ table.getFilteredSelectedRowModel().rows.length }} of
        {{ table.getFilteredRowModel().rows.length }} row(s) selected.
      </div> -->

     
    </div>
  </div>
</template>
