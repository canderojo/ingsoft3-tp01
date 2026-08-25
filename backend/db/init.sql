-- Script de creación de la base de datos del turnero.
-- Se corre una sola vez (o cada vez que recreás la base desde cero).
--
-- Nota para principiante: en vez de un "enum" nativo de Postgres para
-- especialidad/estado, usamos TEXT + CHECK. Es más simple de leer y de
-- modificar (agregar un valor nuevo es un ALTER TABLE chico) que un
-- tipo ENUM de Postgres, que requiere migraciones más pesadas.

-- SERIAL = un entero que se autoincrementa solo (1, 2, 3, ...). Es el
-- equivalente en Postgres de un "autoincrement" / identity column.
CREATE TABLE profesionales (
    id                      SERIAL PRIMARY KEY,
    nombre                  TEXT NOT NULL,
    especialidad            TEXT NOT NULL
        CHECK (especialidad IN ('dermatologia', 'nutricion', 'ecografia', 'endocrinologia')),
    hora_inicio_atencion    TIME NOT NULL,
    hora_fin_atencion       TIME NOT NULL,
    duracion_turno_minutos  INT NOT NULL CHECK (duracion_turno_minutos > 0),
    precio_consulta         NUMERIC(10, 2) NOT NULL CHECK (precio_consulta >= 0)
);

CREATE TABLE pacientes (
    id       SERIAL PRIMARY KEY,
    nombre   TEXT NOT NULL,
    dni      TEXT NOT NULL UNIQUE,
    email    TEXT NOT NULL UNIQUE,
    telefono TEXT
);

CREATE TABLE turnos (
    id                 SERIAL PRIMARY KEY,
    -- REFERENCES crea una foreign key: Postgres no deja insertar un
    -- turno con un profesional_id que no exista en la tabla profesionales.
    profesional_id     INT NOT NULL REFERENCES profesionales(id),
    paciente_id        INT NOT NULL REFERENCES pacientes(id),
    fecha_hora_inicio  TIMESTAMP NOT NULL,
    fecha_hora_fin     TIMESTAMP NOT NULL,
    estado             TEXT NOT NULL DEFAULT 'pendiente'
        CHECK (estado IN ('pendiente', 'confirmado', 'cancelado', 'completado')),
    precio             NUMERIC(10, 2) NOT NULL CHECK (precio >= 0)
);

-- Índices para acelerar las consultas que hacemos seguido: buscar los
-- turnos de un profesional o de un paciente en un rango de fechas
-- (esto es justo lo que necesitamos para detectar superposiciones).
CREATE INDEX idx_turnos_profesional_fecha ON turnos (profesional_id, fecha_hora_inicio);
CREATE INDEX idx_turnos_paciente_fecha ON turnos (paciente_id, fecha_hora_inicio);

-- Datos de ejemplo para poder probar la app sin cargar todo a mano.
INSERT INTO profesionales (nombre, especialidad, hora_inicio_atencion, hora_fin_atencion, duracion_turno_minutos, precio_consulta) VALUES
    ('Dra. Laura Gómez',      'dermatologia',    '09:00', '17:00', 30, 8000.00),
    ('Lic. Marina Ruiz',      'nutricion',       '08:00', '14:00', 45, 6500.00),
    ('Dr. Pablo Sánchez',     'ecografia',       '10:00', '18:00', 20, 9500.00),
    ('Dra. Carla Fernández',  'endocrinologia',  '09:00', '13:00', 30, 8500.00),
    -- Un segundo profesional por especialidad, con horarios bien distintos
    -- (turno vespertino, duraciones cortas/largas), para tener variedad
    -- real al probar filtros y reservas.
    ('Dr. Nicolás Herrera',   'dermatologia',    '14:00', '19:00', 25, 7200.00),
    ('Lic. Valentina Suárez', 'nutricion',       '13:00', '20:00', 30, 5800.00),
    ('Dra. Agustina López',   'ecografia',       '07:00', '12:00', 30, 10200.00),
    ('Dr. Tomás Ibáñez',      'endocrinologia',  '10:00', '16:00', 40, 9000.00);
