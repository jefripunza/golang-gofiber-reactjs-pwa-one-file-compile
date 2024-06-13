import { useState } from "react";
import reactLogo from "./assets/react.svg";
import appLogo from "/favicon.svg";
import PWABadge from "./PWABadge.tsx";
import "./App.css";
import { hostname } from "./consts.ts";

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

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={appLogo} className="logo" alt="App Example logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>App Example</h1>
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
        Click on the Vite and React logos to learn more
      </p>
      <PWABadge />
    </>
  );
}

export default App;
