<script setup lang="ts">
import { AspectRatio } from "@/components/ui/aspect-ratio";
// import { cn } from '@/lib/utils'
</script>

<template>
  <div class="w-[42px] mr-8">
    <AspectRatio :ratio="1 / 1">
      <img
        src="/vite.svg"
        alt="moliaibot.com"
        class="rounded-md object-cover w-full h-full transition-all hover:scale-105"
      />
    </AspectRatio>
  </div>
</template>
