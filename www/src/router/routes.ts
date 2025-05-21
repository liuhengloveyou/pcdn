import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: () => import('layouts/EmptyLayout.vue'),
    children: [{ path: '', component: import('pages/LoginPage.vue') }],
  },
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/IndexPage.vue') },
      { path: 'devices', component: () => import('pages/DeviceManagementPage.vue') },
      { path: 'members', component: () => import('pages/MemberManagementPage.vue') },
      { path: 'providers', component: () => import('pages/ProviderManagementPage.vue') },
      { path: 'member-levels', component: () => import('pages/MemberLevelPage.vue') },
      { path: 'provider-levels', component: () => import('pages/ProviderLevelPage.vue') },
      { path: 'system', component: () => import('pages/SystemManagementPage.vue') },
    ],
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
