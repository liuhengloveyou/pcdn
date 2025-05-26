<script setup lang="ts">
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Command, CommandItem, CommandList } from "@/components/ui/command";
import Toaster from "@/components/ui/toast/Toaster.vue";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { cn } from "@/lib/utils";
import { MartParamModel, MartParamService } from "@/services/MartParamService";
import { useMartParamStore } from "@/stores/martParam-store";
import { CheckIcon, CaretSortIcon } from "@radix-icons/vue";
import { onMounted, ref } from "vue";
import { toast } from "./ui/toast";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const loading = ref(false);
// 当前交易所和交易对
const currMartParam = ref<MartParamModel>({
  id: 0,
  uid: 0,
  domain: "domain",
  symbol: "symbol",
  memo: "memo",
  accessKey: "",
  secretKey: "",
  active: 0,
});
const myMartParams = ref<MartParamModel[]>([]);
const martParamStore = useMartParamStore();
const openAlert = ref(false);

onMounted(() => {
  if (martParamStore.getActiveMartParam === null) {
    loadMyMartParam();
  }
});

function loadMyMartParam() {
  loading.value = true;
  MartParamService.LoadLite().then((resp) => {
    loading.value = false;

    if (resp) {
      if (resp.code == 0) {
        myMartParams.value = resp.data;
      } else if (resp.code) {
        toast({
          variant: "destructive",
          title: t("error." + resp.code),
        });
      }
    }

    if (myMartParams.value.length === 0) {
      openAlert.value = true;
    }

    if (myMartParams.value.length > 0) {
      currMartParam.value = myMartParams.value[0];
      for (let i = 0; i < myMartParams.value.length; i++) {
        if (myMartParams.value[i].active > currMartParam.value.active) {
          currMartParam.value = myMartParams.value[i];
        }
      }

      martParamStore.SetActiveMartParam(currMartParam.value);
    }
  });
}

async function setCurrentMartParam(val: any) {
  loading.value = true;
  currMartParam.value = {
    id: 0,
    uid: 0,
    domain: "domain",
    symbol: "symbol",
    memo: "memo",
    accessKey: "",
    secretKey: "",
    active: 0,
  };
  martParamStore.SetActiveMartParam(null);
  const resp = await MartParamService.Active(val.id);
  loading.value = false;

  if (resp && resp.code == 0) {
    toast({
      title: "Change symbol success.",
    });
    await loadMyMartParam();
  } else if (resp) {
    toast({
      variant: "destructive",
      title: t("error." + resp.code),
    });
  }
}

const open = ref(false);
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <Button
        :disabled="loading || myMartParams.length === 0"
        variant="outline"
        role="combobox"
        aria-expanded="open"
        aria-label="Select a team"
        :class="cn('w-[245px] justify-between mr-6', $attrs.class ?? '')"
      >
        {{ currMartParam?.symbol }}
        <span class="text-xs">@</span>
        {{ currMartParam?.domain }}
        <CaretSortIcon class="ml-auto h-4 w-4 shrink-0 opacity-50" />
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-[245px] p-0">
      <Command @update:model-value="setCurrentMartParam">
        <CommandList>
          <CommandItem
            v-for="mart in myMartParams"
            :key="mart.id"
            :value="mart"
            class="text-sm"
            @select="
              () => {
                currMartParam = mart;
                open = false;
              }
            "
          >
            {{ mart?.symbol }}
            @
            {{ mart?.domain }}
            <CheckIcon
              :class="
                cn(
                  'ml-auto h-4 w-4',
                  currMartParam?.id === mart.id ? 'opacity-100' : 'opacity-0'
                )
              "
            />
          </CommandItem>
        </CommandList>
      </Command>
    </PopoverContent>
  </Popover>

  <AlertDialog v-model:open="openAlert">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>Please configure APIKEY</AlertDialogTitle>
        <AlertDialogDescription>
          You have not configured the APIKEY, you need to configure the APIKEY
          before proceeding.
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <RouterLink to="/user">
          <AlertDialogAction>Go to configure the APIKEY</AlertDialogAction>
        </RouterLink>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>

  <Toaster />
</template>
