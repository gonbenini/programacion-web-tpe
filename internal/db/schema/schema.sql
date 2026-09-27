CREATE TABLE Usuario (
    id_usuario INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    mail VARCHAR(255) UNIQUE NOT NULL,
    contrasenia VARCHAR(255) NOT NULL,
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Mensaje (
    id_mensaje INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id_usuario INT,
    id_respuesta_a_mensaje INT,
    texto TEXT NOT NULL,
    me_gusta INT DEFAULT 0,
    CONSTRAINT fk_mensaje_usuario FOREIGN KEY (id_usuario) REFERENCES Usuario(id_usuario) ON DELETE SET NULL,
    CONSTRAINT fk_mensaje_respuesta FOREIGN KEY (id_respuesta_a_mensaje) REFERENCES Mensaje(id_mensaje) ON DELETE CASCADE
);