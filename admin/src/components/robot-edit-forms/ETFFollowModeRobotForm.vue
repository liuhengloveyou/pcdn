<script setup lang="ts">
import { ScrollArea } from "@/components/ui/scroll-area";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  NumberField,
  NumberFieldContent,
  NumberFieldDecrement,
  NumberFieldIncrement,
  NumberFieldInput,
} from "@/components/ui/number-field";
import Toaster from "@/components/ui/toast/Toaster.vue";
import { useToast } from "@/components/ui/toast/use-toast";
const { toast } = useToast();
// import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import { h, ref } from "vue";
// import * as z from "zod";

const loading = ref(false);

// const formSchema = toTypedSchema(
//   z.object({
//     priceMode: z.enum(["normalMode", "safeMode"], {
//       required_error: "You need to select a notification type.",
//     }),
//     launchMode: z.enum(["immediateStart", "scheduledStart"], {
//       required_error: "You need to select a notification type.",
//     }),
//   })
// );

const { handleSubmit } = useForm({
  // validationSchema: formSchema,
});

const onSubmit = handleSubmit((values) => {
  console.log(">?>>", values);
  toast({
    variant: "destructive",
    title: "You submitted the following values:",
    description: h(
      "pre",
      { class: "mt-2 w-[340px] rounded-md bg-slate-950 p-4" },
      h("code", { class: "text-white" }, JSON.stringify(values, null, 2))
    ),
  });
});
</script>

<template>
  <ScrollArea class="h-[90vh] w-full">
    <form class="space-y-6" @submit="onSubmit">
      <div>
        <FormField v-slot="{ componentField }" name="buyOrderSize">
          <FormLabel>{{ $t("buyOrderSize") }}</FormLabel>
          <div class="flex row justify-between">
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
            <span style="text-align: center; align-content: center"> ~ </span>
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
          </div>
        </FormField>
      </div>

      <div>
        <FormField v-slot="{ componentField }" name="buySheet">
          <FormLabel>{{ $t("buySheet") }}</FormLabel>
          <div class="flex row justify-between">
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
            <span style="text-align: center; align-content: center"> ~ </span>
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
          </div>
        </FormField>
      </div>

      <div>
        <FormField v-slot="{ componentField }" name="sellOrderSize">
          <FormLabel>{{ $t("sellOrderSize") }}</FormLabel>
          <div class="flex row justify-between">
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
            <span style="text-align: center; align-content: center"> ~ </span>
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
          </div>
        </FormField>
      </div>

      <div>
        <FormField v-slot="{ componentField }" name="orderInterval">
          <FormLabel>{{ $t("orderInterval") }}</FormLabel>
          <div class="flex row justify-between">
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
            <span style="text-align: center; align-content: center"> ~ </span>
            <FormItem>
              <NumberField class="gap-2" :min="0" v-bind="componentField">
                <NumberFieldContent>
                  <NumberFieldDecrement />
                  <FormControl>
                    <NumberFieldInput />
                  </FormControl>
                  <NumberFieldIncrement />
                </NumberFieldContent>
              </NumberField>
              <FormMessage />
            </FormItem>
          </div>
        </FormField>
      </div>

      <FormField v-slot="{ componentField }" name="orderType">
        <FormItem>
          <FormLabel>{{ $t("orderMode") }}:</FormLabel>
          <FormControl>
            <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="random" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("random") }}
                </FormLabel>
              </FormItem>
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="buyFirst" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("buyFirst") }}
                </FormLabel>
              </FormItem>
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="sellFirst" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("sellFirst") }}
                </FormLabel>
              </FormItem>
            </RadioGroup>
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="priceMode">
        <FormItem>
          <FormLabel>{{ $t("priceMode") }}:</FormLabel>
          <FormControl>
            <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="normalMode" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("normalMode") }}
                </FormLabel>
              </FormItem>
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="safeMode" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("safeMode") }}
                </FormLabel>
              </FormItem>
            </RadioGroup>
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="username">
        <FormItem>
          <FormLabel>{{ $t("launchMode") }}:</FormLabel>
          <FormControl>
            <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="immediateStart" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("immediateStart") }}
                </FormLabel>
              </FormItem>
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="scheduledStart" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("scheduledStart") }}
                </FormLabel>
              </FormItem>
            </RadioGroup>
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="amountDistributionType">
        <FormItem>
          <FormLabel>{{ $t("amountDistributionType") }}:</FormLabel>
          <FormControl>
            <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="continuousDistribution" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("continuousDistribution") }}
                </FormLabel>
              </FormItem>
              <FormItem class="flex items-center space-y-0 gap-x-3">
                <FormControl>
                  <RadioGroupItem value="randomDistribution" />
                </FormControl>
                <FormLabel class="font-normal">
                  {{ $t("randomDistribution") }}
                </FormLabel>
              </FormItem>
            </RadioGroup>
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="note">
        <FormItem>
          <FormLabel>{{ $t("note") }}</FormLabel>
          <FormControl>
            <Textarea class="w-full" v-bind="componentField"></Textarea>
          </FormControl>
        </FormItem>
      </FormField>

      <Button :disabled="loading" type="submit" class="w-full">
        <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
        {{ $t("submit") }}
      </Button>
    </form>
  </ScrollArea>
  <Toaster />
</template>
