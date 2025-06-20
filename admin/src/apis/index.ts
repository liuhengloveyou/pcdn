export interface HttpResponse<T = never> {
  code: number;
  total?: number;
  msg: string;
  data: T;
}


export interface UserInfo {
  uid: number;
  tenant_id: number;
  nickname: string;
  cellphone: string;
  email: string;
  avatarUrl?: string;
  loginTime?: string;
  password?: string;
  gender?: string;
  disable?: number;
  ext: {
    dep?: string;
    deps?: number[];
    disabled?: number;
    isMarket?: boolean;
  };

  label?: string; // 给select显示用
  checked?: boolean; // 选择营销人员
  tcNum?: number; // 营销人员提成份数
  idx?: number; // 显示序号
}
