import { ESTADOS } from "../api/turnos";

export default function EstadoBadge({ estado }) {
  return (
    <span
      className="badge"
      style={{
        color: `var(--color-${estado})`,
        background: `var(--color-${estado}-bg)`,
      }}
    >
      {ESTADOS[estado] ?? estado}
    </span>
  );
}
