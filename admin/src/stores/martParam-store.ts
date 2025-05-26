import { defineStore } from 'pinia';

import type { MartParamModel } from 'src/services/MartParamService';

export const useMartParamStore = defineStore('martParam', {
  state: () => ({
    activeMartParam: null as null | MartParamModel,
  }),
  getters: {
    getActiveMartParam: (state) => state.activeMartParam,
  },
  actions: {
    SetActiveMartParam(data: MartParamModel | null) {
      this.$state.activeMartParam = data;
    },
  },
});
