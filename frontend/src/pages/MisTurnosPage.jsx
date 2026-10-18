import { useState } from "react";
import { Link } from "react-router-dom";
import { listarTurnosDePaciente, ESTADOS } from "../api/turnos";
import { ApiError } from "../api/client";
import { formatFechaHora, formatPrecio } from "../utils/format";
import EstadoBadge from "../components/EstadoBadge";
import EmptyIcon from "../components/icons/EmptyIcon";
import "./MisTurnosPage.css";

export default function MisTurnosPage() {
  const [modo, setModo] = useState("dni");
  const [valor, setValor] = useState("");
  const [turnos, setTurnos] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  async function handleBuscar(e) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setTurnos(null);
    try {
      const resultado = await listarTurnosDePaciente({ [modo]: valor });
      setTurnos(resultado);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No pudimos buscar tus turnos.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="page">
      <h1>Mis turnos</h1>
      <p className="muted">Buscá tus turnos con tu DNI o tu email.</p>

      <form className="card buscar-form" onSubmit={handleBuscar}>
        <div className="buscar-modo">
          <label>
            <input
              type="radio"
              name="modo"
              checked={modo === "dni"}
              onChange={() => setModo("dni")}
            />
            DNI
          </label>
          <label>
            <input
              type="radio"
              name="modo"
              checked={modo === "email"}
              onChange={() => setModo("email")}
            />
            Email
          </label>
        </div>
        <div className="buscar-input-row">
          <input
            type={modo === "email" ? "email" : "text"}
            placeholder={modo === "dni" ? "Tu DNI" : "Tu email"}
            aria-label="Dato para buscar tus turnos"
            value={valor}
            onChange={(e) => setValor(e.target.value)}
            required
          />
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? "Buscando..." : "Buscar"}
          </button>
        </div>
      </form>

      {error && <div className="alert alert-error" role="alert">{error}</div>}

      {turnos && turnos.length === 0 && (
        <div className="empty-state">
          <EmptyIcon />
          <p>No encontramos turnos con esos datos.</p>
        </div>
      )}

      {turnos && turnos.length > 0 && (
        <div className="turnos-list">
          {turnos.map((t) => {
            const fecha = formatFechaHora(t.fecha_hora_inicio);
            // Mismo texto que muestra EstadoBadge.
            const estado = ESTADOS[t.estado] ?? t.estado;
            return (
              <Link
                to={`/turnos/${t.id}`}
                key={t.id}
                className="turno-item card"
                aria-label={`Turno #${t.id}, ${fecha}, ${estado}`}
              >
                <div>
                  <p className="turno-item-fecha">{fecha}</p>
                  <p className="muted">{formatPrecio(t.precio)}</p>
                </div>
                <EstadoBadge estado={t.estado} />
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}
