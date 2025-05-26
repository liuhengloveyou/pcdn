import {
  createRouter,
  createWebHistory,
} from 'vue-router';
import routes from './routes';
import { useUserSessionStore } from '@/stores/userSessionStore';
import { AccountService } from '@/services/AccountService';

const Router = createRouter({
  scrollBehavior: () => ({ left: 0, top: 0 }),
  routes,

  history: createWebHistory(),
});

Router.beforeEach(async (to, from, next) => {
  console.log('Router.beforeEach:: ', from.fullPath, to.fullPath);

  if (to.fullPath == '/login') {
    console.log('Router.beforeEach need not auth', to.fullPath);
    return next();
  }

  const uerSessionStore = useUserSessionStore();
  const isLoign = uerSessionStore.isLoign;
  console.log('Router.beforeEach:: ', uerSessionStore.getSessionUser);

  if (isLoign == false) {
    const authResp = await AccountService.userAuth();
    console.log('router.beforeEach userAuth: ', authResp);
    if (authResp.code == 0 && authResp.data.uid && authResp.data.uid > 0) {
      uerSessionStore.login(authResp.data);
    } else {
      return next({ path: '/login', replace: true });
    }
  } else if (uerSessionStore.isLoign == true) {
    return next(); // root
  }

  return next();
});

export default Router;
