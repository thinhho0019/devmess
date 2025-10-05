import { BrowserRouter } from "react-router-dom";
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import './i18n';
import { SocketProvider } from "./contexts/SocketContext.tsx";

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <SocketProvider>
      <App />
    </SocketProvider>
  </BrowserRouter>
)
