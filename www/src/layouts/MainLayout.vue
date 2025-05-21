<template>
  <q-layout view="lHh Lpr lFf">
    <q-header>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />

        <q-toolbar-title> PCDN管理系统 </q-toolbar-title>

        <div>v0.1</div>
      </q-toolbar>
    </q-header>

    <q-drawer
      show-if-above
      :mini="miniState"
      :width="255"
      :breakpoint="500"
      style="background-color: #191b20; color: #fafafa"
    >
      <q-card v-if="!miniState" flat style="background-color: #191b20; color: #fafafa">
        <q-card-section class="text-h6 text-center">
          <q-avatar size="56px" class="q-mb-sm">
            <img src="https://cdn.quasar.dev/img/boy-avatar.png" />
          </q-avatar>
          <div>
            <div class="text-weight-bold">刘恒</div>
            <div>@liuheng</div>
          </div>
        </q-card-section>
      </q-card>

      <q-scroll-area style="height: calc(100% - 160px); border-right: 1px solid #191b20">
        <q-list padding>
          <EssentialLink v-for="link in linksList" :key="link.title" v-bind="link" />
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import EssentialLink, { type EssentialLinkProps } from 'components/EssentialLink.vue';

const linksList: EssentialLinkProps[] = [
  {
    title: '设备管理',
    caption: '管理所有设备',
    icon: 'devices',
    link: '#/devices', // 修改为指向设备管理页面的链接
  },
  {
    title: '会员管理',
    caption: '管理系统会员',
    icon: 'people', // 人员/用户相关图标
    link: '#/members', // 修改为指向会员管理页面的链接
  },
  {
    title: '服务商管理',
    caption: '管理系统服务商',
    icon: 'business', // 企业/商业相关图标
    link: '#/providers', // 修改为指向服务商管理页面的链接
  },
  {
    title: '提现管理',
    caption: 'forum.quasar.dev',
    icon: 'payments', // 更改为支付/金钱相关图标
    link: 'https://forum.quasar.dev',
  },
  {
    title: '日结账单',
    caption: '@quasarframework',
    icon: 'receipt_long', // 更改为账单/收据相关图标
    link: 'https://twitter.quasar.dev',
  },
  {
    title: '激活码管理',
    caption: '@QuasarFramework',
    icon: 'key', // 更改为钥匙/激活相关图标
    link: 'https://facebook.quasar.dev',
  },
  {
    title: '上游计费管理',
    caption: 'Community Quasar projects',
    icon: 'account_balance', // 更改为财务/银行相关图标
    link: 'https://awesome.quasar.dev',
  },
  {
    title: '会员级别',
    caption: '配置不同会员级别的提成比例',
    icon: 'stars', // 星级/等级相关图标
    link: '#/member-levels', // 修改为指向会员级别页面的链接
  },
  {
    title: '服务商级别',
    caption: '配置不同服务商级别的提成比例',
    icon: 'workspace_premium', // 高级/优质相关图标
    link: '#/provider-levels', // 修改为指向服务商级别页面的链接
  },
  {
    title: '收益统计',
    caption: 'Community Quasar projects',
    icon: 'bar_chart', // 更改为图表/统计相关图标
    link: 'https://awesome.quasar.dev',
  },
  {
    title: '系统管理',
    caption: '系统设置与权限管理',
    icon: 'settings',
    link: '#/system',
  },
];

const miniState = ref(false);
// const leftDrawerOpen = ref(false);

function toggleLeftDrawer() {
  miniState.value = !miniState.value;
}
</script>
