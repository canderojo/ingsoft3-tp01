import { Link } from "react-router-dom";
import EmptyIcon from "../components/icons/EmptyIcon";
import "./NotFoundPage.css";

export default function NotFoundPage() {
  return (
    <div className="page not-found">
      <EmptyIcon size={56} />
      <p className="not-found-codigo">404</p>
      <h1>No encontramos esta página</h1>
      <p className="muted">
        Puede que el link esté roto o que la página se haya movido.
      </p>
      <Link to="/" className="btn btn-primary">
        Volver a profesionales
      </Link>
    </div>
  );
}
