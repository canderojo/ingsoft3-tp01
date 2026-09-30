package service

import "time"

// AnticipacionMinima es cuánto antes de su hora de inicio un turno deja
// de poder reservarse: un horario de las 19:30 se ofrece y se acepta
// hasta las 19:20 inclusive.
const AnticipacionMinima = 10 * time.Minute

// zonaArgentina es UTC-3 fijo. Argentina no usa horario de verano, así
// que alcanza con un corrimiento fijo y no hace falta la base de husos
// horarios del sistema (que la imagen alpine no trae).
var zonaArgentina = time.FixedZone("ART", -3*60*60)

// AhoraEnArgentina es el reloj que main.go le pasa al servicio en la app
// real. Devuelve la hora actual de Argentina con la misma convención que
// usa toda la app: los turnos se guardan y viajan con la hora de reloj
// de Argentina "etiquetada" como UTC (por ejemplo, las 19:30 de acá
// viajan como "19:30Z"). Si se usara time.Now() a secas, se compararía
// esa hora de reloj contra la hora UTC real, que va 3 horas adelantada,
// y los turnos de la tarde de hoy aparecerían como ya pasados.
func AhoraEnArgentina() time.Time {
	return horaDeRelojArgentina(time.Now())
}

// horaDeRelojArgentina pasa un instante cualquiera a la hora de reloj de
// Argentina, etiquetada como UTC. Está separada de AhoraEnArgentina para
// poder testearla con un instante fijo.
func horaDeRelojArgentina(t time.Time) time.Time {
	local := t.In(zonaArgentina)
	return time.Date(local.Year(), local.Month(), local.Day(),
		local.Hour(), local.Minute(), local.Second(), local.Nanosecond(), time.UTC)
}
