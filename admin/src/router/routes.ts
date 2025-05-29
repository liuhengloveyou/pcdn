import DashboardView from "@/pages/DashboardView.vue";
import EmptyLayout from "../layouts/EmptyLayout.vue";
import MainLayout from "../layouts/MainLayout.vue";
import LoginView from "../pages/LoginView.vue";
import SystemMonitorView from '@/pages/SystemMonitorView.vue'


const routes = [
  {
    path: "/login",
    component: EmptyLayout,
    children: [{ path: "", component: LoginView }],
  },
  {
    path: "/",
    component: MainLayout,
    children: [
      {
        path: '/system-monitor',
        name: 'SystemMonitor',
        component: SystemMonitorView,
        meta: {
          requiresAuth: true,
          layout: 'main'
        }
      },
      { path: "/dashboard", component: DashboardView },
      { path: "/devices", component: () => import("../pages/devices/index.vue")},
      { path: "", redirect: "/dashboard"},

    ],
  },

  // Always leave this as last one,
  // but you can also remove it
];

export default routes;
