const formulario = document.querySelector('form');
const mensajePantalla = document.querySelector('#mensaje-alerta');

// Función que recibe el mensaje del servidor e imprime un mensaje con contenido HTML.
formulario.addEventListener('submit', async (event) => {
  event.preventDefault(); 

  const datosFormulario = new URLSearchParams(new FormData(formulario));

  try {
    const respuesta = await fetch('/register', {
      method: 'POST',
      body: datosFormulario 
    });

    const resultado = await respuesta.json();

    mensajePantalla.textContent = resultado.mensaje;
    mensajePantalla.style.color = "green"; 

  } catch (error) {
    console.error("Error al conectar con el servidor:", error);
    mensajePantalla.textContent = "Hubo un error al enviar el formulario.";
    mensajePantalla.style.color = "red";
  }
});