<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { onMounted, onUnmounted, ref } from "vue";
import MartParamComponent from "@/components/MartParamComponent.vue";
import { MartParamModel, MartParamService } from "@/services/MartParamService";
import { useI18n } from "vue-i18n";
import Toaster from "@/components/ui/toast/Toaster.vue";
import { toast } from "@/components/ui/toast";

defineOptions({
  name: "UserView",
});

const { t } = useI18n();
// const tab = ref("mart");
const martParamDialog = ref(false);
const models = ref<MartParamModel[]>([]);
// const editModelId = ref(0);

onMounted(() => {
  loadMyMartParam();
});

onUnmounted(() => {});

async function loadMyMartParam() {
  const resp = await MartParamService.Load();
  if (resp) {
    if (resp.code == 0) {
      models.value = resp.data;
    } else {
      toast({
        variant: "destructive",
        title: t("error." + resp.code),
      });
      // q.notify({ message: , type: 'warning' });
    }
  }
}

function onMartParamCommited(_data: any) {
  martParamDialog.value = false;
  loadMyMartParam();
}
</script>

<template>
  <div class="flex-1 space-y-4 p-8 pt-6">
    <Tabs default-value="apikeys" class="space-y-4">
      <TabsList>
        <TabsTrigger value="apikeys">
          {{ $t("bindApiKey") }}
        </TabsTrigger>
        <TabsTrigger value="analytics" disabled> Account </TabsTrigger>
      </TabsList>
      <TabsContent value="apikeys" class="space-y-4">
        <div class="flex items-center justify-between space-y-2">
          <h2 class="text-3xl font-bold tracking-tight">
            {{ $t("userView.APIKEYs") }}
          </h2>
          <div class="flex items-center space-x-2">
            <MartParamComponent :id="0" @commited="onMartParamCommited" />
          </div>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead> {{ $t("mart") }} </TableHead>
              <TableHead>{{ $t("symbol") }}</TableHead>
              <TableHead>AccessKey</TableHead>
              <TableHead> MEMO </TableHead>
              <TableHead> {{ $t("edit") }} </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="(mart, idx) in models" :key="idx">
              <TableCell> {{ mart.domain }} </TableCell>
              <TableCell> {{ mart.symbol }} </TableCell>
              <TableCell> {{ mart.accessKey }} </TableCell>
              <TableCell> {{ mart.memo }} </TableCell>
              <TableCell>
                <MartParamComponent :id="mart.id" @commited="onMartParamCommited" />
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </TabsContent>
    </Tabs>
  </div>

  <Toaster />
</template>

<style scoped>
.mytab {
  font-weight: 700 !important;
  font-size: larger !important;
}

.my-card-div {
  padding: 6px;
}

.my-card {
  height: 200px;
  border-radius: 12px;
  font-size: larger;

  .text-caption {
    font-size: 14px;
  }
  .q-item__label {
    font-size: large;
  }
}
</style>
