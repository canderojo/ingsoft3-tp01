import { Routes, Route } from "react-router-dom";
import Header from "./components/Header";
import Footer from "./components/Footer";
import ProfesionalesPage from "./pages/ProfesionalesPage";
import ProfesionalDetailPage from "./pages/ProfesionalDetailPage";
import MisTurnosPage from "./pages/MisTurnosPage";
import TurnoDetailPage from "./pages/TurnoDetailPage";
import NotFoundPage from "./pages/NotFoundPage";

export default function App() {
  return (
    <>
      <Header />
      <main className="app-main">
        <Routes>
          <Route path="/" element={<ProfesionalesPage />} />
          <Route path="/profesionales/:id" element={<ProfesionalDetailPage />} />
          <Route path="/mis-turnos" element={<MisTurnosPage />} />
          <Route path="/turnos/:id" element={<TurnoDetailPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>
      <Footer />
    </>
  );
}
