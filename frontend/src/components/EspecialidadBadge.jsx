import { ESPECIALIDADES } from "../api/profesionales";
import EspecialidadIcon from "./icons/EspecialidadIcon";

const COLOR_VAR = {
  dermatologia: "--color-dermatologia",
  nutricion: "--color-nutricion",
  ecografia: "--color-ecografia",
  endocrinologia: "--color-endocrinologia",
};

export default function EspecialidadBadge({ especialidad }) {
  const label =
    ESPECIALIDADES.find((e) => e.value === especialidad)?.label ?? especialidad;
  const colorVar = COLOR_VAR[especialidad] ?? "--color-secondary";

  return (
    <span
      className="badge"
      style={{
        color: `var(${colorVar})`,
        background: `color-mix(in srgb, var(${colorVar}) 15%, white)`,
      }}
    >
      <EspecialidadIcon especialidad={especialidad} size={14} />
      {label}
    </span>
  );
}
