export type Log = {
    readonly message: string;
    readonly date: number;
    readonly type: LogType;
  };
  export type LogType = "stdout" | "stderr";