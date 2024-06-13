export const isDEV = import.meta.env.DEV;
export const hostname = isDEV ? "http://localhost:1337" : "";
