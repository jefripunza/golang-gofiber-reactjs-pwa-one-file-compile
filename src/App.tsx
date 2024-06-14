import { useEffect, useState } from "react";
import PWABadge from "./PWABadge.tsx";

import { hostname } from "./consts.ts";
import { env } from "./env.ts";

import golangLogo from "./assets/golang.svg";
import appLogo from "/favicon.svg";
import reactLogo from "./assets/react.svg";
import typescriptLogo from "./assets/typescript.svg";

const APP_NAME = import.meta.env.VITE_APP_NAME;

import "./App.css";

function App() {
  const [count, setCount] = useState<number>(0);
  const [response, setResponse] = useState<string>("???");

  async function getHelloWorld() {
    try {
      const response = await fetch(`${hostname}/api/example`);
      const res = await response.json();
      setResponse(JSON.stringify(res));
    } catch (error: any) {
      setResponse(error.message);
    }
  }

  useEffect(() => {
    console.log({ env });
  }, []);

  return (
    <>
      <div>
        <a href="https://go.dev" target="_blank">
          <img src={golangLogo} className="logo" alt="Golang logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
        <a href="https://vitejs.dev" target="_blank">
          <img src={appLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://www.typescriptlang.org" target="_blank">
          <img src={typescriptLogo} className="logo" alt="TypeScript logo" />
        </a>
      </div>
      <h1>Golang + ReactJS + VitePWA + TS</h1>
      <h1>
        <i>{`{ ${APP_NAME} }`}</i>
      </h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <button onClick={getHelloWorld}>response is {response}</button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Golang, Vite and React logos to learn more
      </p>
      <PWABadge />
    </>
  );
}

export default App;
