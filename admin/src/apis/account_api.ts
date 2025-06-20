/* eslint-disable no-console */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import api from "../utils/axios.ts";
import type { HttpResponse, UserInfo } from "./index.ts";

export class AccountApi {
  static async userAuth() {
    const resp = await api.get<HttpResponse<UserInfo>>("/usercenter", {
      headers: {
        "X-API": "user/auth",
      },
    });

    // console.log('AccountService.userAuth:', resp.data);

    return resp.data;

    //   this.store.commit("userAuth", resp.data.data);
    //   router.replace({ path: '/' })
  }

  static async login(email: string, password: string, sms: string) {
    const resp = await api.post<HttpResponse<UserInfo>>(
      "/usercenter",
      {
        cellphone: email,
        password: password,
        sms: sms,
      },
      {
        headers: {
          "X-API": "user/login",
        },
      }
    );

    return resp.data;
  }

  static async logout() {
    const resp = await api.post<HttpResponse<string>>("/usercenter", null, {
      headers: {
        "X-API": "user/logout",
      },
    });

    if (resp.status != 200) {
      return resp;
    }

    if (resp.data.code != 0) {
      //   presentToast(resp.data.msg);
    } else {
      console.log("logout ok:", resp.data);
    }

    return resp;
  }

  static async sendLoginSms(cellphone: string, aliveSec: number) {
    try {
      const resp = await api.post<HttpResponse<string>>(
        "/usercenter",
        {
          cellphone: cellphone,
          aliveSec: aliveSec,
        },
        {
          headers: {
            "X-API": "sms/sendUserLoginSms",
          },
        }
      );

      if (resp.status != 200) {
        return;
      }

      if (resp.data.code != 0) {
        // // Notify.create(resp.data.msg);
      }

      return resp.data;
    } catch (error) {
      console.log("err: ", error);
      return;
    }
  }

  static async sendWxBindSms(cellphone: string, aliveSec: number) {
    try {
      const resp = await api.post<HttpResponse<string>>(
        "/usercenter",
        {
          cellphone: cellphone,
          aliveSec: aliveSec,
        },
        {
          headers: {
            "X-API": "sms/sendWxBindSms",
          },
        }
      );

      if (resp.status != 200) {
        return;
      }

      if (resp.data.code != 0) {
        // Notify.create(resp.data.msg);
      }

      return resp.data;
    } catch (error) {
      console.log("err: ", error);
      return;
    }
  }

  static async update(data: unknown) {
    const resp = await api.post<HttpResponse<string>>("/usercenter", data, {
      headers: {
        "X-API": "user/modify",
      },
    });

    if (resp.status != 200) {
      return resp;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }

  static async info() {
    const resp = await api.get<HttpResponse<UserInfo>>("/usercenter", {
      headers: {
        "X-API": "user/info",
      },
    });

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    }

    return resp.data;
  }

  static async updateAvatarByForm(img: File) {
    const formData = new FormData();

    const names = img.name.split(".");
    formData.set("ext_name", names[names.length - 1]);
    formData.append("file", img);

    const resp = await fetch("/usercenter", {
      method: "POST",
      headers: {
        "Content-Type": "multipart/form-data", // 设置Content-Type为multipart/form-data
        "X-API": "user/modify/avatarForm",
      },
      body: formData,
    });

    // const resp = await api.post<HttpResponse>('/usercenter', data, {
    //   headers: {
    //     'X-API': 'user/modify/avatarForm',
    //     'Content-Type': 'multipart/form-data',
    //   },
    //   maxContentLength: 10000000,
    // });

    if (resp.status != 200) {
      return;
    }

    // if (resp.data.code != 0) {
    //   // Notify.create(resp.data.msg);
    // } else {
    //   // Notify.create('成功');
    // }

    // return resp.data;
  }

  static async updatePwd(n: string, o: string) {
    const resp = await api.post<HttpResponse<string>>(
      "/usercenter",
      { n: n, o: o },
      {
        headers: {
          "X-API": "user/modify/password",
        },
      }
    );

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }

  // 建账号
  static async create(obj: UserInfo) {
    const resp = await api.post<HttpResponse<string>>("/usercenter", obj, {
      headers: {
        "X-API": "tenant/user/add",
      },
    });

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else if (resp.data.data) {
      // Notify.create('成功');
    }

    return resp.data;
  }

  static async delUserByUID(uid: number) {
    const resp = await api.post<HttpResponse<string>>(
      "/usercenter",
      {
        uid: uid,
      },
      {
        headers: {
          "X-API": "tenant/delUser",
        },
      }
    );

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }

  static async userDisableByUID(uid: number, disable: number) {
    const resp = await api.post<HttpResponse<string>>(
      "/usercenter",
      {
        uid: uid,
        disable: disable,
      },
      {
        headers: {
          "X-API": "tenant/userDisableByUID",
        },
      }
    );

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }

  static async setExt(uid: number, k: string, v: unknown) {
    const resp = await api.post<HttpResponse<string>>(
      "/usercenter",
      {
        id: uid,
        k: k,
        v: v,
      },
      {
        headers: {
          "X-API": "tenant/userModifyExtInfo",
        },
      }
    );

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }

  static async setDepartment(uid: number, depId: number) {
    const resp = await api.post<HttpResponse<string>>(
      "/usercenter",
      {
        uid: uid,
        depIds: [depId],
      },
      {
        headers: {
          "X-API": "tenant/user/setDepartment",
        },
      }
    );

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }

  static async addRoleForUser(uid: number, roleVal: string) {
    const resp = await api.post<HttpResponse<string>>(
      "/usercenter",
      {
        uid: uid,
        value: roleVal,
      },
      {
        headers: {
          "X-API": "access/addRoleForUser",
        },
      }
    );

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }

  static async removeRoleForUser(uid: number, roleVal: string) {
    const resp = await api.post<HttpResponse<string>>(
      "/usercenter",
      {
        uid: uid,
        value: roleVal,
      },
      {
        headers: {
          "X-API": "access/removeRoleForUser",
        },
      }
    );

    if (resp.status != 200) {
      return;
    }

    if (resp.data.code != 0) {
      // Notify.create(resp.data.msg);
    } else {
      // Notify.create('成功');
    }

    return resp.data;
  }
}
