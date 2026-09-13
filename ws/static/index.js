const formulario = document.querySelector('form');
const mensajePantalla = document.querySelector('#mensaje-alerta');

formulario.addEventListener('submit', async (event) => {
  event.preventDefault(); 

  // 1. Envolvemos FormData en URLSearchParams para que Go pueda leerlo
  const datosFormulario = new URLSearchParams(new FormData(formulario));

  try {
    const respuesta = await fetch('/register', {
      method: 'POST',
      body: datosFormulario 
    });

    const resultado = await respuesta.json();

    // 2. Cambiamos a resultado.Mensaje (con M mayúscula) coincidiendo con Go
    mensajePantalla.textContent = resultado.mensaje;
    mensajePantalla.style.color = "green"; 

  } catch (error) {
    console.error("Error al conectar con el servidor:", error);
    mensajePantalla.textContent = "Hubo un error al enviar el formulario.";
    mensajePantalla.style.color = "red";
  }
});