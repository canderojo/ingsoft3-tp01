import { useEffect, useState } from "react";
import { Link, useLocation, useParams } from "react-router-dom";
import { obtenerTurno, cambiarEstadoTurno, TRANSICIONES_PERMITIDAS, ESTADOS } from "../api/turnos";
import { obtenerProfesional } from "../api/profesionales";
import { ApiError } from "../api/client";
import { formatFechaHora, formatPrecio } from "../utils/format";
import EstadoBadge from "../components/EstadoBadge";
import EspecialidadBadge from "../components/EspecialidadBadge";
import CheckIcon from "../components/icons/CheckIcon";
import "./TurnoDetailPage.css";

// "completado" es una transición real del turno (TRANSICIONES_PERMITIDAS),
// pero es el centro médico quien marca una consulta como completada, no
// el propio paciente — por eso no se ofrece como botón acá.
const ACCION_POR_ESTADO = {
  confirmado: { texto: "Confirmar turno", clase: "btn-primary" },
  cancelado: { texto: "Cancelar turno", clase: "btn-danger" },
};

export default function TurnoDetailPage() {
  const { id } = useParams();
  const location = useLocation();

  const [turno, setTurno] = useState(null);
  const [profesional, setProfesional] = useState(null);
  const [error, setError] = useState(null);
  const [cambiandoEstado, setCambiandoEstado] = useState(null);
  const [errorEstado, setErrorEstado] = useState(null);

  function cargarTurno() {
    setError(null);
    obtenerTurno(id)
      .then((t) => {
        setTurno(t);
        return obtenerProfesional(t.profesional_id);
      })
      .then(setProfesional)
      .catch((err) =>
        setError(err instanceof ApiError ? err.message : "No pudimos cargar este turno.")
      );
  }

  useEffect(cargarTurno, [id]);

  async function handleCambiarEstado(nuevoEstado) {
    setCambiandoEstado(nuevoEstado);
    setErrorEstado(null);
    try {
      const actualizado = await cambiarEstadoTurno(id, nuevoEstado);
      setTurno(actualizado);
    } catch (err) {
      setErrorEstado(err instanceof ApiError ? err.message : "No pudimos actualizar el turno.");
    } finally {
      setCambiandoEstado(null);
    }
  }

  if (error) {
    return (
      <div className="page">
        <div className="alert alert-error">{error}</div>
        <Link to="/mis-turnos">Volver a mis turnos</Link>
      </div>
    );
  }

  if (!turno || !profesional) {
    return (
      <div className="page">
        <div className="spinner" />
      </div>
    );
  }

  const transiciones = (TRANSICIONES_PERMITIDAS[turno.estado] ?? []).filter(
    (estado) => estado in ACCION_POR_ESTADO
  );

  return (
    <div className="page">
      <Link to="/mis-turnos" className="back-link">
        ← Mis turnos
      </Link>

      {location.state?.recienCreado && (
        <div className="alert alert-success alert-success-icon">
          <CheckIcon />
          <span>Tu turno se reservó con éxito.</span>
        </div>
      )}

      <div className="card turno-detail-card">
        <div className="turno-detail-header">
          <h1>Turno #{turno.id}</h1>
          <EstadoBadge estado={turno.estado} />
        </div>

        <p className="turno-detail-fecha">{formatFechaHora(turno.fecha_hora_inicio)}</p>

        <div className="turno-detail-profesional">
          <div>
            <p className="turno-detail-label muted">Profesional</p>
            <p className="turno-detail-value">{profesional.nombre}</p>
          </div>
          <EspecialidadBadge especialidad={profesional.especialidad} />
        </div>

        <div>
          <p className="turno-detail-label muted">Precio</p>
          <p className="turno-detail-value">{formatPrecio(turno.precio)}</p>
        </div>

        {errorEstado && <div className="alert alert-error">{errorEstado}</div>}

        {transiciones.length > 0 && (
          <div className="turno-detail-acciones">
            {transiciones.map((estado) => (
              <button
                key={estado}
                type="button"
                className={`btn ${ACCION_POR_ESTADO[estado].clase}`}
                disabled={cambiandoEstado !== null}
                onClick={() => handleCambiarEstado(estado)}
              >
                {cambiandoEstado === estado ? "Actualizando..." : ACCION_POR_ESTADO[estado].texto}
              </button>
            ))}
          </div>
        )}

        {transiciones.length === 0 && (
          <p className="muted">Este turno ya está {ESTADOS[turno.estado]?.toLowerCase()} y no admite cambios.</p>
        )}
      </div>
    </div>
  );
}
