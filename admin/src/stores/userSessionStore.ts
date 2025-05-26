import { defineStore } from 'pinia';
import type { UserInfoStruct } from 'src/services';

export const useUserSessionStore = defineStore('UserSessionStore', {
  state: () => ({
    account: undefined as UserInfoStruct | undefined,
  }),
  getters: {
    getSessionUser: (state) => {
      return state.account;
    },
    isLoign: (state) => {
      return (
        state.account != undefined &&
        state.account.uid != undefined &&
        state.account.uid > 0
      );
    },
  },
  actions: {
    login(userInfo: UserInfoStruct) {
      this.$state.account = userInfo;
    },
    logout() {
      this.$state.account = undefined;
    },
  },
});
