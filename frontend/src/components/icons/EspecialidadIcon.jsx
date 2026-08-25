const PATHS = {
  dermatologia: (
    <>
      <path d="M12 3c3 3.5 5 6.6 5 9.2A5 5 0 0 1 7 12.2C7 9.6 9 6.5 12 3Z" />
      <path d="M9.5 13c.3 1.6 1.6 2.8 3 3" strokeLinecap="round" />
    </>
  ),
  nutricion: (
    <>
      <path d="M12 8.2c-1.2-1.3-3.2-1.5-4.4-.2-1.7 1.8-1.5 5.4.6 7.9 1.4 1.6 2.7 2.4 3.8 2.4s2.4-.8 3.8-2.4c2.1-2.5 2.3-6.1.6-7.9-1.2-1.3-3.2-1.1-4.4.2Z" />
      <path d="M12 8.2c0-1.6.6-2.8 1.8-3.6" strokeLinecap="round" />
    </>
  ),
  ecografia: (
    <>
      <path d="M4 15c1.6 0 1.6-3 3.2-3s1.6 3 3.2 3 1.6-6 3.2-6 1.6 6 3.2 6 1.6-3 3.2-3" strokeLinecap="round" strokeLinejoin="round" />
    </>
  ),
  endocrinologia: (
    <>
      <circle cx="7" cy="8" r="2.1" />
      <circle cx="17" cy="8" r="2.1" />
      <circle cx="12" cy="17" r="2.1" />
      <path d="M8.8 9.3 10.5 14.7" strokeLinecap="round" />
      <path d="M15.2 9.3 13.5 14.7" strokeLinecap="round" />
    </>
  ),
};

export default function EspecialidadIcon({ especialidad, size = 18 }) {
  const path = PATHS[especialidad];
  if (!path) return null;

  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.6"
      aria-hidden="true"
    >
      {path}
    </svg>
  );
}
