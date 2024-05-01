
const content = document.querySelector("script#pulseup__config")?.textContent || "{}";

export interface Config {
    hostname: string;
    version: string;
    address: string;
    base: string;
    hosts: { name: string; id: string }[];
    authProvider: "simple" | "none" ;
  }
  
  const jsonConfig = JSON.parse(content);

  const config: Config = {
    ...jsonConfig,
  };
  
  export default Object.freeze(config);

  export const withBase = (path: string) => `${config.base}${path}`;