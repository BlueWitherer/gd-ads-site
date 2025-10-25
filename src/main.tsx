import "./index.css";
import App from "./App.tsx";
import Dashboard from "./Dashboard.tsx";
import Admin from "./dashboard/Admin.tsx";
import NotFound from "./NotFound.tsx";
import MobileWarning from "./MobileWarning";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <MobileWarning />
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/admin" element={<Admin />} />

        {/* 404 page */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>
);