
  export type Container = {
    readonly id: string;
    readonly created: number;
    readonly image: string;
    readonly name: string;
    readonly command: string;
    readonly status: string;
    readonly state: ContainerState;
    readonly host: string;
    readonly labels: Record<string, string>;
    readonly health?: ContainerHealth;
  };
  
  export type ContainerState = "created" | "running" | "exited" | "dead" | "paused" | "restarting" | "stopped";
  export type ContainerHealth = "healthy" | "unhealthy" | "starting";