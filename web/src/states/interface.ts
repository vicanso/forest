export interface IList<T> {
  processing: boolean;
  items: T[];
  count: -1;
}

export interface IStatus {
  desc: string;
  value: number;
}
