/* eslint-disable no-console */
import axios from 'axios';



// Be careful when using SSR for cross-request state pollution
// due to creating a Singleton instance here;
// If any client changes this (global) instance, it might be a
// good idea to move this instance creation inside of the
// "export default () => {}" function below (which runs individually
// for each client)
const api = axios.create({
  //  baseURL: 'https://api.example.com',
  validateStatus: function (status) {
    return status < 500;
  },
});

// 添加请求拦截器
axios.interceptors.request.use(
  function (config) {
    // 在发送请求之前做些什么
    return config;
  },
  function (error) {
    // 对请求错误做些什么
    return Promise.reject(error);
  },
);

// 添加响应拦截器
axios.interceptors.response.use(
  function (response) {
    if (response.status >= 404) {
      console.log( '网络错误');
    }

    if (response.status == 401) {
      // setTimeout(() => {
      //   void Router.push({
      //     path: '/login',
      //   });
      // }, 1000);
      // router.replace({ path: '/login' });
    } else if (response.status == 403) {
      console.log( '权限错误');
    }
    // 对响应数据做点什么
    return response;
  },
  function (error) {
    console.log('axios error: ', error);
    if (error.response) {
      if (error.response.data && error.response.data.message) {
        // Notify.create({
        //   message: error.response.data.message,
        //   type: 'negative',
        // });
      } else {
        // Notify.create({
        //   message: error.response.statusText || error.response.status,
        //   type: 'negative',
        // });
      }
    } else if (error.message.indexOf('timeout') > -1) {
      // Notify.create({
      //   message: 'Network timeout',
      //   type: 'negative',
      // });
    } else if (error.message) {
      // Notify.create({
      //   message: error.message,
      //   type: 'negative',
      // });
    } else {
      // Notify.create({
      //   message: 'http request error',
      //   type: 'negative',
      // });
    }
    // 对响应错误做点什么
    return Promise.reject(error);
  },
);



export default api;
