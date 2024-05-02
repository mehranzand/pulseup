export type Log = {
    readonly message: string | undefined;
    readonly date: number | undefined;
    readonly type: LogType | undefined;
    readonly rowType: RowType
  };
  export type LogType = "stdout" | "stderr";
  export type RowType =  "gutter" | "log";