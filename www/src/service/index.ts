export interface HttpResponse<T = never> {
  code: number;
  total?: number;
  msg: string;
  data: T;
}
