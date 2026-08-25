import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { obtenerProfesional, horariosDisponibles } from "../api/profesionales";
import { crearTurno } from "../api/turnos";
import { ApiError } from "../api/client";
import { formatHora, formatPrecio, formatHoraSlot, hoyISO, claveSlot } from "../utils/format";
import EspecialidadBadge from "../components/EspecialidadBadge";
import "./ProfesionalDetailPage.css";

// Arma la grilla completa de horarios del día (ocupados incluidos) a partir
// del horario de atención del profesional, marcando cuáles siguen libres
// según lo que devolvió la API. Así el turno ya tomado se ve gris en vez de
// simplemente desaparecer de la lista.
function armarGrillaDelDia(profesional, fecha, disponibles) {
  const [horaInicio, minInicio] = profesional.hora_inicio_atencion.split(":").map(Number);
  const [horaFin, minFin] = profesional.hora_fin_atencion.split(":").map(Number);
  const duracion = profesional.duracion_turno_minutos;

  // El backend maneja las horas como reloj de Argentina "sin huso" (las
  // manda con sufijo Z pero sin convertir). Armamos la grilla en UTC
  // también, para comparar dígito a dígito sin que el navegador aplique
  // una conversión de huso horario que no corresponde.
  const inicio = new Date(`${fecha}T00:00:00.000Z`);
  inicio.setUTCHours(horaInicio, minInicio, 0, 0);
  const fin = new Date(`${fecha}T00:00:00.000Z`);
  fin.setUTCHours(horaFin, minFin, 0, 0);

  const disponiblesClaves = new Set(disponibles.map(claveSlot));
  const slots = [];
  for (let t = new Date(inicio); t < fin; t.setUTCMinutes(t.getUTCMinutes() + duracion)) {
    const iso = new Date(t).toISOString();
    slots.push({ iso, disponible: disponiblesClaves.has(claveSlot(iso)) });
  }
  return slots;
}

const FORM_INICIAL = { nombre: "", dni: "", email: "", telefono: "" };

