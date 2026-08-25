import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listarProfesionales, ESPECIALIDADES } from "../api/profesionales";
import { ApiError } from "../api/client";
import { formatHora, formatPrecio } from "../utils/format";
import EspecialidadBadge from "../components/EspecialidadBadge";
import EspecialidadIcon from "../components/icons/EspecialidadIcon";
import EmptyIcon from "../components/icons/EmptyIcon";
import "./ProfesionalesPage.css";

export default function ProfesionalesPage() {
  const [profesionales, setProfesionales] = useState([]);
  const [especialidad, setEspecialidad] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    setError(null);
    listarProfesionales(especialidad)
      .then(setProfesionales)
      .catch((err) => setError(err instanceof ApiError ? err.message : "No pudimos cargar los profesionales."))
      .finally(() => setLoading(false));
  }, [especialidad]);

  return (
    <>
      <section className="hero">
        <div className="hero-inner">
          <p className="hero-eyebrow">Turnos online</p>
          <h1>Reservá tu consulta con el profesional que necesitás</h1>
          <p className="hero-subtitulo">
            Elegí el horario que te quede mejor y reservá en minutos, sin
            llamados ni esperas.
          </p>

          <div className="tabs" role="group" aria-label="Filtrar por especialidad">
            <button
              type="button"
              className={`tab${especialidad === "" ? " tab-active" : ""}`}
              onClick={() => setEspecialidad("")}
            >
              Todas
            </button>
            {ESPECIALIDADES.map((e) => (
              <button
                key={e.value}
                type="button"
                className={`tab${especialidad === e.value ? " tab-active" : ""}`}
                onClick={() => setEspecialidad(e.value)}
              >
                <EspecialidadIcon especialidad={e.value} size={15} />
                {e.label}
              </button>
            ))}
          </div>
        </div>
      </section>

      <div className="page">
        {loading && <div className="spinner" />}
        {error && <div className="alert alert-error">{error}</div>}

        {!loading && !error && profesionales.length === 0 && (
          <div className="empty-state">
            <EmptyIcon />
            <p>No hay profesionales disponibles para este filtro.</p>
          </div>
        )}

        <div className="profesionales-grid">
          {profesionales.map((p) => (
            <Link
              to={`/profesionales/${p.id}`}
              key={p.id}
              className="profesional-card card"
              style={{ "--card-accent": `var(--color-${p.especialidad})` }}
            >
              <div className="profesional-card-header">
                <span className="profesional-card-icon">
                  <EspecialidadIcon especialidad={p.especialidad} size={18} />
                </span>
                <div>
                  <h3>{p.nombre}</h3>
                  <EspecialidadBadge especialidad={p.especialidad} />
                </div>
              </div>
              <p className="muted">
                Atiende de {formatHora(p.hora_inicio_atencion)} a {formatHora(p.hora_fin_atencion)}
              </p>
              <div className="profesional-card-footer">
                <span className="profesional-card-precio">{formatPrecio(p.precio_consulta)}</span>
                <span className="profesional-card-cta">Ver horarios →</span>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </>
  );
}
