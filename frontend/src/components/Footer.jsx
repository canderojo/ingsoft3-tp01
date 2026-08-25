import "./Footer.css";

function PhoneIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92Z" />
    </svg>
  );
}

function MailIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <rect x="2" y="4" width="20" height="16" rx="2" />
      <path d="M22 6 12 13 2 6" />
    </svg>
  );
}

function MapPinIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0Z" />
      <circle cx="12" cy="10" r="3" />
    </svg>
  );
}

function ClockIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <circle cx="12" cy="12" r="10" />
      <path d="M12 6v6l4 2" />
    </svg>
  );
}

export default function Footer() {
  return (
    <footer className="footer">
      <div className="footer-inner">
        <div className="footer-brand">
          <span className="footer-brand-mark">
            <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="M12 4v16M4 12h16" stroke="#fff" strokeWidth="2.4" strokeLinecap="round" />
            </svg>
          </span>
          <div>
            <p className="footer-brand-nombre">Centro de Salud de la Mujer</p>
            <p className="footer-brand-tagline">Turnos online, sin vueltas.</p>
          </div>
        </div>

        <div className="footer-col">
          <p className="footer-col-title">Contacto</p>
          <ul>
            <li>
              <PhoneIcon />
              <span>(0351) 400-1234</span>
            </li>
            <li>
              <MailIcon />
              <span>turnos@centrosaluddelamujer.com.ar</span>
            </li>
            <li>
              <MapPinIcon />
              <span>Nueva Córdoba, Córdoba</span>
            </li>
          </ul>
        </div>

        <div className="footer-col">
          <p className="footer-col-title">Horarios</p>
          <ul>
            <li>
              <ClockIcon />
              <span>Lunes a viernes</span>
            </li>
            <li className="footer-col-nota">El horario varía según el profesional.</li>
          </ul>
        </div>
      </div>
      <p className="footer-nota">
        Proyecto académico — Ingeniería del Software 3, UCC.
      </p>
    </footer>
  );
}