export default function ProfesionalDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [profesional, setProfesional] = useState(null);
  const [errorProfesional, setErrorProfesional] = useState(null);

  const [fecha, setFecha] = useState(hoyISO());
  const [horarios, setHorarios] = useState([]);
  const [loadingHorarios, setLoadingHorarios] = useState(false);
  const [errorHorarios, setErrorHorarios] = useState(null);

  const [slotSeleccionado, setSlotSeleccionado] = useState(null);
  const [form, setForm] = useState(FORM_INICIAL);
  const [enviando, setEnviando] = useState(false);
  const [errorReserva, setErrorReserva] = useState(null);

  const grillaDelDia = useMemo(
    () => (profesional ? armarGrillaDelDia(profesional, fecha, horarios) : []),
    [profesional, fecha, horarios]
  );

  useEffect(() => {
    obtenerProfesional(id)
      .then(setProfesional)
      .catch((err) =>
        setErrorProfesional(err instanceof ApiError ? err.message : "No pudimos cargar este profesional.")
      );
  }, [id]);

  useEffect(() => {
    setSlotSeleccionado(null);
    setLoadingHorarios(true);
    setErrorHorarios(null);
    horariosDisponibles(id, fecha)
      .then(setHorarios)
      .catch((err) =>
        setErrorHorarios(err instanceof ApiError ? err.message : "No pudimos cargar los horarios.")
      )
      .finally(() => setLoadingHorarios(false));
  }, [id, fecha]);

  async function handleReservar(e) {
    e.preventDefault();
    setErrorReserva(null);

    if (!/^\d{8}$/.test(form.dni)) {
      setErrorReserva("El DNI debe tener 8 dígitos, sin puntos ni espacios.");
      return;
    }

    setEnviando(true);
    try {
      const turno = await crearTurno({
        profesional_id: Number(id),
        fecha_hora_inicio: slotSeleccionado,
        ...form,
      });
      navigate(`/turnos/${turno.id}`, { state: { recienCreado: true } });
    } catch (err) {
      setErrorReserva(err instanceof ApiError ? err.message : "No pudimos crear el turno.");
      setSlotSeleccionado(null);
      // alguien más se quedó con el horario justo antes que nosotros: refrescamos
      // la grilla para que se vea gris en vez de dejarla desactualizada
      horariosDisponibles(id, fecha).then(setHorarios).catch(() => {});
    } finally {
      setEnviando(false);
    }
  }

  if (errorProfesional) {
    return (
      <div className="page">
        <div className="alert alert-error">{errorProfesional}</div>
        <Link to="/">Volver a profesionales</Link>
      </div>
    );
  }

  if (!profesional) {
    return (
      <div className="page">
        <div className="spinner" />
      </div>
    );
  }

  return (
    <div className="page">
      <Link to="/" className="back-link">
        ← Profesionales
      </Link>

      <div className="profesional-detail-header">
        <h1>{profesional.nombre}</h1>
        <EspecialidadBadge especialidad={profesional.especialidad} />
      </div>
      <p className="muted">
        Atiende de {formatHora(profesional.hora_inicio_atencion)} a{" "}
        {formatHora(profesional.hora_fin_atencion)} · Consultas de{" "}
        {profesional.duracion_turno_minutos} min
      </p>
      <p className="profesional-detail-precio">{formatPrecio(profesional.precio_consulta)}</p>

      <div className="card">
        <div className="field" style={{ maxWidth: 220 }}>
          <label htmlFor="fecha">Elegí una fecha</label>
          <input
            id="fecha"
            type="date"
            min={hoyISO()}
            value={fecha}
            onChange={(e) => setFecha(e.target.value)}
          />
        </div>

        {loadingHorarios && <div className="spinner" />}
        {errorHorarios && <div className="alert alert-error">{errorHorarios}</div>}

        {!loadingHorarios && !errorHorarios && grillaDelDia.length === 0 && (
          <p className="muted">Este profesional no atiende ese día.</p>
        )}

        {!loadingHorarios && !errorHorarios && grillaDelDia.length > 0 && (
          <div className="slots-grid">
            {grillaDelDia.map(({ iso, disponible }) => (
              <button
                key={iso}
                type="button"
                disabled={!disponible}
                title={disponible ? undefined : "Ya reservado"}
                className={`slot-btn${slotSeleccionado === iso ? " selected" : ""}${
                  disponible ? "" : " slot-btn-ocupado"
                }`}
                onClick={() => setSlotSeleccionado(iso)}
              >
                {formatHoraSlot(iso)}
              </button>
            ))}
          </div>
        )}
      </div>

      {slotSeleccionado && (
        <form className="card reserva-form" onSubmit={handleReservar}>
          <h3>Completá tus datos para reservar</h3>
          <p className="muted">
            Turno el {formatHoraSlot(slotSeleccionado)} hs — {fecha}
          </p>

          <div className="field">
            <label htmlFor="nombre">Nombre y apellido</label>
            <input
              id="nombre"
              required
              value={form.nombre}
              onChange={(e) => setForm({ ...form, nombre: e.target.value })}
            />
          </div>
          <div className="field">
            <label htmlFor="dni">DNI</label>
            <input
              id="dni"
              required
              inputMode="numeric"
              maxLength={8}
              placeholder="Sin puntos, 8 dígitos"
              value={form.dni}
              onChange={(e) => setForm({ ...form, dni: e.target.value.replace(/\D/g, "").slice(0, 8) })}
            />
          </div>
          <div className="field">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              required
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
            />
          </div>
          <div className="field">
            <label htmlFor="telefono">Teléfono (opcional)</label>
            <input
              id="telefono"
              value={form.telefono}
              onChange={(e) => setForm({ ...form, telefono: e.target.value })}
            />
          </div>

          {errorReserva && <div className="alert alert-error">{errorReserva}</div>}

          <button type="submit" className="btn btn-primary" disabled={enviando}>
            {enviando ? "Reservando..." : "Confirmar reserva"}
          </button>
        </form>
      )}
    </div>
  );
}
