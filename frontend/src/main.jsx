import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import * as Sentry from "@sentry/react";
import App from "./App.jsx";
import "./styles/global.css";

// El DSN de un front es público por naturaleza: sólo sirve para enviar datos.
Sentry.init({
  dsn: "https://f4d1ffa13d678356040775028525b0bc@o4512313576914944.ingest.us.sentry.io/4512313640878080",
  environment: window.location.hostname, // de qué sitio vino: QA, producción o localhost
  integrations: [Sentry.browserTracingIntegration()],
  tracesSampleRate: 1.0,
  // GET /turnos?dni=...&email=... lleva datos del paciente en el query string.
  dataCollection: { urlQueryParams: false },
});

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>
);
