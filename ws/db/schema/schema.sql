CREATE TABLE Usuario (
    id_usuario INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    mail VARCHAR(255) UNIQUE NOT NULL,
    contrasenia VARCHAR(255) NOT NULL,
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Tablero (
    id_tablero INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    descripcion TEXT,
    id_tablero_padre INT,
    CONSTRAINT fk_tablero_padre FOREIGN KEY (id_tablero_padre) REFERENCES Tablero(id_tablero) ON DELETE CASCADE
);

CREATE TABLE Discusion (
    id_discusion INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    id_usuario_creador INT,
    id_tablero INT,
    id_primer_mensaje INT,
    CONSTRAINT fk_discusion_usuario FOREIGN KEY (id_usuario_creador) REFERENCES Usuario(id_usuario) ON DELETE SET NULL,
    CONSTRAINT fk_discusion_tablero FOREIGN KEY (id_tablero) REFERENCES Tablero(id_tablero) ON DELETE CASCADE
);

CREATE TABLE Mensaje (
    id_mensaje INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id_usuario INT,
    id_respuesta_a_mensaje INT,
    texto TEXT NOT NULL,
    me_gusta INT DEFAULT 0,
    id_discusion INT,
    CONSTRAINT fk_mensaje_usuario FOREIGN KEY (id_usuario) REFERENCES Usuario(id_usuario) ON DELETE SET NULL,
    CONSTRAINT fk_mensaje_respuesta FOREIGN KEY (id_respuesta_a_mensaje) REFERENCES Mensaje(id_mensaje) ON DELETE CASCADE,
    CONSTRAINT fk_mensaje_discusion FOREIGN KEY (id_discusion) REFERENCES Discusion(id_discusion) ON DELETE CASCADE
);

ALTER TABLE Discusion ADD CONSTRAINT fk_discusion_primer_mensaje FOREIGN KEY (id_primer_mensaje) REFERENCES Mensaje(id_mensaje) ON DELETE SET NULL;