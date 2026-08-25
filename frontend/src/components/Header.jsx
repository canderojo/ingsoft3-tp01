import { useEffect, useState } from "react";
import { NavLink } from "react-router-dom";
import "./Header.css";

export default function Header() {
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    function onScroll() {
      setScrolled(window.scrollY > 4);
    }
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <header className={`header${scrolled ? " header-scrolled" : ""}`}>
      <div className="header-inner">
        <NavLink to="/" className="header-brand">
          <span className="header-brand-mark">
            <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="M12 4v16M4 12h16" stroke="#fff" strokeWidth="2.4" strokeLinecap="round" />
            </svg>
          </span>
          Centro de Salud de la Mujer
        </NavLink>
        <nav className="header-nav">
          <NavLink
            to="/"
            end
            className={({ isActive }) => (isActive ? "header-link active" : "header-link")}
          >
            Profesionales
          </NavLink>
          <NavLink
            to="/mis-turnos"
            className={({ isActive }) => (isActive ? "header-link active" : "header-link")}
          >
            Mis turnos
          </NavLink>
        </nav>
      </div>
    </header>
  );
}
