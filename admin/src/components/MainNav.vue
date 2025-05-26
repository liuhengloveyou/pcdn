<script setup lang="ts">
import { cn } from "@/lib/utils";
import { onMounted } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();

onMounted(() => {
  console.log(route.path);
});
</script>

<template>
  <nav
    :class="cn('flex items-center space-x-4 lg:space-x-6', $attrs.class ?? '')"
  >
    <router-link
      to="/"
      :class="
        cn(
          'text-sm font-medium text-muted-foreground transition-colors hover:text-primary',
          route.path === '/' ? 'text-primary' : ''
        )
      "
    >
      {{ $t("menu.ControlPannel") }}
    </router-link>
    <router-link
      to="/botdeploy"
      :class="
        cn(
          'text-sm font-medium text-muted-foreground transition-colors hover:text-primary',
          route.path === '/botdeploy' ? 'text-primary' : ''
        )
      "
    >
      {{ $t("menu.RobotDeployment") }}
    </router-link>

    <router-link
      to="/workpanel"
      :class="
        cn(
          'text-sm font-medium text-muted-foreground transition-colors hover:text-primary',
          route.path === '/workpanel' ? 'text-primary' : ''
        )
      "
    >
      {{ $t("menu.WorkSurface") }}
    </router-link>

    <router-link
      to="/riskman"
      :class="
        cn(
          'text-sm font-medium text-muted-foreground transition-colors hover:text-primary',
          route.path === '/riskman' ? 'text-primary' : ''
        )
      "
    >
      {{ $t("menu.RiskControl") }}
    </router-link>

    <router-link
      to="/dataana"
      :class="
        cn(
          'text-sm font-medium text-muted-foreground transition-colors hover:text-primary',
          route.path === '/dataana' ? 'text-primary' : ''
        )
      "
    >
      {{ $t("menu.DataAnalysis") }}
    </router-link>
  </nav>
</template>
